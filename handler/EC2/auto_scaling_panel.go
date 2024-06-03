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
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/spf13/cobra"
)

type AutoScalingGroupDetails struct {
	AutoScalingGroupName    string `json:"AutoScalingGroupName"`
	LaunchConfigurationName string `json:"LaunchConfigurationName"`
	InstanceType            string `json:"InstanceType"`
	MinSize                 int64  `json:"MinSize"`
	MaxSize                 int64  `json:"MaxSize"`
	DesiredCapacity         int64  `json:"DesiredCapacity"`
	HealthCheckType         string `json:"HealthCheckType"`
}

var AwsxAutoScalingGroupsCmd = &cobra.Command{
	Use:   "autoscaling_groups",
	Short: "get autoscaling groups details",
	Long:  `command to get autoscaling groups details`,

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
			jsonResp, autoScalingGroupsResp, err := GetAutoScalingGroupsDetails(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting autoscaling groups details: ", err)
				return
			}
			if responseType == "frame" {
				for _, detail := range autoScalingGroupsResp {
					fmt.Printf("%+v\n", *detail)
				}
			} else {
				fmt.Println(jsonResp)
			}
		}
	},
}

func GetAutoScalingGroupsDetails(cmd *cobra.Command, clientAuth *model.Auth, autoScalingClient *autoscaling.AutoScaling) (string, []*AutoScalingGroupDetails, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Specify the region
	})
	if err != nil {
		return "", nil, fmt.Errorf("error creating AWS session: %v", err)
	}
	if autoScalingClient == nil {
		autoScalingClient = autoscaling.New(sess)
	}

	autoScalingGroupsInput := &autoscaling.DescribeAutoScalingGroupsInput{}
	autoScalingGroupsOutput, err := autoScalingClient.DescribeAutoScalingGroups(autoScalingGroupsInput)
	if err != nil {
		return "", nil, fmt.Errorf("error describing Auto Scaling groups: %v", err)
	}

	var autoScalingGroupDetailsList []*AutoScalingGroupDetails

	for _, group := range autoScalingGroupsOutput.AutoScalingGroups {
		for _, instance := range group.Instances {
			details := &AutoScalingGroupDetails{
				AutoScalingGroupName:    aws.StringValue(group.AutoScalingGroupName),
				LaunchConfigurationName: aws.StringValue(group.LaunchConfigurationName),
				InstanceType:            aws.StringValue(instance.InstanceType),
				MinSize:                 aws.Int64Value(group.MinSize),
				MaxSize:                 aws.Int64Value(group.MaxSize),
				DesiredCapacity:         aws.Int64Value(group.DesiredCapacity),
				HealthCheckType:         aws.StringValue(group.HealthCheckType),
			}
			autoScalingGroupDetailsList = append(autoScalingGroupDetailsList, details)
		}
	}

	jsonString, err := json.Marshal(autoScalingGroupDetailsList)
	if err != nil {
		log.Println("Error in marshalling JSON to string: ", err)
		return "", nil, err
	}

	return string(jsonString), autoScalingGroupDetailsList, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxAutoScalingGroupsCmd)

}
