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

var AwsxRDSDBLoadNonCPUCmd = &cobra.Command{
	Use:   "db_load_non_cpu_panel",
	Short: "get non-cpu load in database operations",
	Long:  `command to get non-cpu load in database operations`,

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
			jsonResp, cloudwatchMetricResp, err, _ := GetRDSDBLoadNonCPU(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting non-cpu load data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// default case. it prints json
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetRDSDBLoadNonCPU(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
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

	// Fetch raw data for non-cpu load metric
	rawData, err := GetMetricData(clientAuth, elementType, startTime, endTime, "DBLoadNonCPU", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting non-cpu load data: ", err)
		return "", "", nil, err
	}
	cloudwatchMetricData["DBLoadNonCPU"] = rawData

	// Process raw data
	result := processedRawDataa(rawData)
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json for non-cpu load data: ", err)
		return "", "", nil, err
	}

	return string(jsonData), "", cloudwatchMetricData, nil
}

func GetMetricData(clientAuth *model.Auth, elementType string, startTime, endTime *time.Time, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace AWS/RDS from %v to %v", elementType, startTime, endTime)

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{},
						MetricName: aws.String(metricName),
						Namespace:  aws.String("AWS/RDS"),
					},
					Period: aws.Int64(60),
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

func processedRawDataa(result *cloudwatch.GetMetricDataOutput) []struct {
	Timestamp time.Time
	Value     float64
} {
	var processedData []struct {
		Timestamp time.Time
		Value     float64
	}

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		value := *result.MetricDataResults[0].Values[i]
		processedData = append(processedData, struct {
			Timestamp time.Time
			Value     float64
		}{Timestamp: *timestamp, Value: value})
	}

	return processedData
}



func init() {
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}

