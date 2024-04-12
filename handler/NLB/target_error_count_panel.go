package NLB

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	// "github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type NlbTargetErrorCountTime struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"target_error_count_panel"`
}

var AwsxNlbTargetErrorCountCmd = &cobra.Command{
	Use:   "target_error_count_panel",
	Short: "get target error count metrics data",
	Long:  `command to get target count metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetTargetErrorCountData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting nlb target response data: ", err)
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

func GetTargetErrorCountData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
		// cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		// if err != nil {
		// 	return "", nil, err
		// }
		// instanceId = cmdbData.InstanceId

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
	TargetErrorCount, err := GetNlbTargetErrorCountMetricValue(clientAuth, instanceId, elementType, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NLB active connections data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["TargetTLSErrorCount"] = TargetErrorCount

	result := ProcessTargetResponseRawData(TargetErrorCount)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetNlbTargetErrorCountMetricValue(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("LoadBalancer"),
								Value: aws.String("net/a0affec9643ca40c5a4e837eab2f07fb/f623f27b6210158f"),
							},
						},
						Namespace:  aws.String("AWS/NetworkELB"),
						MetricName: aws.String("ClientTLSNegotiationErrorCount"),
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

func ProcessTargetResponseRawData(result *cloudwatch.GetMetricDataOutput) NlbTargetErrorCountTime {
	var rawData NlbTargetErrorCountTime
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
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("query", "", "query")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxNlbTargetErrorCountCmd.PersistentFlags().String("ApiName", "", "api name")
}
