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
import Util

#判断字符串是否为json可以解析的字符串
def checkStrIsJson(s:str)->bool:
    try:
        json.loads(s)
        return True
    except (json.JSONDecodeError, ValueError, TypeError):
        return False
#
import asyncio
import json
import aiofiles
import protoc_pb2
import protoc_pb2_grpc


# ============================================================
# 判断退出码是否为程序自身原因退出
# ============================================================
def IsProgramSelfReason(code: int) -> bool:
    # 程序主动退出：正常退出、通用错误、命令执行失败、找不到命令、程序内部触发的abort/非法指令/除零/段错误
    self_exit_codes = {0, 1, 125, 126, 127, 132, 134, 136, 139}
    # 外部信号干掉：SIGINT、SIGKILL、SIGTERM
    external_kill_codes = {130, 137, 143}
    
    if code in self_exit_codes:
        return True
    elif code in external_kill_codes:
        return False
    else:
        # 不在定义的字典内，可按你的需求处理，这里默认返回True或者抛异常，当前返回True
        return True

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
    return exitReasons.get(code, f"Unknown Exit Code ({code})")

def CodeRun(
    id: str,
    indices: int,
    dir: str,
    crashDir: str,
    logFile: str,
    recordFile: str,
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
    container_name=f"{id}_{indices}_{Util.GetRandomStr()}"
    # ============================================================
    # 挂载文件
    # ============================================================
    fileMap = {
        f"{new_aoe_folder}/res.rcc": f"{workdir}/res.rcc",
        f"{new_aoe_folder}/config.json": f"{workdir}/config.json:ro",
        f"{dir}/newAOE": f"{workdir}/newAOE:ro",
        f"{resultFile}": f"{workdir}/{RunResultFileName}",
        f"{recordFile}": f"{workdir}/{RecordFileName}",
    }
    # 挂载所有地图文件
    mountArg = []
    for localP, containerP in fileMap.items():
        mountArg.append("-v")
        mountArg.append(f"{localP}:{containerP}")

    # 所有地图文件
    for x in MapFiles:
        mountArg.append("-v")
        mountArg.append(f"{new_aoe_folder}/{x}:{workdir}/{x}:ro")
    #读取bash脚本
    with open(f"bash/coderun.sh", "r",encoding="utf-8") as f:
        bashScript = f.read()
    bashScript = bashScript.format(
        id=id,indices=indices,
        workdir=workdir,RunResultFileName=RunResultFileName,
        RecordOutputFileName=RecordFileName,
        ).strip()
    # ============================================================
    # Docker 命令
    # ============================================================
    docker_cmd = [
        "docker", "run",
        # 保存崩溃文件
        "-v", f"{crashDir}:{workdir}/crash",
        "--name",container_name,
        # 内存限制
        "-m", RunMemoryLimit,
        "--memory-swap", RunMemoryLimit,
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
        bashScript,
    ]
    
    # ============================================================
    # 启动 Docker
    # ============================================================
    async def CodeRunProcess(logFile, docker_cmd):
        log = await aiofiles.open(logFile, "w", errors="replace")
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
    async def CodeStatusPost(resultFile, process, id, indices):
        # 保存尚未写完整的一行
        preInfo = ""
        # 最后一个合法 JSON
        last_post_data = ""

        async with aiofiles.open(resultFile, "r", errors="replace") as result:
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
                while allLines and allLines[-1].strip() == "":
                    allLines.pop()

                # 最后一行可能没有写完整
                if allLines and not checkStrIsJson(allLines[-1]):
                    preInfo = allLines.pop()

                # =================================================
                # 找最后一条完整 JSON
                # =================================================
                if allLines:
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
                        try:
                            currentJson = json.loads(currentJson)
                        except:
                            continue
                        # =========================================
                        # 上传 Running 状态
                        # =========================================
                        PostRunStatus(
                            server=server,
                            data=protoc_pb2.CodeStatusUpdateRequest(
                                auth=GRPCAuth,
                                indices=indices,
                                id=id,
                                status=PostRunStatusEnum.Code_Status_Running.value,
                                data=CodeRunStatusInfo(
                                    status=PostRunStatusEnum.Code_Status_Running.value,
                                    gold=currentJson.get("gold", 0),
                                    stone=currentJson.get("stone", 0),
                                    wood=currentJson.get("wood", 0),
                                    food=currentJson.get("food", 0),
                                    frame=currentJson.get("frame", 0),
                                    win=currentJson.get("win", False),
                                    score=currentJson.get("score", 0),
                                    data=last_post_data
                                    ).tostr()
                            )
                        ).Response()

                # =================================================
                # 判断 Docker 是否已经退出
                # =================================================
                if process.returncode is not None:
                    break

                await asyncio.sleep(CodeRunStatusUploadInterval)

            # ====================================================
            # Docker 已退出
            # ====================================================
            await process.wait()
            returnCode = process.returncode
            # ====================================================
            # 检查是否 OOM
            # ====================================================
            docker_result = subprocess.check_output(
                [
                    "docker",
                    "inspect",
                    container_name
                ]
            )
            info = json.loads(docker_result)[0]
            oom = info["State"]["OOMKilled"]
            if oom:
                returnCode = 137
            # ====================================================
            # 删除容器
            # ====================================================
            subprocess.run(
                [
                    "docker",
                    "rm",
                    container_name
                ]
            )
            
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
            while lineParts:
                last_line = lineParts.pop().strip()
                if not last_line:
                    continue
                if checkStrIsJson(last_line):
                    try:
                        finalJson = json.loads(last_line)
                        break
                    except:
                        continue

            # 如果最后残留数据中没找到，使用运行过程中读到的最后一条 JSON
            if finalJson is None:
                if last_post_data and checkStrIsJson(last_post_data):
                    try:
                        finalJson = json.loads(last_post_data)
                    except:
                        finalJson = None
            # 如果最后残留数据中也没找到，使用默认 JSON
            if finalJson is None:
                finalJson = CodeRunStatusInfo().tojson()
            # 获取崩溃文件名
            crashLogFileName = Util.GetFolerRandomFileIfExist(crashDir)
            need_Log = crashLogFileName!="" and IsProgramSelfReason(returnCode)
            # ====================================================
            # 情况 1：returncode != 0，只要非 0，就无条件 Crash
            # ====================================================
            if returnCode != 0:
                crashReason = GetExitReason(returnCode)
                Log(f"[CodeRun] Program crashed: id={id}, \
                    indices={indices}, \
                    returncode={returnCode}, \
                    reason={crashReason}")
                # 上传 Crash
                resp=PostRunStatus(
                    server=server,
                    data=protoc_pb2.CodeStatusUpdateRequest(
                        auth=GRPCAuth,
                        indices=indices,
                        id=id,
                        status=PostRunStatusEnum.Code_Status_Crash.value,
                        data=CodeRunStatusInfo(
                            status=PostRunStatusEnum.Code_Status_Crash.value,
                            gold=finalJson.get("gold", 0),
                            stone=finalJson.get("stone", 0),
                            wood=finalJson.get("wood", 0),
                            food=finalJson.get("food", 0),
                            frame=finalJson.get("frame", 0),
                            win=finalJson.get("win", False),
                            score=finalJson.get("score", 0),
                            data=json.dumps(
                                {
                                    "crash_reason": crashReason,
                                    "needlog": need_Log
                                }, 
                                ensure_ascii=False
                                )
                        ).tostr()
                    )
                ).Response()
                if need_Log and resp:
                    Util.UploadData(resp.data.encode(), Util.read_any_text(f"{crashDir}/{crashLogFileName}"))
                else:
                    Log(f"上传CrashStatus失败，错误信息:{resp.error}")
                return

            # ====================================================
            # 情况 2：returncode == 0 + 有合法最终 JSON
            # ====================================================
            # win 判断
            if finalJson.get("win", False):
                status = PostRunStatusEnum.Code_Status_Success
            else:
                status = PostRunStatusEnum.Code_Status_Fail

            # 上传最终正常结果
            resp=PostRunStatus(
                server=server,
                data=protoc_pb2.CodeStatusUpdateRequest(
                    auth=GRPCAuth,
                    indices=indices,
                    id=id,
                    status=status.value,
                    data=CodeRunStatusInfo(
                        status=status.value,
                        gold=finalJson.get("gold", 0),
                        stone=finalJson.get("stone", 0),
                        wood=finalJson.get("wood", 0),
                        food=finalJson.get("food", 0),
                        frame=finalJson.get("frame", 0),
                        win=finalJson.get("win", False),
                        score=finalJson.get("score", 0),
                        data=""
                    ).tostr()
                )
            ).Response()
            if resp :
                Util.UploadData(resp.data.encode(), Util.read_any_text(recordFile))
    
    # ============================================================
    # 定时器，超时直接Crash
    # ============================================================
    async def DockerWatcher(process: asyncio.subprocess.Process,timeout:int):
        """
        监听进程结束，超时直接Crash
        :param process: 进程对象
        """
        await asyncio.sleep(timeout)
        # 判断进程是否结束
        if process.returncode is None:
            Log(f"[CodeRun] Program timeout: id={id},indices={indices},timeout={timeout}")
            # 杀死docker容器
            subprocess.run(
                [
                    "docker",
                    "kill",
                    container_name
                ]
            )
    # ============================================================
    # 总运行函数
    # ============================================================
    async def run(logFile, docker_cmd, resultFile, id, indices):
        process = None
        logHandle = None
        try:
            process, logHandle = await CodeRunProcess(logFile, docker_cmd)
            t1 = asyncio.create_task(CodeStatusPost(resultFile, process, id, indices))
            t2 = asyncio.create_task(DockerWatcher(process, RunTimeout))
            # 等待任意一个任务完成
            done, pending = await asyncio.wait(
                [t1, t2],
                return_when=asyncio.FIRST_COMPLETED
            )
            # 把剩下还在跑的任务全部取消
            for task in pending:
                task.cancel()
                try:
                    await task
                except asyncio.CancelledError:
                    # 被取消属于正常，直接吞掉
                    pass
            
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
    exitReason = GetExitReason(process.returncode)
    with open(logFile, "a", errors="replace") as f:
        f.write("\n========================================\n")
        f.write(f"Process exited with return code {process.returncode}\n")
        f.write(f"Exit reason: {exitReason}\n")
        f.write("========================================\n")
