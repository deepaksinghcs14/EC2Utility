use below command after setting up go and aws cli

CREATE ONE lambda and set its hook as EC2Utility 
run below command

``env GOOS=linux GOARCH=amd64 go build -o /tmp/EC2Utility . && zip -j /tmp/EC2Utility.zip /tmp/EC2Utility && aws lambda update-function-code --function-name Ec2Utility --zip-file fileb:///tmp/EC2Utility.zip``