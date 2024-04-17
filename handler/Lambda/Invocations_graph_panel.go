package Lambda

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type InvocationsGraph struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"invocations_graph_panel"`
}

var AwsxLambdaInvocationsGraphCmd = &cobra.Command{
	Use:   "Invocations_graph_panel",
	Short: "get Invocations count graph metrics data",
	Long:  `command to get Invocations count graph metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaInvocationsGraphData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda Invocations response data: ", err)
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

func GetLambdaInvocationsGraphData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		if err != nil {
			return "", nil, err
		}
		instanceId = cmdbData.InstanceId

	}

	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

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

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	InvocationsCount, err := GetLambdaInvocationsCountMetricValue(clientAuth, instanceId, elementType, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting lambda throttles count data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Invocations"] = InvocationsCount

	result := ProcessLambdaInvocationsRawData(InvocationsCount)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}
	fmt.Println(jsonString)

	return string(jsonString), cloudwatchMetricData, nil
}

func GetLambdaInvocationsCountMetricValue(clientAuth *model.Auth, instanceId string, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("invocations"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{

						Namespace:  aws.String("AWS/Lambda"),
						MetricName: aws.String("Invocations"),
					},
					Period: aws.Int64(300),
					Stat:   aws.String(statistic),
				},
			},
		},
		StartTime: startTime,
		EndTime:   endTime,
	}

	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ProcessLambdaInvocationsRawData(result *cloudwatch.GetMetricDataOutput) InvocationsGraph {
	var rawData InvocationsGraph
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
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxLambdaInvocationsGraphCmd.PersistentFlags().String("ApiName", "", "api name")
}
