package Lambda

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxLambdaIdleFunctionCmd = &cobra.Command{
	Use:   "idle_function_panel",
	Short: "get idle function metrics data",
	Long:  `Command to get idle function metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")
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
			idleFunctionCount, cloudwatchMetricResp, err := GetLambdaIdleFunctionData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting idle function data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
				// Print cloudwatchMetricResp if necessary
			} else {
				fmt.Println("Idle Function Count:", idleFunctionCount)
			}
		}
	},
}

func GetLambdaIdleFunctionData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	//func GetLambdaIdleFunctionData(clientAuth *model.Auth, lambdaClient *lambda.Lambda) (int, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = comman_function.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}
	idleFunctionCount, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/Lambda", "Invocations", startTime, endTime, "Sum", "FunctionName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting idle function count: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Idle_Functions"] = idleFunctionCount
	return "", cloudwatchMetricData, nil
}

//     idleFunctionCount, err := GetIdleLambdaFunctionCount(clientAuth, lambdaClient)
// 	if err != nil {
// 		log.Println("Error in getting idle function count: ", err)
// 		return 0, err
// 	}
// 	return idleFunctionCount, nil
// }

// func GetIdleLambdaFunctionCount(clientAuth *model.Auth, lambdaClient *lambda.Lambda) (int, error) {
// 	if lambdaClient == nil {
// 		lambdaClient = awsclient.GetClient(*clientAuth, awsclient.LAMBDA_CLIENT).(*lambda.Lambda)
// 	}

// 	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

// 	input := &lambda.ListFunctionsInput{}

// 	idleFunctionCount := 0

// 	for {
// 		resp, err := lambdaClient.ListFunctions(input)
// 		if err != nil {
// 			return 0, err
// 		}

// 		for _, function := range resp.Functions {
// 			// Get the invocation count for the last 5 minutes
// 			input := &cloudwatch.GetMetricStatisticsInput{
// 				Namespace:  aws.String("AWS/Lambda"),
// 				MetricName: aws.String("Invocations"),
// 				Dimensions: []*cloudwatch.Dimension{
// 					{
// 						Name:  aws.String("FunctionName"),
// 						Value: function.FunctionName,
// 					},
// 				},
// 				StartTime:  aws.Time(time.Now().Add(-5 * time.Minute)),
// 				EndTime:    aws.Time(time.Now()),
// 				Period:     aws.Int64(300), // 5 minutes
// 				Statistics: []*string{aws.String("Sum")},
// 			}
// 			// fmt.Println(function)
// 			resp, err := cloudWatchClient.GetMetricStatistics(input)
// 			if err != nil {
// 				log.Println("Error getting metric statistics: ", err)
// 				continue
// 			}

// 			// If there are no invocations in the last 5 minutes, consider the function idle
// 			if len(resp.Datapoints) == 0 || *resp.Datapoints[0].Sum == 0 {
// 				idleFunctionCount++
// 			}
// 		}

// 		// Check if there are more functions to fetch
// 		if resp.NextMarker != nil {
// 			input.Marker = resp.NextMarker
// 		} else {
// 			break // No more functions to fetch
// 		}
// 	}

// 	return idleFunctionCount, nil
// }

func init() {
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaIdleFunctionCmd.PersistentFlags().String("endTime", "", "end time")
}
