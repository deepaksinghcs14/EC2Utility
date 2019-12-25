package actions

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"log"
	"net/http"
	"Ec2Utility/models"
	"encoding/json"
	"fmt"
)

func OnServer(svc *ec2.EC2,instanceId string)  (events.APIGatewayProxyResponse,error){
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
		DryRun: aws.Bool(true),
	}
	result, err := svc.StartInstances(input)
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		// Let's now set dry run to be false. This will allow us to start the instances
		input.DryRun = aws.Bool(false)
		result, err = svc.StartInstances(input)
		if err != nil {
			log.Println("Error", err)
			return events.APIGatewayProxyResponse{
				Body:CreateResponse(false,err),
				StatusCode:http.StatusInternalServerError,
			},nil
		} else {
			log.Println("Success", result.StartingInstances)
			return events.APIGatewayProxyResponse{
				Body:CreateResponse(true,fmt.Errorf("")),
				StatusCode:http.StatusOK,
			},nil
		}
	} else { // This could be due to a lack of permissions
		log.Println("Error", err)
		return events.APIGatewayProxyResponse{
			Body:CreateResponse(false,err),
			StatusCode:http.StatusInternalServerError,
		},nil
	}

}


func OffServer(svc *ec2.EC2,instanceId string)  (events.APIGatewayProxyResponse,error){
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
		DryRun: aws.Bool(true),
	}
	result, err := svc.StopInstances(input)

	awsErr,ok := err.(awserr.Error)
	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		result, err = svc.StopInstances(input)
		if err != nil {
			log.Println("Error", err)
			return events.APIGatewayProxyResponse{
				Body:CreateResponse(false,err),
				StatusCode:http.StatusInternalServerError,
			},nil
		} else {
			log.Println("Success", result.StoppingInstances)
			return events.APIGatewayProxyResponse{
				Body:CreateResponse(true,fmt.Errorf("")),
				StatusCode:http.StatusOK,
			},nil
		}
	} else {
		log.Println("Error", err)
		return events.APIGatewayProxyResponse{
			Body:CreateResponse(false,err),
			StatusCode:http.StatusInternalServerError,
		},nil
	}
}

func RebootServer(svc *ec2.EC2,instanceId string)  (events.APIGatewayProxyResponse,error){
	input := &ec2.RebootInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
		DryRun: aws.Bool(true),
	}
	result, err := svc.RebootInstances(input)
	awsErr, ok := err.(awserr.Error)

	// If the error code is `DryRunOperation` it means we have the necessary
	// permissions to Start this instance
	if ok && awsErr.Code() == "DryRunOperation" {
		// Let's now set dry run to be false. This will allow us to reboot the instances
		input.DryRun = aws.Bool(false)
		result, err = svc.RebootInstances(input)
		if err != nil {
			log.Println("Error", err)
			return events.APIGatewayProxyResponse{
				Body:CreateResponse(false,err),
				StatusCode:http.StatusInternalServerError,
			},nil
		} else {
			log.Println("Success", result)
			return events.APIGatewayProxyResponse{
				Body:CreateResponse(true,fmt.Errorf("")),
				StatusCode:http.StatusOK,
			},nil
		}
	} else { // This could be due to a lack of permissions
		log.Println("Error", err)
		return events.APIGatewayProxyResponse{
			Body:CreateResponse(false,err),
			StatusCode:http.StatusInternalServerError,
		},nil
	}
}

func ToogleServer(svc *ec2.EC2,instanceId string)  (events.APIGatewayProxyResponse,error){
	params := &ec2.DescribeInstancesInput{
		InstanceIds:[]*string{
			aws.String(instanceId),
		},
	}
	resp,err:= svc.DescribeInstances(params)
	for idx, res := range resp.Reservations {
		log.Println("  > Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
		for _, inst := range resp.Reservations[idx].Instances {
			log.Println("    - Instance ID: ", *inst.InstanceId)
			if *inst.State.Name == "running"{
				return OffServer(svc,instanceId)
			}else{
				return OnServer(svc,instanceId)
			}
		}
	}
	return events.APIGatewayProxyResponse{
		Body:CreateResponse(false,err),
		StatusCode:http.StatusInternalServerError,
	},nil
}

func CreateResponse(success bool,error error)  string{
	response:= models.Response_DTO{
		Success: success,
		Error:   error.Error(),
	}
	res,_ := json.Marshal(response)
	return string(res)
}
