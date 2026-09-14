from CodeGet import base_url
import subprocess
import os
from config import *
import threading
import Util

compileLock=threading.Lock()

def CodeCompile(buildDir: str, logfile: str) -> tuple[bool, str]:
    """
    编译用户代码。

    不使用全局 compileLock，
    不同 buildDir 可以并发编译。
    """

    docker_cmd = [
        "docker", "run",
        "--rm",

        # 公共预编译文件，只读
        "-v", f"{new_aoe_folder}:/app/project:ro",

        # 每个用户独立的编译目录
        "-v", f"{buildDir}:/app/build",

        "-w", "/app/",

        new_aoe_docker_img,

        "bash", "-c",

        r'''
        
        export PATH="$PATH:/opt/qt5.9.2/bin/" &&
        export QT="/opt/qt5.9.2" &&
        export QTINCLUDE="/opt/qt5.9.2/include" &&
        #先copy一下newaoe_copy目录
        mkdir -p newaoe_copy
        cd project
        find ./ \
            \( -name "*.h" -o -name "*.cpp" -o -name "*.hpp" \) \
            -exec cp --parents {} ../newaoe_copy \;
        cd ../
        #把UsrAI.h和UsrAI.cpp copy到newaoe_copy目录下
        cp build/UsrAI.h build/UsrAI.cpp newaoe_copy/
        #切换到newaoe_copy目录下
        cd newaoe_copy
        #修复大小写问题
        fixcase -f ./ >/dev/null
        # 编译用户代码
        g++ -c UsrAI.cpp \
            -O2 \
            -fPIC \
            -I./ \
            -I${QTINCLUDE} \
            -I${QTINCLUDE}/QtCore \
            -I${QTINCLUDE}/QtMultimedia \
            -I${QTINCLUDE}/QtWidgets \
            -I${QTINCLUDE}/QtGui \
            -I${QTINCLUDE}/QtNetwork &&

        # 链接公共 .o 文件
        g++ UsrAI.o ../project/release/*.o \
            -o newAOE \
            -L/opt/qt5.9.2/lib \
            -lQt5Widgets \
            -lQt5Gui \
            -lQt5Core \
            -lQt5Multimedia \
            -lQt5Network
        #把newAOE copy回buildDir目录下
        cp newAOE ../build/  >/dev/null
        '''
    ]

    with open(
        logfile,
        "w",
        errors="replace"
    ) as f:

        result = subprocess.run(
            docker_cmd,
            stdout=f,
            stderr=f,
            text=True,
            check=False
        )

    if result.returncode != 0:
        return (
            False,
            Util.read_any_text(logfile)
        )

    return (True, "")