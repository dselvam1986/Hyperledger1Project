// Created by Dinesh  & Milan
package contracts

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SimplePaymentContract contract for handling writing and reading from the world state
type InsuranceContract struct {
	contractapi.Contract
}

// Account : User account
type Account struct {
	AccountID         string `json:"accountID"`
	Name              string `json:"name"`
	Address         bool   `json:"address"`
}

// Policy : Hold policy data
type Policy struct {
	PolicyID string `json:"policyID"`
}

