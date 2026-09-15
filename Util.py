import chardet
import os
import requests
from config import Log
import random
import time

def read_any_text(path):
    # 先读二进制探测编码
    with open(path, "rb") as f:
        raw = f.read()
    detected = chardet.detect(raw)
    encoding = detected["encoding"] or "utf-8"
    
    # 用探测出的编码解码
    return raw.decode(encoding, errors="replace")


def UploadData(url: str, data: bytes) :
    """
    传入预签名上传链接，把二进制data PUT上传到该url
    :param url: minio预签名PUT链接
    :param data: 二进制字节数据（日志、文本转bytes）
    """
    resp = requests.put(url, data=data)
    try:
        resp.raise_for_status()  # 如果http状态不是2xx，直接抛异常
    except Exception as e:
        Log(f"上传数据到{url}失败，错误信息:{e}")

def FileExist(path: str) -> bool:
    """
    检查文件是否存在
    :param path: 文件路径
    :return: 是否存在
    """
    return os.path.exists(path)

def GetFolerRandomFileIfExist(path:str)->str:
    """
    从文件夹中随机获取一个文件
    :param path: 文件夹路径
    :return: 随机文件路径
    """
    if not os.path.exists(path):
        return ""
    files = os.listdir(path)
    if not files:
        return ""
    return files[0]

def GetRandomStr():
    ns_total = time.time_ns()
    sec = ns_total // 1_000_000_000
    ns = ns_total % 1_000_000_000
    tm = time.localtime(sec)
    # 样例：2026_09_16_15_40_22_123456789
    return time.strftime("%Y_%m_%d_%H_%M_%S", tm) + f"_{ns:09d}_f{random.randint(0,999999)}"