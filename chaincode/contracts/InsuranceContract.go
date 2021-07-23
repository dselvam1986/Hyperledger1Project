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
type User struct {
	UserID string `json:"userID"`
	UserName      	string `json:"username"`
	UserAddress   	bool   `json:"useraddress"`
	ProviderID 		string `json:"providerID"`
}

// Provider: Account
type Provider struct {
	ProviderID 		string `json:"providerID"`
	ProviderName 	string `json:"providername"`

}

// Policy : Hold policy data
type Policy struct { 
	PolicyID 		string `json:"policyID"`
	ProviderID 		string `json:"providerID"`
	PolicyName 		string `json:"policyname"`
	PolicyPlan		[]		Plans

}

type Plans struct{
	PlanName		string `json:"planname"`
	Deductible		int    `json:"deductible"`
	OOPLimitPerson	int    `json:"ooplimitperson"`
	OOPLimitfamily	int    `json:"ooplimitfamily"`
}

// Init and Creator Functions for User, Organization, Policy and Plan
func (spc *InsuranceContract) InitInsurance(ctx contractapi.TransactionContextInterface) error {

	// possible function to pre create policy and then create plans. then add the plans to the policy array. 
	

	return nil;
}

