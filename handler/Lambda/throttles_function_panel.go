package Lambda

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/cobra"
)

var AwsxLambdaThrottlesFunctionCmd = &cobra.Command{
	Use:   "throttles_function_panel",
	Short: "get throttles function metrics data",
	Long:  `Command to get throttles function metrics data`,

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
			throttlesFunctionCount, err := GetLambdaThrottlesFunctionData(clientAuth)
			if err != nil {
				log.Println("Error getting throttles function data: ", err)
				return
			}
			if responseType == "frame" {
				// Print cloudwatchMetricResp if necessary
			} else {
				fmt.Println("Throttles Function Count:", throttlesFunctionCount)
			}
		}
	},
}

func GetLambdaThrottlesFunctionData(clientAuth *model.Auth) (int, error) {
	lambdaClient := awsclient.GetClient(*clientAuth, awsclient.LAMBDA_CLIENT).(*lambda.Lambda)
	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	input := &lambda.ListFunctionsInput{}

	throttlesFunctionCount := 0
	err := lambdaClient.ListFunctionsPages(input,
		func(page *lambda.ListFunctionsOutput, lastPage bool) bool {
			for _, function := range page.Functions {
				throttles, err := GetFunctionThrottlesCount(cloudWatchClient, function.FunctionName)
				if err != nil {
					log.Printf("Error getting throttles count for function %s: %v", *function.FunctionName, err)
					continue
				}
				throttlesFunctionCount += throttles
			}
			return !lastPage
		})
	if err != nil {
		return 0, err
	}

	return throttlesFunctionCount, nil
}

func GetFunctionThrottlesCount(cloudWatchClient *cloudwatch.CloudWatch, functionName *string) (int, error) {
	now := time.Now()
	// Get the throttles count for the last 5 minutes
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Throttles"),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("FunctionName"),
				Value: functionName,
			},
		},
		StartTime:  aws.Time(now.Add(-5 * time.Minute)),
		EndTime:    aws.Time(now),
		Period:     aws.Int64(300), // 5 minutes
		Statistics: []*string{aws.String("Sum")},
	}

	resp, err := cloudWatchClient.GetMetricStatistics(input)
	if err != nil {
		return 0, err
	}

	throttlesCount := 0
	for _, dp := range resp.Datapoints {
		if dp.Sum != nil {
			throttlesCount += int(*dp.Sum)
		}
	}

	return throttlesCount, nil
}

func init() {
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaThrottlesFunctionCmd.PersistentFlags().String("endTime", "", "end time")
}
