package subrouter

import (
	AssessmentSubmitController "newaoe/Src/codeAssessmentSubmit/controller"
	CommonSubmitController "newaoe/Src/codeCommonSubmit/controller"
	"newaoe/Src/filter"

	"github.com/gin-gonic/gin"
)

func CodeSubmitGroup_Init(group *gin.RouterGroup) {
	g := group.Group("codesubmit")
	g.Use(filter.FilterCookieCheck())
	//普通提交（默认普通提交都会运行，所以要FilterCodeRun一下)
	g.GET("codecommonsubmit", filter.FilterTourist(), filter.FilterLimitCodeSubmit() /* filter.FilterCodeRun(), */, CommonSubmitController.CodeCommonSubmit)
	g.POST("codecommonsubmitACK", CommonSubmitController.CodeCommonSubmitACK)
	//考核提交
	g.GET("codeassessmentsubmit", filter.FilterTourist(), filter.FilterLimitAssessmentSubmission(), AssessmentSubmitController.AssessmentSubmissionSubmit)
	g.POST("codeassessmentsubmitACK", AssessmentSubmitController.AssessmentSubmissionSubmitACK)
	g.GET("fetchteacher", AssessmentSubmitController.StudentGetTeacher)
}
