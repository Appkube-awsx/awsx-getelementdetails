package EC2

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var AwsxEC2DiskReadOpsPerInstanceTypeCommmand = &cobra.Command{
	Use:   "disk_read_ops_per_type",
	Short: "get disk read ops per instance type metrics data",
	Long:  `command to disk read ops per instance type metrics data`,
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
			jsonResp, resp, err := DiskReadOpsPerInstanceType(cmd, clientAuth, nil, nil)
			if err != nil {
				log.Println("Error getting disk read ops per instance type data : ", err)
				return
			}
			if responseType == "json" {
				fmt.Println(jsonResp)
			} else {
				fmt.Println(resp)
			}
		}
	},
}

func DiskReadOpsPerInstanceType(cmd *cobra.Command, clientAuth *model.Auth, ec2Client *ec2.EC2, cloudWatchClient *cloudwatch.CloudWatch) (string, []Ec2DiskReadOpsResult, error) {
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
	if ec2Client == nil {
		ec2Client = awsclient.GetClient(*clientAuth, awsclient.EC2_CLIENT).(*ec2.EC2)
	}
	ec2Input := ec2.DescribeInstancesInput{}
	instancesResult, err := ec2Client.DescribeInstances(&ec2Input)
	if err != nil {
		return "", nil, fmt.Errorf("Error getting disk read ops per instance type data")
	}
	var instances []Ec2InstanceOutputData
	for _, reserv := range instancesResult.Reservations {
		for _, instance := range reserv.Instances {
			temp := Ec2InstanceOutputData{
				InstanceType: *instance.InstanceType,
				InstanceId:   *instance.InstanceId,
			}
			instances = append(instances, temp)
		}
	}
	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}
	var wg sync.WaitGroup

	ch := make(chan Ec2DiskReadOpsResult)

	for _, instance := range instances {
		wg.Add(1)
		go getDiskReadOps(cloudWatchClient, instance, startTime, endTime, &wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var data []Ec2DiskReadOpsResult
	for result := range ch {
		data = append(data, result)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", nil, fmt.Errorf("Error getting disk read ops per instance type data")
	}
	return string(jsonData), data, nil
}

type Ec2DiskReadOpsResult struct {
	InstanceType string
	items        interface{}
}

func getDiskReadOps(cloudWatchClient *cloudwatch.CloudWatch, instance Ec2InstanceOutputData, startTime, endTime *time.Time, wg *sync.WaitGroup, ch chan<- Ec2DiskReadOpsResult) {
	defer wg.Done()

	cwInput := cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("diskReadOps"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("AWS/EC2"),
						MetricName: aws.String("DiskReadOps"),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("InstanceId"),
								Value: aws.String(instance.InstanceId),
							},
						},
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Average"),
				},
				ReturnData: aws.Bool(true),
			},
		},
		StartTime: aws.Time(*startTime),
		EndTime:   aws.Time(*endTime),
	}
	result, err := cloudWatchClient.GetMetricData(&cwInput)
	if err != nil {
		log.Printf("internal server error : %w", err)
	}
	dataMap := make(map[*time.Time]*float64)
	for i := 0; i < len(result.MetricDataResults[0].Timestamps); i++ {
		k := result.MetricDataResults[0].Timestamps[i]
		v := result.MetricDataResults[0].Values[i]
		dataMap[k] = v
	}
	ch <- Ec2DiskReadOpsResult{
		InstanceType: instance.InstanceType,
		items:        dataMap,
	}
}
