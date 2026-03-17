import requests
import json
from config import *

def GetOneStudentCode()->map:
    '''
    分别对应：
        是否获取成功
        学号
        代码编号
        头文件
        源文件
    '''
    resp=requests.get(url=CodeGetURL,headers=header)
    #如果获取不成功直接返回
    if resp.status_code!=200:
        return {"status":False,"msg":"获取失败"}
    #获取成功直接返回内容
    try:
        jsData=json.loads(resp.content)
        jsData=jsData["data"]
        return {"status":True,"id":jsData["id"],"indices":int(jsData["indices"]),"header":jsData["header"],"source":jsData["source"]}
    except Exception as e:
        return {"status":False,"msg":f"出现解析异常:{e}"}
    
    