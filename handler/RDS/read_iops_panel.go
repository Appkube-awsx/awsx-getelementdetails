package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type ReadIOPS struct {
	ReadIOPS []struct {
		Timestamp time.Time
		Value     float64
	} `json:"read_iops"`
}

var AwsxRDSReadIOPSCmd = &cobra.Command{
	Use:   "read_iops_panel",
	Short: "Get read IOPS metrics data",
	Long:  `Command to get read IOPS metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")
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
			jsonResp, _, err := GetRDSReadIOPSPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting read IOPS data: ", err)
				return
			}
			fmt.Println("Read IOPS Data:")
			fmt.Println(jsonResp)
		}
	},
}

func GetRDSReadIOPSPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	instanceID, _ := cmd.PersistentFlags().GetString("instanceId")
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

	rawReadIOPSData, err := GetReadIOPSMetricData(clientAuth, instanceID, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting read IOPS data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["ReadIOPS"] = rawReadIOPSData

	result := processRawReadIOPSData(rawReadIOPSData)
	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetReadIOPSMetricData(clientAuth *model.Auth, instanceID string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace AWS/RDS from %v to %v", instanceID, startTime, endTime)

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("readIOPS"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{},
						MetricName: aws.String("ReadIOPS"),
						Namespace:  aws.String("AWS/RDS"),
					},
					Period: aws.Int64(300), // 5 minutes (in seconds)
					Stat:   aws.String("Sum"),
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

	return result, nil
}

func processRawReadIOPSData(result *cloudwatch.GetMetricDataOutput) ReadIOPS {
	var rawData ReadIOPS
	rawData.ReadIOPS = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.ReadIOPS[i].Timestamp = *timestamp
		rawData.ReadIOPS[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}

func init() {
	AwsxRDSReadIOPSCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSReadIOPSCmd.PersistentFlags().String("endTime", "", "endcl time")
}
