import os
from CodeGet import base_url,header
import subprocess
from CodeCompile import new_aoe_folder,new_aoe_docker_img
from config import *
import asyncio
import aiofiles
import json
#
def CodeRun(id:str,indices:int,dir:str,logFile:str,resultFile:str):
    '''
    运行 实时结果上传
    '''
    workdir="/tmp/project"
    #代挂在文件/目录
    fileMap={
        f"{new_aoe_folder}/res":f"{workdir}/res:ro",#挂载资源路径
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
        ./newAOE --offscreen --exam --freq=8 --ResultLogFile={RunResultFileName}
        ''' 
    ]
    #异步运行
    async def CodeRunProcess(logFile,docker_cmd):
        async with aiofiles.open(logFile, "w") as log:
            process = await asyncio.create_subprocess_exec(
                *docker_cmd,
                stdout=log,
                stderr=log
            )
        return process
    async def CodeStatusPost(resultFile,process,id,indices):
        pre_line=""
        frame=0
        async with aiofiles.open(resultFile, "r") as result:
            while True:
                #跳到最后一行读取数据
                await result.seek(0, os.SEEK_END)  # 跳到文件末尾
                line = await result.readline()  # 读取最后一行
                line=line.strip()
                if line and line != pre_line:
                    frame+=1
                    #上传数据
                    pre_line=line
                    data={"frame": frame,"data": line}
                    PostRunStatus(id=id,indices=indices,status=PostRunStatusEnum.Code_Status_Running,data=data).Response()
                    await asyncio.sleep(CodeRunStatusUploadInterval)
    
                #判断进程是否结束
                if process.returncode is not None:
                    break
    async def run(logFile,docker_cmd,resultFile,id,indices):
        process=await CodeRunProcess(logFile,docker_cmd)
        await CodeStatusPost(resultFile,process,id,indices)
        return process

    #运行
    process=asyncio.run(run(logFile,docker_cmd,resultFile,id,indices))
    #记录return code
    with open(logFile,"a") as f:
        f.write(f"\nProcess exited with return code {process.returncode}\n")
    #把最后一次记录的数据上传给服务器（如果可以的话）
    with open(resultFile, "r") as result:
        result.seek(0, os.SEEK_END)  # 跳到文件末尾
        last_line=result.readline()  # 读取最后一行
        last_line=last_line.strip()
        if last_line:
            status=PostRunStatusEnum.Code_Status_Success if js["status"] else PostRunStatusEnum.Code_Status_Fail
            js=json.loads(last_line)
            data={"frame": 200501190650,"data": last_line}
            PostRunStatus(id=id,indices=indices,status=status,data=data).Response()
        else:
            data={"frame": 200501190650,"data": "程序崩溃"}
            PostRunStatus(id=id,indices=indices,status=PostRunStatusEnum.Code_Status_Crash,data=data).Response()