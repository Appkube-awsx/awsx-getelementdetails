package EC2

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	comman_function "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var AwsxEc2MemoryUtilizationForAllInstancesCmd = &cobra.Command{
	Use:   "memory_utilization_panel",
	Short: "get memory utilization metrics data",
	Long:  `command to get memory utilization metrics data`,
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
			jsonResp, cloudwatchMetricResp, err := GetMemoryUtilizationForAllInstancesPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting memory utilization: ", err)
				fmt.Println("null")
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

func GetMemoryUtilizationForAllInstancesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	// elementType, _ := cmd.PersistentFlags().GetString("elementType")

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(clientAuth.Region),
	})
	if err != nil {
		return "", nil, fmt.Errorf("error creating AWS session: %v", err)
	}

	ec2Svc := ec2.New(sess)

	// Get the list of instances
	instances, err := getAllInstances(ec2Svc)
	fmt.Println("instances", instances)
	if err != nil {
		return "", nil, fmt.Errorf("error listing instances: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}
	jsonOutput := make(map[string]float64)

	for _, instanceId := range instances {
		// Get average utilization
		averageUsage, err := comman_function.GetMetricData(clientAuth, instanceId, "CWAgent", "disk_used_percent", startTime, endTime, "Average", "InstanceId", cloudWatchClient)
		if err != nil {
			log.Printf("Error in getting average for instance %s: %v\n", instanceId, err)
			continue
		}
		if len(averageUsage.MetricDataResults) > 0 && len(averageUsage.MetricDataResults[0].Values) > 0 {
			cloudwatchMetricData[instanceId] = averageUsage
			jsonOutput[instanceId] = *averageUsage.MetricDataResults[0].Values[0]
		} else {
			log.Printf("No data found for average usage for instance %s\n", instanceId)
		}
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func getAllInstances(svc *ec2.EC2) ([]string, error) {
	input := &ec2.DescribeInstancesInput{}
	instanceIds := []string{}

	err := svc.DescribeInstancesPages(input, func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				instanceIds = append(instanceIds, *instance.InstanceId)
			}
		}
		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return instanceIds, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxEc2MemoryUtilizationForAllInstancesCmd)
}
