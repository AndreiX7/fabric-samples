#!/bin/bash -e
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# exit on first error

export BASE_FOLDER=$WORKSPACE/gopath/src/github.com/hyperledger
export NEXUS_URL=nexus3.hyperledger.org:10001
export ORG_NAME="hyperledger/fabric"

Parse_Arguments() {
      while [ $# -gt 0 ]; do
              case $1 in
                      --env_Info)
                            env_Info
                            ;;
                      --pull_Docker_Images)
                            pull_Docker_Images
                            ;;
                      --pull_Fabric_CA_Images)
                            pull_Fabric_CA_Images
                            ;;
                      --clean_Environment)
                            clean_Environment
                            ;;
                      --byfn_eyfn_Tests)
                            byfn_eyfn_Tests
                            ;;
                      --fabcar_Tests)
                            fabcar_Tests
                            ;;
                      --pull_Thirdparty_Images)
                            pull_Thirdparty_Images
                            ;;
              esac
              shift
      done
}

clean_Environment() {

echo "-----------> Clean Docker Containers & Images, unused/lefover build artifacts"
function clearContainers () {
        CONTAINER_IDS=$(docker ps -aq)
        if [ -z "$CONTAINER_IDS" ] || [ "$CONTAINER_IDS" = " " ]; then
                echo "---- No containers available for deletion ----"
        else
                docker rm -f $CONTAINER_IDS || true
                docker ps -a
        fi
}

function removeUnwantedImages() {

        for i in $(docker images | grep none | awk '{print $3}'); do
                docker rmi ${i};
        done

        for i in $(docker images | grep -vE ".*baseimage.*(0.4.13|0.4.14)" | grep -vE ".*baseos.*(0.4.13|0.4.14)" | grep -vE ".*couchdb.*(0.4.13|0.4.14)" | grep -vE ".*zoo.*(0.4.13|0.4.14)" | grep -vE ".*kafka.*(0.4.13|0.4.14)" | grep -v "REPOSITORY" | awk '{print $1":" $2}'); do
                docker rmi ${i};
        done
}

# Remove /tmp/fabric-shim
docker run -v /tmp:/tmp library/alpine rm -rf /tmp/fabric-shim || true

# remove tmp/hfc and hfc-key-store data
rm -rf /home/jenkins/.nvm /home/jenkins/npm /tmp/fabric-shim /tmp/hfc* /tmp/npm* /home/jenkins/kvsTemp /home/jenkins/.hfc-key-store

rm -rf /var/hyperledger/*

rm -rf gopath/src/github.com/hyperledger/fabric-ca/vendor/github.com/cloudflare/cfssl/vendor/github.com/cloudflare/cfssl_trust/ca-bundle || true
# yamllint disable-line rule:line-length
rm -rf gopath/src/github.com/hyperledger/fabric-ca/vendor/github.com/cloudflare/cfssl/vendor/github.com/cloudflare/cfssl_trust/intermediate_ca || true

clearContainers
removeUnwantedImages
}

env_Info() {
	# This function prints system info

	#### Build Env INFO
	echo "-----------> Build Env INFO"
	# Output all information about the Jenkins environment
	uname -a
	cat /etc/*-release
	env
	gcc --version
	docker version
	docker info
	docker-compose version
	pgrep -a docker
}

# Pull Thirdparty Docker images (Kafka, couchdb, zookeeper)
pull_Thirdparty_Images() {
            echo "------> BASE_IMAGE_TAG:" $BASE_IMAGE_TAG
            for IMAGES in kafka couchdb zookeeper; do
                 echo "-----------> Pull $IMAGES image"
                 echo
                 docker pull $ORG_NAME-$IMAGES:${BASE_IMAGE_TAG} > /dev/null 2>&1
                 if [ $? -ne 0 ]; then
                       echo -e "\033[31m FAILED to pull docker images" "\033[0m"
                       exit 1
                 fi
                 docker tag $ORG_NAME-$IMAGES:${BASE_IMAGE_TAG} $ORG_NAME-$IMAGES
            done
                 echo
                 docker images | grep hyperledger/fabric
}
# pull fabric images from nexus
pull_Docker_Images() {
            pull_Fabric_CA_Image
            for IMAGES in peer orderer tools ccenv; do
                 echo "-----------> pull $IMAGES image"
                 echo
                 docker pull $NEXUS_URL/$ORG_NAME-$IMAGES:$IMAGE_TAG > /dev/null 2>&1
                 if [ $? -ne 0 ]; then
                       echo -e "\033[31m FAILED to pull docker images" "\033[0m"
                       exit 1
                 fi
                 docker tag $NEXUS_URL/$ORG_NAME-$IMAGES:$IMAGE_TAG $ORG_NAME-$IMAGES
                 docker tag $NEXUS_URL/$ORG_NAME-$IMAGES:$IMAGE_TAG $ORG_NAME-$IMAGES:$ARCH-$VERSION
                 docker rmi -f $NEXUS_URL/$ORG_NAME-$IMAGES:$IMAGE_TAG
            done
                 echo
                 docker images | grep hyperledger/fabric
}
# pull fabric-ca images from nexus
pull_Fabric_CA_Image() {
            echo "------> IMAGE_TAG:" $IMAGE_TAG
            for IMAGES in ca ca-peer ca-orderer ca-tools; do
                 echo "-----------> pull $IMAGES image"
                 echo
                 docker pull $NEXUS_URL/$ORG_NAME-$IMAGES:$IMAGE_TAG > /dev/null 2>&1
                 if [ $? -ne 0 ]; then
                       echo -e "\033[31m FAILED to pull docker images" "\033[0m"
                       exit 1
                 fi
                 docker tag $NEXUS_URL/$ORG_NAME-$IMAGES:$IMAGE_TAG $ORG_NAME-$IMAGES
	         docker tag $NEXUS_URL/$ORG_NAME-$IMAGES:$IMAGE_TAG $ORG_NAME-$IMAGES:$ARCH-$VERSION
                 docker rmi -f $NEXUS_URL/$ORG_NAME-$IMAGES:$IMAGE_TAG
            done
                 echo
                 docker images | grep hyperledger/fabric-ca
}

# run byfn,eyfn tests
byfn_eyfn_Tests() {
                 echo
                 echo "-----------> Execute Byfn and Eyfn Tests"
                 ./byfn_eyfn.sh
}
# run fabcar tests
fabcar_Tests() {
                 echo
                 echo "npm version ------> $(npm -v)"
                 echo "node version ------> $(node -v)"
                 echo "-----------> Execute FabCar Tests"
                 ./fabcar.sh
}
Parse_Arguments $@
