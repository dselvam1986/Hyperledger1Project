package main

import (
	"insurance-application-chaincode/contracts"

	// "simple-payment-application-chaincode/contracts"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	insuranceContract := new(contracts.InsuranceContract)

	cc, err := contractapi.NewChaincode(insuranceContract)

	if err != nil {
		panic(err.Error())
	}

	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}
