
package main

import (
	"crypto/aes"
	"crypto/cipher"

	"encoding/base64"

	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("secret", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "setCCID" {
		return t.write(stub, args)
	} else if function == "setPrice" {
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "process" {
		return t.process(stub, args)
	}

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) setCCID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. name of the key and value to set")
	}

	key = "CCID"
	value = args[0]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) setPrice(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var uuid, key1, key2, value1, value2 string
	var err error
	fmt.Println("running write()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. name of the key and value to set")
	}

	uuid = args[0]
	key1 = uuid + "guid"
	value1 = args[1]
	err = stub.PutState(key1, []byte(value1)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	key2 = uuid + "price"
	value2 = args[2]
	err = stub.PutState(key2, []byte(value2)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *SimpleChaincode) process(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {


	//transferBean
	/*{
	"ccid" : "eiwquriewuripoq",
 "sendAddr" : "supernova27",
 "recvAddr": "8EBCFBE6C46F6145DC5A6E341265E60D455684541CDEA539A47B07E061404A54",
 "beanAmount" : 5
	}*/

	var jsonResp string
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(args[0])
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to Decode enckey\"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Printf("%q\n", ciphertext)

	secret := []byte("abcdefghijklmnop")
	iv := []byte("abcdefgh")

	block, err := aes.NewCipher(secret)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	fmt.Printf("%s\n", ciphertext)
	// Output: exampleplaintext
	return ciphertext, nil
}
