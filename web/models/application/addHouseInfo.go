//在这里调用链码，读写区块链

package application

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	//"github.com/pkg/errors"
)

func (this *Application) AddHouseItem(args []string) (string, error) {
	len1 := len(args)
	fmt.Printf("len of args:%f\n", len1)

	request := channel.Request{ChaincodeID: this.FabricSetup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(args[5]), []byte(args[6]), []byte(args[7])}}
	//进行数据的上传
	response, err := this.FabricSetup.Client.Execute(request)
	if err != nil {
		// 资产转移失败
		return "", err
	}

	return string(response.TransactionID), nil
}