from CodeGet import *
from CodeRun import *
from CodeCompile import *
import shutil
from concurrent.futures import ThreadPoolExecutor
import os
import time
from config import *
import atexit
import psutil
#################################################编译结果回复

#################################################任务流程
def TaskProcess():
    #首先获取代码
    res0=GetOneStudentCode()
    if res0["status"]==False:
        if 'id' in res0:
            Log(f"{res0['id']}/{res0['indices']}/解析错误，原因:{res0['msg']}，休眠{JudgeSleepTimeWhenGetCodeFailed}秒后重试")
        else:
            Log(f"无法获取代码，原因:{res0['msg']}，休眠{JudgeSleepTimeWhenGetCodeFailed}秒后重试")
        time.sleep(JudgeSleepTimeWhenGetCodeFailed)#休眠一段时间
        return
    else:
        Log(f"{res0['id']}/{res0['indices']}/成功获取代码!")
        PostRunStatus(id=res0["id"],indices=res0["indices"],status=PostRunStatusEnum.Code_Status_ServerGet,data="成功获取代码").Response()
    #创建运行目录和编译目录
    rundir=f"{RunDir}/{res0["id"]}_{res0["indices"]}"
    buildDir=f"{rundir}/build"
    os.makedirs(rundir,exist_ok=True)
    os.makedirs(buildDir,exist_ok=True)
    #保存代码为文件到编译目录
    with open(f"{buildDir}/UsrAI.h","w")as f:
        f.write(res0["header"])
    with open(f"{buildDir}/UsrAI.cpp","w") as f:
        f.write(res0["source"])
    #编译代码
    logfile=f"{buildDir}/{CompileLogFileName}"
    res1=CodeCompile(buildDir,logfile)
    if res1[0]==False:
        Log(f"{res0['id']}/{res0['indices']}/编译失败!")
        PostRunStatus(id=res0["id"],indices=res0["indices"],status=PostRunStatusEnum.Code_Status_Compile_Error,data=res1[1]).Response()
        return
    else:
        Log(f"{res0['id']}/{res0['indices']}/编译成功!")
        PostRunStatus(id=res0["id"],indices=res0["indices"],status=PostRunStatusEnum.Code_Status_Compile_Success,data="编译成功").Response()
    #将可执行文件移动到运行目录
    shutil.move(f"{buildDir}/newAOE",f"{rundir}/newAOE")
    #创建运行日志文件和运行实时结果文件
    with open(f"{rundir}/{RunLogFileName}","w")as f:
        pass
    with open(f"{rundir}/{RunResultFileName}","w")as f:
        pass
    #运行代码
    Log(f"{res0['id']}/{res0['indices']}/开始运行!")
    CodeRun(res0["id"],res0["indices"],rundir,f"{rundir}/{RunLogFileName}",f"{rundir}/{RunResultFileName}")
    Log(f"{res0['id']}/{res0['indices']}/运行结束!")
#################################################任务逻辑
def Task():
    #
    while True:
        try:
            TaskProcess()
        except Exception as e:
            print(f"出现异常:{e}")
##################################################静态编译
def PreCompile():
    #
    docker_cmd = [
        "docker", "run",
        "--rm", # 运行完自动删除容器
        "-v", f"{new_aoe_folder}:/app",  # 挂载当前项目目录到容器 /app
        "-w", "/app",  # 工作目录
        new_aoe_docker_img,
        # 编译命令（把所有输出打出来）
        "bash","-c",
        ''' export PATH+=":/opt/qt5.9.2/bin/" && qmake && make -j$(nproc) && rm UsrAI.o newAOE '''
    ]
    result = subprocess.run(
            docker_cmd,
            stdout=subprocess.PIPE,     # 输出捕获到变量
            stderr=LogFile,                   # 报错也写文件
            text=True,
            check=False
        )
    if result.returncode != 0:
        Log(f"预编译失败，错误信息:{result.stdout}")
        exit(1)
##################################################退出清理
def CleanUp():
    Log("正在清理子进程...")
    parent = psutil.Process(os.getpid())
    for child in parent.children(recursive=True):
        Log(f"杀掉子进程: {child.pid}")
        child.kill()
    Log("清理完成!")
##################################################初始化
def Init():
    #保证退出时杀掉所有子进程
    atexit.register(CleanUp)
    #创建基础环境
    os.makedirs(RunDir,exist_ok=True)
    #静态编译一次，让docker把环境准备好
    PreCompile()
    #清空log日志
    LogFile.truncate(0)
    LogFile.seek(0)
    Log("初始化完成!")
########################################################
if __name__=="__main__":
    #初始化
    Init()
    #运行多个线程执行任务
    with ThreadPoolExecutor(max_workers=JudgeMaxPayLoad) as executor:
        for _ in range(JudgeMaxPayLoad):
            executor.submit(Task)
    