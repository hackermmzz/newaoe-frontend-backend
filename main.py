import functools
import json
import os
import queue
import shutil
import subprocess
import sys
import threading
import time
from concurrent.futures.thread import ThreadPoolExecutor
from datetime import datetime

import requests

############################################
base_url = "http://localhost:8080/api/admin"   #服务器网址
LegalValidateHeaderKey = "mmzz"                 #合法的值
codeGet_url = f"{base_url}/CodeGet"             #获取代码
codeRunPost_url=f"{base_url}/CodeRunStatusPost" #上传运行结果
newaoe_path="C:\\QtHomeWork\\new-aoe"      #newaoe项目所在目录
header_file=f"{newaoe_path}/UsrAI.h"
source_file=f"{newaoe_path}/UsrAI.cpp"
qmake_path="C:\\Qt\\Qt5.9.2\\5.9.2\mingw53_32\\bin\\qmake.exe" #qmake路径
WorkDir=os.getcwd()+"/WorkDir"
MaxPayload=20                                   #机器支持同时运行的最多进程个数
CodeRunPool=queue.Queue(MaxPayload//2)             #待运行队列(暂设为最大负载的一半)
RecordFile="record.txt"                         #newaoe运行存放结果的文件名称
CodeRunArgument={
    "--exam":"true",
    "--id":"",
    "--indices":"",
    "--api":LegalValidateHeaderKey
}                              #程序传入的参数
############################################线程运行函数
def ThreadRunMain():
    while True:
        try:
            task=CodeRunPool.get()
            try:
                task()
            except Exception as e:
                print(e)
        except Exception as e:
            print(e)
############################################初始化
def Init():
    #创建线程
    for _ in range(MaxPayload):
        thread=threading.Thread(target=ThreadRunMain,daemon=True)
        thread.start()
    #创建文件夹
    if not os.path.exists(WorkDir):
        os.mkdir(WorkDir)
    #检查newae目录是否存在MakeFile
    if not os.path.exists(f"{newaoe_path}/Makefile"):
        #执行Qmake
        cmd=[qmake_path]
        res=subprocess.run(cmd)
        if res!=0:
            print("qmake失败")
            exit(0)

    print("初始化成功")
############################################获取代码
def GetCodeFile():
    res = requests.get(codeGet_url, headers={"api":LegalValidateHeaderKey})
    data = res.content.decode("utf-8")
    data = json.loads(data)
    if data["status"] == False:
        return False,"",""
    #获取信息
    data=data["data"]
    id=data["id"]
    indices=data["indices"]
    header_code=data["header"]
    source_code=data["source"]
    #写入指定文件
    with open(header_file,"w",encoding="utf-8") as f:
        f.write(header_code)
    with open(source_file,"w",encoding="utf-8") as f:
        f.write(source_code)
    return True,id,indices
    #
#################################################编译代码
def CompileCode(error_fd):
    cmd=["make","-j4"]
    result = subprocess.run(
        cmd,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,  # 输出为字符串（而非 bytes）,
        cwd=newaoe_path
    )
    #
    if result.returncode==0:
        return True,None
    else:
        return False,result.stderr
#################################################将运行结果返回给服务器
def CodeResultPost(id:str, indices:int,status:int,data:any):
    data={
        "id":id,
        "indices":indices,
        "status":status,
        "data":data
    }
    resp=requests.post(url=codeRunPost_url,json=data,headers={"api":LegalValidateHeaderKey})
    data=json.loads(resp.content)
    return data["status"]
#################################################
def OJRun():
    print("OJ启动")
    while True:
        try:
            #下载代码
            flag, id, indices = GetCodeFile()
            if flag == False:
                print("获取代码失败")
                time.sleep(10)
                continue
            print("获取代码成功")
            #编译文件
            print("正在编译...")
            CodeResultPost(id,indices,7,"编译中")
            flag,errPipe=CompileCode(sys.stdout)
            if flag==False:
                print("编译失败")
                CodeResultPost(id, indices, 9,f"{errPipe}")
                continue
            print("编译成功")
            CodeResultPost(id,indices,10,f"编译成功")
            #创建文件夹
            current_time = datetime.now().strftime("%Y_%m_%d_%H_%M_%S")
            target_dir=f"{WorkDir}/{id}_{indices}_{current_time}"
            exe_finame=f"{target_dir}/{id}.exe"
            os.mkdir(target_dir)
            shutil.copy(f"{newaoe_path}/release/newAOE.exe", exe_finame)
            CodeRunArgument["--id"]=id
            CodeRunArgument["--indices"]=str(indices)

            #运行文件，放入运行队列
            def RunTask(argument:dict):
                cmd=[exe_finame]
                for key,val in argument.items():
                    cmd.append(key)
                    cmd.append(val)
                res=subprocess.run(
                    cmd,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                    cwd=newaoe_path
                )
                print(res.stderr)
            print(f"正在运行{id}的代码")
            CodeRunPool.put(functools.partial(RunTask,argument=CodeRunArgument.copy()))
        except Exception as e:
            print(e)
#################################################运行
def main():
    Init()
    OJRun()
#####################################################
if __name__=="__main__":
    main()