import os
import sys
import subprocess
import time  
# ====================== 依赖配置区：在这里填需要的pip包 ======================
REQUIRED_PACKAGES = [
    # 示例: "requests>=2.31.0",
    # 本脚本本身不需要任何第三方库，留空即可
]
# ==========================================================================


def install_packages(packages):
    if not packages:
        return
    # 处理 break‑system‑packages 问题
    pip_args = [sys.executable, "-m", "pip", "install"]
    # linux非windows加上break标记，容器专用
    if sys.platform != "win32":
        pip_args.append("--break-system-packages")
    # 清华源
    pip_args.extend(["-i", "https://pypi.tuna.tsinghua.edu.cn/simple"])
    pip_args.extend(packages)

    print(f"准备安装依赖: {packages}")
    proc = subprocess.run(pip_args)
    if proc.returncode != 0:
        print("依赖安装失败，退出")
        sys.exit(1)


def get_bin_name():
    if "win" in sys.platform:
        return "new-aoe-judge.exe"
    else:
        return "new-aoe-judge.out"

def build_project():
    bin_name = get_bin_name()
    # 编译judge
    print("编译judge")
    build_ret = subprocess.run(["go", "build", "-o", bin_name, "./"])
    if build_ret.returncode != 0:
        raise Exception("编译程序失败")
    
    
def run_project():
    bin_name = get_bin_name()
    # 启动judge
    print("启动judge")
    res=subprocess.run(
        [os.path.join(".", bin_name)],
        stdin=sys.stdin,
        stdout=sys.stdout,
        stderr=sys.stderr
        )
    if res.returncode != 0:
        raise Exception("judge运行异常")
    else:
        print("judge正常退出")
        
def git_pull(dir:str,branch: str, retry:int=3):
    # 拉取最新代码
    for i in range(retry):
        try:
            print(f"拉取最新代码，分支: {branch}，重试次数: {i}/{retry}")
            subprocess.run(["git", "pull", "origin", branch], cwd=dir)
            print("拉取最新代码成功")
            break
        except Exception as e:
            print(f"拉取最新代码失败: {e}")
            if i == retry - 1:
                raise e
    
def main():
    # 自动安装缺失依赖
    begin_time = time.time()
    install_packages(REQUIRED_PACKAGES)
    end_time = time.time()
    print(f"依赖安装时间: {end_time - begin_time:.2f}秒")
    
    while True:
        try:
            # 拉取最新代码(judge)的judge分支
            print("开始拉取judge代码")
            begin_time = time.time()
            git_pull("./", "judge")
            end_time = time.time()
            print(f"judge拉取最新代码时间: {end_time - begin_time:.2f}秒")
            print(f"judge拉取最新代码成功")
            
            # 编译judge
            begin_time = time.time()
            try:
                build_project()
            except Exception as e:
                print(f"judge编译异常: {e}")
                sys.exit(1)
            build_end_time = time.time()
            print(f"judge编译时间: {build_end_time - begin_time:.2f}秒")
            
            # 启动judge
            begin_time = time.time()
            try:
                run_project()
            except Exception as e:
                print(f"judge运行异常: {e}")
            run_end_time = time.time()
            print(f"judge运行时间: {run_end_time - begin_time:.2f}秒")
        except Exception as e:
            print(f"judge进程异常退出: {e}")
        finally:
            print("进程退出，3秒后重启...")    
            time.sleep(3)


if __name__ == "__main__":
    main()
