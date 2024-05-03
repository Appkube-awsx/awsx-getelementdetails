package Lambda

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/cobra"
)

var AwsxLambdaFunctionsByRegionCmd = &cobra.Command{
	Use:   "functions_by_region_panel",
	Short: "get Lambda functions by region",
	Long:  `command to get Lambda functions by region`,

	Run: func(cmd *cobra.Command, args []string) {
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
			jsonResp, functionCounts, err := GetLambdaFunctionsByRegion(clientAuth)
			if err != nil {
				log.Println("Error getting Lambda functions by region: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(functionCounts)
			} else {
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetLambdaFunctionsByRegion(clientAuth *model.Auth) (string, map[string]interface{}, error) {
	cloudwatchMetricData := make(map[string]interface{})

	if clientAuth == nil {
		log.Println("Error: Authentication failed. Client authentication credentials are nil.")
		return "", nil, errors.New("authentication failed: clientAuth is nil")
	}

	// List of AWS regions
	regions := []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "eu-west-1", "eu-west-2", "ap-northeast-1"}

	// Initialize total functions count
	totalFunctions := 0

	for _, region := range regions {
		newAuth := model.Auth{
			AccessKey:           clientAuth.AccessKey,
			SecretKey:           clientAuth.SecretKey,
			CrossAccountRoleArn: clientAuth.CrossAccountRoleArn,
			ExternalId:          clientAuth.ExternalId,
			Region:              region, // Use the region from the list
		}

		// Get Lambda client for the current region
		lambdaClient := awsclient.GetClient(newAuth, awsclient.LAMBDA_CLIENT).(*lambda.Lambda)

		// Get total Lambda functions count for the current region
		count, err := getTotalLambdaFunctions(lambdaClient)
		if err != nil {
			log.Printf("Error getting total functions in region %s: %v", region, err)
			continue
		}
		cloudwatchMetricData[region] = count
		totalFunctions += count
	}

	// Add total functions count to the map
	cloudwatchMetricData["TotalFunctions"] = totalFunctions

	// Marshal the function counts map to JSON string
	jsonString, err := json.Marshal(cloudwatchMetricData)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	// Construct JSON response
	jsonResp := string(jsonString)

	return jsonResp, cloudwatchMetricData, nil
}

// Function to get total Lambda functions in a region
func getTotalLambdaFunctions(lambdaClient *lambda.Lambda) (int, error) {
	totalFunctions := 0

	input := &lambda.ListFunctionsInput{}

	for {
		resp, err := lambdaClient.ListFunctions(input)
		if err != nil {
			return 0, err
		}

		totalFunctions += len(resp.Functions)

		if resp.NextMarker != nil {
			input.Marker = resp.NextMarker
		} else {
			break
		}
	}

	return totalFunctions, nil
}

func init() {
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxLambdaFunctionsByRegionCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
