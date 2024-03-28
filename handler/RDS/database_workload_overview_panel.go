package RDS

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

type DatabaseWorkloadOverview struct {
	Timestamp time.Time
	Value     float64
}

var AwsxRDSDBLoadCmd = &cobra.Command{
	Use:   "db_load_panel",
	Short: "get database workload overview metrics data",
	Long:  `Command to get database workload overview metrics data`,

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
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			jsonResp, cloudwatchMetricResp, err, _ := GetRDSDBLoadPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting database workload overview data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// Default case: print JSON
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetRDSDBLoadPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")

	if elementId != "" {
		log.Println("Getting cloud-element data from CMDB")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("Using default CMDB URL")
			apiUrl = config.CmdbUrl
		}
		log.Println("CMDB URL: " + apiUrl)
	}

	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return "", "", nil, err
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
			return "", "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data for database workload overview metric
	rawData, err := GetDBLoadMetricData(clientAuth, elementType, startTime, endTime, "DBLoad", cloudWatchClient)
	if err != nil {
		log.Println("Error getting database workload overview data: ", err)
		return "", "", nil, err
	}
	cloudwatchMetricData["DBLoad"] = rawData

	// Process raw data
	result := processedRawDatabaseWorkloadOverviewData(rawData)
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshalling JSON for database workload overview data: ", err)
		return "", "", nil, err
	}

	return string(jsonData), "", cloudwatchMetricData, nil
}

func GetDBLoadMetricData(clientAuth *model.Auth, elementType string, startTime, endTime *time.Time, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace AWS/RDS from %v to %v", elementType, startTime, endTime)

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String  ("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{},
						MetricName: aws.String(metricName),
						Namespace:  aws.String("AWS/RDS"),
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

	return result, nil
}

func processedRawDatabaseWorkloadOverviewData(result *cloudwatch.GetMetricDataOutput) []DatabaseWorkloadOverview {
	var processedData []DatabaseWorkloadOverview

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		value := *result.MetricDataResults[0].Values[i]
		processedData = append(processedData, DatabaseWorkloadOverview{
			Timestamp: *timestamp,
			Value:     value,
		})
	}

	return processedData
}

func init() {
	AwsxRDSDBLoadCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSDBLoadCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSDBLoadCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSDBLoadCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSDBLoadCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSDBLoadCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSDBLoadCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSDBLoadCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSDBLoadCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSDBLoadCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSDBLoadCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSDBLoadCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSDBLoadCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSDBLoadCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSDBLoadCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSDBLoadCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}

