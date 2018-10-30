#!/bin/bash 

STAGE=demo
CLIENTID=$1
CLIENTSECRET=$2
AWSACCOUNT=$3


echo setup SecretsManager secret for Auth0 secrets....
AUTH0_SECRET=auth0-secrets-$STAGE
if aws secretsmanager list-secrets --region eu-west-1 2>&1 | grep -q $AUTH0_SECRET
then 
  echo secret $AUTH0_SECRET found
else
  echo secret $AUTH0_SECRET not found. Creating it
  aws secretsmanager create-secret --region eu-west-1 \
      --name $AUTH0_SECRET \
      --description "Store Auth0 connection sensitive data '$STAGE'" \
      --secret-string '{{"CLIENT_ID":'$CLIENTID'},{"CLIENT_SECRET":'$CLIENTSECRET'}}'
fi
echo finished
sleep 3
##################################

TOKEN_SECRET=service-token-$STAGE

echo setup SecretsManager secret for Token....
if aws secretsmanager list-secrets --region eu-west-1 2>&1 | grep -q $TOKEN_SECRET
then 
  echo secret $TOKEN_SECRET found
else
  echo secret $TOKEN_SECRET not found. Creating it
  aws secretsmanager create-secret --region eu-west-1 --name $TOKEN_SECRET \
      --description "Service bearer token rotated on daily basis '$STAGE'" \
      --secret-string [{"token":""},{"rotatedAt":""}]
fi

echo grant AWS SecretsManager permission to call Lambda
aws lambda add-permission --region eu-west-1\
          --function-name medium-secrets-$STAGE-rotate \
          --principal secretsmanager.amazonaws.com \
          --action lambda:InvokeFunction \
          --statement-id SecretsManagerAccess
{
    "Statement": "{\"Sid\":\"SecretsManagerAccess\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"secretsmanager.amazonaws.com\"},\"Action\":\"lambda:InvokeFunction\",\"Resource\":\"arn:aws:lambda:eu-west-1:$AWSACCOUNT:function:medium-secrets-$STAGE-rotate\"}"
} 

echo bind rotating Lambda to to secret "$TOKEN_SECRET" and push rotation...
aws secretsmanager rotate-secret --region eu-west-1 --secret-id service-token-$STAGE \
  --rotation-lambda-arn arn:aws:lambda:eu-west-1:$AWSACCOUNT:function:medium-secrets-$STAGE-rotate \
  --rotation-rules AutomaticallyAfterDays=1
