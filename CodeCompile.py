from CodeGet import base_url
import subprocess
import os
from config import *
#
def CodeCompile(buildDir:str,logfile:str)->tuple[bool,str]:
    '''
    编译代码
    '''
    #编译代码
    docker_cmd = [
        "docker", "run",
        "--rm", # 运行完自动删除容器
        "-v", f"{new_aoe_folder}:/app/project:ro",  # 挂载当前项目目录到容器 /app
        "-v",f"{buildDir}:/app/build", #挂载编译输出目录(用户源码在这)
        "-w", "/app/build",  # 工作目录
        new_aoe_docker_img,
        # 编译命令（把所有输出打出来）
        "bash","-c",
        #不用make,直接g++编译，减少编译时间，同时也能避免makefile被改了导致的编译问题
        '''
        export PATH+=":/opt/qt5.9.2/bin/" &&
        export QT="/opt/qt5.9.2" &&
        export QTINCLUDE="/opt/qt5.9.2/include" && 
        g++ -c UsrAI.cpp -fPIC -I./ -I../project/ -I${QTINCLUDE} -I${QTINCLUDE}/QtCore -I${QTINCLUDE}/QtMultimedia -I${QTINCLUDE}/QtWidgets -I${QTINCLUDE}/QtGui -I${QTINCLUDE}/QtNetwork &&
        g++ UsrAI.o ../project/*.o -o newAOE  -L/opt/qt5.9.2/lib -lQt5Widgets -lQt5Gui -lQt5Core -lQt5Multimedia -lQt5Network 
        '''
    ]
    #编译并且保存编译日志
    with open(logfile,"w")as f:
        result = subprocess.run(
            docker_cmd,
            stdout=subprocess.PIPE,     # 输出捕获到变量
            stderr=f,                   # 报错也写文件
            text=True,
            check=False
        )
    #判断是否编译成功
    if result.returncode != 0:
        with open(logfile,"r")as f:
            return (False,f.read())
    #编译好的文件为newAOE
    return (True,"")
    
    
    
    