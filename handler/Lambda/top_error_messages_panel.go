package Lambda

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	comman_function "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxTopErrorsMessagesPanelCmd = &cobra.Command{
	Use:   "top_errors_messages_panel",
	Short: "Get top errors messages events",
	Long:  `Command to retrieve top errors events`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running  panel command")

		var authFlag bool
		var clientAuth *model.Auth
		var err error
		authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)

		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			err := cmd.Help()
			if err != nil {
				return
			}
			return
		}
		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			jsonResp, resp, err := GetLambdaTopErrorsMessagesEvents(cmd, clientAuth)
			if err != nil {
				log.Println("Error getting top lambda zones data : ", err)
				return
			}
			if responseType == "json" {
				fmt.Println(jsonResp)
			} else {
				fmt.Println(resp)

			}

		}
	},
}

type ResultData struct {
	EventTime    string `json:"eventTime"`
	EventVersion string `json:"eventVersion"`
	Frequency    string `json:"frequency"`
	FunctionName string `json:"functionName"`
}

func GetLambdaTopErrorsMessagesEvents(cmd *cobra.Command, clientAuth *model.Auth) (string, []ResultData, error) {

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}
	logClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	input := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String("CloudTrail/DefaultLogGroup"),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`
		fields @timestamp, @message, eventVersion, eventTime, requestParameters
			| filter eventSource = "lambda.amazonaws.com"
			| filter @message like /ERROR|Exception|Failed/
			| stats count(*) as frequency by eventTime, requestParameters.functionName as functionName, eventVersion
			| sort frequency desc
			| limit 10`),
	}
	res, err := logClient.StartQuery(input)
	if err != nil {
		return "", nil, fmt.Errorf("failed to start query: %v", err)
	}
	fmt.Println("---------", res)
	queryId := res.QueryId
	var queryResults []*cloudwatchlogs.GetQueryResultsOutput // Declare queryResults outside the loop

	for {
		// Check query status
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}

		queryResult, err := logClient.GetQueryResults(queryStatusInput)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get query results: %v", err)
		}
		fmt.Println("queryResult", queryResult)
		queryResults = append(queryResults, queryResult) // Append each query result to queryResults
		if *queryResult.Status != "Complete" {
			time.Sleep(5 * time.Second) // wait before querying again
			continue
		}
		break // exit loop if query is complete
	}
	resArrMap := make([]ResultData, 0)
	for i := 0; i < len(queryResults); i++ {
		fmt.Println(i)
		if *queryResults[i].Status == "Complete" {
			res := queryResults[i].Results
			for _, resFields := range res {
				tempStruct := ResultData{}
				for _, resField := range resFields {
					switch *resField.Field {
					case "eventTime":
						tempStruct.EventTime = *resField.Value
					case "eventVersion":
						tempStruct.EventVersion = *resField.Value
					case "frequency":
						tempStruct.Frequency = *resField.Value
					case "functionName":
						tempStruct.FunctionName = *resField.Value
					}
				}
				resArrMap = append(resArrMap, tempStruct)
			}
		}
	}
	// fmt.Println("resArrMap", resArrMap)
	jsonData, err := json.Marshal(resArrMap)
	if err != nil {
		return "", nil, fmt.Errorf("error paring json: %v", err)
	}
	return string(jsonData), resArrMap, nil

	
}

func init() {
	AwsxTopErrorsMessagesPanelCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxTopErrorsMessagesPanelCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxTopErrorsMessagesPanelCmd.PersistentFlags().String("endTime", "", "end time")
}
