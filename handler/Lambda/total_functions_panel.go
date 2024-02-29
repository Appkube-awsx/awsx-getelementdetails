package Lambda

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/cobra"
)

type TotalFunctionResult struct {
	Value float64 `json:"Value"`
}

var AwsxLambdaTotalFunctionCmd = &cobra.Command{
	Use:   "total_function_panel",
	Short: "get total function metrics data",
	Long:  `command to get total function metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaTotalFunctionData(clientAuth, nil)
			if err != nil {
				log.Println("Error getting total function data : ", err)
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

func GetLambdaTotalFunctionData(clientAuth *model.Auth, lambdaClient *lambda.Lambda) (string, map[string]float64, error) {
	cloudwatchMetricData := map[string]float64{}

	totalFunctions, err := GetTotalLambdaFunctions(clientAuth, lambdaClient)
	if err != nil {
		log.Println("Error in getting total functions: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["TotalFunctions"] = float64(totalFunctions)

	log.Printf("Total Functions: %d", totalFunctions)

	jsonString, err := json.Marshal(TotalFunctionResult{Value: float64(totalFunctions)})
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetTotalLambdaFunctions(clientAuth *model.Auth, lambdaClient *lambda.Lambda) (int, error) {
	if lambdaClient == nil {
		lambdaClient = awsclient.GetClient(*clientAuth, awsclient.LAMBDA_CLIENT).(*lambda.Lambda)
	}

	totalFunctions := 0

	input := &lambda.ListFunctionsInput{}

	for {
		resp, err := lambdaClient.ListFunctions(input)
		if err != nil {
			return 0, err
		}

		totalFunctions += len(resp.Functions)

		// If there are more functions to fetch, update the marker
		if resp.NextMarker != nil {
			input.Marker = resp.NextMarker
		} else {
			break
		}
	}

	return totalFunctions, nil
}

func init() {
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}


// package Lambda

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/Appkube-awsx/awsx-common/authenticate"
// 	"github.com/Appkube-awsx/awsx-common/awsclient"
// 	"github.com/Appkube-awsx/awsx-common/model"
// 	// "github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/lambda"
// 	"github.com/spf13/cobra"
// )

// type TotalFunctionResult struct {
// 	Value float64 `json:"Value"`
// }

// var AwsxLambdaTotalFunctionCmd = &cobra.Command{
// 	Use:   "total_function_panel",
// 	Short: "get total function metrics data",
// 	Long:  `command to get total function metrics data`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("running from child command")
// 		var authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)
// 		if err != nil {
// 			log.Printf("Error during authentication: %v\n", err)
// 			err := cmd.Help()
// 			if err != nil {
// 				return
// 			}
// 			return
// 		}
// 		if authFlag {
// 			responseType, _ := cmd.PersistentFlags().GetString("responseType")
// 			jsonResp, cloudwatchMetricResp, err := GetLambdaTotalFunctionData(cmd, clientAuth, nil)
// 			if err != nil {
// 				log.Println("Error getting total function data : ", err)
// 				return
// 			}
// 			if responseType == "frame" {
// 				fmt.Println(cloudwatchMetricResp)
// 			} else {
// 				fmt.Println(jsonResp)
// 			}
// 		}

// 	},
// }

// func GetLambdaTotalFunctionData(cmd *cobra.Command, clientAuth *model.Auth, lambdaClient *lambda.Lambda) (string, map[string]float64, error) {
// 	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
// 	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

// 	var startTime, endTime *time.Time

// 	// Parse start time if provided
// 	if startTimeStr != "" {
// 		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
// 		if err != nil {
// 			log.Printf("Error parsing start time: %v", err)
// 			return "", nil, err
// 		}
// 		startTime = &parsedStartTime
// 	} else {
// 		defaultStartTime := time.Now().Add(-5 * time.Minute)
// 		startTime = &defaultStartTime
// 	}

// 	if endTimeStr != "" {
// 		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
// 		if err != nil {
// 			log.Printf("Error parsing end time: %v", err)
// 			return "", nil, err
// 		}
// 		endTime = &parsedEndTime
// 	} else {
// 		defaultEndTime := time.Now()
// 		endTime = &defaultEndTime
// 	}

// 	// Debug prints
// 	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

// 	cloudwatchMetricData := map[string]float64{}

// 	// Fetch raw data
// 	totalFunctions, err := GetTotalLambdaFunctions(clientAuth, lambdaClient)
// 	if err != nil {
// 		log.Println("Error in getting total functions: ", err)
// 		return "", nil, err
// 	}
// 	cloudwatchMetricData["TotalFunctions"] = float64(totalFunctions)

// 	// Debug prints
// 	log.Printf("Total Functions: %d", totalFunctions)

// 	jsonString, err := json.Marshal(TotalFunctionResult{Value: float64(totalFunctions)})
// 	if err != nil {
// 		log.Println("Error in marshalling json in string: ", err)
// 		return "", nil, err
// 	}

// 	return string(jsonString), cloudwatchMetricData, nil
// }

// func GetTotalLambdaFunctions(clientAuth *model.Auth, lambdaClient *lambda.Lambda) (int, error) {
// 	if lambdaClient == nil {
// 		lambdaClient = awsclient.GetClient(*clientAuth, awsclient.LAMBDA_CLIENT).(*lambda.Lambda)
// 	}

// 	input := &lambda.ListFunctionsInput{}

// 	resp, err := lambdaClient.ListFunctions(input)
// 	if err != nil {
// 		return 0, err
// 	}

// 	totalFunctions := len(resp.Functions)

// 	return totalFunctions, nil
// }

// func init() {
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("elementId", "", "element id")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("elementType", "", "element type")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("query", "", "query")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("vaultToken", "", "vault token")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("zone", "", "aws region")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("accessKey", "", "aws access key")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("secretKey", "", "aws secret key")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("externalId", "", "aws external id")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("instanceId", "", "instance id")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("startTime", "", "start time")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("endTime", "", "end time")
// 	AwsxLambdaTotalFunctionCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
// }
