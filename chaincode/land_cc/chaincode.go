package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode is the definition of the chaincode structure.
type Chaincode struct {
}

// Defintion of transfer record
type transfer struct {
	PreviousOwner     string `json:"PreviousOwner"`
	CurrentOwner      string `json:"CurrentOwner"`
	TransferDate      int    `json:"TransferDate"`
	TransferRequestID string `json:"TransferRequest"`
	BLRO              string `json:"BLRO"`
}

// Definition of the Land structure
type land struct {
	ID      string     `json:"ID"`
	Address string     `json:"Address"`
	Owner   string     `json:"Owner"`
	History []transfer `json:"History"`
	Type    string     `json:"Type"`
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Invoke()", fcn, params)

	if fcn == "createLand" {
		return cc.createLand(stub, params)
	} else if fcn == "readLand" {
		return cc.readLand(stub, params)
	} else if fcn == "transferLand" {
		return cc.transferLand(stub, params)
	} else if fcn == "getLands" {
		return cc.getLands(stub, params)
	} else {
		fmt.Println("Invoke() did not find func: " + fcn)
		return shim.Error("Received unknown function invocation!")
	}
}

// Function to create new land (C of CRUD)
func (cc *Chaincode) createLand(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	creatorOrg, creatorCertIssuer, creator, err := getTxCreatorInfo(stub)
	if !authenticateBLRO(creatorOrg, creatorCertIssuer) {
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

	key := "land-" + params[0]
	ID := params[0]
	Address := params[1]
	Owner := params[2]
	Date := params[3]
	DateI, err := strconv.Atoi(Date)
	if err != nil {
		return shim.Error("Error: Invalid Date!")
	}

	// Check if Land exists with Key => key
	landAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to check if Land exists!")
	} else if landAsBytes != nil {
		return shim.Error("Land Already Exists!")
	}

	// Generate Initial Transfer Record
	var History []transfer
	initialHistory := transfer{"BLRO", Owner, DateI, "Land Created By BLRO", creator}
	History = append(History, initialHistory)

	// Generate Land from params provided
	land := &land{ID, Address, Owner, History, "LAND"}
	landJSONasBytes, err := json.Marshal(land)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Put State of newly generated Land with Key => key
	err = stub.PutState(key, landJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Returned on successful execution of the function
	return shim.Success(nil)
}

// Function to read an land (R of CRUD)
func (cc *Chaincode) readLand(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	// Check if sufficient Params passed
	if len(params) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Check if Params are non-empty
	if len(params[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	key := "land-" + params[0]

	// Get State of Land with Key => key
	landAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + params[0] + "\"}"
		return shim.Error(jsonResp)
	} else if landAsBytes == nil {
		jsonResp := "{\"Error\":\"Land does not exist!\"}"
		return shim.Error(jsonResp)
	}

	// Returned on successful execution of the function
	return shim.Success(landAsBytes)
}

// Function to update an land's owner (U of CRUD)
func (cc *Chaincode) transferLand(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	creatorOrg, creatorCertIssuer, creator, err := getTxCreatorInfo(stub)
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

	key := "land-" + params[0]
	CurrentOwner := params[1]
	TransferDate := params[2]
	TransferRequestID := params[3]

	DateI, err := strconv.Atoi(TransferDate)
	if err != nil {
		return shim.Error("Error: Invalid Date!")
	}

	// Get State of Land with Key => key
	landAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + params[0] + "\"}"
		return shim.Error(jsonResp)
	} else if landAsBytes == nil {
		jsonResp := "{\"Error\":\"Land does not exist!\"}"
		return shim.Error(jsonResp)
	}

	// Create new Land Variable
	landToTransfer := land{}
	err = json.Unmarshal(landAsBytes, &landToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

	// Append Transfer History
	initialHistory := transfer{landToTransfer.Owner, CurrentOwner, DateI, TransferRequestID, creator}
	landToTransfer.History = append(landToTransfer.History, initialHistory)

	// Update land.Owner => params[1]
	landToTransfer.Owner = CurrentOwner

	// Convert to Byte[]
	landJSONasBytes, err := json.Marshal(landToTransfer)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Put updated State of the Land with Key => key
	err = stub.PutState(key, landJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Returned on successful execution of the function
	return shim.Success(nil)
}

// Function to Delete an land (D of CRUD)
func (cc *Chaincode) getLands(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	// Check if sufficient Params passed
	if len(params) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Check if Params are non-empty
	if len(params[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	Owner := params[0]

	regex := "(?i:.*%s.*)"
	search := "{\"selector\": {\"$and\": [{\"Type\": \"LAND\" },{\"Owner\": { \"$regex\": \"%s\" }}]}}"

	search = fmt.Sprintf(search, fmt.Sprintf(regex, Owner))

	queryResults, err := getQueryResultForQueryString(stub, search)
	if err != nil {

	}

	return shim.Success(queryResults)
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

// Query Helpers
// +++++++++++++

// Construct Query Response from Iterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return &buffer, nil
}

// Get Query Result for Query String
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
