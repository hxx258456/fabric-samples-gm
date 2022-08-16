/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package functions

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gitee.com/zhaochuninhefei/fabric-sdk-go-gm/pkg/core/config"
	"gitee.com/zhaochuninhefei/fabric-sdk-go-gm/pkg/gateway"
)

// ManyUpdates allows you to push many cuncurrent updates to a variable
func ManyUpdates(function, variableName, change, sign string) ([]byte, error) {

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		return nil, fmt.Errorf("error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err := populateWallet(wallet)
		if err != nil {
			return nil, fmt.Errorf("failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return nil, fmt.Errorf("failed to get network: %v", err)
	}

	contract := network.GetContract("bigdatacc")

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() ([]byte, error) {
			defer wg.Done()
			result, err := contract.SubmitTransaction(function, variableName, change, sign)
			if err != nil {
				return result, fmt.Errorf("failed to evaluate transaction: %v", err)
			}
			return result, nil
		}()
	}

	wg.Wait()

	result, err := contract.EvaluateTransaction("get", variableName)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %v", err)
	}
	return result, err
}
