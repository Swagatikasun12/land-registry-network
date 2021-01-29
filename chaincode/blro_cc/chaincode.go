package main

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode is the definition of the chaincode structure.
type Chaincode struct {
}

// Definition of the BLRO structure
type blro struct {
	ID             string   `json:"ID"`
	Name           string   `json:"Name"`
	Description    string   `json:"Description"`
	CompletedCases []string `json:"CompletedCases"`
	ActiveCases    []string `json:"ActiveCases"`
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Invoke()", fcn, params)

	if fcn == "createBLRO" {
		return cc.createBLRO(stub, params)
	} else if fcn == "readBLRO" {
		return cc.readBLRO(stub, params)
	} else if fcn == "addCase" {
		return cc.addCase(stub, params)
	} else if fcn == "completeCase" {
		return cc.completeCase(stub, params)
	} else {
		fmt.Println("Invoke() did not find func: " + fcn)
		return shim.Error("Received unknown function invocation!")
	}
}

// Function to create new blro (C of CRUD)
func (cc *Chaincode) createBLRO(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	creatorOrg, creatorCertIssuer, _, err := getTxCreatorInfo(stub)
	if !authenticateBLRO(creatorOrg, creatorCertIssuer) {
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

	key := "blro-" + params[0]
	ID := params[0]
	Name := params[1]
	Description := params[2]
	var CompletedCases []string
	var ActiveCases []string

	// Check if BLRO exists with Key => key
	blroAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to check if BLRO exists!")
	} else if blroAsBytes != nil {
		return shim.Error("BLRO Already Exists!")
	}

	// Generate BLRO from params provided
	blro := &blro{ID, Name, Description, CompletedCases, ActiveCases}
	blroJSONasBytes, err := json.Marshal(blro)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Put State of newly generated BLRO with Key => key
	err = stub.PutState(key, blroJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Returned on successful execution of the function
	return shim.Success(nil)
}

// Function to read an blro (R of CRUD)
func (cc *Chaincode) readBLRO(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	// Check if sufficient Params passed
	if len(params) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Check if Params are non-empty
	if len(params[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	key := "blro-" + params[0]

	// Get State of BLRO with Key => key
	blroAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + params[0] + "\"}"
		return shim.Error(jsonResp)
	} else if blroAsBytes == nil {
		jsonResp := "{\"Error\":\"BLRO does not exist!\"}"
		return shim.Error(jsonResp)
	}

	// Returned on successful execution of the function
	return shim.Success(blroAsBytes)
}

// Function to add new active case (U of CRUD)
func (cc *Chaincode) addCase(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	creatorOrg, creatorCertIssuer, _, err := getTxCreatorInfo(stub)
	if !authenticateRegistryOffice(creatorOrg, creatorCertIssuer) {
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

	key := "blro-" + params[0]
	CaseID := params[1]

	// Get State of BLRO with Key => key
	blroAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + params[0] + "\"}"
		return shim.Error(jsonResp)
	} else if blroAsBytes == nil {
		jsonResp := "{\"Error\":\"BLRO does not exist!\"}"
		return shim.Error(jsonResp)
	}

	// Create new BLRO Variable
	blroToUpdate := blro{}
	err = json.Unmarshal(blroAsBytes, &blroToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

	// Update blro.Owner => params[1]
	blroToUpdate.ActiveCases = append(blroToUpdate.ActiveCases, CaseID)

	// Convert to Byte[]
	blroJSONasBytes, err := json.Marshal(blroToUpdate)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Put updated State of the BLRO with Key => key
	err = stub.PutState(key, blroJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Returned on successful execution of the function
	return shim.Success(nil)
}

// Function to complete a case, add to CompletedCases, remove from ActiveCases (U of CRUD)
func (cc *Chaincode) completeCase(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	creatorOrg, creatorCertIssuer, _, err := getTxCreatorInfo(stub)
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

	key := "blro-" + params[0]
	CaseID := params[1]

	// Get State of BLRO with Key => key
	blroAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + params[0] + "\"}"
		return shim.Error(jsonResp)
	} else if blroAsBytes == nil {
		jsonResp := "{\"Error\":\"BLRO does not exist!\"}"
		return shim.Error(jsonResp)
	}

	// Create new BLRO Variable
	blroToUpdate := blro{}
	err = json.Unmarshal(blroAsBytes, &blroToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

	// Remove from ActiveCases
	for i, v := range blroToUpdate.ActiveCases {
		if v == CaseID {
			blroToUpdate.ActiveCases = append(blroToUpdate.ActiveCases[:i], blroToUpdate.ActiveCases[i+1:]...)
		}
	}

	// Append to CompletedCases
	blroToUpdate.CompletedCases = append(blroToUpdate.CompletedCases, CaseID)

	// Convert to Byte[]
	blroJSONasBytes, err := json.Marshal(blroToUpdate)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Put updated State of the BLRO with Key => key
	err = stub.PutState(key, blroJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
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

// Authenticate => BLRO
func authenticateBLRO(mspID string, certCN string) bool {
	return (mspID == "BLROMSP") && (certCN == "ca.blro.lran.com")
}

// Authenticate => RegistryOffice
func authenticateRegistryOffice(mspID string, certCN string) bool {
	return (mspID == "RegistryOfficeMSP") && (certCN == "ca.registryoffice.lran.com")
}
