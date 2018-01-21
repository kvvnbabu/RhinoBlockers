package main


import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("distributionsmartcontract")

// Define the Smart Contract structure
type SaplingDistributionSmartContract struct {
	DistributionId  string `json:"distributionid"`
	BeneficiaryId  string `json:"beneficiaryid"`
	Land  string `json:"land"`
	EligibleSaplings string `json:"eligiblesaplings"`
	ActiveSaplings  string `json:"activesaplings"`
	RfIds string `json:"rfids"`
}

func (s *SaplingDistributionSmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SaplingDistributionSmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Info("########### distributionsmartcontract Invoke ###########")
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "createDistribution" {
		return s.createDistribution(APIstub, args)
	}else if function == "getDistribution"{
		return s.queryDistribution(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SaplingDistributionSmartContract) createDistribution(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var distribution = SaplingDistributionSmartContract {
			DistributionId		:args[0],
			BeneficiaryId  		:args[0],
			Land  				:args[1],
			EligibleSaplings 	:args[2],
			ActiveSaplings  	:args[3],
			RfIds 				:args[4],
	}

	distributionAsBytes, _ := json.Marshal(distribution)
	APIstub.PutState(args[0], distributionAsBytes)

	return shim.Success(nil)
}

func (s *SaplingDistributionSmartContract) queryDistribution(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var err error
	var distributionId = args[0]

	Avalbytes, err := APIstub.GetState(distributionId)

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + distributionId + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + distributionId + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Id\":\"" + distributionId + "\",\"Data\":\"" + string(Avalbytes) + "\"}"
	logger.Infof("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}


func main() {
	err := shim.Start(new(SaplingDistributionSmartContract))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
