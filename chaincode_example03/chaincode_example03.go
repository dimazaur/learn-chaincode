/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This program is an erroneous chaincode program that attempts to put state in query context - query should return error
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init takes a string and int. These are stored as a key/value pair in the state
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))

	var A, Aval string // Entity
	var X int // Asset holding
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	// Initialize the chaincode
	A = args[0]
	Aval = args[1]
	X, err = strconv.Atoi(Aval)
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d\n", Aval)

	// Write the state to the ledger - this put is legal within Run
	err = stub.PutState(A, []byte(strconv.Itoa(X)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is a no-op
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var A, Aval string
	var X int
	var err error

	A = args[0]
	Aval = args[1]
	X, err = strconv.Atoi(Aval)
	if err != nil {
		return nil, errors.New("Invalid transaction amount, expecting a integer value")
	}

	err = stub.PutState(A, []byte(strconv.Itoa(X)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var A string // Entity
	//var X string // Asset holding
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	A = args[0]

	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}
	if Avalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	//X, err = strconv.Atoi(string(Avalbytes))

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	return []byte(jsonResp), nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
