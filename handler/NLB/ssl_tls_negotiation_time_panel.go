package NLB

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

type SSLTLSNegotiationData struct {
	Timestamp             time.Time
	SSLTLSNegotiationData float64
}

type SSLTLSNegotiationDataa struct {
	SSLTLSNegotiationData []SSLTLSNegotiationData `json:"SSLTLSNegotiationData"`
}

var AwsxNLBSSLTLSNegotiationCmd = &cobra.Command{
	Use:   "nlb_ssl_tls_negotiation_panel",
	Short: "Get NLB SSL/TLS negotiation time metrics data",
	Long:  `Command to get NLB SSL/TLS negotiation time metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command..")
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
			rawData, calculatedData, err := GetSSLTLSNegotiationDataData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting NLB SSL/TLS negotiation time metrics: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(calculatedData)
			} else {
				fmt.Println(rawData)
			}
		}

	},
}

func GetSSLTLSNegotiationDataData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	//lbID, _ := cmd.PersistentFlags().GetString("lbID")
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
			return "", nil, err
		}
		instanceId = cmdbData.InstanceId

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

	// Fetch raw data
	rawData, err := GetSSLTLSNegotiationDataMetricData(clientAuth, instanceId, elementType, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}

	// Process the raw data if needed
	result := processssltlsrawdata(rawData)

	// Collect all timestamps and values separately
	timestamps := make([]time.Time, len(result.SSLTLSNegotiationData))
	values := make([]float64, len(result.SSLTLSNegotiationData))

	// Populate the slices with actual data
	for i, data := range result.SSLTLSNegotiationData {
		// Assigning values directly to slices without taking their addresses
		timestamps[i] = data.Timestamp
		values[i] = data.SSLTLSNegotiationData
	}

	// Initialize the MetricDataResults slice
	metricDataResults := make([]*cloudwatch.MetricDataResult, len(result.SSLTLSNegotiationData))

	// Populate the MetricDataResults with actual data
	for i := range result.SSLTLSNegotiationData {
		metricDataResults[i] = &cloudwatch.MetricDataResult{
			Timestamps: []*time.Time{&timestamps[i]},
			Values:     []*float64{&values[i]},
		}
	}

	// Assign the processed data to cloudwatchMetricData
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{
		"SSLTLSNegotiationData": {
			MetricDataResults: metricDataResults,
		},
	}

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetSSLTLSNegotiationDataMetricData(clientAuth *model.Auth, instanceId string, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {

	log.Printf("Getting SSL/TLS negotiation metric data for NLB %s from %v to %v", instanceId, startTime, endTime)

	// Replace this block with actual code to fetch SSL/TLS negotiation metrics from CloudWatch
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("clientTLSNegotiation"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("LoadBalancer"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("ClientTLSNegotiationErrorCount"),
						Namespace:  aws.String("AWS/NetworkELB"),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Average"),
				},
			},
			{
				Id: aws.String("targetTLSNegotiation"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("LoadBalancer"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("TargetTLSNegotiationErrorCount"),
						Namespace:  aws.String("AWS/NetworkELB"),
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

func processssltlsrawdata(result *cloudwatch.GetMetricDataOutput) SSLTLSNegotiationDataa {
	var rawData SSLTLSNegotiationDataa
	rawData.SSLTLSNegotiationData = make([]SSLTLSNegotiationData, len(result.MetricDataResults[0].Timestamps))

	// Assuming the two metrics have the same number of data points
	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.SSLTLSNegotiationData[i].Timestamp = *timestamp

		clientTLSNegotiationErrorCount := *result.MetricDataResults[0].Values[i]
		targetTLSNegotiationErrorCount := *result.MetricDataResults[1].Values[i]

		SSLTLSNegotiationData := clientTLSNegotiationErrorCount + targetTLSNegotiationErrorCount

		rawData.SSLTLSNegotiationData[i].SSLTLSNegotiationData = SSLTLSNegotiationData
	}
	return rawData
}

func init() {
	AwsxNLBSSLTLSNegotiationCmd.PersistentFlags().String("instanceId", "", "instanceId")
	AwsxNLBSSLTLSNegotiationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxNLBSSLTLSNegotiationCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxNLBSSLTLSNegotiationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
