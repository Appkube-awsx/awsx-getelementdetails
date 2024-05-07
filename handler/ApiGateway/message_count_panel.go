package ApiGateway

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxApiMessageCountCmd = &cobra.Command{

	Use: "message_count_panel",

	Short: "get message count metrics data",

	Long: `command to get message count data`,

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

			panel, err := GetMessageCountPanel(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)

		}

	},
}

func GetMessageCountPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	events, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message| filter eventSource="apigateway.amazonaws.com" | parse @message /"name":\s*"(?<ApiName>[^"]+)"/| stats count(@message) as MessageCount`, cloudWatchLogs)
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		// handle error
	}
	for _, event := range events {
		fmt.Println(event)
	}
	
	return nil, err
}

func init() {
	AwsxApiMessageCountCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	AwsxApiMessageCountCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	AwsxApiMessageCountCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	AwsxApiMessageCountCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApiMessageCountCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApiMessageCountCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApiMessageCountCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApiMessageCountCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxApiMessageCountCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApiMessageCountCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApiMessageCountCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApiMessageCountCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApiMessageCountCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApiMessageCountCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApiMessageCountCmd.PersistentFlags().String("ServiceName", "", "Service Name")
	AwsxApiMessageCountCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApiMessageCountCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApiMessageCountCmd.PersistentFlags().String("clusterName", "", "cluster name")
	AwsxApiMessageCountCmd.PersistentFlags().String("query", "", "query")
	AwsxApiMessageCountCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiMessageCountCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxApiMessageCountCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApiMessageCountCmd.PersistentFlags().String("logGroupName", "", "log group name")
}
