package main
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	sc "github.com/hyperledger/fabric/protos/peer"
)
type SmartContract struct {
}
type Invoice struct {
	InvoiceFrom string `json:"invoiceFrom"`
	BilledTo  string `json:"billedTo"`
	InvoiceDate string `json:"invoiceDate"`
	InvoiceAmount  string `json:"invoiceAmount"`
	ItemDescription  string `json:"itemDescription"`
	GoodReceived  string `json:"goodReceived"`
	IsPaid  string `json:"isPaid"`
	PaidAmount  string `json:"paidAmount"`
	Repaid  string `json:"repaid"`
	RepaymentAmount  string `json:"repaymentAmount"`
}
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	 // Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	 // Route to the appropriate handler function to interact with the ledger appropriately
	 
	if function == "initLedger" {
		return s.initLedger(APIstub)
	}else if function == "raiseInvoice" {
		return s.raiseInvoice(APIstub, args)
	}else if function == "goodReceived" {
		return s.goodReceived(APIstub, args) 
	}else if function == "bankPayment" {
		return s.bankPayment(APIstub, args) 
	}else if function == "oemPayment" {
		return s.oemPayment(APIstub, args) 
	}else if function == "getHistoryForInvoice" {
		return s.getHistoryForInvoice(APIstub, args) 
	}else if function == "queryInvoiceByType" {
		return s.queryInvoiceByType(APIstub, args)
	}else if function == "queryAllInvoices" {
		return s.queryAllInvoices(APIstub)
	}
	return shim.Error("Invalid Smart Contract function name.")
}
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	
	// invoices := []Invoice{
	// 	Invoice{BilledTo: "OEM", InvoiceDate: "02/07/19",InvoiceAmount: "200000", ItemDescription: "Processor", GoodReceived: "False", IsPaid: "False", PaidAmount: "0", Repaid: "False", RepaymentAmount: "0"},
	// 	Invoice{BilledTo: "OEM", InvoiceDate: "02/08/19",InvoiceAmount: "22000", ItemDescription: "SSD", GoodReceived: "False", IsPaid: "False", PaidAmount: "0", Repaid: "False", RepaymentAmount: "0"},
	// 	Invoice{BilledTo: "OEM", InvoiceDate: "02/09/19",InvoiceAmount: "14000", ItemDescription: "RAM", GoodReceived: "False", IsPaid: "False", PaidAmount: "0", Repaid: "False", RepaymentAmount: "0"},
	// 	Invoice{BilledTo: "OEM", InvoiceDate: "02/10/19",InvoiceAmount: "50000", ItemDescription: "HDD", GoodReceived: "False", IsPaid: "False", PaidAmount: "0" , Repaid: "False", RepaymentAmount: "0"},
	// }

	// i := 0
	// for i < len(invoices) {
	// 	fmt.Println("i is ", i)
	// 	invoiceAsBytes, _ := json.Marshal(invoices[i])
	// 	APIstub.PutState("INVOICE"+strconv.Itoa(i), invoiceAsBytes)
	// 	fmt.Println("Added", invoices[i])
	// 	i = i + 1
	// }

	return shim.Success(nil)
}
func (s *SmartContract) raiseInvoice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	attr, _, _ := cid.GetAttributeValue(APIstub, "type")

	if(attr != "Supplier"){
		return shim.Error("You are not allowed to invoke raiseInvoice as " + attr)
	}
	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments. Expecting 10")
	}
	fmt.Printf("InvoiceFrom: " + attr)
	var invoice = Invoice{InvoiceFrom: attr, BilledTo: args[1], InvoiceDate: args[2] ,InvoiceAmount: args[3], ItemDescription: args[4], GoodReceived: args[5], IsPaid: args[6], PaidAmount: args[7], Repaid: args[8], RepaymentAmount: args[9]}

	invoiceAsBytes, _:= json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)
}
func (s *SmartContract) goodReceived(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	attr, _, _ := cid.GetAttributeValue(APIstub, "type")
	if(attr != "OEM"){
		return shim.Error("You are not allowed to invoke goodReceived " + attr) 
	}
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	json.Unmarshal(invoiceAsBytes, &invoice)
	invoice.GoodReceived = args[1]

	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)
	
	return shim.Success(nil)
}
func (s *SmartContract) bankPayment(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	attr, _, _ := cid.GetAttributeValue(APIstub, "type")
	if(attr != "Bank"){
		return shim.Error("You are not allowed to invoke bankPayment " + attr)
	}

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	json.Unmarshal(invoiceAsBytes, &invoice)
	invoice.PaidAmount = args[1]

	paid, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		// do something sensible
	}
	invoiceAmount, err := strconv.ParseFloat(invoice.InvoiceAmount, 32)
	if err != nil {
		// do something sensible
	}

	if(paid > invoiceAmount) {
		return shim.Error("Paid should be less than Invoice Amount!")
	}

	invoice.IsPaid = "True"
	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)
}
func (s *SmartContract) oemPayment(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	attr, _, _ := cid.GetAttributeValue(APIstub, "type")
	if(attr != "OEM"){
		return shim.Error("You are not allowed to invoke oemPayment " + attr)
	}

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	json.Unmarshal(invoiceAsBytes, &invoice)
	invoice.RepaymentAmount = args[1]

	repaymentAmount, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		// do something sensible
	}
	paidAmount, err := strconv.ParseFloat(invoice.PaidAmount, 32)
	if err != nil {
		// do something sensible
	}

	if(paidAmount > repaymentAmount) {
		return shim.Error("Repayment should be greater than Paid Amount!")
	}

	invoice.Repaid = "True"
	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)
}
func (s *SmartContract) getHistoryForInvoice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	invoiceNumber := args[0]

	resultsIterator, err := APIstub.GetHistoryForKey(invoiceNumber)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the car
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
		buffer.WriteString(string(response.Value))

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
func (s *SmartContract) queryInvoiceByType(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var queryType = ""
	
	attr, _, _ := cid.GetAttributeValue(APIstub, "type")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	if (attr == "Supplier") {
		queryType = "invoiceFrom"
	}else if (attr == "OEM") {
		queryType = "billedTo"
	} else {
		return shim.Error("You are not allowed to Invoke")
	}

	queryString := fmt.Sprintf("{\"selector\":{\"%s\":\"%s\"}}", queryType, attr)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
func (s *SmartContract) queryAllInvoices(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "INVOICE0"
	endKey := "INVOICE999"

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

	fmt.Printf("- queryAllInvoices:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
func getQueryResultForQueryString(APIstub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {


	resultsIterator, err := APIstub.GetQueryResult(queryString)
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
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")	

	return buffer.Bytes(), nil
}
func main() {
 
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}