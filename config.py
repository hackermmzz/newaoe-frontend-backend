import os
from enum import Enum
import requests
import time
import threading
import protoc_pb2
import protoc_pb2_grpc
#服务器网址
BaseIP="114.66.62.156"
GRPCHost=f"{BaseIP}:50051"
Cookie='''nxd_tooken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbF9iaW5kIjoiMjA0OTk4MzQ3NEBxcS5jb20iLCJleHBpcmVfdGltZSI6IjIxMjYtMDMtMThUMTM6NDQ6NDAuNDAxMzI4MjM3WiIsImlkIjoiOTIzMTA2ODQwNDI5IiwiaXAiOiIxMjcuMC4wLjEiLCJsb2dpbl90aW1lIjoiMjAyNi0wNC0xMVQxMzo0NDo0MC40MDEzMjk2MDhaIiwicmFuZG9tIjoiSXVtdms3WFdJTmZqczBCQ09yeDBJaWhvckFrMDA0MGciLCJyZWdpc3RfZGF0ZSI6IjIwMjYtMDQtMTFUMjE6MzY6MjMrMDg6MDAiLCJ3bGhfdG9feW91Ijoi5Li65LuA5LmI5LiN546p5Y6f56WePyEifQ.MBW4WjnyuSu5-wu2FNWPRxjXA5xaU4uUbx--b9JLgX0'''
GRPCAuth='''5rGq56uL5rSq5piv5YWo5LiW55WM5pyA5biF55qE55S355Sf'''

BaseHost=f"http://{BaseIP}:8080"
base_url=f"{BaseHost}/api/code"
CodeGetURL=f"{base_url}/CodeGet"
CodeStatusPostURL=f"{base_url}/CodeRunStatusPost"
#定义api_key
api_key="new-aoe"
#定义请求头
header={
    "api":api_key
}
#AOE项目目录
new_aoe_folder=f"{os.getcwd()}/new-aoe"
#AOE镜像名称
new_aoe_docker_img="hackermmzz/qt-env"
#编译的输出日志
CompileLogFileName="CompileLog.txt"
#编译是否开启多线程编译模式
CompileUseMultiThread=False
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
#本地测试
DebugLocal=False
#是否单线程测试
IsSingleThreadTest=False
#测试得轮数
Epochs=3

###########################################################
def Log(msg):
    with LogLock:
        LogFile.write(f"[{time.strftime('%Y-%m-%d %H:%M:%S')}] Thread {threading.current_thread().ident}: {msg}"+"\n")
        LogFile.flush()
###########################################################

class PostRunStatusEnum(Enum):
    Code_Status_Wait            = 1 #在等待队列里面
    Code_Status_Compile         = 2 #编译中
    Code_Status_Compile_Error   = 3 #编译错误
    Code_Status_Compile_Success = 4 #编译成功
    Code_Status_Running         = 5 #正在运行出结果
    Code_Status_Success         = 6 #运行胜利
    Code_Status_Crash           = 7 #游戏崩溃
    Code_Status_Fail            = 8 #游戏失败

class PostRunStatus():
    def __init__(self,server:protoc_pb2_grpc.CodeStub,data:protoc_pb2.CodeStatusUpdateRequest):
        self.data=data
        self.server=server
    def Response(self)->bool:
        if DebugLocal:
            return
        try:
            resp=self.server.CodeStatusUpdate(self.data)
            if resp==None:
                Log(f"{self.id}/{self.indices}/后端异常!")
        except Exception as e:
            Log(f"{self.id}/{self.indices}/状态上传出现异常:{e}")