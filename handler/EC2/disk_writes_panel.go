package EC2

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"log"
)

//type DiskWritePanelData struct {
//	RawData []struct {
//		Timestamp time.Time
//		Value     float64
//	} `json:"Disk_Writes"`
//}

var AwsxEc2DiskWriteCmd = &cobra.Command{
	Use:   "disk_write_panel",
	Short: "get disk write metrics data",
	Long:  `command to get disk write metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetDiskWritePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu utilization: ", err)
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

func GetDiskWritePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = commanFunction.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "DiskWriteBytes", startTime, endTime, "Average", "InstanceId", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting disk write data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Disk_Writes"] = rawData

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxEc2DiskWriteCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
