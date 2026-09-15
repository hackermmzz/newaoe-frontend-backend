
// 代码运行状态配置
let Code_Status_Error           = 0 //服务器异常
let Code_Status_Wait            = 1 //在等待队列里面
let Code_Status_Compile         = 2 //编译中
let Code_Status_Compile_Success = 3 //编译成功
let Code_Status_Compile_Fail    = 4 //编译错误
let Code_Status_Running         = 5 //正在运行出结果
let Code_Status_Success         = 6 //运行胜利
let Code_Status_Fail            = 7 //游戏失败
let Code_Status_Crash           = 8 //游戏崩溃

//特权配置
let VIP_NONE  = 0 //没有特权
let VIP_SUPER = 1 //特权用户

//api配置
let base_url = "your domain here"; // 默认值
let download_url="your download domain here";
let upload_url="your upload domain here"
let code_url="your admin domain here"
let uploadConfirm_url="your uploadConfirm_url domain here"
let ranking_url="your ranking domain here"
let manager_url="your manager domain here"
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
ranking_url=base_url+"/rank"
manager_url=base_url+"/vip"
// 导出配置
module.exports = {
    ranking_url,
    manager_url,
    manager_student_url: manager_url + "/fetchStudentInfos",
    manager_history_url: manager_url + "/getstudenthistory",
    ManagerStudentRecordPerPage: 10,
    RankingRecordPerPage: 10,
    VIP_NONE,
    VIP_SUPER,
    base_url,
    download_url,
    uploadConfirm_url,
    upload_url,
    code_url: code_url,
    HistoryRecordPerPage,
    Code_Status_Error,
    Code_Status_Wait,
    Code_Status_Compile,
    Code_Status_Compile_Success,
    Code_Status_Compile_Fail,
    Code_Status_Running,
    Code_Status_Success,
    Code_Status_Fail,
    Code_Status_Crash,
};
