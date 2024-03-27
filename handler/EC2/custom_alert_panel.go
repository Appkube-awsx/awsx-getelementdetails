package EC2

import (
    "fmt"
    "log"
    "time"

    "github.com/Appkube-awsx/awsx-common/authenticate"
    "github.com/Appkube-awsx/awsx-common/awsclient"
    "github.com/Appkube-awsx/awsx-common/cmdb"
    "github.com/Appkube-awsx/awsx-common/config"
    "github.com/Appkube-awsx/awsx-common/model"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
    "github.com/spf13/cobra"
)

var AwsxEc2CustomAlertPanelCmd = &cobra.Command{
    Use:   "custom_alert_panel",
    Short: "get custom alerts for EC2 security group changes",
    Long:  `command to get custom alerts for EC2 security group changes`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("running from custom alert panel")

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
            cloudwatchMetric, err := GetEc2CustomAlertPanel(cmd, clientAuth, nil)
            if err != nil {
                log.Println("Error getting custom alerts: ", err)
                return
            }
            fmt.Println(cloudwatchMetric)
        }
    },
}

func GetEc2CustomAlertPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
    elementId, _ := cmd.PersistentFlags().GetString("elementId")
    cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
    logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
    if elementId != "" {
        log.Println("getting cloud-element data from cmdb")
        apiUrl := cmdbApiUrl
        if cmdbApiUrl == "" {
            log.Println("using default cmdb url")
            apiUrl = config.CmdbUrl
        }
        log.Println("cmdb url: " + apiUrl)
        cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
        if err != nil {
            return nil, err
        }
        logGroupName = cmdbData.LogGroup

    }
    startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
    endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
    var startTime, endTime *time.Time

    // Parse start time if provided
    if startTimeStr != "" {
        parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
        if err != nil {
            log.Printf("Error parsing start time: %v", err)
            return nil, err
        }
        startTime = &parsedStartTime
    } else {
        defaultStartTime := time.Now().Add(-24 * time.Hour)
        startTime = &defaultStartTime
    }

    if endTimeStr != "" {
        parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
        if err != nil {
            log.Printf("Error parsing end time: %v", err)
            return nil, err
        }
        endTime = &parsedEndTime
    } else {
        defaultEndTime := time.Now()
        endTime = &defaultEndTime
    }

    results, err := filtercloudWatchLogs(clientAuth, startTime, endTime, logGroupName)
    if err != nil {
        log.Println("Error in getting custom alert data: ", err)
        return nil, err
    }

    return results, nil
}

func filtercloudWatchLogs(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
    // Initialize CloudWatch Logs client
    cloudWatchLogs := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)

    // Construct input parameters
    params := &cloudwatchlogs.StartQueryInput{
        LogGroupName: aws.String(logGroupName),
        StartTime:    aws.Int64(startTime.Unix() * 1000),
        EndTime:      aws.Int64(endTime.Unix() * 1000),
        QueryString: aws.String(`fields @timestamp, requestParameters.groupId AS SecurityGroupID,
        if (eventName = 'AuthorizeSecurityGroupIngress' OR eventName = 'AuthorizeSecurityGroupEgress', 'Added', 'Removed') AS Action,
        userIdentity.sessionContext.sessionIssuer.userName AS UserName
        | filter eventSource = 'ec2.amazonaws.com' AND (eventName = 'AuthorizeSecurityGroupIngress' OR eventName = 'RevokeSecurityGroupIngress' OR eventName = 'AuthorizeSecurityGroupEgress' OR eventName = 'RevokeSecurityGroupEgress')
        | sort @timestamp desc`),
    }

    // Start the query
    queryResult, err := cloudWatchLogs.StartQuery(params)
    if err != nil {
        return nil, fmt.Errorf("failed to start query: %v", err)
    }

    queryId := queryResult.QueryId
    queryStatus := ""
    var queryResults []*cloudwatchlogs.GetQueryResultsOutput // Declare queryResults outside the loop
    for queryStatus != "Complete" {
        // Check query status
        queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
            QueryId: queryId,
        }

        queryResult, err := cloudWatchLogs.GetQueryResults(queryStatusInput) // Assign value to queryResults
        if err != nil {
            return nil, fmt.Errorf("failed to get query results: %v", err)
        }

        queryStatus = aws.StringValue(queryResult.Status)
        time.Sleep(1 * time.Second) // Wait for a second before checking status again

        queryResults = append(queryResults, queryResult)
    }

    return queryResults, nil
}

func init() {
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("secretKey", "", "aws secret key")
}
