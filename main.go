package main

import (

	"github.com/aws/aws-lambda-go/lambda"
	"Ec2Utility/handler"
)

/** this is the main start of our application */
func main() {
	lambda.Start(handler.Handler)
}
