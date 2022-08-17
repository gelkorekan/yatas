package apigateway

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfStagesCloudwatchLogsExist(checkConfig yatas.CheckConfig, stages []types.Stage, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Apigateway Cloudwatch Logs enabled", "Check if all cloudwatch logs are enabled for all stages", testName)
	for _, stage := range stages {
		if stage.AccessLogSettings != nil && stage.AccessLogSettings.DestinationArn != nil {
			Message := "Cloudwatch logs are enabled on stage" + *stage.StageName
			result := results.Result{Status: "OK", Message: Message, ResourceID: *stage.StageName}
			check.AddResult(result)
		} else {
			Message := "Cloudwatch logs are not enabled on " + *stage.StageName
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *stage.StageName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
