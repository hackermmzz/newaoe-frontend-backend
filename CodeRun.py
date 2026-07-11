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
def CodeRun(id:str,indices:int,dir:str,logFile:str,resultFile:str,server:protoc_pb2_grpc.CodeStub):
    '''
    运行 实时结果上传
    '''
    workdir="/tmp/project"
    #代挂在文件/目录
    fileMap={
        f"{new_aoe_folder}/res.rcc":f"{workdir}/res.rcc",
        f"{new_aoe_folder}/config.json":f"{workdir}/config.json:ro",#挂在配置文件
        f"{dir}/newAOE":f"{workdir}/newAOE:ro",   #挂载可执行文件
        f"{resultFile}":f"{workdir}/{RunResultFileName}", #挂载结果文件
    }
    mountArg=[]
    for localP,containterP in fileMap.items():
        mountArg.append("-v")
        mountArg.append(f"{localP}:{containterP}")
    for x in MapFiles:#所有地图文件
        mountArg.append("-v")
        mountArg.append(f"{new_aoe_folder}/{x}:{workdir}/{x}:ro")
    #组合指令
    docker_cmd = [
        "docker", "run",
        "--rm", # 运行完自动删除容器
        "-m", RunMemoryLimit, #限制内存使用
        "--cpus", RunCPULimit, #限制CPU使用
        "--tmpfs", f"/tmp:rw,size={RunDiskLimit}", #限制磁盘使用
    ]+mountArg+[
        "-w", f"{workdir}",  # 工作目录
        new_aoe_docker_img,
        #运行
        "bash","-c",
        #注意这里我加了一个离屏渲染的参数QT_QPA_PLATFORM=offscreen，防止运行时因为没有显示设备而崩溃
        f'''
        export QT_QPA_PLATFORM=offscreen &&
        export LD_LIBRARY_PATH=/opt/qt5.9.2/lib:$LD_LIBRARY_PATH &&
        ./newAOE --offscreen --exam --freq=MAX --ResultLogFile={RunResultFileName}
        ''' 
    ]
    #异步运行
    async def CodeRunProcess(logFile,docker_cmd):
        async with aiofiles.open(logFile, "w",errors="replace") as log:
            process = await asyncio.create_subprocess_exec(
                *docker_cmd,
                stdout=log,
                stderr=log
            )
        return process
    async def CodeStatusPost(resultFile,process,id,indices):
        preInfo=""
        last_post_data=""
        async with aiofiles.open(resultFile, "r",errors="replace") as result:
            while True:
                #读取所有新增的数据
                curInfo=await result.read()
                info=preInfo+curInfo
                allLines=info.split("\n")
                #保证最后一个元素肯定不是空字符串
                while len(allLines)!=0:
                    if allLines[-1].strip()=="":
                        allLines.pop()
                    else:
                        break
                #如果只有一行，那么不管是否完整，都放到preInfo里面
                if len(allLines)>=1 and not checkStrIsJson(allLines[-1]):
                    preInfo=allLines.pop()
                #当有大于等于2行的json数据选择倒数第二个数据上传
                if len(allLines)>=1:
                    last_post_data=allLines[-1]
                    #上传数据
                    if not DebugLocal:
                        PostRunStatus(server=server,
                                    data=protoc_pb2.CodeStatusUpdateRequest(
                                        auth=GRPCAuth,
                                        indices=indices,
                                        id=id,
                                        status=PostRunStatusEnum.Code_Status_Running.value,
                                        data=last_post_data
                                    )
                                    ).Response()
                if not DebugLocal:
                    await asyncio.sleep(CodeRunStatusUploadInterval)
                #判断进程是否结束
                if process.returncode is not None:
                    break
            #把最后一次记录的数据上传给服务器（如果可以的话）
            if not DebugLocal:
                leftInfo=await result.read()
                preInfo+=leftInfo
                lienParts=preInfo.split("\n")
                flag=False
                while len(lienParts)!=0:
                    last_line=lienParts.pop()
                    if checkStrIsJson(last_line):
                        flag=True
                        js=json.loads(last_line)
                        status=PostRunStatusEnum.Code_Status_Success if js["win"] else PostRunStatusEnum.Code_Status_Fail
                        PostRunStatus(server=server,data=protoc_pb2.CodeStatusUpdateRequest(
                            auth=GRPCAuth,
                            indices=indices,
                            id=id,
                            status=status.value,
                            data=last_line
                            )).Response()
                if not flag:
                    PostRunStatus(server=server,data=protoc_pb2.CodeStatusUpdateRequest(
                        auth=GRPCAuth,
                        indices=indices,
                        id=id,
                        status=PostRunStatusEnum.Code_Status_Crash.value,
                        data=last_post_data
                        )).Response()
    #######################
    async def run(logFile,docker_cmd,resultFile,id,indices):
        process=await CodeRunProcess(logFile,docker_cmd)
        await CodeStatusPost(resultFile,process,id,indices)
        return process
    #运行
    process=asyncio.run(run(logFile,docker_cmd,resultFile,id,indices))
    #记录return code
    with open(logFile,"a",errors="replace") as f:
        f.write(f"\nProcess exited with return code {process.returncode}\n")