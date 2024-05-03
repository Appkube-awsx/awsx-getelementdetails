package EKS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type DiskIOPerformanceResult struct {
//     TotalOps []struct {
//         Timestamp time.Time
//         Value     float64
//     } `json:"total_ops"`
// }

var AwsxEKSDiskIOPerformanceCmd = &cobra.Command{
	Use:   "disk_io_performance_panel",
	Short: "get disk I/O performance metrics data",
	Long:  `command to get disk I/O performance metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetEKSDiskIOPerformancePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting disk I/O performance metrics: ", err)
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

func GetEKSDiskIOPerformancePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = comman_function.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_diskio_io_serviced_total", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error fetching total operations raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["TotalOps"] = rawData

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
