import chardet
import os
def read_any_text(path):
    # 先读二进制探测编码
    with open(path, "rb") as f:
        raw = f.read()
    detected = chardet.detect(raw)
    encoding = detected["encoding"] or "utf-8"
    
    # 用探测出的编码解码
    return raw.decode(encoding, errors="replace")

        
