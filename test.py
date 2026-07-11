import os 
import sys
import shutil
dir="AOE"
for file in os.listdir(dir):
    if file.endswith(".cpp"):
        id=file.split(".cpp")[0]
        dd=f"{dir}/{id}"
        os.makedirs(dd)
        shutil.copy(f"{dir}/{file}",f"{dd}/{file}")
        shutil.copy(f"{dir}/{id}.h",f"{dd}/{id}.h")
