//常量配置

const Code_Status_IDLE         = 3 //闲置
const Code_Status_Wait         = 1 //在等待队列里面
const Code_Status_Running      = 2 //正在运行出结果
const Code_Status_Run_Timeout  = 4 //运行出错/超时
const Code_Status_Wait_Timeout = 5 //等待超时/出错
const Code_Status_Success      = 6 //运行成功
// 配置文件
let base_url = "your domain here"; // 默认值
let download_url="your download domain here";
let upload_url="your upload domain here"
let code_url="your admin domain here"
const isTestMode =true;
const HistoryRecordPerPage=10;
// 如果不是测试模式，则使用生产环境的base_url
if (isTestMode) {
    base_url = "http://localhost:8080/api";
    
}else{
    base_url = "http://www.adhn.asia:8080/api";
}
download_url=base_url+"/download";
upload_url=base_url+"/upload"
code_url=base_url+"/code"
// 导出配置
module.exports = {
    base_url,
    download_url,
    upload_url,
    code_url: code_url,
    Code_Status_IDLE  ,      
    Code_Status_Wait     ,    
    Code_Status_Running     ,
    Code_Status_Run_Timeout,  
    Code_Status_Wait_Timeout ,
    Code_Status_Success,
    HistoryRecordPerPage,
};
    