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
	//
	g.GET("codecommonsubmit", filter.FilterLimitCodeSubmit(), CommonSubmitController.CodeCommonSubmit)
	g.POST("codecommonsubmitACK", CommonSubmitController.CodeCommonSubmitACK)
	g.GET("codeassessmentsubmit", filter.FilterLimitAssessmentSubmission(), AssessmentSubmitController.AssessmentSubmissionSubmit)
	g.POST("codeassessmentsubmitACK", AssessmentSubmitController.AssessmentSubmissionSubmitACK)
	g.GET("fetchteacher", AssessmentSubmitController.StudentGetTeacher)
}
