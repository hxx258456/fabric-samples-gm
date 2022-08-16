/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"gitee.com/zhaochuninhefei/fabric-contract-api-go-gm/contractapi"
	"gitee.com/zhaochuninhefei/fabric-samples-gm/asset-transfer-basic/chaincode-go/chaincode"
)

func main() {
	assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}
