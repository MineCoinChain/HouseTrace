package main

import (
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type HouseChainCode struct {
	// 房屋溯源链代码
}

func (hcc *HouseChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (hcc *HouseChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fn, parameters := stub.GetFunctionAndParameters()

	//添加房源
	if fn == "setHouseInfo" {
		return hcc.setHouseInfo(stub, parameters)
	} else if fn == "getHouseInfo" {
		return hcc.getHouseInfo(stub, parameters)
	}

	//其他的，TODO

	return shim.Success(nil)
}

func main() {
	shim.Start(new(HouseChainCode))
}

func (hcc *HouseChainCode) setHouseInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var houseInfo RentingHouseInfo
	if len(args) != 7 {
		return shim.Error("输入的参数不足！")
	}
	houseInfo.RentingID = args[0]
	if houseInfo.RentingID == "" {
		return shim.Error("出租id不应为空")
	}

	houseInfo.HouseID = args[1]
	houseInfo.HouseOwner = args[2]
	houseInfo.RegDate = args[3]
	houseInfo.HouseArea = args[4]
	houseInfo.HouseUsed = args[5]
	houseInfo.IsMortgage = args[6]

	fmt.Printf("houseInfo：", houseInfo)

	houseInfoBytes, err := json.Marshal(houseInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(houseInfo.RentingID, houseInfoBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("ok"))
}

func (hcc *HouseChainCode) getHouseInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("输入参数只需1个rentingID")
	}

	rentingId := args[0]

	iterator, err := stub.GetHistoryForKey(rentingId)

	defer iterator.Close()
	var houseInfo HouseInfo

	// 使用迭代器遍历查询结果集
	for iterator.HasNext() {
		var rentingInfo RentingHouseInfo
		// 取一条结果
		response, err := iterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		err = json.Unmarshal(response.Value, &rentingInfo)
		if err != nil {
			return shim.Error(err.Error())
		}

		if rentingInfo.HouseOwner != "" {
			houseInfo = rentingInfo.HouseInfo
			continue
		}
	}

	jsonAsBytes, err := json.Marshal(houseInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	// 返回数据查询结果
	return shim.Success(jsonAsBytes)
}
