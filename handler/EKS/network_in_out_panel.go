package EKS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type NetworkInOutResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Network in and Network out"`
// }

var AwsxEKSNetworkInOutCmd = &cobra.Command{
	Use:   "Network_in_out_panel",
	Short: "get Network in out graph metrics data",
	Long:  `command to get Network in out graph metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkInOutData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting Network in out data: ", err)
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

func GetNetworkInOutData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = commanFunction.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_network_total_bytes", startTime, endTime, "Sum", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Network in and Network out"] = rawData

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
