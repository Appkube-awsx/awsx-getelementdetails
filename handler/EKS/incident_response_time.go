package EKS

import (
	"encoding/json"
	"errors"
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

type IncidentResponseResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"incident response time"`
}

var AwsxEKSIncidentResponseTimeCmd = &cobra.Command{
	Use:   "incident_response_time_panel",
	Short: "get incident response time metrics data",
	Long:  `command to get incident response time metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from child command")
		var authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			err := cmd.Help()
			if err != nil {
				log.Println("Error displaying help: ", err)
				return
			}
			return
		}
		if !authFlag {
			log.Println("Authentication failed.")
			return
		}

		responseType, _ := cmd.PersistentFlags().GetString("responseType")
		jsonResp, cloudwatchMetricResp, err := GetIncidentResponseTimeData(cmd, clientAuth, nil)
		if err != nil {
			log.Println("Error getting incident response time data: ", err)
			return
		}
		if responseType == "frame" {
			fmt.Println(cloudwatchMetricResp)
		} else {
			fmt.Println(jsonResp)
		}
	},
}

func GetIncidentResponseTimeData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		if err != nil {
			return "", nil, errors.New("error retrieving cloud element data from CMDB: " + err.Error())
		}
		instanceId = cmdbData.InstanceId
	}

	startTime, endTime, err := parseTimeRange(startTimeStr, endTimeStr)
	if err != nil {
		return "", nil, err
	}

	cloudwatchMetricData, err := GetIncidentResponseTimeMetricData(clientAuth, instanceId, elementType, startTime, endTime, cloudWatchClient)
	if err != nil {
		return "", nil, errors.New("error retrieving incident response time metric data: " + err.Error())
	}

	result := processIncidentResponseTimeRawData(cloudwatchMetricData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		return "", nil, errors.New("error marshalling JSON response: " + err.Error())
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func parseTimeRange(startTimeStr, endTimeStr string) (*time.Time, *time.Time, error) {
	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return nil, nil, errors.New("error parsing start time: " + err.Error())
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return nil, nil, errors.New("error parsing end time: " + err.Error())
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	return startTime, endTime, nil
}

func GetIncidentResponseTimeMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (map[string]*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceId, elementType, startTime, endTime)
	elmType := "ContainerInsights"
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("pod_status_failed"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Average"),
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

	return map[string]*cloudwatch.GetMetricDataOutput{"IncidentResponse": result}, nil
}

func processIncidentResponseTimeRawData(result map[string]*cloudwatch.GetMetricDataOutput) IncidentResponseResult {
	rawData := IncidentResponseResult{}
	for _, metricData := range result {
		for i, timestamp := range metricData.MetricDataResults[0].Timestamps {
			rawData.RawData = append(rawData.RawData, struct {
				Timestamp time.Time
				Value     float64
			}{*timestamp, *metricData.MetricDataResults[0].Values[i]})
		}
	}
	return rawData
}

func init() {
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSIncidentResponseTimeCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
