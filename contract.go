package seed

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
	"strconv"
	"syscall"
)

//ETH ...
type ETH struct {
	conn            *ethclient.Client
	key             string
	ContractAddress string
	DialAddress     string
}

func getSeedKey() string {
	key, _ := syscall.Getenv("SEED_KEY")
	return key
}

// NewETH ...
func NewETH(key string) *ETH {
	// Create an IPC based RPC connection to a remote node and instantiate a contract binding
	conn, err := ethclient.Dial("https://ropsten.infura.io/QVsqBu3yopMu2svcHqRj")
	if err != nil {
		logrus.Fatalf("Failed to connect to the Ethereum client: %v", err)
		return nil
	}
	return &ETH{
		conn: conn,
		key:  key,
	}
}

// InfoInput ...
func (eth *ETH) InfoInput(video *model.Video) (e error) {
	for idx := range video.VideoGroupList {
		e = infoInput(eth, video, idx)
		if e != nil {
			return e
		}
	}
	return nil
}

// CheckExist ...
func (eth *ETH) CheckExist(ban string) (e error) {
	token, e := eth.ConnectToken()
	if e != nil {
		return e
	}
	hash, e := token.QueryHash(&bind.CallOpts{Pending: true}, ban)
	if e != nil {
		return e
	}
	logrus.Println("hash:", hash)
	if hash == "" {
		return xerrors.New(ban + " hash is not found!")
	}
	return nil
}

// Close ...
func (eth *ETH) Close() {
	if eth.conn == nil {
		return
	}
	eth.conn.Close()
}

// Contract ...
func Contract(key string) (e error) {
	var videos = new([]*model.Video)

	if e = model.DB().Find(videos); e != nil {
		return e
	}
	if key == "" {
		key = getSeedKey()
	}
	eth := NewETH(key)
	if eth == nil {
		return xerrors.New("nil eth")
	}
	for _, v := range *videos {

		e = eth.InfoInput(v)
		if e != nil {
			logrus.Error("contract err:", v.Bangumi, e)
			return e
		}
	}
	return
}

// ConnectToken ...
func (eth *ETH) ConnectToken() (*BangumiData, error) {
	token, err := NewBangumiData(common.HexToAddress("0xb5eb6bf5eab725e9285d0d27201603ecf31a1d37"), eth.conn)
	if err != nil {
		logrus.Fatalf("Failed to instantiate a Token contract: %v", err)
		return &BangumiData{}, nil
	}
	logrus.Info(token)

	return token, nil
}

// InfoInput ...
func infoInput(eth *ETH, video *model.Video, index int) (e error) {
	token, e := eth.ConnectToken()
	if e != nil {
		return e
	}
	logrus.Info(token)
	privateKey, err := crypto.HexToECDSA(eth.key)
	if err != nil {
		logrus.Fatal(err)
		return err
	}

	opt := bind.NewKeyedTransactor(privateKey)
	name := video.Bangumi
	max := len(video.VideoGroupList[index].Object)
	maxv := strconv.FormatInt(int64(max), 10)
	for i := 0; i < max; i++ {
		idxv := strconv.FormatInt(int64(i+1), 10)
		e = eth.CheckExist(name + "@" + idxv)
		if e == nil {
			continue
		}
		transaction, err := token.InfoInput(opt,
			name+"@"+idxv,
			video.Poster,
			video.Role[0],
			video.VideoGroupList[index].Object[i].Link.Hash,
			video.Alias[0],
			video.VideoGroupList[index].Sharpness,
			idxv,
			maxv,
			video.VideoGroupList[index].Season,
			video.VideoGroupList[index].Output,
			"",
			"")
		if err != nil {
			return err
		}
		ctx := context.Background()
		receipt, err := bind.WaitMined(ctx, eth.conn, transaction)
		if err != nil {
			//logrus.Fatalf("tx mining error:%v\n", err)
			return err
		}

		//fmt.Printf("tx is :%+v\n", transaction)
		fmt.Printf("receipt is :%x\n", string(receipt.TxHash[:]))
	}

	return nil
}
