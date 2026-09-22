from CodeGet import base_url
import subprocess
import os
from config import *
import threading
import Util
import config
compileLock=threading.Lock()

def CodeCompile(buildDir: str, logfile: str,runtype:int) -> tuple[bool, str]:
    """
    编译用户代码。

    不使用全局 compileLock，
    不同 buildDir 可以并发编译。
    """
    # 读取bash脚本
    with open(f"bash/codecompile.sh", "r",encoding="utf-8") as f:
        bashScript = f.read()
    bashScript = bashScript.format(
                buildDir=buildDir,
                DebugMode=runtype==config.CodeRunTypeEnum.CodeRunType_DebugRun.value,
                ).strip()
    
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
        bashScript,
    ]

    with open(logfile,"w",errors="replace") as f:
        result = subprocess.run(
            docker_cmd,
            stdout=f,
            stderr=f,
            text=True,
            check=False
        )
    if result.returncode != 0:
        Log(f"编译失败：{Util.read_any_text(logfile)}")
        return (False,Util.read_any_text(logfile))
    return (True, "")