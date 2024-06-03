package EC2

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/spf13/cobra"
)

type AutoScalingCounts struct {
	AutoScalingActiveCounts int
	LaunchConfig            int
}

var autoScalingConfigPanelCmd = &cobra.Command{
	Use:   "auto_scaling_config_panel",
	Short: "gets auto scaling active counts and launch config count",
	Long:  `Command to get auto scaling active counts and launch config count `,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")
		authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
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
			jsonResp, cloudwatchMetricResp, err := GetAutoScalingInfo(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting  CLI data: ", err)
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

func GetAutoScalingInfo(cmd *cobra.Command, clientAuth *model.Auth, autoScalingClient *autoscaling.AutoScaling) (string, []AutoScalingCounts, error) {
	if autoScalingClient == nil {
		autoScalingClient = awsclient.GetClient(*clientAuth, awsclient.AUTOSCALING_CLIENT).(*autoscaling.AutoScaling)
	}
	// sess := session.Must(session.NewSession(&aws.Config{
	// 	Region: aws.String("us-east-1"), // Specify your AWS region
	// }))
	// ec2Client = autoscaling.New(sess)
	// autoscalinginput := &autoscaling.DescribeAutoScalingGroupsInput{}
	launchconfigInput := &autoscaling.DescribeLaunchConfigurationsInput{}

	// result, err := ec2Client.DescribeAutoScalingGroups(autoscalinginput)
	// if err != nil {
	// 	log.Println("Error describing auto scaling groups:", err)
	// 	return "", nil, err
	// }
	// autoscalingstatuscount, err := json.Marshal(result)
	// if err != nil {
	// 	fmt.Println("Error marshalling result:", err)
	// }
	launchconfigcount, err := autoScalingClient.DescribeLaunchConfigurations(launchconfigInput)
	if err != nil {
		log.Println("Error describing launch configurations:", err)
		return "", nil, err
	}
	//log.Println("launchconfigcount: ", launchconfigcount)

	// activeCount := 0
	// for _, asg := range autoscalingstatuscount.AutoScalingGroups {
	// 	if asg.Status != nil && *asg.Status == "Active" {
	// 		activeCount++
	// 	}
	// }

	instanceData := AutoScalingCounts{
		// AutoScalingActiveCounts: activeCount,
		LaunchConfig: len(launchconfigcount.LaunchConfigurations),
	}

	//instances := []AutoScalingCounts{instanceData}
	jsonResp, _ := json.Marshal(instanceData)
	return string(jsonResp), nil, err
}

func init() {
	autoScalingConfigPanelCmd.PersistentFlags().String("region", "", "AWS region to use")
}
