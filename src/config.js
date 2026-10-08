
// 运行类型（Release运行或者重度Debug运行）
let	Code_ReleaseRun = 0 //release版本运行
let	Code_DebugRun   = 1 //重度调试模式运行

//下载链接类型
let DownloadURLType_UnKnown = 0 //未知
let DownloadURLType_OSS     = 1 //oss下载
let DownloadURLType_CDN     = 2 //cdn下载

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
let VIP_TOURIST = -1 //游客用户(只看查看，不可提交)
let VIP_NONE  = 0 //没有特权
let VIP_SUPER = 1 //特权用户

// 上传的代码类别
let Code_CommonSubmit     = 1 //普通提交
let Code_AssessmentSubmit = 2 //考核提交
let Code_ReRunSubmit= 3      //重新运行提交

//api配置
let base_url = "your domain here"; // 默认值
let download_url="your download domain here";
let ranking_url="your ranking domain here"
let manager_url="your manager domain here"
let history_url="your history domain here"
let feedback_url="your feedback domain here"
let avatarUpdate_url="your avatarUpdate domain here"
let codeSubmit_url="your codeSubmit domain here"
let codeRun_url="your codeRun domain here"
let uploadConfirm_url="your uploadConfirm domain here"
let teacherFetch_url="your teacherFetch domain here"
let teacherAdd_url="your teacherAdd domain here"
const isTestMode =false;
const HistoryRecordPerPage=10;
// 如果不是测试模式，则使用生产环境的base_url
if (isTestMode) {
    base_url = "http://localhost:8080/api";
    
}else{
    base_url = "http://114.66.62.156:8080/api";
}
download_url=base_url+"/download";
codeSubmit_url=base_url+"/codesubmit"
codeRun_url=base_url+"/coderun"
uploadConfirm_url=base_url+"/uploadconfirm"
ranking_url=codeRun_url+"/fetchrank"
manager_url=base_url+"/vip"
history_url=base_url+"/coderun/fetchhistory"
feedback_url=base_url+"/home"
avatarUpdate_url=base_url+"/home/avatarUpdate"
teacherFetch_url=base_url+"/codesubmit/fetchteacher"
teacherAdd_url=manager_url+"/teacherAdd"
// 导出配置
module.exports = {
    Code_ReleaseRun,
    Code_DebugRun,
    DownloadURLType_UnKnown,
    DownloadURLType_OSS,
    DownloadURLType_CDN,
    Code_CommonSubmit,
    Code_AssessmentSubmit,
    Code_ReRunSubmit,
    ranking_url,
    manager_url,
    teacherFetch_url,
    teacherAdd_url,
    manager_student_url: manager_url + "/fetchStudentInfos",
    manager_student_search_url: manager_url + "/searchstudentbyid",
    manager_history_url: manager_url + "/getstudenthistory",
    manager_all_history_url: manager_url + "/fetchSubmitRecord",
    manager_assessment_history_url: manager_url + "/getassessmenthistory",
    manager_run_all_assessment_submit_url: manager_url + "/runAllAssessmentSubmit",
    manager_rerun_anomal_record_url: manager_url + "/rerunAnomalRecord",
    manager_block_all_submit_url: manager_url + "/blockAllSubmit",
    manager_cancel_submit_block_url: manager_url + "/cancelSubmitBlock",
    manager_feedback_url: manager_url + "/fetchfeedback",
    manager_info_export_url: manager_url + "/infoExport",
    manager_reset_common_submit_time_url: manager_url + "/resetcommonsubmittime",
    // “重置所有人次数”接口每次请求处理的学生数量
    resetcommonsubmitTimeBatch: 100,
    ManagerStudentRecordPerPage: 10,
    ManagerHistoryRecordPerPage: 10,
    ManagerFeedbackRecordPerPage: 10,
    RankingRecordPerPage: 10,
    VIP_NONE,
    VIP_TOURIST,
    VIP_SUPER,
    avatarUpdate_url,
    codeSubmit_url,
    codeRun_url,
    uploadConfirm_url,
    base_url,
    download_url,
    feedback_url,
    history_url,
    fetch_announcement_url: base_url + "/home/fetchAnnouncement",
    publish_announcement_url: manager_url + "/publishAnnouncement",
    manager_oj_version_update_url: manager_url + "/ojVersionUpdate",
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
