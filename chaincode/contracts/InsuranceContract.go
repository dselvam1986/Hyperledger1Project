// Created by Dinesh & Milan
package contracts

import (
	//"encoding/json"
	"encoding/json"
	"fmt"

	//"log"
	//"time"

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
	LatestTransaction string `json:"transaction"`
	Users             []User
}

// Account : User account
type User struct {
	UserID   string `json:"userID"` //is this memberID? How does it defer from UserName?
	UserName string `json:"username"`
	OwnerRel string `json:"rel"` //SELF, SPOUSE, DEPENDANT
}

// Policy : Hold policy data
type Plans struct {
	PlanID      string `json:"policyID"`
	PlanName    string `json:"policyname"`
	PlanOptions []Policy
}

type Policy struct {
	PolicyID       string `json:"planID"`
	PolicyName     string `json:"planname"`
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
func (spc *InsuranceContract) RegisterUserAccount(ctx contractapi.TransactionContextInterface, name string, provider string) (*Account, *User, error) {
	id, _ := ctx.GetClientIdentity().GetID()
	//check if there is any error returning the worldstate of user certificate ID
	accountBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	//check if ID already exists (return the state of the ID by checking the world state)
	if accountBytes != nil {
		return nil, nil, fmt.Errorf("the account already exists for user %s", name)
	}
	//defince structs
	account := Account{
		DocType:           "Account",
		AccountID:         id,
		OwnerName:         name,
		LatestTransaction: ctx.GetStub().GetTxID(),
	}

	//convert Golang to jSon format (JSON Byte Array)
	accountBytes, err = json.Marshal(account)
	if err != nil {
		return nil, nil, err
	}

	//put account data unto the Ledger (key value pair)
	err = ctx.GetStub().PutState(id, accountBytes)
	if err != nil {
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, err
	}

	//declare user variable to save registered user
	var user User
	/*
		user := User{

		}*/
	//register the owner as a user
	//user = RegisterUser(ctx, name, "SELF")

	/* I DISABLED THIS CODE BLOCK BECAUSE THERE IS NO TRANSACTION STRUCT ABOVE! ALSO WHAT IS THE POINT OF THIS? TO RECORD ALL TRANSACTIONS?
	IF THIS WAS FOR MONEY TRANSACTION IT WOULD BE USEFULL, BUT FOR USER CREATION I DONT THINK THATS THIS IS NEEDED. LETS DISCUSS TOMROROW.*/
	// transaction := Transaction{
	// 	DocType:       "Transaction",
	// 	TransactionID: ctx.GetStub().GetTxID(),
	// 	Beneficiary:   id,
	// 	Remitter:      provider,
	// 	Amount:        0,
	// }

	// var transactionBytes []byte
	// transactionBytes, err = json.Marshal(transaction)
	// if err != nil {
	// 	return nil, err
	// }entIdentity().GetID()
	//userBytes, err := ctx.GetStub().GetState(id)

	// //write info to the ledger
	// err = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), transactionBytes)
	// if err != nil {
	// 	return nil, err
	// }

	return &account, &user, nil
}

func (spc *InsuranceContract) RegisterUser(ctx contractapi.TransactionContextInterface, name string, relation string) (*User, error) {
	// checks to see if user already exists
	id, _ := ctx.GetClientIdentity().GetID()
	userBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if userBytes != nil {
		return nil, fmt.Errorf("the account already exists for user %s", name)
	}

	user := User{
		UserID:   id,
		UserName: name,
		OwnerRel: relation,
	}

	userBytes, err = json.Marshal(user)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState(id, userBytes)
	if err != nil {
		return nil, err

	}
	return &user, nil
}

//Getter Functions
func (spc *InsuranceContract) GetUser(ctx contractapi.TransactionContextInterface, id string) (*User, error) {
	userbytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userbytes == nil {
		return nil, fmt.Errorf("user not found")
	}

	User := User{}
	json.Unmarshal(userbytes, &User)

	return &User, err
}

// FetchID : to check the owner's AccountID
func (spc *InsuranceContract) FetchID(ctx contractapi.TransactionContextInterface) (string, error) {

	id, _ := ctx.GetClientIdentity().GetID()
	accountBytes, err := ctx.GetStub().GetState(id)
	//check if there is any error returning the worldstate of user certificate ID
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}

	var account Account
	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		return "", err
	}

	return account.AccountID, nil
}

// DeleteUserAccount deletes an given asset from the world state.
func (spc *InsuranceContract) DeleteUserAccount(ctx contractapi.TransactionContextInterface, id string) error {

	_, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	return ctx.GetStub().DelState(id)
}
