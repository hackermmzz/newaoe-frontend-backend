
// 配置文件
let base_url = "your domain here"; // 默认值
let download_url="your download domain here";
let upload_url="your upload domain here"
let code_url="your admin domain here"
let uploadConfirm_url="your uploadConfirm_url domain here"
const isTestMode =false;
const HistoryRecordPerPage=10;
// 如果不是测试模式，则使用生产环境的base_url
if (isTestMode) {
    base_url = "http://localhost:8080/api";
    
}else{
    base_url = "http://114.66.62.156:8080/api";
}
download_url=base_url+"/download";
upload_url=base_url+"/upload"
uploadConfirm_url=base_url+"/uploadconfirm"
code_url=base_url+"/code"
// 导出配置
module.exports = {
    base_url,
    download_url,
    uploadConfirm_url,
    upload_url,
    code_url: code_url,
    HistoryRecordPerPage,

};
    