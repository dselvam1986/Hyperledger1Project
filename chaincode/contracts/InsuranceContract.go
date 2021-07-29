// Created by Dinesh & Milan
package contracts

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

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
	Users             []User `json:"user"`
	PlanId            string `json:"planid"`
}

// Account : User account
// DOCTYPE: USER
type User struct {
	DocType  string `json:"docType"`
	UserID   string `json:"userID"` //is this memberID? How does it defer from UserName?
	UserName string `json:"username"`
	OwnerRel string `json:"rel"` //SELF, SPOUSE, DEPENDANT
}

// Policy : Hold policy data
//DOCTYPE: PLANS
type Plans struct {
	DocType     string   `json:"docType"`
	PlanID      string   `json:"planID"`
	PlanName    string   `json:"planname"`
	FamilyPlan  bool     `json:"familyplan"`
	PlanOptions []string `json:"planoptions"`
}

//DOCTYPEL POLICY
type Policy struct {
	DocType        string `json:"docType"`
	PolicyID       string `json:"policyID"`
	PolicyName     string `json:"policyname"` /* Medical, Vision, Dental */
	Deductible     int    `json:"deductible"`
	IsFamily       bool   `json:"isfamily"`
	OOPLimitSingle int    `json:"ooplimitsingle"`
	OOPLimitFamily int    `json:"ooplimitfamily"`
	FSA            bool   `json:"fsa"`
	FSABalance     int    `json:"fsabalance"`
}

var accountCount int
var userCount int
var planCount int
var policyCount int
var defaultFSABalance int

var contract InsuranceContract

// Init and Creator Functions for User, Organization, Policy and Plan
func (spc *InsuranceContract) InitInsurance(ctx contractapi.TransactionContextInterface) error {

	defaultFSABalance = 1000

	// possible function to pre create policy and then create plans. then add the plans to the policy array.
	accountCount = 0
	userCount = 0
	planCount = 0
	policyCount = 0

	return nil
}

// RegisterUserAccount : User registers his account - WORKS FINE
func (spc *InsuranceContract) RegisterAccount(ctx contractapi.TransactionContextInterface, name string) (*Account, error) {

	id, _ := contract.IDGenerator("account", name, accountCount)
	//check if there is any error returning the worldstate of user certificate ID
	accountBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	//check if ID already exists (return the state of the ID by checking the world state)
	if accountBytes != nil {
		return nil, fmt.Errorf("the account already exists for user %s", name)
	}

	//declare user variable to save registered user,  declare contract var to call func within func
	var user *User
	user, _ = contract.RegisterUser(ctx, name, "SELF", true, "")

	// save user data in account User array
	usrArry := []User{}

	usrArry = append(usrArry, *user)

	//defince structs
	account := Account{
		DocType:           "Account",
		AccountID:         id,
		OwnerName:         name,
		LatestTransaction: ctx.GetStub().GetTxID(),
		Users:             usrArry,
		PlanId:            "",
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

	accountCount += 1

	return &account, nil
}

// Create a new User ( ORG pov )
func (spc *InsuranceContract) RegisterUser(ctx contractapi.TransactionContextInterface, name string, relation string, isSelf bool, accountID string) (*User, error) {

	// checks to see if user already exists
	id, _ := contract.IDGenerator("user", name, userCount)
	userBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if userBytes != nil {
		return nil, fmt.Errorf("the account already exists for user %s", name)
	}
	log.Println("creating user")

	user := User{
		DocType:  "User",
		UserID:   id,
		UserName: name,
		OwnerRel: relation,
	}

	userBytes, err = json.Marshal(user)
	if err != nil {
		return nil, err
	}

	log.Println("marshal user")

	err = ctx.GetStub().PutState(id, userBytes)
	if err != nil {
		return nil, err

	}
	log.Println("Put State user")

	if !isSelf {
		log.Println("Getting account from ID - getState")
		accountBytes, err := ctx.GetStub().GetState(accountID)
		if err != nil {
			return nil, fmt.Errorf("failed to read from world state: %v", err)
		}
		//check if ID already exists (return the state of the ID by checking the world state)
		if accountBytes == nil {
			return nil, fmt.Errorf("the account not found for user %s", name)
		}

		var account Account
		err = json.Unmarshal(accountBytes, &account)
		if err != nil {
			return nil, fmt.Errorf("line 153: %v", err)
		}

		account.Users = append(account.Users, user)

		//convert Golang to jSon format (JSON Byte Array)
		accountBytes, err = json.Marshal(account)
		if err != nil {
			return nil, err
		}

		err = ctx.GetStub().PutState(accountID, accountBytes)
		if err != nil {
			return nil, err

		}
	}

	userCount += 1

	return &user, nil
}

// Create a new Plan ( ORG pov )
func (spc *InsuranceContract) RegisterPlan(ctx contractapi.TransactionContextInterface, name string, isFamilyPlan bool) (*Plans, error) {

	// checks to see if user already exists
	id, _ := contract.IDGenerator("plan", name, planCount)
	planBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if planBytes != nil {
		return nil, fmt.Errorf("the account already exists for user %s", name)
	}
	// var planArray []string
	plan := Plans{
		DocType:     "Plan",
		PlanID:      id,
		PlanName:    name,
		FamilyPlan:  isFamilyPlan,
		PlanOptions: []string{},
	}

	planBytes, err = json.Marshal(plan)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState(id, planBytes)
	if err != nil {
		return nil, err

	}

	planCount += 1

	return &plan, nil
}

// Create a new Policy ( ORG pov )
func (spc *InsuranceContract) RegisterPolicy(ctx contractapi.TransactionContextInterface, name string, deductible int, isFamilyPolicy bool, OOPlimit int, FSAbal int, isFSA bool) (*Policy, error) {

	// checks to see if user already exists
	id, _ := contract.IDGenerator("policy", name, policyCount)
	policyBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if policyBytes != nil {
		return nil, fmt.Errorf("the account already exists for user %s", name)
	}

	if isFSA && FSAbal == 0 {
		FSAbal = defaultFSABalance
	}

	if !isFSA {
		FSAbal = 0
	}

	var OOPFamily int
	var OOPSingle int
	if isFamilyPolicy {
		OOPFamily = OOPlimit
		OOPSingle = 0
	} else {
		OOPFamily = 0
		OOPSingle = OOPlimit
	}

	policy := Policy{
		DocType:        "Policy",
		PolicyID:       id,
		PolicyName:     name,
		Deductible:     deductible,
		IsFamily:       isFamilyPolicy,
		OOPLimitSingle: OOPSingle,
		OOPLimitFamily: OOPFamily,
		FSA:            isFSA,
		FSABalance:     FSAbal,
	}

	policyBytes, err = json.Marshal(policy)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState(id, policyBytes)
	if err != nil {
		return nil, err

	}
	policyCount += 1

	return &policy, nil
}

//Getter Functions
// Get User from ID - WORKS FINE
func (spc *InsuranceContract) GetUser(ctx contractapi.TransactionContextInterface, id string) (*User, error) {
	userbytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userbytes == nil {
		return nil, fmt.Errorf("user not found")
	}

	var user User
	err = json.Unmarshal(userbytes, &user)
	if err != nil {
		return nil, err
	}

	return &user, err
}

// FetchID : to check the owner's AccountID - WORKS FINE
func (spc *InsuranceContract) FetchID(ctx contractapi.TransactionContextInterface, accountID string) (*Account, error) {

	accountBytes, err := ctx.GetStub().GetState(accountID)
	//check if there is any error returning the worldstate of user certificate ID
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	var account Account
	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
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
func (spc *InsuranceContract) LinkPolicyToPlan(ctx contractapi.TransactionContextInterface, policyID string, planID string) (string, error) {
	// get plan info
	var plan Plans
	planBytes, err := ctx.GetStub().GetState(planID)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	//check if ID already exists (return the state of the ID by checking the world state)
	if planBytes == nil {
		return "", fmt.Errorf("plan does not exists for planID %s", planID)
	}

	err = json.Unmarshal(planBytes, &plan)
	if err != nil {
		return "", fmt.Errorf("unmarshall planbytes error: %v", err)
	}
	// // get policy info
	// var policy Policy
	// policybytes, err := ctx.GetStub().GetState(policyID)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to read from world state: %v", err)
	// }
	// //check if ID already exists (return the state of the ID by checking the world state)
	// if policybytes != nil {
	// 	return "", fmt.Errorf("confirmed the plan already exists for planID %s", planID)
	// }

	// json.Unmarshal(policybytes, &policy)

	plan.PlanOptions = append(plan.PlanOptions, policyID)
	planBytes, err = json.Marshal(plan)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(planID, planBytes)
	if err != nil {
		return "", err

	}

	return "Policy addded to Plan", nil
}

//Show all available plan to User by Account ID
func (spc *InsuranceContract) ShowAvailablePlans(ctx contractapi.TransactionContextInterface) ([]*Plans, error) {

	queryString := `{"selector":{"docType":"Plan"}}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var plansArray []*Plans

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var plans Plans
		err = json.Unmarshal(queryResult.Value, &plans)
		if err != nil {
			return nil, err
		}
		plansArray = append(plansArray, &plans)
	}

	return plansArray, nil

}

//Customer Support / Provider perspective show all info regarding account - plan and policy : Param - accountID

// RegisterPolicy : User subscribes to a policy
func (spc *InsuranceContract) LinkPlanToAccount(ctx contractapi.TransactionContextInterface, accountID string, planID string) (string, error) {

	//check if there is any error returning the worldstate of user certificate ID
	accountBytes, err := ctx.GetStub().GetState(accountID)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	//check if ID already exists (return the state of the ID by checking the world state)
	if accountBytes == nil {
		return "", fmt.Errorf("the account is not found: %v", accountID)
	}

	// unmarshal the policy to update value
	var account Account
	json.Unmarshal(accountBytes, &account)

	account.PlanId = planID

	//convert Golang to jSon format (JSON Byte Array)
	accountBytes, err = json.Marshal(account)
	if err != nil {
		return "", err
	}
	//put policy data unto the Ledger (key value pair)
	err = ctx.GetStub().PutState(accountID, accountBytes)
	if err != nil {
		return "nil", err
	}
	return "Account linked to Plan", nil
}

// Get All available plans for User: Should return account is and Array of Plans (full details) - User will choose plan and pass accountId and PlanID to Function LInkPlanToAccount

// Helper function
func (spc *InsuranceContract) IDGenerator(doctype string, name string, count int) (string, error) {

	docSubstring := doctype[0:3]
	nameSubString := name[0:3]

	s := []string{docSubstring, nameSubString, strconv.Itoa(count)}

	return strings.Join(s, ""), nil
}
