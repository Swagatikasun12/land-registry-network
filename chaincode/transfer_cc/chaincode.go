package main

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	sc "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

// Chaincode is the definition of the chaincode structure.
type Chaincode struct {
}

// Definition of status of TransferRequest
type statusHistory struct {
	Status        string `json:"Status"`
	StatusCreator string `json:"StatusCreator"`
	Date          int    `json:"Date"`
}

// Definition of the TransferRequest structure
type transferRequest struct {
	ID              string          `json:"ID"`
	LandID          string          `json:"LandID"`
	Lawyer          string          `json:"Lawyer"`
	RegistryOfficer string          `json:"RegistryOfficer"`
	BLRO            string          `json:"BLRO"`
	Stage           string          `json:"Stage"`
	StatusHistory   []statusHistory `json:"StatusHistory"`
	Complete        bool            `json:"Complete"`
}

var stage = [...]string{"lawyer", "registry", "blro"}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Invoke()", fcn, params)

	if fcn == "createTransferRequest" {
		return cc.createTransferRequest(stub, params)
	} else if fcn == "readTransferRequest" {
		return cc.readTransferRequest(stub, params)
	} else if fcn == "transfer2RegistryOfficer" {
		return cc.transfer2RegistryOfficer(stub, params)
	} else if fcn == "transfer2BLRO" {
		return cc.transfer2BLRO(stub, params)
	} else if fcn == "approveTransferRequest" {
		return cc.approveTransferRequest(stub, params)
	} else {
		fmt.Println("Invoke() did not find func: " + fcn)
		return shim.Error("Received unknown function invocation!")
	}
}

// Function to create new transferRequest (C of CRUD)
func (cc *Chaincode) createTransferRequest(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	creatorOrg, creatorCertIssuer, creator, err := getTxCreatorInfo(stub)
	if !authenticateCitizen(creatorOrg, creatorCertIssuer) {
		return shim.Error("{\"Error\":\"Access Denied!\",\"Payload\":{\"MSP\":\"" + creatorOrg + "\",\"CA\":\"" + creatorCertIssuer + "\"}}")
	}

	// Check if sufficient Params passed
	if len(params) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Check if Params are non-empty
	for a := 0; a < 4; a++ {
		if len(params[a]) <= 0 {
			return shim.Error("Arguments must be a non-empty string")
		}
	}

	key := "transferRequest-" + params[0]
	ID := params[0]
	LandID := params[1]
	Lawyer := params[2]
	Date := params[3]
	var StatusHistory []statusHistory
	Complete := false
	DateI, err := strconv.Atoi(Date)
	if err != nil {
		return shim.Error("Error: Invalid Date!")
	}

	// Check if TransferRequest exists with Key => key
	transferRequestAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to check if TransferRequest exists!")
	} else if transferRequestAsBytes != nil {
		return shim.Error("TransferRequest Already Exists!")
	}

	// Generate StatusHistory
	status := statusHistory{"Transfer Request Created.", creator, DateI}
	StatusHistory = append(StatusHistory, status)

	// Generate TransferRequest from params provided
	transferRequest := &transferRequest{ID, LandID, Lawyer, "", "", stage[0], StatusHistory, Complete}
	transferRequestJSONasBytes, err := json.Marshal(transferRequest)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Put State of newly generated TransferRequest with Key => key
	err = stub.PutState(key, transferRequestJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Add TransferRequestID to The Lawyer Profile with LawyerID
	args := util.ToChaincodeArgs("addCase", Lawyer, ID)
	response := stub.InvokeChaincode("lawyer_cc", args, "mainchannel")
	if response.Status != shim.OK {
		return shim.Error(response.Message)
	}

	// Returned on successful execution of the function
	return shim.Success(nil)
}

// Function to read an transferRequest (R of CRUD)
func (cc *Chaincode) readTransferRequest(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	// Check if sufficient Params passed
	if len(params) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Check if Params are non-empty
	if len(params[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	key := "transferRequest-" + params[0]

	// Get State of TransferRequest with Key => key
	transferRequestAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + params[0] + "\"}"
		return shim.Error(jsonResp)
	} else if transferRequestAsBytes == nil {
		jsonResp := "{\"Error\":\"TransferRequest does not exist!\"}"
		return shim.Error(jsonResp)
	}

	// Returned on successful execution of the function
	return shim.Success(transferRequestAsBytes)
}

// Function to add new active case (U of CRUD)
func (cc *Chaincode) transfer2RegistryOfficer(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	creatorOrg, creatorCertIssuer, creator, err := getTxCreatorInfo(stub)
	if !authenticateLawyer(creatorOrg, creatorCertIssuer) {
		return shim.Error("{\"Error\":\"Access Denied!\",\"Payload\":{\"MSP\":\"" + creatorOrg + "\",\"CA\":\"" + creatorCertIssuer + "\"}}")
	}

	// Check if sufficient Params passed
	if len(params) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Check if Params are non-empty
	for a := 0; a < 3; a++ {
		if len(params[a]) <= 0 {
			return shim.Error("Arguments must be a non-empty string")
		}
	}

	key := "transferRequest-" + params[0]
	ID := params[0]
	RegistryOfficer := params[1]
	Date := params[2]
	DateI, err := strconv.Atoi(Date)
	if err != nil {
		return shim.Error("Error: Invalid Date!")
	}

	// Get State of TransferRequest with Key => key
	transferRequestAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + params[0] + "\"}"
		return shim.Error(jsonResp)
	} else if transferRequestAsBytes == nil {
		jsonResp := "{\"Error\":\"TransferRequest does not exist!\"}"
		return shim.Error(jsonResp)
	}

	// Create new TransferRequest Variable
	transferRequestToUpdate := transferRequest{}
	err = json.Unmarshal(transferRequestAsBytes, &transferRequestToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

	// Generate StatusHistory
	status := statusHistory{"Request forwarded to Registry Officer.", creator, DateI}
	transferRequestToUpdate.StatusHistory = append(transferRequestToUpdate.StatusHistory, status)

	// Update transferRequest.RegistryOfficer => params[1]
	transferRequestToUpdate.RegistryOfficer = RegistryOfficer

	// Update Stage
	transferRequestToUpdate.Stage = stage[1]

	// Convert to Byte[]
	transferRequestJSONasBytes, err := json.Marshal(transferRequestToUpdate)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Put updated State of the TransferRequest with Key => key
	err = stub.PutState(key, transferRequestJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Add TransferRequestID to The RegistryOfficer Profile with RegistryOfficerID
	args := util.ToChaincodeArgs("addCase", RegistryOfficer, ID)
	response := stub.InvokeChaincode("registryoffice_cc", args, "mainchannel")
	if response.Status != shim.OK {
		return shim.Error(response.Message)
	}

	// Returned on successful execution of the function
	return shim.Success(nil)
}

// Function to complete a case, add to StatusHistory, remove from Complete (U of CRUD)
func (cc *Chaincode) transfer2BLRO(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	creatorOrg, creatorCertIssuer, creator, err := getTxCreatorInfo(stub)
	if !authenticateRegistryOffice(creatorOrg, creatorCertIssuer) {
		return shim.Error("{\"Error\":\"Access Denied!\",\"Payload\":{\"MSP\":\"" + creatorOrg + "\",\"CA\":\"" + creatorCertIssuer + "\"}}")
	}

	// Check if sufficient Params passed
	if len(params) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// Check if Params are non-empty
	for a := 0; a < 3; a++ {
		if len(params[a]) <= 0 {
			return shim.Error("Arguments must be a non-empty string")
		}
	}

	key := "transferRequest-" + params[0]
	ID := params[0]
	BLRO := params[1]
	Date := params[2]
	DateI, err := strconv.Atoi(Date)
	if err != nil {
		return shim.Error("Error: Invalid Date!")
	}

	// Get State of TransferRequest with Key => key
	transferRequestAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + params[0] + "\"}"
		return shim.Error(jsonResp)
	} else if transferRequestAsBytes == nil {
		jsonResp := "{\"Error\":\"TransferRequest does not exist!\"}"
		return shim.Error(jsonResp)
	}

	// Create new TransferRequest Variable
	transferRequestToUpdate := transferRequest{}
	err = json.Unmarshal(transferRequestAsBytes, &transferRequestToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

	// Generate StatusHistory
	status := statusHistory{"Request forwarded to BLRO.", creator, DateI}
	transferRequestToUpdate.StatusHistory = append(transferRequestToUpdate.StatusHistory, status)

	// Update transferRequest.RegistryOfficer => params[1]
	transferRequestToUpdate.BLRO = BLRO

	// Update Stage
	transferRequestToUpdate.Stage = stage[2]

	// Convert to Byte[]
	transferRequestJSONasBytes, err := json.Marshal(transferRequestToUpdate)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Put updated State of the TransferRequest with Key => key
	err = stub.PutState(key, transferRequestJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Add TransferRequestID to The RegistryOfficer Profile with RegistryOfficerID
	args := util.ToChaincodeArgs("addCase", BLRO, ID)
	response := stub.InvokeChaincode("blro_cc", args, "mainchannel")
	if response.Status != shim.OK {
		return shim.Error(response.Message)
	}

	// Returned on successful execution of the function
	return shim.Success(nil)
}

// Function to complete a case, add to StatusHistory, remove from Complete (U of CRUD)
func (cc *Chaincode) approveTransferRequest(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	creatorOrg, creatorCertIssuer, creator, err := getTxCreatorInfo(stub)
	if !authenticateBLRO(creatorOrg, creatorCertIssuer) {
		return shim.Error("{\"Error\":\"Access Denied!\",\"Payload\":{\"MSP\":\"" + creatorOrg + "\",\"CA\":\"" + creatorCertIssuer + "\"}}")
	}

	// Check if sufficient Params passed
	if len(params) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Check if Params are non-empty
	for a := 0; a < 2; a++ {
		if len(params[a]) <= 0 {
			return shim.Error("Arguments must be a non-empty string")
		}
	}

	key := "transferRequest-" + params[0]
	ID := params[0]
	Date := params[1]
	DateI, err := strconv.Atoi(Date)
	if err != nil {
		return shim.Error("Error: Invalid Date!")
	}

	// Get State of TransferRequest with Key => key
	transferRequestAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + params[0] + "\"}"
		return shim.Error(jsonResp)
	} else if transferRequestAsBytes == nil {
		jsonResp := "{\"Error\":\"TransferRequest does not exist!\"}"
		return shim.Error(jsonResp)
	}

	// Create new TransferRequest Variable
	transferRequestToUpdate := transferRequest{}
	err = json.Unmarshal(transferRequestAsBytes, &transferRequestToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

	// Generate StatusHistory
	status := statusHistory{"Transfer Request Approved by BLRO.", creator, DateI}
	transferRequestToUpdate.StatusHistory = append(transferRequestToUpdate.StatusHistory, status)

	// Update transferRequest.Complete => true
	transferRequestToUpdate.Complete = true

	// Convert to Byte[]
	transferRequestJSONasBytes, err := json.Marshal(transferRequestToUpdate)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Put updated State of the TransferRequest with Key => key
	err = stub.PutState(key, transferRequestJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Set complete to case with ID
	args0 := util.ToChaincodeArgs("completeCase", transferRequestToUpdate.BLRO, ID)
	response0 := stub.InvokeChaincode("blro_cc", args0, "mainchannel")
	if response0.Status != shim.OK {
		return shim.Error(response0.Message)
	}

	args1 := util.ToChaincodeArgs("completeCase", transferRequestToUpdate.RegistryOfficer, ID)
	response1 := stub.InvokeChaincode("registryoffice_cc", args1, "mainchannel")
	if response1.Status != shim.OK {
		return shim.Error(response1.Message)
	}

	args2 := util.ToChaincodeArgs("completeCase", transferRequestToUpdate.Lawyer, ID)
	response2 := stub.InvokeChaincode("lawyer_cc", args2, "mainchannel")
	if response2.Status != shim.OK {
		return shim.Error(response2.Message)
	}

	// Returned on successful execution of the function
	return shim.Success(nil)
}

// ---------------------------------------------
// Helper Functions
// ---------------------------------------------

// Authentication
// ++++++++++++++

// Get Tx Creator Info
func getTxCreatorInfo(stub shim.ChaincodeStubInterface) (string, string, string, error) {
	var mspid string
	var err error
	var cert *x509.Certificate
	mspid, err = cid.GetMSPID(stub)

	if err != nil {
		fmt.Printf("Error getting MSP identity: %sn", err.Error())
		return "", "", "", err
	}

	cert, err = cid.GetX509Certificate(stub)
	if err != nil {
		fmt.Printf("Error getting client certificate: %sn", err.Error())
		return "", "", "", err
	}

	return mspid, cert.Issuer.CommonName, cert.Subject.CommonName, nil
}

// Authenticate => TransferRequest
func authenticateCitizen(mspID string, certCN string) bool {
	return (mspID == "CitizenMSP") && (certCN == "ca.citizen.lran.com")
}

// Authenticate => TransferRequest
func authenticateLawyer(mspID string, certCN string) bool {
	return (mspID == "LawyerMSP") && (certCN == "ca.lawyer.lran.com")
}

// Authenticate => TransferRequest
func authenticateRegistryOffice(mspID string, certCN string) bool {
	return (mspID == "RegistryOfficeMSP") && (certCN == "ca.registryoffice.lran.com")
}

// Authenticate => TransferRequest
func authenticateBLRO(mspID string, certCN string) bool {
	return (mspID == "BLROMSP") && (certCN == "ca.blro.lran.com")
}
