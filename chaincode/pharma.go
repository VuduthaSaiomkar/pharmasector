
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Asset Declaration */

/* Blockchain /Couchdb - TransactionID - Onchain
OffChain - ipfs
front-end - angular js - data can be pulled from database..

AssetID  AssetName   Distributor    WS	Ret transactionid
1		abc				A			     	    12345
2       abc             A 			aa		    123456
3		abc				A			aa   bb     1234567

Fabric CA server
Register distributors, Wholesalers and Retailers on CA server as per their role
*/

	//DistributorID        int    `json:"distributorid"`
	//WholesalerID  int  `json:"wholesalerid"`
	//RetailerID       int    `json:"retailerid"`
	//CreatedBy  string  `json:"createdby"`
	//TransactionID string `json:"transactionid"`
	//transactionID not needed as fabric creates on its own

	//ExpDate time.Time `json:"ExpDate"`
	//1 cretedatetime is needed as separate transactions will be generated
	//during the change of ownership/sale..even this 1 createdatetime also
	//not needed as fabric creates this value during the creation of each transaction

	//Manu_CreatedDatetime time.Time `json:"Manu_CreatedDatetime"`
	//Dist_CreatedDatetime time.Time `json:"Dist_CreatedDatetime"`
	//WS_CreatedDatetime time.Time `json:"WS_CreatedDatetime"`
	//Retailer_CreatedDatetime time.Time `json:"Retailer_CreatedDatetime"`
//}

/*
We will store this in RDBMS and retrieve data to the interface, selected data will be stored on blockchain
Interface usually developed in angular which can connect to Database and pull the data

2) Create Manufacturer table off chain
ManufacturerID	int
ManufacturerName int
3) Create Distributor table offchain
DistributorID
DistributorName
Address
City
State
Zip
4) Create Wholesaler table offchain
WholesalerID
WholesalerName
Address
City
State
Zip
5) Create Retailer table offchain
RetailerID
RetailerName
Address
City
State
Zip

I see here a challenge, I don’t expect retailers or wholesalers data will be at
Manufacturer’s place.. In that case, how to read this data? Need to discuss..

If this data is manufacturer’s side then there won’t be an issue..
If it is at different places then read that source and integrate it with this application.
*/




type Asset struct {
	AssetID            string  `json:"assetid"`
	AssetName          string  `json:"assetname"`
	QRID               string  `json:"qrid"`
	BoxID      		   string   `json:"boxid"`
	ConsignmentID      string    `json:"consignmentid"`
	OwnerID            string  `json:"ownerid"`
	OwnerRole          string  `json:"ownerrole"`
	MfgDate 		   string `json:"MfgDate"`
	CustomerPhone      int   `json:"customerphone",omitempty`
	ManufactureOwner string `json:"manuowner"`
	DistributorOwner string `json:"distowner",omitempty`
	WholesalerOwner string  `json:"wholeowner",omitempty`
	RetailerOwner  string `json:retailerowner",omitempty`
	CustomerOwner  string 	`json:customerowner",omitempty`
}




func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()

	if function == "manu_saletransaction" {
		/* On the interface manufacturer will see
		AssetName,  QRID, BoxID ,ConsignmentID, OwnerID
		So total 5 argunments to the function
		AssetID will be unique id.. after submit create this unique id and store it
		OwnerRole will be default value of "Distributor" (make sure to update it for all assets belonging to respective
			boxid, consignmentid) because
		manufacturer will be assigning it to distributor
		*/

		return s.manu_saletransaction(stub, args)
	} else if function == "manu_creation"{
		return s.manu_creation(stub,args)
	}else if function == "manu_consignments_list" {
		/*
		On the interface manufacturer will see a list of consignments (like consignementID, DistributorID/OwnerID)
		So no arguments just pull all consignments
		***
		Make sure once sale has been made that should not be available for sale again..
		*/
		return s.manu_consignments_list(stub, args)
	} else if function == "consignment_detail" {
		/*
		On the consignments list interface if manufacturer clicks on given consignment

		So 1 argument consignmentid will be passed..
		for example, if consignment has 10 boxes and each box has 1000 assets then his next screen
		should present after selecting the consignment
		Box1 ID – assets count
		Box2 ID – assets count
		Box3 ID – assets count
		*/
		return s.consignment_detail(stub, args)
	} else if function == "dist_sale_transaction" {
		/*
		Upon clicking Sale transaction button on above consignment detail, Create new transactions in
		the assets table with selected WholesalerID/OwnerID and boxes in assets table

		On the interface manufacturer will see
		So 2 arguments wholesalerID, list of BoxIDs will be passed..(no need assetIDs because assetids are already tied to boxes)
		(one box may have 1000 assets and one consignment may have 1000 boxes)
		OwnerRole will be default value of "Wholesaler" because Distributor will be assigning
		it to wholesaler
		*/
		return s.dist_sale_transaction(stub, args)
	}  else if function == "ws_sale_transaction" {
		/*
		Upon clicking Sale transaction button on above ws_carton_list, create new transactions
		in the assets table with selected retailerid/OwnerID and boxes

		So 2 arguments retailerid, list of BoxIDs will be passed..(no need assetIDs because assetids are already tied to boxes)
		(note that one box may have 1000 assets and one consignment may have 1000 boxes)
		OwnerRole will be default value of "Retailer" because Wholesaler will be assigning
		it to retailer
		*/
		return s.ws_sale_transaction(stub, args)
	} else if function == "ws_box_detail" {
		/*

		On the above ws_carton_list interface if wholesaler clicks on a given box

		So 1 argument boxID will be passed..
		for example, if Box has 1000 assets his next screen
		should present after selecting the consignment

		a) Pre-populated Retailer IDs
		b) submit button with  a label "Submit Sale transaction"

		Upon submit, invoke  ws_sale_transaction by passing 2 arguments retailerID, boxid
		OwnerRole will be default value of "Retailer" because Wholesaler will be assigning
		it to retailer

		*** we may need a status attribute to hold a value of whether particular asset is sold or not ??
		Make sure once sale has been made that should not be available for sale again..

		*/
		return s.ws_box_detail(stub, args)
	} else if function == "ret_sale_transaction_single" {
		/*
		Retailer can sell a single asset or a box
		but he may scan a single asset or just perform a sale without even scanning it..
		if he scans pass assetID and phone number of customer and create a transaction for that asset
		with his phone number and
		OwnerRole will be default value of "Customer" because retailer will be selling it to customer

		as soon as he scans we should not invoke sale transaction because sometimes it would be just a lookup too
		may be present him a "submit to sell" and cancel button on the interace.. invoke ret_sale_transaction only
		when he clicks on "submit to sell" button

		*/
		return s.ret_sale_transaction_single(stub, args)
	}else if function == "ret_sale_transaction_box"{


		return s.ret_sale_transaction_box(stub, args)

	} else if function == "scan_asset" {
		/*
		it is more like a history of a Asset or asset detail screen
		On the interface, we will give a assetid and hit submit

		show him the history of the asset starting from manufacturer to him
		so 1 argument (assetID) needed for this function to work
		*/
		return s.scan_asset(stub,args)
	} else if function == "ret_box_list" {
		/*TODO*/
		return s.ret_box_list(stub,args)
	} else if function =="transferassetcontainsbox"{
		return s.transferassetcontainsbox(stub,args)
	} else if function =="transferconsignment"{

		return s.transferconsignment(stub,args)
	}else if function =="transferassetbydealer"{
		
		return s.transferassetbydealer(stub,args)
	}else if function =="dist_consignments_list"{
		return s.dist_consignments_list(stub,args)
	}else if function =="SaleHistory_Ratail"{
		return s.SaleHistory_Ratail(stub,args)
	}else if function =="SaleHistory_Manu"{
		return s.SaleHistory_Manu(stub,args)
	}else if function =="SaleHistory_Dist"{
		return s.SaleHistory_Dist(stub,args)
	}else if function=="SaleHistory_Ratail"{
		return s.SaleHistory_Ratail(stub,args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) ret_sale_transaction_single(stub shim.ChaincodeStubInterface,args []string) sc.Response{

	if len(args) !=3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	assetid := args[0]
	assetnewowner := strings.ToLower(args[1])
	number, err := strconv.Atoi(args[2])
	assetAsBytes, err := stub.GetState(assetid)
	if err != nil {
		return shim.Error("Failed to get box:" + err.Error())
	} else if assetAsBytes == nil {
		return shim.Error("box does not exist")
	}

	modfiyasset := Asset{}
	err = json.Unmarshal(assetAsBytes, &modfiyasset) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	modfiyasset.OwnerID = assetnewowner //change the owner
	modfiyasset.OwnerRole="customer"
	modfiyasset.CustomerPhone=number
	modfiyasset.CustomerOwner=assetnewowner

	assetJSONasBytes, _ := json.Marshal(modfiyasset)
	err = stub.PutState(assetid, assetJSONasBytes) //rewrite the asset
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "consignmentid~ownerid~assetid" //query by cosingmentid
	consignmentwithowner, err := stub.CreateCompositeKey(indexName, []string{modfiyasset.ConsignmentID,modfiyasset.OwnerID,modfiyasset.AssetID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the box.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(consignmentwithowner, value)
	fmt.Println("- end of assetid")
	indexName1 := "boxid~ownerid~assetid" //query by box id
	boxidwithowner, err := stub.CreateCompositeKey(indexName1, []string{modfiyasset.BoxID,modfiyasset.OwnerID,modfiyasset.AssetID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the box.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	stub.PutState(boxidwithowner, value)
	fmt.Println("- end of assetid")


	fmt.Println("- end transfer of asset (success)")

	return shim.Success(nil)
}



func (s *SmartContract) ret_sale_transaction_box(stub shim.ChaincodeStubInterface,args []string) sc.Response{


	if len(args) != 4{
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	newownerid:=strings.ToLower(args[0])
	assetowner:=strings.ToLower(args[1])
	boxid:=strings.ToLower(args[2])
	number:=args[3]
	boxidassetResultsIterator, err := stub.GetStateByPartialCompositeKey("boxid~ownerid~assetid", []string{boxid,assetowner})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer boxidassetResultsIterator.Close()

	// Iterate through result set and for each box found, transfer to newOwner
	var i int
	for i = 0; boxidassetResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the box name from the composite key
		responseRange, err := boxidassetResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedboxid := compositeKeyParts[0]
		returnedownerName := compositeKeyParts[1]
		returnboxassetid:=compositeKeyParts[2]
		fmt.Printf("- found a box from index:%s boxid:%s name:%s\n", objectType, returnedboxid, returnedownerName,returnboxassetid)

		// Now call the transfer function for the found box.
		// Re-use the same function that is used to transfer individual consigment

		response := s.ret_sale_transaction_single(stub, []string{returnboxassetid, newownerid,number})


		// if the transfer failed break out of loop and return error
		if response.Status != shim.OK {
			return shim.Error("Transfer failed: " + response.Message)
		}
	}

	responsePayload := fmt.Sprintf("Transferred %d %s consigment to %s", i,boxid, assetowner)
	fmt.Println("- end transferconsigmentBasedOnbox: " + responsePayload)
	return shim.Success([]byte(responsePayload))

}


func (s *SmartContract) ws_box_detail(stub shim.ChaincodeStubInterface,args []string) sc.Response{


	if len(args) !=2{
		return shim.Error("sorry you missed ownerid oops!!")
	}

	ownerid:=strings.ToLower(args[0])
	boxid:=strings.ToLower(args[1])
	ownerrole:="wholesaler"

    queryString := fmt.Sprintf("{\"selector\":{\"ownerid\":\"%s\",\"boxid\":\"%s\",\"ownerrole\":\"%s\"}}",ownerid,boxid,ownerrole)
    resultsIterator, err := stub.GetQueryResult(queryString)
    defer resultsIterator.Close()
    if err != nil {
        return shim.Error(err.Error())
    }
    // buffer is a JSON array containing QueryRecords
    var buffer bytes.Buffer
    buffer.WriteString("[")
    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        queryResponse,
        err := resultsIterator.Next()
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
    fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
    return shim.Success(buffer.Bytes())


}


func (s *SmartContract) ws_sale_transaction(stub shim.ChaincodeStubInterface,args []string) sc.Response{


	if len(args) != 3{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	newownerid:=strings.ToLower(args[0])
	assetowner:=strings.ToLower(args[1])
	boxid:=strings.ToLower(args[2])
	boxidassetResultsIterator, err := stub.GetStateByPartialCompositeKey("boxid~ownerid~assetid", []string{boxid,assetowner})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer boxidassetResultsIterator.Close()

	// Iterate through result set and for each box found, transfer to newOwner
	var i int
	for i = 0; boxidassetResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the box name from the composite key
		responseRange, err := boxidassetResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedboxid := compositeKeyParts[0]
		returnedownerName := compositeKeyParts[1]
		returnboxassetid:=compositeKeyParts[2]
		fmt.Printf("- found a box from index:%s boxid:%s name:%s\n", objectType, returnedboxid, returnedownerName,returnboxassetid)

		// Now call the transfer function for the found box.
		// Re-use the same function that is used to transfer individual consigment

			response := s.transferassetcontainsbox(stub, []string{returnboxassetid, newownerid})


		// if the transfer failed break out of loop and return error
		if response.Status != shim.OK {
			return shim.Error("Transfer failed: " + response.Message)
		}
	}

	responsePayload := fmt.Sprintf("Transferred %d %s consigment to %s", i, boxid, assetowner)
	fmt.Println("- end transferconsigmentBasedOnbox: " + responsePayload)
	return shim.Success([]byte(responsePayload))


}

func (t *SmartContract) transferassetcontainsbox(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	assetid := args[0]
	assetnewowner := strings.ToLower(args[1])

	assetAsBytes, err := stub.GetState(assetid)
	if err != nil {
		return shim.Error("Failed to get box:" + err.Error())
	} else if assetAsBytes == nil {
		return shim.Error("box does not exist")
	}

	modfiyasset := Asset{}
	err = json.Unmarshal(assetAsBytes, &modfiyasset) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	modfiyasset.OwnerID = assetnewowner //change the owner
	modfiyasset.OwnerRole="retailer"
	modfiyasset.RetailerOwner=assetnewowner

	assetJSONasBytes, _ := json.Marshal(modfiyasset)
	err = stub.PutState(assetid, assetJSONasBytes) //rewrite the asset
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "consignmentid~ownerid~assetid" //query by cosingmentid
	consignmentwithowner, err := stub.CreateCompositeKey(indexName, []string{modfiyasset.ConsignmentID,modfiyasset.OwnerID,modfiyasset.AssetID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the box.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(consignmentwithowner, value)
	fmt.Println("- end of assetid")
	indexName1 := "boxid~ownerid~assetid" //query by box id
	boxidwithowner, err := stub.CreateCompositeKey(indexName1, []string{modfiyasset.BoxID,modfiyasset.OwnerID,modfiyasset.AssetID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the box.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	stub.PutState(boxidwithowner, value)
	fmt.Println("- end of assetid")

	fmt.Println("- end transfer of asset (success)")
	return shim.Success(nil)
}


func (s *SmartContract) ret_box_list(stub shim.ChaincodeStubInterface,args []string) sc.Response{

	if len(args) !=2{
		return shim.Error("sorry you missed ownerid oops!!")
	}

	ownerid:=strings.ToLower(args[0])
	boxid:=strings.ToLower(args[1])
	ownerrole:="retailer"

    queryString := fmt.Sprintf("{\"selector\":{\"boxid\":\"%s\",\"ownerid\":\"%s\",\"ownerrole\":\"%s\"}}",boxid,ownerid,ownerrole)
    resultsIterator, err := stub.GetQueryResult(queryString)
    defer resultsIterator.Close()
    if err != nil {
        return shim.Error(err.Error())
    }
    // buffer is a JSON array containing QueryRecords
    var buffer bytes.Buffer
    buffer.WriteString("[")
    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        queryResponse,
        err := resultsIterator.Next()
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
    fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
    return shim.Success(buffer.Bytes())

}

func (s *SmartContract) dist_consignments_list(stub shim.ChaincodeStubInterface, args []string) sc.Response{
	if len(args) !=2{
		return shim.Error("sorry you missed ownerid oops!!")
	}

	ownerid:=strings.ToLower(args[0])
	consigmnetid:=strings.ToLower(args[1])
	ownerrole:="dealer"
    queryString := fmt.Sprintf("{\"selector\":{\"ownerid\":\"%s\",\"consignmentid\":\"%s\",\"ownerrole\":\"%s\"}}",ownerid,consigmnetid,ownerrole)
    resultsIterator, err := stub.GetQueryResult(queryString)
    if err != nil {
        return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

    // buffer is a JSON array containing QueryRecords
    var buffer bytes.Buffer
    buffer.WriteString("[")
    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        queryResponse,
        err := resultsIterator.Next()
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
    fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
    return shim.Success(buffer.Bytes())


}

func (s *SmartContract) manu_consignments_list(stub shim.ChaincodeStubInterface,args []string) sc.Response{

	if len(args) !=1{
		return shim.Error("sorry you missed ownerid oops!!")
	}

	queryString:=args[0]
	//queryString := fmt.Sprintf("{\"selector\":{\"OwnerID\":\"%s\"}}",ownerid)
	fmt.Printf("%s",queryString);
	resultsIterator, err := stub.GetQueryResult(queryString)
	defer resultsIterator.Close()

    if err != nil {
        return shim.Error(err.Error())
	}
    // buffer is a JSON array containing QueryRecords
    var buffer bytes.Buffer
    buffer.WriteString("[")
    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        queryResponse,
        err := resultsIterator.Next()
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
    fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
    return shim.Success(buffer.Bytes())

}

func (s *SmartContract) consignment_detail(stub shim.ChaincodeStubInterface,args []string) sc.Response{

	if len(args) !=2{
		return shim.Error("sorry you missed consignment id oops!!")
	}

	ownerid:=strings.ToLower(args[0])
	consignementid:=strings.ToLower(args[1])
	ownerrole:="manufacture"
    queryString := fmt.Sprintf("{\"selector\":{\"consignmentid\":\"%s\",\"ownerid\":\"%s\",\"ownerrole\":\"%s\"}}",consignementid,ownerid,ownerrole)
    queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) SaleHistory_Manu(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) !=1{
		return shim.Error("sorry you missed consignment id oops!!")
	}

	ownerid:=strings.ToLower(args[0])
    queryString := fmt.Sprintf("{\"selector\":{\"manuowner\":\"%s\"}}",ownerid)
    queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)

}


func (s *SmartContract) SaleHistory_Dist(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) !=1{
		return shim.Error("sorry you missed consignment id oops!!")
	}

	ownerid:=strings.ToLower(args[0])
    queryString := fmt.Sprintf("{\"selector\":{\"distowner\":\"%s\"}}",ownerid)
    queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)

}

func (s *SmartContract) SaleHistory_Wsale(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) !=1{
		return shim.Error("sorry you missed consignment id oops!!")
	}

	ownerid:=strings.ToLower(args[0])
    queryString := fmt.Sprintf("{\"selector\":{\"wholeowner\":\"%s\"}}",ownerid)
    queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)

}
func (s *SmartContract) SaleHistory_Ratail(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) !=1{
		return shim.Error("sorry you missed consignment id oops!!")
	}

	ownerid:=strings.ToLower(args[0])
    queryString := fmt.Sprintf("{\"selector\":{\"retailerowner\":\"%s\"}}",ownerid)
    queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)

}


func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
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

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (s *SmartContract) manu_creation(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	assetname:=strings.ToLower(args[0])
	qrid:=strings.ToLower(args[1])
	boxid:=strings.ToLower(args[2])
	consignementid:=strings.ToLower(args[3])
	ownerid:=strings.ToLower(args[4])
	assetid:=assetname+""+consignementid
	
	Asset := Asset{AssetID:assetid,AssetName:assetname, QRID: qrid, BoxID: boxid, ConsignmentID: consignementid, OwnerID: ownerid,OwnerRole:"manufacture",MfgDate:args[5],ManufactureOwner:ownerid}

	AssetAsBytes, _ := json.Marshal(Asset)
	stub.PutState(assetid, AssetAsBytes)
	indexName := "consignmentid~ownerid~assetid" //query by cosingmentid
	consignmentwithowner, err := stub.CreateCompositeKey(indexName, []string{Asset.ConsignmentID,Asset.OwnerID,Asset.AssetID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the box.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(consignmentwithowner, value)
	fmt.Println("- end of assetid")

	indexName1 := "boxid~ownerid~assetid" //query by box id
	boxidwithowner, err := stub.CreateCompositeKey(indexName1, []string{Asset.BoxID,Asset.OwnerID,Asset.AssetID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the box.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	stub.PutState(boxidwithowner, value)
	fmt.Println("- end of assetid")

	return shim.Success(nil)
}


func (s *SmartContract) transferconsignment(stub shim.ChaincodeStubInterface,args []string) sc.Response{


	if len(args) != 2{
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	assetid:=strings.ToLower(args[0])
	assestasbytes, _ := stub.GetState(assetid)
	asset := Asset{}

	json.Unmarshal(assestasbytes, &asset)
	
	 
	asset.OwnerID = strings.ToLower(args[1])
	asset.OwnerRole="dealer"
	asset.DistributorOwner=strings.ToLower(args[1])
	asset_modify, _ := json.Marshal(asset)
	
	stub.PutState(args[0], asset_modify)
	indexName := "consignmentid~ownerid~assetid" //query by cosingmentid
	consignmentwithowner, err := stub.CreateCompositeKey(indexName, []string{asset.ConsignmentID,asset.OwnerID,asset.AssetID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the box.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(consignmentwithowner, value)
	fmt.Println("- end of assetid")
	indexName1 := "boxid~ownerid~assetid" //query by box id
	boxidwithowner, err := stub.CreateCompositeKey(indexName1, []string{asset.BoxID,asset.OwnerID,asset.AssetID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the box.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	stub.PutState(boxidwithowner, value)
	fmt.Println("- end of assetid")



	return shim.Success(nil)

}

func (s *SmartContract) manu_saletransaction(stub shim.ChaincodeStubInterface,args []string) sc.Response {


	if len(args) != 3{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	newownerid:=strings.ToLower(args[0])
	assetowner:=strings.ToLower(args[1])
	consignmentid:=strings.ToLower(args[2])
	boxidassetResultsIterator, err := stub.GetStateByPartialCompositeKey("consignmentid~ownerid~assetid", []string{consignmentid,assetowner})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer boxidassetResultsIterator.Close()

	// Iterate through result set and for each box found, transfer to newOwner
	var i int
	for i = 0; boxidassetResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the box name from the composite key
		responseRange, err := boxidassetResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedboxid := compositeKeyParts[0]
		returnedownerName := compositeKeyParts[1]
		returnboxassetid:=compositeKeyParts[2]
		fmt.Printf("- found a box from index:%s boxid:%s name:%s\n", objectType, returnedboxid, returnedownerName,returnboxassetid)

		// Now call the transfer function for the found box.
		// Re-use the same function that is used to transfer individual consigment

			response := s.transferconsignment(stub, []string{returnboxassetid, newownerid})

		// if the transfer failed break out of loop and return error
		if response.Status != shim.OK {
			return shim.Error("Transfer failed: " + response.Message)
		}
	}

	responsePayload := fmt.Sprintf("Transferred %d %s consigment to %s", i, consignmentid, newownerid)
	fmt.Println("- end transferconsigmentBasedOnbox: " + responsePayload)
	return shim.Success([]byte(responsePayload))

}


func (s *SmartContract) transferassetbydealer(stub shim.ChaincodeStubInterface,args []string) sc.Response{


	if len(args) != 2{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	assetid:=strings.ToLower(args[0])
	assestasbytes, _ := stub.GetState(assetid)
	asset := Asset{}

	json.Unmarshal(assestasbytes, &asset)
	asset.OwnerID = strings.ToLower(args[1])
	asset.OwnerRole ="wholesaler"
	asset.WholesalerOwner=strings.ToLower(args[1])
	asset_modify, _ := json.Marshal(asset)
	stub.PutState(args[0], asset_modify)
	indexName := "consignmentid~ownerid~assetid" //query by cosingmentid
	consignmentwithowner, err := stub.CreateCompositeKey(indexName, []string{asset.ConsignmentID,asset.OwnerID,asset.AssetID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the box.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(consignmentwithowner, value)
	fmt.Println("- end of assetid")
	indexName1 := "boxid~ownerid~assetid" //query by box id
	boxidwithowner, err := stub.CreateCompositeKey(indexName1, []string{asset.BoxID,asset.OwnerID,asset.AssetID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the box.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	stub.PutState(boxidwithowner, value)
	fmt.Println("- end of assetid")


	return shim.Success(nil)
}

func (s *SmartContract) dist_sale_transaction(stub shim.ChaincodeStubInterface, args []string) sc.Response{

	if len(args) != 3{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	newownerid:=strings.ToLower(args[0])
	assetowner:=strings.ToLower(args[1])
	consignmentid:=strings.ToLower(args[2])
	boxidassetResultsIterator, err := stub.GetStateByPartialCompositeKey("consignmentid~ownerid~assetid", []string{consignmentid,assetowner})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer boxidassetResultsIterator.Close()

	// Iterate through result set and for each box found, transfer to newOwner
	var i int
	for i = 0; boxidassetResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the box name from the composite key
		responseRange, err := boxidassetResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedboxid := compositeKeyParts[0]
		returnedownerName := compositeKeyParts[1]
		returnboxassetid:=compositeKeyParts[2]
		fmt.Printf("- found a box from index:%s boxid:%s name:%s\n", objectType, returnedboxid, returnedownerName,returnboxassetid)

		// Now call the transfer function for the found box.
		// Re-use the same function that is used to transfer individual consigment

			response := s.transferassetbydealer(stub, []string{returnboxassetid, newownerid})

		// if the transfer failed break out of loop and return error
		if response.Status != shim.OK {
			return shim.Error("Transfer failed: " + response.Message)
		}
	}

	responsePayload := fmt.Sprintf("Transferred %d %s consigment to %s", i, consignmentid, newownerid)
	fmt.Println("- end transferconsigmentBasedOnbox: " + responsePayload)
	return shim.Success([]byte(responsePayload))

}



func (s *SmartContract) scan_asset(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	assetid := args[0]
	fmt.Printf("- start getHistoryForasset: %s\n", assetid)
	resultsIterator, err := stub.GetHistoryForKey(assetid)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the box
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON asset)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForAsset returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
