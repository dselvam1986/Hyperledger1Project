// Created by Dinesh & Milan
package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// InsuranceContract contract for handling writing and reading from the world state
type InsuranceContract struct {
	contractapi.Contract
}

// Account : The asset being tracked on the chain
type Account struct {
	DocType           string `json:"docType"`
	AccountID         string `json:"accountID"`
	OwnerName         string `json:"name"`
	ProviderName      string `json:"provider"`
	LatestTransaction string `json:"transaction"`
}

// Account : User account
type User struct {
	UserID      string `json:"userID"` //is this memberID? How does it defer from UserName?
	UserName    string `json:"username"`
	UserAddress string `json:"useraddress"`
	OwnerRel    string `json:"rel"` //SELF, SPOUSE, DEPENDANT
	ProviderID  string `json:"providerID"`
}

// Provider: Account
type Provider struct {
	ProviderID      string `json:"providerID"`
	ProviderName    string `json:"providername"`
	ProviderAddress string `json:"providerAddr"` //store the state in which the provider is operating
}

// Policy : Hold policy data
type Policy struct {
	PolicyID   string `json:"policyID"`
	ProviderID string `json:"providerID"`
	PolicyName string `json:"policyname"`
	PolicyPlan []Plans
}

type Plans struct {
	PlanName       string `json:"planname"`
	Deductible     int    `json:"deductible"`
	OOPLimitPerson int    `json:"ooplimitperson"`
	OOPLimitfamily int    `json:"ooplimitfamily"`
}

// Init and Creator Functions for User, Organization, Policy and Plan
func (spc *InsuranceContract) InitInsurance(ctx contractapi.TransactionContextInterface) error {

	// possible function to pre create policy and then create plans. then add the plans to the policy array.

	return nil
}

// RegisterUserAccount : User registers his account
func (spc *InsuranceContract) RegisterUserAccount(ctx contractapi.TransactionContextInterface, name string, provider string) (*Account, error) {

	id, _ := ctx.GetClientIdentity().GetID()
	//check if there is any error returning the worldstate of user certificate ID
	accountBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	//check if ID already exists (return the state of the ID by checking the world state)
	if accountBytes != nil {
		return nil, fmt.Errorf("the account already exists for user %s", name)
	}
	//defince structs
	account := Account{
		DocType:           "Account",
		AccountID:         id,
		OwnerName:         name,
		ProviderName:      provider,
		LatestTransaction: ctx.GetStub().GetTxID(),
	}
	//convert Golang to jSon format (JSON Byte Array)
	accountBytes, err = json.Marshal(account)
	if err != nil {
		return nil, err//defince structs
		account := Account{
			DocType:           "Account",
			AccountID:         id,
			OwnerName:         name,
			ProviderName:      provider,
			LatestTransaction: ctx.GetStub().GetTxID(),
		}
		//convert Golang to jSon format (JSON Byte Array)
		accountBytes, err = json.Marshal(account)
		if err != nil {
			return nil, err
		}
		//put account data unto the Ledger (key value pair)
		err = ctx.GetStub().PutState(id, accountBytes)
		if err != nil {
			return nil, err
		}
	if err != nil {
		return nil, err
	}

	transaction := Transaction{
		DocType:       "Transaction",
		TransactionID: ctx.GetStub().GetTxID(),
		Beneficiary:   id,
		Remitter:      provider,
		Amount:        0,
	}

	var transactionBytes []byte
	transactionBytes, err = json.Marshal(transaction)
	if err != nil {
		return nil, err
	}
	//write info to the ledger
	err = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), transactionBytes)
	if err != nil {
		return nil, err
	}

	return &account, nil
}