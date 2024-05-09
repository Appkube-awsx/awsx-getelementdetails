package comman_function

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/spf13/cobra"
)

func ParseTimes(cmd *cobra.Command) (*time.Time, *time.Time, error) {
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return nil, nil, err
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			return nil, nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	return startTime, endTime, nil
}

func GetCmdbData(cmd *cobra.Command) (string, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")

	if elementId != "" {
		log.Println("Getting cloud-element data from CMDB")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("Using default CMDB URL")
			apiUrl = config.CmdbUrl
		}
		log.Println("CMDB URL: " + apiUrl)
		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		if err != nil {
			return "", fmt.Errorf("error getting cloud element data: %v", err)
		}
		return cmdbData.InstanceId, nil
	}

	return "", errors.New("element ID is required")
}

func GetCmdbLogsData(cmd *cobra.Command) (string, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")

	if elementId != "" {
		log.Println("Getting cloud-element data from CMDB")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("Using default CMDB URL")
			apiUrl = config.CmdbUrl
		}
		log.Println("CMDB URL: " + apiUrl)
		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		if err != nil {
			return "", fmt.Errorf("error getting cloud element data: %v", err)
		}
		return cmdbData.LogGroup, nil
	}

	return "", errors.New("element ID is required")
}

func InitAwsCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	cmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	cmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	cmd.PersistentFlags().String("elementId", "", "element id")
	cmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	cmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	cmd.PersistentFlags().String("vaultToken", "", "vault token")
	cmd.PersistentFlags().String("accountId", "", "aws account number")
	cmd.PersistentFlags().String("zone", "", "aws region")
	cmd.PersistentFlags().String("accessKey", "", "aws access key")
	cmd.PersistentFlags().String("secretKey", "", "aws secret key")
	cmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	cmd.PersistentFlags().String("externalId", "", "aws external id")
	cmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	cmd.PersistentFlags().String("ServiceName", "", "Service Name")
	cmd.PersistentFlags().String("elementType", "", "element type")
	cmd.PersistentFlags().String("instanceId", "", "instance id")
	cmd.PersistentFlags().String("clusterName", "", "cluster name")
	cmd.PersistentFlags().String("query", "", "query")
	cmd.PersistentFlags().String("startTime", "", "start time")
	cmd.PersistentFlags().String("endTime", "", "endcl time")
	cmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	cmd.PersistentFlags().String("logGroupName", "", "log group name")
	cmd.PersistentFlags().String("ApiName", "", "api name")
	cmd.PersistentFlags().String("FunctionName", "", "function name")
	cmd.PersistentFlags().String("LoadBalancer", "", "loadbalancer name")
	cmd.PersistentFlags().String("DBInstanceIdentifier", "", "dbinstance identifier name")

}
