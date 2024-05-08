package Lambda

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxLambdaTopLambdaWarningsCommmand = &cobra.Command{
	Use:   "top_lambda_warnings",
	Short: "get top lambda warnings data",
	Long:  `command to get top lambda warnings metrics data`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from child command")
		var authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)
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
			jsonResp, resp, err := GetLambdaTopLambdaWarningsData(cmd, clientAuth)
			if err != nil {
				log.Println("Error getting top lambda zones data : ", err)
				return
			}
			if responseType == "json" {
				fmt.Println(jsonResp)
				}else{
				fmt.Println(resp)

			}
		}
	},
}

func GetLambdaTopLambdaWarningsData(cmd *cobra.Command, clientAuth *model.Auth) (string, []ResData ,error) {
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			err := cmd.Help()
			if err != nil {

			}
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			err := cmd.Help()
			if err != nil {
				// handle error
			}
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}
	logClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	input := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String("CloudTrail/DefaultLogGroup"),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`fields @timestamp, @message, eventVersion, eventTime, requestParameters
		| filter @message like /LAMBDA_WARNING/
		| stats count(*) as frequency by eventTime, requestParameters.functionName as functionName, eventVersion
		| sort frequency desc
		| limit 10`),
	}
	res, err := logClient.StartQuery(input)
	if err != nil {
		return "", nil, fmt.Errorf("failed to start query: %v", err)
	}
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

		queryResults = append(queryResults, queryResult)
		if *queryResult.Status != "Complete" {
			time.Sleep(5 * time.Second) // wait before querying again
			continue
		}
		break // exit loop if query is complete
	}
	resArrMap := make([]ResData, 0)
	for i := 0; i < len(queryResults); i++ {
		fmt.Println(i)
		if *queryResults[i].Status == "Complete" {
			res := queryResults[i].Results
			temMap := make(map[string]string)
			for _, resFileds := range res {
				for _, resField := range resFileds {
					temMap[*resField.Field] = *resField.Value
				}
			}
			tempStruct := &ResData{
				EventTime:    temMap["eventTime"],
				EventVersion: temMap["eventVersion"],
				Frequency:    temMap["frequency"],
				FunctionName: temMap["functionName"],
			}
			resArrMap = append(resArrMap, *tempStruct)
		}
	}
	jsonData, err := json.Marshal(resArrMap)
	if err != nil {
		return "", nil,  fmt.Errorf("error paring json: %v", err)
	}
	return string(jsonData), resArrMap, nil

}

type ResData struct {
	EventTime    string `json:"eventTime"`
	EventVersion string `json:"eventVersion"`
	Frequency    string `json:"frequency"`
	FunctionName string `json:"functionName"`
}
