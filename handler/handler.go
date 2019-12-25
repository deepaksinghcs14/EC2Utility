package handler

import (
	"github.com/aws/aws-lambda-go/events"
	"context"
	"log"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"
	"Ec2Utility/actions"
	"net/http"
	"fmt"
)

func Handler(context context.Context,request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse,error) {
	aws_user := os.Getenv("ACCESS_KEY_ID")
	aws_pass := os.Getenv("SECRET_ACCESS_KEY")
	aws_region := os.Getenv("REGION")
	sess,err := session.NewSession(&aws.Config{
		Region:aws.String(aws_region),
		Credentials: credentials.NewStaticCredentials(aws_user,aws_pass,""),
	})
	if err!=nil{
		log.Println(err)
	}
	svc := ec2.New(sess)
	instanceId := request.QueryStringParameters["instanceid"]
	switch request.QueryStringParameters["action"] {
	case "on":{
		return  actions.OnServer(svc,instanceId)
		}
	case "off":{
		return actions.OffServer(svc,instanceId)
	}
	case "toogle":{
		return actions.ToogleServer(svc,instanceId)
	}

	case "reboot":{
		return actions.RebootServer(svc,instanceId)
	}

	default:{
		return events.APIGatewayProxyResponse{
			Body: actions.CreateResponse(false,fmt.Errorf("Wrong Options")),
			StatusCode: http.StatusNotFound},nil
		}
	}



}

