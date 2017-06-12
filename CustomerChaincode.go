package main

import (
	"errors"
	"fmt"
	//"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/golang/protobuf/ptypes/timestamp"
	"io"
        "io/ioutil"
        "log"
        "os"
)

// Customer Chaincode implementation
type CustomerChaincode struct {
}

var customerIndexTxStr = "_customerIndexTxStr"

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

type InsuranceClientInformation struct{
	clientId string `json:"clientId"`
	personalInfo InsurancePersonalInfo
	residenceAddress Address
	permanentAddress Address 
	officeAddress Address
	contactDetails ContactDetails
	employmentDetails EmploymentDetails
	personalAssets PersonalAssets
	bankAccountDetails BankDetails
	kycDocuments []KYCDocuments
	
}

type InsurancePersonalInfo struct{
	applicantNumber string `json:"applicantNumber"`
	firstName string `json:"firstName"`
	middleName string `json:"middleName"`
	lastName string `json:"lastName"`
	dateOfBirth string `json:"dateOfBirth"`
	panNumber string `json:"panNumber"`
	passportNumber string `json:"passportNumber"`
	residentStatus string `json:"residentStatus"`
	residencePlaceOwnership string `json:"residencePlaceOwnership"`
	numberofDependents string `json:"numberofDependents"`
	qualification string `json:"qualification"`
	annualIncome string `json:"annualIncome"`
	gender string `json:"gender"`
	maritalStatus string `json:"maritalStatus"`
	criminalRecordDetails string `json:"criminalRecordDetails"`
	financialStability string `json:"financialStability"`
	creditScore string `json:"creditScore"`
}

type Address struct{
	addrLine1 string `json:"addrLine1"`
	addrLine2 string `json:"addrLine2"`
	city string `json:"city"`
	province string `json:"province"`
	country string `json:"country"`
	postalCode string `json:"postalCode"`
	addressType string `json:"addressType"`
	searchLocation string `json:"searchLocation"`
}

type ContactDetails struct{
	homeNumber string `json:"homeNumber"`
	officeNumber string `json:"officeNumber"`
	mobileNumber string `json:"mobileNumber"`
	emailId string `json:"emailId"`
}

type EmploymentDetails struct{
	nameOfEmployer string `json:"nameOfEmployer"`
	designation string `json:"designation"`
	title string `json:"title"`
	noOfYearsExperience string `json:"noOfYearsExperience"`	
}


type PersonalAssets struct{
	assetType string `json:"assetType"`
	assetName string `json:"assetName"`
	details string `json:"details"`
	valueOfAsset string `json:"valueOfAsset"`
	asOnDate string `json:"asOnDate"`
	
}

type BankDetails struct{
	bankName string `json:"bankName"`
	bankBranch string `json:"bankBranch"`
	accountNo string `json:"accountNo"`
	swiftCode string `json:"swiftCode"`		
}

type KYCDocuments struct{
	documentName string `json:"documentName"`
	base64String string `json:"base64String"`	
}

func (t *CustomerChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	// Initialize the chaincode

	var insuranceClientInformationTxs []InsuranceClientInformation
	jsonAsBytes, _ := json.Marshal(insuranceClientInformationTxs)
	err = stub.PutState(customerIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	fmt.Printf("Deployment of Customer ChainCode is completed\n")
	
	return nil, nil
}

// Add customer data for the policy
func (t *CustomerChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	var TAX_IDENTIFIER string // Entities
	var UNIQUE_IDENTIFIER string

    	var err error
	
    	fmt.Printf("********Invoke Call with args length :%s\n", len(args))
	
	if len(args) < 56 {
	    	fmt.Printf("********Inside Invoke length:%s\n", len(args))
		return nil, errors.New("Incorrect number of arguments. Need 31 arguments")
	}
	TAX_IDENTIFIER = args[6]
	UNIQUE_IDENTIFIER = args[7]
	if (TAX_IDENTIFIER == "" || UNIQUE_IDENTIFIER == ""){
		return nil, errors.New(" Tax Identifier and Unique Identifier are mandatory")
	}
	
	//var requiredObj CustomerData
	var objFound bool
	CustomerTxsAsBytes, err := stub.GetState(customerIndexTxStr)
	if err != nil {
		return CustomerTxsAsBytes, errors.New("Failed to get Customer Records")
	}
	var CustomerTxObjects []InsuranceClientInformation
	var CustomerTxObjects1 []InsuranceClientInformation
	json.Unmarshal(CustomerTxsAsBytes, &CustomerTxObjects)
	length := len(CustomerTxObjects)
	fmt.Printf("Output from chaincode: %s\n", CustomerTxsAsBytes)

	
	objFound = false
	var counter int
	// iterate
	for i := 0; i < length; i++ {
		obj := CustomerTxObjects[i]
		//if ((customer_id == obj.CUSTOMER_ID) && (customer_name == obj.CUSTOMER_NAME) && (customer_dob == obj.CUSTOMER_DOB)) 
		
		if (((obj.personalInfo.panNumber) == TAX_IDENTIFIER) && ((obj.personalInfo.passportNumber) == UNIQUE_IDENTIFIER)){			
			CustomerTxObjects1 = append(CustomerTxObjects1,obj)
			//requiredObj = obj
			objFound = true
			counter = i
			break;
		} 
		if ((((obj.personalInfo.panNumber) == TAX_IDENTIFIER) && ((obj.personalInfo.passportNumber) != UNIQUE_IDENTIFIER))||((((obj.personalInfo.panNumber) != TAX_IDENTIFIER) && ((obj.personalInfo.passportNumber) == UNIQUE_IDENTIFIER)  ))){
			return nil, errors.New("Bad Request : Tax Identifier or Unique Identifier mapped for different Customer")
		}
	
	}
	
	if objFound {
		
		//Update CustomerTxObjects1 with new values from args 
		
		CustomerTxObjects[counter].clientId = args[0]
		CustomerTxObjects[counter].personalInfo.applicantNumber = args[1]
		CustomerTxObjects[counter].personalInfo.firstName = args[2]
		CustomerTxObjects[counter].personalInfo.middleName = args[3]
		CustomerTxObjects[counter].personalInfo.lastName = args[4]
		CustomerTxObjects[counter].personalInfo.dateOfBirth= args[5]
		CustomerTxObjects[counter].personalInfo.panNumber = args[6]
		CustomerTxObjects[counter].personalInfo.passportNumber = args[7]
		CustomerTxObjects[counter].personalInfo.residentStatus = args[8]
		CustomerTxObjects[counter].personalInfo.residencePlaceOwnership = args[9]
		CustomerTxObjects[counter].personalInfo.numberofDependents = args[10]
		CustomerTxObjects[counter].personalInfo.qualification   = args[11]
		CustomerTxObjects[counter].personalInfo.annualIncome = args[12]
		CustomerTxObjects[counter].personalInfo.gender = args[13]
		CustomerTxObjects[counter].personalInfo.maritalStatus   = args[14]
		CustomerTxObjects[counter].personalInfo.criminalRecordDetails = args[15]
		CustomerTxObjects[counter].personalInfo.financialStability = args[16]
		CustomerTxObjects[counter].personalInfo.creditScore   = args[17]
		
		
		CustomerTxObjects[counter].residenceAddress.addrLine1 = args[18]
		CustomerTxObjects[counter].residenceAddress.addrLine2 = args[19]
		CustomerTxObjects[counter].residenceAddress.city   = args[20]
		CustomerTxObjects[counter].residenceAddress.province = args[21]
		CustomerTxObjects[counter].residenceAddress.country = args[22]
		CustomerTxObjects[counter].residenceAddress.postalCode   = args[23]
		CustomerTxObjects[counter].residenceAddress.addressType = args[24]
		CustomerTxObjects[counter].residenceAddress.searchLocation = args[25]
		
		
		CustomerTxObjects[counter].permanentAddress.addrLine1 = args[26]
		CustomerTxObjects[counter].permanentAddress.addrLine2 = args[27]
		CustomerTxObjects[counter].permanentAddress.city   = args[28]
		CustomerTxObjects[counter].permanentAddress.province = args[29]
		CustomerTxObjects[counter].permanentAddress.country = args[30]
		CustomerTxObjects[counter].permanentAddress.postalCode   = args[31]
		CustomerTxObjects[counter].permanentAddress.addressType = args[32]
		CustomerTxObjects[counter].permanentAddress.searchLocation = args[33]		
		
		CustomerTxObjects[counter].officeAddress.addrLine1 = args[34]
		CustomerTxObjects[counter].officeAddress.addrLine2 = args[35]
		CustomerTxObjects[counter].officeAddress.city   = args[36]
		CustomerTxObjects[counter].officeAddress.province = args[37]
		CustomerTxObjects[counter].officeAddress.country = args[38]
		CustomerTxObjects[counter].officeAddress.postalCode   = args[39]
		CustomerTxObjects[counter].officeAddress.addressType = args[40]
		CustomerTxObjects[counter].officeAddress.searchLocation = args[41]
	
		CustomerTxObjects[counter].contactDetails.homeNumber = args[42]
		CustomerTxObjects[counter].contactDetails.officeNumber = args[43]
		CustomerTxObjects[counter].contactDetails.mobileNumber = args[44]
		CustomerTxObjects[counter].contactDetails.emailId = args[45]
		
		CustomerTxObjects[counter].employmentDetails.nameOfEmployer = args[46]
		CustomerTxObjects[counter].employmentDetails.designation = args[47]
		CustomerTxObjects[counter].employmentDetails.title = args[48]
		CustomerTxObjects[counter].employmentDetails.noOfYearsExperience = args[49]		
		
		CustomerTxObjects[counter].personalAssets.assetType = args[50]
		CustomerTxObjects[counter].personalAssets.assetName = args[51]
		CustomerTxObjects[counter].personalAssets.details = args[52]
		CustomerTxObjects[counter].personalAssets.valueOfAsset = args[53]
		CustomerTxObjects[counter].personalAssets.asOnDate = args[54]
		
		CustomerTxObjects[counter].bankAccountDetails.bankName = args[55]
		CustomerTxObjects[counter].bankAccountDetails.bankBranch = args[56]
		CustomerTxObjects[counter].bankAccountDetails.accountNo = args[57]
		CustomerTxObjects[counter].bankAccountDetails.swiftCode = args[58]		
		
		//Code for the Document Process	
		fmt.Printf("******** CUSTOMER_DOC:%s\n", args[4])
		var number_of_docs int
		number_of_docs = (len(args)-59)/2
		var CustomerDocObjects1 []KYCDocuments
		for i := 0; i < number_of_docs; i++ {
			var CustomerDocObj KYCDocuments
			fmt.Printf("******** CustomerDocObj[i].DOCUMENT_NAMEC:%d\n",i)
			fmt.Printf("******** CustomerDocObj[i].DOCUMENT_NAMEC:%d\n",number_of_docs)
			//CustomerDocObj[i] := CustomerDoc{DOCUMENT_NAME: args[27+(i*2)], DOCUMENT_STRING: args[27+(i*2)]}
			CustomerDocObj.documentName = args[59+(i*2)]
			//fmt.Printf("********pankaj CustomerDocObj[i].DOCUMENT_NAMEC:%s\n", CustomerDocObj[i].DOCUMENT_NAME)
			CustomerDocObj.base64String = args[60+(i*2)]
			CustomerDocObjects1 = append(CustomerDocObjects1,CustomerDocObj)
		}
		CustomerTxObjects[counter].kycDocuments = CustomerDocObjects1
		
		jsonAsBytes, _ := json.Marshal(CustomerTxObjects)
		fmt.Printf("======json print ====:%s\n", jsonAsBytes)

		err = stub.PutState(customerIndexTxStr, jsonAsBytes)
		if err != nil {
			return nil, err
		}
	    return nil, nil
	} else{
		if err != nil {
			return nil, err
		}
		return t.RegisterCustomer(stub ,args)
	
	}
}


func (t *CustomerChaincode)  RegisterCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var CustomerDataObj InsuranceClientInformation
	var CustomerDataList []InsuranceClientInformation
	var err error
   	fmt.Printf("******** CUSTOMER_DOC:%d\n", len(args))
	
	if len(args) < 8 {
		return nil, errors.New("Incorrect number of arguments. Need 8 arguments")
	}

	// Initialize the chaincode
	
		CustomerDataObj.clientId = args[0]
		CustomerDataObj.personalInfo.applicantNumber = args[1]
		CustomerDataObj.personalInfo.firstName = args[2]
		CustomerDataObj.personalInfo.middleName = args[3]
		CustomerDataObj.personalInfo.lastName = args[4]
		CustomerDataObj.personalInfo.dateOfBirth= args[5]
		CustomerDataObj.personalInfo.panNumber = args[6]
		CustomerDataObj.personalInfo.passportNumber = args[7]
		CustomerDataObj.personalInfo.residentStatus = args[8]
		CustomerDataObj.personalInfo.residencePlaceOwnership = args[9]
		CustomerDataObj.personalInfo.numberofDependents = args[10]
		CustomerDataObj.personalInfo.qualification   = args[11]
		CustomerDataObj.personalInfo.annualIncome = args[12]
		CustomerDataObj.personalInfo.gender = args[13]
		CustomerDataObj.personalInfo.maritalStatus   = args[14]
		CustomerDataObj.personalInfo.criminalRecordDetails = args[15]
		CustomerDataObj.personalInfo.financialStability = args[16]
		CustomerDataObj.personalInfo.creditScore   = args[17]	
		
		CustomerDataObj.residenceAddress.addrLine1 = args[18]
		CustomerDataObj.residenceAddress.addrLine2 = args[19]
		CustomerDataObj.residenceAddress.city   = args[20]
		CustomerDataObj.residenceAddress.province = args[21]
		CustomerDataObj.residenceAddress.country = args[22]
		CustomerDataObj.residenceAddress.postalCode   = args[23]
		CustomerDataObj.residenceAddress.addressType = args[24]
		CustomerDataObj.residenceAddress.searchLocation = args[25]
		
		
		CustomerDataObj.permanentAddress.addrLine1 = args[26]
		CustomerDataObj.permanentAddress.addrLine2 = args[27]
		CustomerDataObj.permanentAddress.city   = args[28]
		CustomerDataObj.permanentAddress.province = args[29]
		CustomerDataObj.permanentAddress.country = args[30]
		CustomerDataObj.permanentAddress.postalCode   = args[31]
		CustomerDataObj.permanentAddress.addressType = args[32]
		CustomerDataObj.permanentAddress.searchLocation = args[33]
		
		CustomerDataObj.officeAddress.addrLine1 = args[34]
		CustomerDataObj.officeAddress.addrLine2 = args[35]
		CustomerDataObj.officeAddress.city   = args[36]
		CustomerDataObj.officeAddress.province = args[37]
		CustomerDataObj.officeAddress.country = args[38]
		CustomerDataObj.officeAddress.postalCode   = args[39]
		CustomerDataObj.officeAddress.addressType = args[40]
		CustomerDataObj.officeAddress.searchLocation = args[41]
	    
		CustomerDataObj.contactDetails.homeNumber = args[42]
		CustomerDataObj.contactDetails.officeNumber = args[43]
		CustomerDataObj.contactDetails.mobileNumber = args[44]
		CustomerDataObj.contactDetails.emailId = args[45]
		
		CustomerDataObj.employmentDetails.nameOfEmployer = args[46]
		CustomerDataObj.employmentDetails.designation = args[47]
		CustomerDataObj.employmentDetails.title = args[48]
		CustomerDataObj.employmentDetails.noOfYearsExperience = args[49]		
		
		CustomerDataObj.personalAssets.assetType = args[50]
		CustomerDataObj.personalAssets.assetName = args[51]
		CustomerDataObj.personalAssets.details = args[52]
		CustomerDataObj.personalAssets.valueOfAsset = args[53]
		CustomerDataObj.personalAssets.asOnDate = args[54]
		
		CustomerDataObj.bankAccountDetails.bankName = args[55]
		CustomerDataObj.bankAccountDetails.bankBranch = args[56]
		CustomerDataObj.bankAccountDetails.accountNo = args[57]
		CustomerDataObj.bankAccountDetails.swiftCode = args[58]	
	
	//Code for the Document Process	
	fmt.Printf("********RegisterCustomer CUSTOMER_DOC Proceesing :%s\n", args[4])
	var number_of_docs int
	number_of_docs = (len(args)-59)/2
	var CustomerDocObjects1 []KYCDocuments
	for i := 0; i < number_of_docs; i++ {
		var CustomerDocObj KYCDocuments
		fmt.Printf("******** CustomerDocObj[i].DOCUMENT_NAMEC:%d\n",i)
		fmt.Printf("******** CustomerDocObj[i].DOCUMENT_NAMEC:%d\n",number_of_docs)
		//CustomerDocObj[i] := CustomerDoc{DOCUMENT_NAME: args[27+(i*2)], DOCUMENT_STRING: args[27+(i*2)]}
		CustomerDocObj.documentName = args[59+(i*2)]
		//fmt.Printf("********pankaj CustomerDocObj[i].DOCUMENT_NAMEC:%s\n", CustomerDocObj[i].DOCUMENT_NAME)
		CustomerDocObj.base64String = args[60+(i*2)]
		CustomerDocObjects1 = append(CustomerDocObjects1,CustomerDocObj)
	}
	
	CustomerDataObj.kycDocuments = CustomerDocObjects1
	
	customerTxsAsBytes, err := stub.GetState(customerIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get customer transactions")
	}
	json.Unmarshal(customerTxsAsBytes, &CustomerDataList)

	CustomerDataList = append(CustomerDataList, CustomerDataObj)
	jsonAsBytes, _ := json.Marshal(CustomerDataList)

	err = stub.PutState(customerIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, err
}

// Query callback representing the query of a chaincode
func (t *CustomerChaincode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {

   	var CUSTOMER_FIRST_NAME  string 
	var CUSTOMER_MIDDLE_NAME string 
	var CUSTOMER_LAST_NAME  string 
	var CUSTOMER_DOB string
	var TAX_IDENTIFIER string // Entities
	var UNIQUE_IDENTIFIER string
	var err error
	var resAsBytes []byte

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6 parameters to query")
	}
	CUSTOMER_FIRST_NAME= args[2]
	CUSTOMER_MIDDLE_NAME= args[3]
	CUSTOMER_LAST_NAME = args[4]
	CUSTOMER_DOB = args[5]
	TAX_IDENTIFIER = args[6]
	UNIQUE_IDENTIFIER = args[7]
	
	resAsBytes, err = t.GetCustomerDetails(stub,CUSTOMER_FIRST_NAME,CUSTOMER_MIDDLE_NAME,CUSTOMER_LAST_NAME,CUSTOMER_DOB, TAX_IDENTIFIER, UNIQUE_IDENTIFIER)
        InitLogs(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
        Trace.Println("I have something standard to say")
        Info.Println("Special Information")
        Warning.Println("There is something you need to know about")
        Error.Println("Something has failed")
	
	file, err := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil {
		Error.Println("Failed to open log file : ", err)
         }
        var MyFile  *log.Logger
        MyFile = log.New(file,
        "PREFIX: ",
        log.Ldate|log.Ltime|log.Lshortfile)
	MyFile.Println("Special Information in Myfile")
	fmt.Printf("Query Response:%s\n", resAsBytes)

	if err != nil {
		return nil, err
	}

	return resAsBytes, nil
}

func InitLogs(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func (t *CustomerChaincode)  GetCustomerDetails(stub shim.ChaincodeStubInterface,CUSTOMER_FIRST_NAME string,CUSTOMER_MIDDLE_NAME string,CUSTOMER_LAST_NAME string,CUSTOMER_DOB string, TAX_IDENTIFIER string, UNIQUE_IDENTIFIER string) ([]byte, error) {

	//var requiredObj CustomerData
	var objFound bool
	CustomerTxsAsBytes, err := stub.GetState(customerIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get Customer Records")
	}
	var CustomerTxObjects []InsuranceClientInformation
	var CustomerTxObjects1 []InsuranceClientInformation
	json.Unmarshal(CustomerTxsAsBytes, &CustomerTxObjects)
	length := len(CustomerTxObjects)
	fmt.Printf("Output from chaincode: %s\n", CustomerTxsAsBytes)

	if CUSTOMER_FIRST_NAME == "" && CUSTOMER_MIDDLE_NAME == "" && CUSTOMER_LAST_NAME == "" && CUSTOMER_DOB == "" && TAX_IDENTIFIER == "" && UNIQUE_IDENTIFIER == ""{
		res, err := json.Marshal(CustomerTxObjects)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}

	objFound = false
	// iterate
	for i := 0; i < length; i++ {
		obj := CustomerTxObjects[i]
		//if ((customer_id == obj.CUSTOMER_ID) && (customer_name == obj.CUSTOMER_NAME) && (customer_dob == obj.CUSTOMER_DOB)) 
		
		if (((CUSTOMER_FIRST_NAME != "" && (obj.personalInfo.firstName == CUSTOMER_FIRST_NAME)) || CUSTOMER_FIRST_NAME == "" ) && 
		((CUSTOMER_MIDDLE_NAME != "" && (obj.personalInfo.middleName == CUSTOMER_MIDDLE_NAME)) || CUSTOMER_MIDDLE_NAME == "" ) && 
		((CUSTOMER_LAST_NAME != "" && (obj.personalInfo.lastName == CUSTOMER_LAST_NAME)) || CUSTOMER_LAST_NAME == "" ) &&
		((CUSTOMER_DOB != "" && (obj.personalInfo.dateOfBirth == CUSTOMER_DOB)) || CUSTOMER_DOB == "" ) &&
		((TAX_IDENTIFIER != "" && (obj.personalInfo.panNumber == TAX_IDENTIFIER)) || TAX_IDENTIFIER == "" ) &&
		((UNIQUE_IDENTIFIER != "" && (obj.personalInfo.passportNumber == UNIQUE_IDENTIFIER)) || UNIQUE_IDENTIFIER == "" )){
				
			fmt.Printf("condition matched\n")
			
			CustomerTxObjects1 = append(CustomerTxObjects1,obj)
			//requiredObj = obj
			objFound = true
				
			
		} else {
			fmt.Printf("no condition matched\n")
		}
	}

	if objFound {
		res, err := json.Marshal(CustomerTxObjects1)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	} else {
		res, err := json.Marshal("No Data found")
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}
}




func main() {
	err := shim.Start(new(CustomerChaincode))
	if err != nil {
		fmt.Printf("Error starting Customer Simple chaincode: %s", err)
	}
}
