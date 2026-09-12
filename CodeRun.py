import os
from CodeGet import base_url,header
import subprocess
from CodeCompile import new_aoe_folder,new_aoe_docker_img
from config import *
import asyncio
import aiofiles
import json
import protoc_pb2_grpc
import protoc_pb2

#判断字符串是否为json可以解析的字符串
def checkStrIsJson(s:str)->bool:
    try:
        json.loads(s)
        return True
    except (json.JSONDecodeError, ValueError, TypeError):
        return False
#
def CodeRun(
    id: str,
    indices: int,
    dir: str,
    logFile: str,
    resultFile: str,
    server: protoc_pb2_grpc.CodeStub
):
    """
    运行程序，并实时上传结果。

    判定规则：
    1. returncode != 0
       -> Crash
       -> 在最后一条 JSON 中加入：
          crash_code
          crash_reason

    2. returncode == 0，并且有最终 JSON
       -> 根据 win 判断 Success / Fail

    3. returncode == 0，但是没有最终 JSON
       -> Crash
    """

    workdir = "/tmp/project"

    # ============================================================
    # 挂载文件
    # ============================================================

    fileMap = {
        f"{new_aoe_folder}/res.rcc":
            f"{workdir}/res.rcc",

        f"{new_aoe_folder}/config.json":
            f"{workdir}/config.json:ro",

        f"{dir}/newAOE":
            f"{workdir}/newAOE:ro",

        f"{resultFile}":
            f"{workdir}/{RunResultFileName}",
    }

    mountArg = []

    for localP, containerP in fileMap.items():
        mountArg.append("-v")
        mountArg.append(f"{localP}:{containerP}")

    # 所有地图文件
    for x in MapFiles:
        mountArg.append("-v")
        mountArg.append(
            f"{new_aoe_folder}/{x}:{workdir}/{x}:ro"
        )

    # ============================================================
    # Docker 命令
    # ============================================================

    docker_cmd = [
        "docker", "run",

        "--rm",

        # 内存限制
        "-m", RunMemoryLimit,

        # CPU 限制
        "--cpus", RunCPULimit,

        # /tmp 磁盘限制
        "--tmpfs", f"/tmp:rw,size={RunDiskLimit}",

    ] + mountArg + [

        # 工作目录
        "-w", workdir,

        # Docker 镜像
        new_aoe_docker_img,

        "bash",
        "-c",

        f"""
        export QT_QPA_PLATFORM=offscreen &&
        export LD_LIBRARY_PATH=/opt/qt5.9.2/lib:$LD_LIBRARY_PATH &&
        ./newAOE \
            --offscreen \
            --exam \
            --freq=MAX \
            --record \
            --ResultLogFile={RunResultFileName}
        """
    ]

    # ============================================================
    # 根据退出码获取崩溃原因
    # ============================================================

    def GetExitReason(code: int) -> str:

        exitReasons = {
            0: "Normal Exit",

            1: "General Error",

            125: "Docker Run Error",

            126: "Command Cannot Execute",

            127: "Command Not Found",

            130: "SIGINT",

            132: "SIGILL (Illegal Instruction)",

            134: "SIGABRT (Abort / Assertion Failed)",

            136: "SIGFPE (Arithmetic Exception / Divide By Zero)",

            137: "SIGKILL (Possibly OOM / Memory Limit)",

            139: "SIGSEGV (Segmentation Fault)",

            143: "SIGTERM",
        }

        return exitReasons.get(
            code,
            f"Unknown Exit Code ({code})"
        )

    # ============================================================
    # 构造 Crash JSON
    # ============================================================

    def BuildCrashData(
        last_post_data: str,
        returnCode: int
    ) -> str:

        crashReason = GetExitReason(returnCode)

        try:
            # 如果之前已经存在合法 JSON
            if (
                last_post_data
                and
                checkStrIsJson(last_post_data)
            ):
                crashData = json.loads(last_post_data)

            else:
                # 如果程序刚启动就崩了，没有任何结果
                crashData = {}

            # 添加崩溃信息
            crashData["crash_code"] = returnCode
            crashData["crash_reason"] = crashReason

            return json.dumps(
                crashData,
                ensure_ascii=False
            )

        except Exception as e:

            print(
                f"[CodeRun] BuildCrashData failed: {e}"
            )

            # 即使原 JSON 有问题，也保证 crash 信息能上传
            return json.dumps(
                {
                    "crash_code": returnCode,
                    "crash_reason": crashReason
                },
                ensure_ascii=False
            )

    # ============================================================
    # 启动 Docker
    # ============================================================

    async def CodeRunProcess(
        logFile,
        docker_cmd
    ):

        log = await aiofiles.open(
            logFile,
            "w",
            errors="replace"
        )

        try:

            process = await asyncio.create_subprocess_exec(
                *docker_cmd,
                stdout=log,
                stderr=log
            )

            return process, log

        except Exception:

            await log.close()
            raise

    # ============================================================
    # 状态上传
    # ============================================================

    async def CodeStatusPost(
        resultFile,
        process,
        id,
        indices
    ):

        # 保存尚未写完整的一行
        preInfo = ""

        # 最后一个合法 JSON
        last_post_data = ""

        async with aiofiles.open(
            resultFile,
            "r",
            errors="replace"
        ) as result:

            # ====================================================
            # 程序运行中
            # ====================================================

            while True:

                # 读取新增数据
                curInfo = await result.read()

                info = preInfo + curInfo

                allLines = info.split("\n")

                # 每一轮重新保存未完成行
                preInfo = ""

                # 去除末尾空行
                while len(allLines) != 0:

                    if allLines[-1].strip() == "":
                        allLines.pop()

                    else:
                        break

                # 最后一行可能没有写完整
                if (
                    len(allLines) >= 1
                    and
                    not checkStrIsJson(allLines[-1])
                ):

                    preInfo = allLines.pop()

                # =================================================
                # 找最后一条完整 JSON
                # =================================================

                if len(allLines) >= 1:

                    currentJson = None

                    for line in reversed(allLines):

                        line = line.strip()

                        if not line:
                            continue

                        if checkStrIsJson(line):

                            currentJson = line
                            break

                    if currentJson is not None:

                        last_post_data = currentJson

                        # =========================================
                        # 上传 Running 状态
                        # =========================================

                        if not DebugLocal:

                            PostRunStatus(
                                server=server,

                                data=protoc_pb2.CodeStatusUpdateRequest(
                                    auth=GRPCAuth,

                                    indices=indices,

                                    id=id,

                                    status=(
                                        PostRunStatusEnum
                                        .Code_Status_Running
                                        .value
                                    ),

                                    data=last_post_data
                                )
                            ).Response()

                # =================================================
                # 判断 Docker 是否已经退出
                # =================================================

                if process.returncode is not None:
                    break

                if DebugLocal:

                    await asyncio.sleep(0.05)

                else:

                    await asyncio.sleep(
                        CodeRunStatusUploadInterval
                    )

            # ====================================================
            # Docker 已退出
            # ====================================================

            await process.wait()

            returnCode = process.returncode

            # ====================================================
            # 再读取一次最后残留的数据
            # ====================================================

            leftInfo = await result.read()

            finalInfo = preInfo + leftInfo

            # ====================================================
            # 寻找最终合法 JSON
            # ====================================================

            finalJson = None

            lineParts = finalInfo.split("\n")

            while len(lineParts) != 0:

                last_line = lineParts.pop().strip()

                if not last_line:
                    continue

                if checkStrIsJson(last_line):

                    finalJson = last_line
                    break

            # 如果最后残留数据中没找到，
            # 使用运行过程中读到的最后一条 JSON
            if finalJson is None:

                if (
                    last_post_data
                    and
                    checkStrIsJson(last_post_data)
                ):

                    finalJson = last_post_data

            # ====================================================
            # 情况 1：
            # returncode != 0
            #
            # 只要非 0，就无条件 Crash
            # ====================================================

            if returnCode != 0:

                crashReason = GetExitReason(
                    returnCode
                )

                print(
                    f"[CodeRun] Program crashed: "
                    f"id={id}, "
                    f"indices={indices}, "
                    f"returncode={returnCode}, "
                    f"reason={crashReason}"
                )

                # 如果退出前又产生了一条合法 JSON，
                # 优先使用最终的那条
                if finalJson is not None:
                    last_post_data = finalJson

                # 给 JSON 添加 crash_code + crash_reason
                last_post_data = BuildCrashData(
                    last_post_data,
                    returnCode
                )

                # ===============================================
                # 上传 Crash
                # ===============================================

                if not DebugLocal:

                    PostRunStatus(
                        server=server,

                        data=protoc_pb2.CodeStatusUpdateRequest(
                            auth=GRPCAuth,

                            indices=indices,

                            id=id,

                            status=(
                                PostRunStatusEnum
                                .Code_Status_Crash
                                .value
                            ),

                            data=last_post_data
                        )
                    ).Response()

                return

            # ====================================================
            # 情况 2：
            # returncode == 0
            # 说明 Docker / newAOE 正常退出
            # ====================================================

            # ----------------------------------------------------
            # 正常退出，但是没有最终 JSON
            # ----------------------------------------------------

            if finalJson is None:

                crashReason = (
                    "Program exited normally "
                    "but no final result JSON"
                )

                print(
                    f"[CodeRun] {crashReason}: "
                    f"id={id}, "
                    f"indices={indices}"
                )

                # 这种情况虽然 returncode = 0，
                # 但从判题逻辑上仍然属于异常
                crashData = {
                    "crash_code": 0,
                    "crash_reason": crashReason
                }

                last_post_data = json.dumps(
                    crashData,
                    ensure_ascii=False
                )

                if not DebugLocal:

                    PostRunStatus(
                        server=server,

                        data=protoc_pb2.CodeStatusUpdateRequest(
                            auth=GRPCAuth,

                            indices=indices,

                            id=id,

                            status=(
                                PostRunStatusEnum
                                .Code_Status_Crash
                                .value
                            ),

                            data=last_post_data
                        )
                    ).Response()

                return

            # ====================================================
            # 情况 3：
            # returncode == 0 + 有合法最终 JSON
            # ====================================================

            try:

                js = json.loads(
                    finalJson
                )

            except Exception as e:

                crashReason = (
                    f"Invalid final JSON: {e}"
                )

                print(
                    f"[CodeRun] {crashReason}"
                )

                last_post_data = json.dumps(
                    {
                        "crash_code": 0,
                        "crash_reason": crashReason
                    },
                    ensure_ascii=False
                )

                if not DebugLocal:

                    PostRunStatus(
                        server=server,

                        data=protoc_pb2.CodeStatusUpdateRequest(
                            auth=GRPCAuth,

                            indices=indices,

                            id=id,

                            status=(
                                PostRunStatusEnum
                                .Code_Status_Crash
                                .value
                            ),

                            data=last_post_data
                        )
                    ).Response()

                return

            # ====================================================
            # win 判断
            # ====================================================

            if js.get("win", False):

                status = (
                    PostRunStatusEnum
                    .Code_Status_Success
                )

            else:

                status = (
                    PostRunStatusEnum
                    .Code_Status_Fail
                )

            # ====================================================
            # 上传最终正常结果
            # ====================================================

            if not DebugLocal:

                PostRunStatus(
                    server=server,

                    data=protoc_pb2.CodeStatusUpdateRequest(
                        auth=GRPCAuth,

                        indices=indices,

                        id=id,

                        status=status.value,

                        data=finalJson
                    )
                ).Response()

    # ============================================================
    # 总运行函数
    # ============================================================

    async def run(
        logFile,
        docker_cmd,
        resultFile,
        id,
        indices
    ):

        process = None
        logHandle = None

        try:

            process, logHandle = await CodeRunProcess(
                logFile,
                docker_cmd
            )

            await CodeStatusPost(
                resultFile,
                process,
                id,
                indices
            )

            await process.wait()

            return process

        finally:

            if logHandle is not None:

                await logHandle.flush()
                await logHandle.close()
    # ============================================================
    # 开始运行
    # ============================================================
    process = asyncio.run(
        run(
            logFile,
            docker_cmd,
            resultFile,
            id,
            indices
        )
    )
    # ============================================================
    # 写最终退出日志
    # ============================================================
    exitReason = GetExitReason(
        process.returncode
    )
    with open(
        logFile,
        "a",
        errors="replace"
    ) as f:
        f.write(
            "\n"
            "========================================\n"
        )
        f.write(
            f"Process exited with return code "
            f"{process.returncode}\n"
        )
        f.write(
            f"Exit reason: "
            f"{exitReason}\n"
        )
        f.write(
            "========================================\n"
        )