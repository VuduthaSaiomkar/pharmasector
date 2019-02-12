#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

jq --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
	echo
	exit 1
fi

echo "POST request Enroll on Org1  ..."
echo
export manu_Token=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Jim&orgName=manufacture')
echo $manu_Token
export manu_Token=$(echo $manu_Token | jq ".token" | sed "s/\"//g")
echo
echo "company token is $manu_Token"
echo



echo
export TRX_ID=$(curl -s -X POST \
  http://localhost:4000/manu_creation \
  -H "authorization: Bearer $manu_Token" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.manu.example.com"],
	"fcn":"manu_creation",
	"args":["drug1","qrid1","box1","consignment1","krishna","12-02-2001"]
}')
echo "Transacton ID is $TRX_ID"
echo
echo

echo "GET query chaincode on  of "
echo
curl -s -X GET \
  http://localhost:4000/consignment_detail?peer=peer0.retailer.example.com&fcn=consignment_detail&args=%5B%22krishna%22%2C%22consignment1%22%5D\
  -H "authorization: Bearer $manu_Token" \
  -H "content-type: application/json"
echo


echo "GET query Transaction by TransactionID"
echo
curl -s -X GET http://localhost:4000/channels/mychannel/transactions/$TRX_ID?peer=peer0.org1.example.com \
  -H "authorization: Bearer $manu_Token" \
  -H "content-type: application/json"
echo
echo
