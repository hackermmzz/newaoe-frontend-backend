import os
from enum import Enum
import requests
import time
import threading
#服务器网址
base_ip="http://localhost:8080"
base_url=f"{base_ip}/api/admin"
CodeGetURL=f"{base_url}/CodeGet"
CodeStatusPostURL=f"{base_url}/CodeRunStatusPost"
#定义api_key
api_key="mmzz"
#定义请求头
header={
    "api":api_key
}
#AOE项目目录
new_aoe_folder=f"{os.getcwd()}/new-aoe"
#AOE镜像名称
new_aoe_docker_img="newaoe.img:latest"
#编译的输出日志
CompileLogFileName="CompileLog.txt"
#编译是否开启多线程编译模式
CompileUseMultiThread=True
#每个可执行程序的运行根目录
RunDir=f"{os.getcwd()}/RunDir"
#运行时日志文件
RunLogFileName="RunLog.txt"
#运行时实时作战结果输出文件
RunResultFileName="RunResult.txt"
#所有的地图文件
MapFiles=[x for x in os.listdir(new_aoe_folder) if x.endswith(".njust")]
#限制运行内存
RunMemoryLimit="256m"
#限制运行CPU
RunCPULimit="2.5"
#限制磁盘使用
RunDiskLimit="20m"
#判题机的负载
JudgeMaxPayLoad=10
#获取任务失败进行休眠时长n秒
JudgeSleepTimeWhenGetCodeFailed=5
#代码运行上传实时数据间隔n秒
CodeRunStatusUploadInterval=1
#日志线程锁
LogLock=threading.Lock()
#日志文件
LogFile=open(f"{os.getcwd()}/JudgeLog.txt","w")
###########################################################
def Log(msg):
    with LogLock:
        LogFile.write(f"[{time.strftime('%Y-%m-%d %H:%M:%S')}] Thread {threading.current_thread().ident}: {msg}"+"\n")
        LogFile.flush()
###########################################################

class PostRunStatusEnum(Enum):
	Code_Status_Wait            = 1  #在等待队列里面
	Code_Status_Running         = 2  #正在运行出结果
	Code_Status_IDLE            = 3  #等待排队
	Code_Status_Success         = 4  #运行胜利
	Code_Status_Wait_Timeout    = 5  #等待超时/出错
	Code_Status_ServerGet       = 6  #代码发送给了服务器
	Code_Status_Compile         = 7  #编译中
	Code_Status_Compile_Timeout = 8  #编译超时
	Code_Status_Compile_Error   = 9  #编译错误
	Code_Status_Compile_Success = 10 #编译成功
	Code_Stastus_Fail           = 11 #游戏失败
	Code_Status_Crash           = 12 #游戏崩溃

class PostRunStatus():
    def __init__(self,id:str,indices:int,status:PostRunStatusEnum,data):
        self.id=id
        self.indices=indices
        self.status=status
        self.data=data
    
    def Response(self)->bool:
        js={
            "id":self.id,
            "indices":self.indices,
            "status":self.status.value,
            "data":self.data
        }
        try:
            resp=requests.post(url=CodeStatusPostURL,headers=header,json=js)
            if resp.status_code!=200:
                Log(f"{self.id}/{self.indices}/后端异常!")
        except Exception as e:
            Log(f"{self.id}/{self.indices}/状态上传出现异常:{e}")