/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	//"strconv"

	/*"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"*/
	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Fruit struct {
	Name   string `json:"name"`
	Size  string `json:"size"`
	Colour string `json:"colour"`
	Season  string `json:"season"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryFruit" {
		return s.queryFruit(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub, args)
	} else if function == "createFruit" {
		return s.createFruit(APIstub)
	} else if function == "queryAllFruits" {
		return s.queryAllFruits(APIstub)
	/*} else if function == "changeName" {
		return s.changeName(APIstub, args)*/
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryFruit(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	fruitAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(fruitAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Entered the init function")
	fruits := []Fruit{
		Fruit{Name: "Apple", Size: "Medium", Colour: "Red", Season: "Monsoon"},
		Fruit{Name: "Orange", Size: "Medium", Colour: "Orange", Season: "Winter"},
		Fruit{Name: "Watermelon", Size: "Large", Colour: "Green", Season: "Summer"},
		Fruit{Name: "Grapes", Size: "Small", Colour: "Purple", Season: "Summer"},
		Fruit{Name: "Pear", Size: "Medium", Colour: "Green", Season: "Monsoon"},
	}

	i := 0
	for i < len(fruits) {
		fmt.Println("i is ", i)
		fruitAsBytes, _ := json.Marshal(fruits[i])
		APIstub.PutState(fruits[i].Name, fruitAsBytes)
		fmt.Println("Added", fruits[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createFruit(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var fruit = Fruit{Name: args[0], Size: args[1], Colour: args[2], Season: args[3]}

	fruitAsBytes, _ := json.Marshal(fruit)
	APIstub.PutState(args[0], fruitAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllFruits(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := ""
	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllFruits:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*func (s *SmartContract) changeCarOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	carAsBytes, _ := APIstub.GetState(args[0])
	car := Car{}

	json.Unmarshal(carAsBytes, &car)
	car.Owner = args[1]

	carAsBytes, _ = json.Marshal(car)
	APIstub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}*/

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
