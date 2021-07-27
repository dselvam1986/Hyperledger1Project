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
	PlanId 			  string `json:"planid"`
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
	PolicyName     string `json:"planname"` /* Medical, Vision, Dental */
	Deductible     int    `json:"deductible"`
	OOPLimitPerson int    `json:"ooplimitperson"`
	OOPLimitfamily int    `json:"ooplimitfamily"`
	FSA			   bool   `json:"fsa"`	
	FSABalance     int    `json:"fsabalance"` 
}

// Init and Creator Functions for User, Organization, Policy and Plan
func (spc *InsuranceContract) InitInsurance(ctx contractapi.TransactionContextInterface) error {

	// possible function to pre create policy and then create plans. then add the plans to the policy array.
	return nil
}

// RegisterUserAccount : User registers his account
func (spc *InsuranceContract) RegisterAccount(ctx contractapi.TransactionContextInterface, name string, provider string) (*Account, *User, error) {
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

	//declare user variable to save registered user,  declare contract var to call func within func
	var user *User
	var contract InsuranceContract
	user, _ = contract.RegisterUser(ctx, name, "SELF", true, "")
	
	// save user data in account User array
	usrArry := [] User{}

	usrArry = append(usrArry, *user)

	//defince structs
	account := Account{
		DocType:           "Account",
		AccountID:         id,
		OwnerName:         name,
		LatestTransaction: ctx.GetStub().GetTxID(),
		Users:				usrArry,
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

	return &account, user, nil
}

func (spc *InsuranceContract) RegisterUser(ctx contractapi.TransactionContextInterface, name string, relation string, isSelf bool, accountID string) (*User, error) {
	
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

	if(!isSelf){
		accountBytes, err := ctx.GetStub().GetState(accountID)
		if err != nil {
			return nil, fmt.Errorf("failed to read from world state: %v", err)
		}
		//check if ID already exists (return the state of the ID by checking the world state)
		if accountBytes != nil {
			return nil, fmt.Errorf("the account already exists for user %s", name)
		}
		var account Account
		err = json.Unmarshal(accountBytes, &account)
		if err != nil {
			return nil, err
		}

		account.Users = append(account.Users, user)

		//convert Golang to jSon format (JSON Byte Array)
		accountBytes, err = json.Marshal(account)
		fmt.Print(accountBytes)

		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (spc *InsuranceContract) RegisterPlan(ctx contractapi.TransactionContextInterface, name string ) (*Plans, error) {
	
	// checks to see if user already exists
	id, _ := ctx.GetClientIdentity().GetID()
	planBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if planBytes != nil {
		return nil, fmt.Errorf("the account already exists for user %s", name)
	}

	plan := Plans{
		PlanID: id, 
		PlanName: name,
	}

	planBytes, err = json.Marshal(plan)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState(id, planBytes)
	if err != nil {
		return nil, err

	}

	return &plan, nil
}

func (spc *InsuranceContract) RegisterPolicy(ctx contractapi.TransactionContextInterface, name string, deductible int) (*Policy, error) {
	
	// checks to see if user already exists
	id, _ := ctx.GetClientIdentity().GetID()
	policyBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if policyBytes != nil {
		return nil, fmt.Errorf("the account already exists for user %s", name)
	}

	policy := Policy{
		PolicyID: id,
		PolicyName: name, 
		Deductible: deductible,
		// OOPLimitPerson: personLimit,
		// OOPLimitfamily: familylimit,
	}

	policyBytes, err = json.Marshal(policy)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState(id, policyBytes)
	if err != nil {
		return nil, err

	}

	return &policy, nil
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
func (spc *InsuranceContract) DeleteAccount(ctx contractapi.TransactionContextInterface, id string) error {

	_, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	return ctx.GetStub().DelState(id)
}


//Link Functionality
// PolicyPlan : Link a policy to an exisiting plan
func (spc *InsuranceContract) PolicyPlan(ctx contractapi.TransactionContextInterface, policyID string, planID string) (string, error) {
	// get plan info
	var plan Plans
	planBytes, err := ctx.GetStub().GetState(planID)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	//check if ID already exists (return the state of the ID by checking the world state)
	if planBytes != nil {
		return "", fmt.Errorf("confirmed the plan already exists for planID %s", planID)
	}

	err = json.Unmarshal(planBytes, &plan)
	if err != nil {
		return "", err
	}
	// get policy info
	var policy Policy
	policybytes, err := ctx.GetStub().GetState(policyID)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	//check if ID already exists (return the state of the ID by checking the world state)
	if policybytes != nil {
		return "", fmt.Errorf("confirmed the plan already exists for planID %s", planID)
	}

	json.Unmarshal(policybytes, &policy)

	plan.PlanOptions= append(plan.PlanOptions, policy)   
	return "Policy addded to Plan",nil
}

// RegisterPolicy : User subscribes to a policy
func (spc *InsuranceContract) RegisterPlanToAccount(ctx contractapi.TransactionContextInterface, ptype string, deduct int, isFSA bool, single bool, family bool) (*Policy, error) {

	id, _ := ctx.GetClientIdentity().GetID()
	//check if there is any error returning the worldstate of user certificate ID
	policyBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	//check if ID already exists (return the state of the ID by checking the world state)
	if policyBytes != nil {
		return nil, fmt.Errorf("the policy already exists for policyID %s", id)
	}

	// unmarshal the policy to update value
	var policy Policy
	json.Unmarshal(policyBytes, &policy)

	var OOP int
	var bal int

	if (single){
		OOP = 2600
	}else {
		OOP = 5000
	}
	if(isFSA){
		bal =1000
	}

	// Update values 
	policy.OOPLimitPerson = OOP
	policy.OOPLimitfamily = OOP
	policy.FSA = isFSA
	policy.FSABalance = bal


	// //define structs
	// policy := Policy{		
	// 	PolicyID:        id,
	//     PolicyName:      ptype,
	//     Deductible:      deduct,
	// 	OOPLimitPerson:  OOP,
	// 	OOPLimitfamily:  OOP,
	// 	FSA:             isFSA,
	// 	FSABalance:       bal,
	// }
	//convert Golang to jSon format (JSON Byte Array)
	policyBytes, err = json.Marshal(policy)
	if err != nil {
		return nil, err
	}
	//put policy data unto the Ledger (key value pair)
	err = ctx.GetStub().PutState(id, policyBytes)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}


