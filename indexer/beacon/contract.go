package beacon

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"math/big"
	"strings"
)

func parseFunction(txDataHex string, log logrus.FieldLogger) (functionName *string, value *big.Int, to *common.Address) {
	log.Infof("try to parse: %v", txDataHex)
	contractABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		log.Warnf("Failed to parse ABI: %v", err)
		return nil, nil, nil
	}

	txData, err := hex.DecodeString(txDataHex)
	if err != nil {
		log.Warnf("Failed to decode hex: %v", err)
		return nil, nil, nil
	}

	if len(txData) < 4 {
		log.Warnf("Transaction data too short")
		return nil, nil, nil
	}

	methodSig := txData[:4]

	method, err := contractABI.MethodById(methodSig)
	if err != nil {
		log.Warnf("Failed to find method: %v", err)
		return nil, nil, nil
	}

	args := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(args, txData[4:]); err != nil {
		log.Warnf("Failed to unpack arguments: %v", err)
		return nil, nil, nil
	}

	fmt.Printf("Method called: %s\n", method.Name)
	fmt.Printf("Signature: 0x%x\n", methodSig)

	if method.Name == "transfer" || method.Name == "mint" {
		toAddress := args["to"].(common.Address)
		value := args["amount"].(*big.Int)
		return &method.Name, value, &toAddress
	}

	for name, value := range args {
		switch v := value.(type) {
		case common.Address:
			fmt.Printf("%s: %s (address)\n", name, v.Hex())
		case *big.Int:
			fmt.Printf("%s: %s (uint256)\n", name, v.String())
		case string:
			fmt.Printf("%s: %s (string)\n", name, v)
		case bool:
			fmt.Printf("%s: %t (bool)\n", name, v)
		default:
			fmt.Printf("%s: %v (type: %T)\n", name, v, v)
		}
	}
	return &method.Name, nil, nil
}

const erc20ABI = `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`
