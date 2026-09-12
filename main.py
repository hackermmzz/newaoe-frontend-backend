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
import grpc
import ExportExcel
import Util
##################################################创建一个全局的channel
channel=None
if not DebugLocal:
    channel = grpc.insecure_channel(
        GRPCHost,
        # 关键参数：自动重连配置（断了自动重试）
        options=[
            ('grpc.max_receive_message_length', 1024 * 1024 * 100),  # 100MB
            ('grpc.max_send_message_length', 1024 * 1024 * 100),
            ('grpc.keepalive_time_ms', 10000),    # 每10秒发心跳
            ('grpc.keepalive_timeout_ms', 5000),   # 心跳超时5秒
            ('grpc.keepalive_permit_without_calls', True),  # 无请求也保活
            ('grpc.http2.max_pings_without_data', 0),
        ]
    )
#################################################任务流程
def TaskProcess():
    #获取服务
    server=None
    if not DebugLocal:
        server=protoc_pb2_grpc.CodeStub(channel)
    #首先获取代码
    res0=GetOneStudentCode(server)
    if not res0["ok"]:
        if not DebugLocal:
            Log(res0["msg"])
            time.sleep(JudgeSleepTimeWhenGetCodeFailed)#休眠一段时间
            return False
        else:
            #说明全测完了
            return True
    #告诉服务器处于编译状态
    id=res0["id"]
    indices=res0["indices"]
    Log(f"{id}/{indices}/成功获取代码!")
    if not DebugLocal:
        PostRunStatus(server=server,data=protoc_pb2.CodeStatusUpdateRequest(
            auth=GRPCAuth,
            indices=indices,
            id=id,
            status=PostRunStatusEnum.Code_Status_Compile.value,
            data=""
            )).Response()
    #创建运行目录和编译目录
    rundir=f'{RunDir}/{res0["id"]}_{res0["indices"]}' if indices !="" else f'{RunDir}/{res0["id"]}'
    buildDir=f"{rundir}/build"
    os.makedirs(rundir,exist_ok=True)
    os.makedirs(buildDir,exist_ok=True)
    #保存代码为文件到编译目录
    with open(f"{buildDir}/UsrAI.h","w",errors="replace")as f:
        f.write(res0["header"])
    with open(f"{buildDir}/UsrAI.cpp","w",errors="replace") as f:
        f.write(res0["source"])
    #编译代码
    logfile=f"{buildDir}/{CompileLogFileName}"
    res1=CodeCompile(buildDir,logfile)
    compileError=not res1[0]
    if compileError:
        Log(f"{id}/{indices}/编译失败!")
        if not DebugLocal:
            PostRunStatus(server=server,data=protoc_pb2.CodeStatusUpdateRequest(
                auth=GRPCAuth,
                indices=indices,
                id=id,
                status=PostRunStatusEnum.Code_Status_Compile_Error.value,
                data=res1[1]
                )).Response()
            return False
    else:
        Log(f"{id}/{indices}/编译成功!")
        if not DebugLocal:
            PostRunStatus(server=server,
                        data=protoc_pb2.CodeStatusUpdateRequest(
                            auth=GRPCAuth,
                            indices=indices,
                            id=id,
                            status=PostRunStatusEnum.Code_Status_Compile_Success.value,
                            data=""  
                )).Response()
    if not compileError:
        #将可执行文件移动到运行目录
        shutil.move(f"{buildDir}/newAOE",f"{rundir}/newAOE")
        #创建运行日志文件和运行实时结果文件
        with open(f"{rundir}/{RunLogFileName}","w",errors="replace")as f:
            pass
        with open(f"{rundir}/{RunResultFileName}","w",errors="replace")as f:
            pass
        #运行代码
        Log(f"{res0['id']}/{res0['indices']}/开始运行!")
        CodeRun(res0["id"],res0["indices"],rundir,f"{rundir}/{RunLogFileName}",f"{rundir}/{RunResultFileName}",server)
        Log(f"{res0['id']}/{res0['indices']}/运行结束!")
    #导出到excel表
    if DebugLocal:
        linesStr=[]
        linesData=[]
        if compileError:
            linesStr=['{"food":50,"frame":4288,"gold":150,"score":16,"stone":0,"time":171520,"win":false,"wood":20,"compile":false}']
        else:
            #读取运行实时结果文件
            with open(f"{rundir}/{RunResultFileName}","r") as f:
                linesStr=[line.strip() for line in f if line.strip()]
        for line in linesStr:
            try:
                js=json.loads(line)
                if "compile" not in js:
                    js["compile"]=True
                linesData.append(js)
            except:
                pass
        #导出到excel表
        ExportExcel.ExportExcel(id,linesData)
    return False
#################################################任务逻辑
def Task():
    #
    while True:
        try:
            fine=TaskProcess()
            if fine:
                break
        except Exception as e:
            print(f"出现异常:{e}")
##################################################静态编译
def PreCompile():
    docker_cmd = [
        "docker", "run",
        "--rm",
        "-v", f"{new_aoe_folder}:/app",
        "-w", "/app",
        new_aoe_docker_img,

        "bash", "-c",
        r'''
        # 给所有 .h / .cpp 文件建立“首字母大小写相反”的软链接
        find . -type f \( -name "*.h" -o -name "*.cpp" \) -print0 |
        while IFS= read -r -d '' file; do
            dir=$(dirname "$file")
            name=$(basename "$file")

            first="${name:0:1}"
            rest="${name:1}"

            # 首字母小写 -> 大写
            if [[ "$first" =~ [a-z] ]]; then
                new_first=$(echo "$first" | tr '[:lower:]' '[:upper:]')

            # 首字母大写 -> 小写
            elif [[ "$first" =~ [A-Z] ]]; then
                new_first=$(echo "$first" | tr '[:upper:]' '[:lower:]')

            else
                continue
            fi

            link_name="${new_first}${rest}"
            link_path="$dir/$link_name"

            # 如果目标名字不存在，才创建软链接
            if [ ! -e "$link_path" ] && [ ! -L "$link_path" ]; then
                echo "创建软链接: $link_path -> $name"
                ln -s "$name" "$link_path"
            fi
        done

        export PATH="$PATH:/opt/qt5.9.2/bin/"

        qmake &&
        make -j$(nproc) &&
        rm -f UsrAI.o newAOE
        '''
    ]

    result = subprocess.run(
        docker_cmd,
        stdout=subprocess.PIPE,
        stderr=LogFile,
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
    #
    for i in range(Epochs):
        print(f"第{i+1}/{Epochs}轮测试开始...")
        if not IsSingleThreadTest:
            #运行多个线程执行任务
            with ThreadPoolExecutor(max_workers=JudgeMaxPayLoad) as executor:
                for _ in range(JudgeMaxPayLoad):
                    executor.submit(Task)
        else:
                Task()
        #压缩一轮数据
        dir=f"data{i+1}"
        os.makedirs(dir,exist_ok=True)
        shutil.move(f"{RunDir}", f"{dir}/")
        shutil.move(f"result.xlsx", f"{dir}/")
        import CodeGet
        CodeGet.idx=0
    