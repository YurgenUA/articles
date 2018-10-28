#!/bin/bash 

STAGE=demo
CLIENTID=$1
CLIENTSECRET=$2


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
