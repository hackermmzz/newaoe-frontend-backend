import requests
import json
from config import *
if not DebugLocal:
    import protoc_pb2_grpc
    import protoc_pb2
import sys
import threading
import Util
LocalFiles=[]
idx=0
LocalFilesLock=threading.Lock()

def GetAllFileFromDir(dir:str)->list:
    global LocalFiles
    for d in os.listdir(dir):
        pth=os.path.join(dir,d)
        if os.path.isdir(pth):
            header=f"{pth}/{d}.h"
            source=f"{pth}/{d}.cpp"
            LocalFiles.append((d,header,source))
     
def GetOneStudentCode(server:protoc_pb2_grpc.CodeStub)->map:
    with LocalFilesLock:
        '''
        分别对应：
            是否获取成功
            学号
            代码编号
            头文件
            源文件
        '''
        if not DebugLocal:
            req=protoc_pb2.CodeRequest(auth=GRPCAuth)
            resp=None
            resp=server.GetCode(req)
            if resp is None :
                return {"ok":False,"msg":"服务器异常"}
            if resp.ok is False:
                return {"ok":False,"msg":resp.msg}
            #下载代码
            headerResp=requests.get(resp.headerUrl)
            sourceResp=requests.get(resp.sourceUrl)
            #
            if headerResp.status_code !=200 or sourceResp.status_code!=200:
                return {"ok":False,"msg":"下载文件失败!"}
            headerContent=headerResp.text
            sourceContent=sourceResp.text
            return {
                "ok":True,"msg":"获取代码成功!",
                    "header":headerContent,
                    "source":sourceContent,
                    "id":resp.id,
                    "indices":resp.indices
                    }
        else:
            global LocalFiles,idx
            if idx==len(LocalFiles):
                if len(LocalFiles)==0:
                    GetAllFileFromDir("AOE")
                else:
                    return {"ok":False,"msg":"下载文件失败!"}
            id,headerFile,sourceFile=LocalFiles[idx]
            idx+=1
            headerContent=Util.read_any_text(headerFile)
            sourceContent=Util.read_any_text(sourceFile)
            return {
                "ok":True,"msg":"获取代码成功!",
                    "header":headerContent,
                    "source":sourceContent,
                    "id":id,
                    "indices":""
                    }
                