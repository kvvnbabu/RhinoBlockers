#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
ORG1_TOKEN 

function dkcl(){
        CONTAINER_IDS=$(docker ps -aq)
	echo
        if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" = " " ]; then
                echo "========== No containers available for deletion =========="
        else
                docker rm -f $CONTAINER_IDS
        fi
	echo
}

function dkrm(){
        DOCKER_IMAGE_IDS=$(docker images | grep "dev\|none\|test-vp\|peer[0-9]-" | awk '{print $3}')
	echo
        if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" = " " ]; then
		echo "========== No images available for deletion ==========="
        else
                docker rmi -f $DOCKER_IMAGE_IDS
        fi
	echo
}

function restartNetwork() {
	echo

        #teardown the network and clean the containers and intermediate images
	cd artifacts
	docker-compose down
	dkcl
	dkrm

	#Cleanup the material
	rm -rf /tmp/hfc-test-kvs_peerOrg* $HOME/.hfc-key-store/ /tmp/fabric-client-kvs_peerOrg*

	#Start the network
	docker-compose up -d
	cd -
	echo
}

function installNodeModules() {
	echo
	if [ -d node_modules ]; then
		echo "============== node modules installed already ============="
	else
		echo "============== Installing node modules ============="
		npm install
	fi
	echo
}

function installChainCodes(){
	echo
	initAuthToken
	installCreateDistributionChainCode
	echo
}

function installCreateDistributionChainCode(){
	echo "Install chaincode on CreateDistribution"
	echo
	curl -s -X POST \
	http://localhost:4000/chaincodes \
	-H "authorization: Bearer $ORG1_TOKEN" \
	-H "content-type: application/json" \
	-d '{
		"peers": ["peer1", "peer2"],
		"chaincodeName":"createdistribution",
		"chaincodePath":"github.com/distributionsmartcontract",
		"chaincodeVersion":"v0"
	}'
	echo
	echo
}

function initAuthToken(){
	echo "POST request Enroll on Org1  ..."
	echo
	ORG1_TOKEN=$(curl -s -X POST \
	http://localhost:4000/users \
	-H "content-type: application/x-www-form-urlencoded" \
	-d 'username=Jim&orgName=org1')
	echo $ORG1_TOKEN
	ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".token" | sed "s/\"//g")
	echo
	echo "ORG1 token is $ORG1_TOKEN"
	echo
	echo "POST request Enroll on Org2 ..."
	echo
}

restartNetwork

installNodeModules

PORT=4000 node app

installChainCodes