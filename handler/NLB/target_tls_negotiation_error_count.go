package NLB

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	 "github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type NlbTargetTlsErrorCountTime struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"target_tls_negotiation_error_count_panel"`
}

var AwsxNlbTargetTlsErrorCountCmd = &cobra.Command{
	Use:   "target_tls_negotiation_error_count_panel",
	Short: "get target tls error count metrics data",
	Long:  `command to get target tls count metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetTargetTlsErrorCountData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting nlb target tls response data: ", err)
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

func GetTargetTlsErrorCountData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	TargetTlsErrorCount, err := GetNlbTargetTlsErrorCountMetricValue(clientAuth, instanceId, elementType, startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NLB active connections data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["TargetTLSErrorCount"] = TargetTlsErrorCount

	result := ProcessTargetResponseRawData(TargetTlsErrorCount)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetNlbTargetTlsErrorCountMetricValue(clientAuth *model.Auth, instanceId string, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("TargetTLSNegotiationErrorCount"),
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

func ProcessTargetTlsResponseRawData(result *cloudwatch.GetMetricDataOutput) NlbTargetTlsErrorCountTime {
	var rawData NlbTargetTlsErrorCountTime
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
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("query", "", "query")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxNlbTargetTlsErrorCountCmd.PersistentFlags().String("ApiName", "", "api name")
}
