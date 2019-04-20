package seed

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

//ETH ...
type ETH struct {
	key             string
	ContractAddress string
	DialAddress     string
}

func NewETH(key string) *ETH {
	return &ETH{
		key: key,
	}
}

// ConnectToken ...
func ConnectToken() *BangumiData {
	// Create an IPC based RPC connection to a remote node and instantiate a contract binding
	conn, err := ethclient.Dial("https://ropsten.infura.io/QVsqBu3yopMu2svcHqRj")
	if err != nil {
		logrus.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer conn.Close()

	token, err := NewBangumiData(common.HexToAddress("0xb5eb6bf5eab725e9285d0d27201603ecf31a1d37"), conn)
	if err != nil {
		logrus.Fatalf("Failed to instantiate a Token contract: %v", err)
	}
	logrus.Info(token)

	return token
}

// InfoInput ...
func InfoInput() {
	// Create an IPC based RPC connection to a remote node and instantiate a contract binding
	conn, err := ethclient.Dial("https://ropsten.infura.io/QVsqBu3yopMu2svcHqRj")
	if err != nil {
		logrus.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer conn.Close()

	token, err := NewBangumiData(common.HexToAddress("0xb5eb6bf5eab725e9285d0d27201603ecf31a1d37"), conn)
	if err != nil {
		logrus.Fatalf("Failed to instantiate a Token contract: %v", err)
	}
	logrus.Info(token)

	bytes := "key"
	privateKey, err := crypto.HexToECDSA(bytes)
	if err != nil {
		logrus.Fatal(err)
	}

	opt := bind.NewKeyedTransactor(privateKey)
	logrus.Info(opt)
	transaction, err := token.InfoInput(opt,
		"test",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"")
	if err != nil {
		return
	}
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, conn, transaction)
	if err != nil {
		logrus.Fatalf("tx mining error:%v\n", err)
	}
	fmt.Printf("tx is :%+v\n", transaction)
	fmt.Printf("receipt is :%x\n", string(receipt.TxHash[:]))

}
