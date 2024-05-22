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

var AwsxCpuUtilizationPerInstanceTypeCommmand = &cobra.Command{
	Use:   "cpu_utilization_per_instance_type",
	Short: "get cpu utilization per instance type metrics data",
	Long:  `command to cpu utilization per instance type metrics data`,
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
			jsonResp, resp, err := CpuUtilizationPerInstanceType(cmd, clientAuth, nil, nil)
			if err != nil {
				log.Println("Error getting cpu utilization per instance type data : ", err)
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

func CpuUtilizationPerInstanceType(cmd *cobra.Command, clientAuth *model.Auth, ec2Client *ec2.EC2, cloudWatchClient *cloudwatch.CloudWatch) (string, string, error) {
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return "", "", err
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
			return "", "", err
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
		log.Printf("Error getting cpu utilization per instance type data")
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
	// data := make(map[string]int)
	// data["full_concurrency"] = fullConcurrency
	fmt.Println("instances : ", instances)
	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}
	var wg sync.WaitGroup

	ch := make(chan Ec2CpuUtilizationResult)

	for _, instance := range instances {
		wg.Add(1)
		go getCpuUtilization(cloudWatchClient, instance, startTime, endTime, &wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var data []Ec2CpuUtilizationResult
	for result := range ch {
		data = append(data, result)
	}

	jsonData, err := json.Marshal(data)
	fmt.Println("data : ", jsonData)
	if err != nil {
		log.Printf("error parsing data: %s", err)
		return "", "", err
	}
	return string(jsonData), string(jsonData), nil
}

type Ec2InstanceOutputData struct {
	InstanceType string
	InstanceId   string
}

type Ec2CpuUtilizationResult struct {
	InstanceType string
	items        interface{}
}

func getCpuUtilization(cloudWatchClient *cloudwatch.CloudWatch, instance Ec2InstanceOutputData, startTime, endTime *time.Time, wg *sync.WaitGroup, ch chan<- Ec2CpuUtilizationResult) {
	defer wg.Done()

	cwInput := cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("CPUUtilization"),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("InstanceId"),
				Value: aws.String(instance.InstanceId),
			},
		},
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     aws.Int64(300), // 5-minute intervals
		Statistics: []*string{aws.String("Average")},
	}
	result, err := cloudWatchClient.GetMetricStatistics(&cwInput)
	if err != nil {
		log.Printf("internal server error : %w", err)
	}
	ch <- Ec2CpuUtilizationResult{
		InstanceType: instance.InstanceType,
		items:        result.Datapoints,
	}
}
