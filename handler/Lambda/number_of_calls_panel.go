package Lambda

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type NumberOfCallsResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"number_of_calls_panel"`
}

var AwsxLambdaNumberOfCallsCmd = &cobra.Command{
	Use:   "number of calls panel",
	Short: "get number of calls metrics data",
	Long:  `command to get number of calls metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaNumberOfCallsPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting number of calls data : ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetLambdaNumberOfCallsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)

	}

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return "", nil, err
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
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	// Debug prints
	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	totalInvocations, err := GetLambdaNumberOfCallsMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Number Of Calls"] = totalInvocations

	// Process the raw data if needed
	result := processRawData(totalInvocations)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetLambdaNumberOfCallsMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {

	input := &cloudwatch.GetMetricDataInput{
		StartTime: startTime,
		EndTime:   endTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("AWS/Lambda"),
						MetricName: aws.String("Invocations"),
						Dimensions: []*cloudwatch.Dimension{
							// {
							//     Name:  aws.String("FunctionName"),
							//     Value: aws.String(instanceId),
							// },
						},
					},
					Period: aws.Int64(300), // Adjust period as needed (e.g., 5 minutes)
					Stat:   aws.String(statistic),
				},
			},
		},
	}

	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	if len(result.MetricDataResults) == 0 {
		return nil, fmt.Errorf("no data available for the specified time range")
	}

	return result, nil
}

func processRawData(result *cloudwatch.GetMetricDataOutput) NumberOfCallsResult {
	var rawData NumberOfCallsResult
	rawData.RawData = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.RawData[i].Timestamp = *timestamp
		rawData.RawData[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}

func init() {
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxLambdaNumberOfCallsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
