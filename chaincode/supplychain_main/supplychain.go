package main

import (
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "fmt"
    "errors"
)

// SupplychainChaincode
type SupplychainChaincode struct {}

var ACCESS_DENIED_RESPONSE = []byte("{\"failure\": \"ACCESS_DENIED\"}")
var NOT_FOUND_RESPONSE = []byte("{\"failure\": \"NOT_FOUND\"}")

func main() {
    err := shim.Start(new(SupplychainChaincode))
    if err != nil {
        fmt.Printf("Error starting Org Registrar chaincode: %s", err)
    }
}

func (t *SupplychainChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("initialising")
    InitDao(stub, args[0])

    return nil, nil
}

// Invoke is the entry point to invoke a chaincode function
func (t *SupplychainChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    callerDetails, err := GetCallerDetails(stub)
    if err != nil { fmt.Println("An error occured whilst obtaining the caller details"); return nil, err}

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function == "addOrder" {
        return t.processAddOrder(stub, callerDetails, args)
    }
    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SupplychainChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    callerDetails, err := GetCallerDetails(stub)
    if err != nil { fmt.Println("An error occured whilst obtaining the caller details"); return nil, err }

    if function == "getOrder" {
        return t.processGetOrder(stub, callerDetails, args)
    } else if function == "getAllOrders" {
        return t.processGetAllOrders(stub, callerDetails, args)
    }

    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query: " + function)
}

//=================================================================================================================================
//	 processAddOrder - Processes an addOrganistion request.
//          args -  Recipient,
//                  Address,
//                  SourceWarehouse,
//                  DeliveryCompany
//                  Items
//=================================================================================================================================
func (t *SupplychainChaincode) processAddOrder(stub shim.ChaincodeStubInterface, callerDetails CallerDetails, args[]string) ([]byte, error) {

    fmt.Println("running processAddOrder)")

    if len(args) != 5 {
        return nil, errors.New("Incorrect number of arguments. Expecting (Recipient, Address, SourceWarehouse, DeliveryCompany, Items)")
    }

    items, err := MarshallItems(args[4])

    if err != nil { return nil, LogAndError("Invalid items: " + args[4] + ", error: " + err.Error()) }

    order := NewOrder(stub.GetTxID(), args[0], args[1], args[2], args[3], items)

    return nil, AddOrder(stub, callerDetails, order)
}

//=================================================================================================================================
//	 processGetOrder - Processes a getOrder request.
//          args -  orderId
//=================================================================================================================================
func (t *SupplychainChaincode) processGetOrder(stub shim.ChaincodeStubInterface, callerDetails CallerDetails, args[]string) ([]byte, error) {

    fmt.Println("running processGetOrder()")

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting (OrderId)")
    }

    order, accessDenied, err := GetOrder(stub, callerDetails, args[0])

    if accessDenied {
        return ACCESS_DENIED_RESPONSE, nil
    }

    if err != nil {
        return NOT_FOUND_RESPONSE, nil
    }

    return marshall(order)
}

//=================================================================================================================================
//	 processGetAllOrders - Processes a getAllOrders request.
//
//=================================================================================================================================
func (t *SupplychainChaincode) processGetAllOrders(stub shim.ChaincodeStubInterface, callerDetails CallerDetails, args[]string) ([]byte, error) {

    fmt.Println("running processGetAllOrders()")

    if len(args) != 0 {
        return nil, errors.New("Incorrect number of arguments. Expecting NONE")
    }

    orders, err := GetAllOrders(stub, callerDetails)

    if err != nil { return nil, LogAndError("Error when retrieving orders")}

    return marshall(orders)
}