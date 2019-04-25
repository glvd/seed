package seed

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
	"strconv"
	"strings"
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
	e = infoInput(eth, video)
	if e != nil {
		return e
	}
	return nil
}

// CheckExist ...
func (eth *ETH) CheckExist(ban string) (e error) {
	tk, e := eth.ConnectToken()
	if e != nil {
		return e
	}
	hash, e := tk.QueryHash(&bind.CallOpts{Pending: true}, ban)
	if e != nil {
		return e
	}
	logrus.Info(ban+" checking hash:", hash)
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
	logrus.Debug("key is: ", key)
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
	tk, err := NewBangumiData(common.HexToAddress("0xb5eb6bf5eab725e9285d0d27201603ecf31a1d37"), eth.conn)
	if err != nil {
		logrus.Fatalf("Failed to instantiate a Token contract: %v", err)
		return &BangumiData{}, nil
	}
	return tk, nil
}

func singleInput(eth *ETH, video *model.Video) (e error) {
	tk, e := eth.ConnectToken()
	if e != nil {
		return e
	}
	privateKey, err := crypto.HexToECDSA(eth.key)
	if err != nil {
		logrus.Fatal(err)
		return err
	}

	opt := bind.NewKeyedTransactor(privateKey)
	name := video.Bangumi
	list := video.VideoGroupList[0]
	objMax := len(list.Object)
	objMaxStr := strconv.FormatInt(int64(objMax), 10)
	for i := 0; i < objMax; i++ {
		idxStr := strconv.FormatInt(int64(i+1), 10)
		upperName := strings.ToUpper(name + "@" + idxStr)
		e = eth.CheckExist(upperName)
		if e == nil {
			continue
		}
		transaction, err := tk.InfoInput(opt,
			strings.ToUpper(upperName),
			video.Poster,
			video.Role[0],
			list.Object[i].Link.Hash,
			video.Alias[0],
			list.Sharpness,
			idxStr,
			objMaxStr,
			list.Season,
			list.Output,
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
		logrus.Info(name + "@" + idxStr + " success")
		logrus.Debugf("receipt is :%x\n", string(receipt.TxHash[:]))
	}
	return nil
}

func multipleInput(eth *ETH, video *model.Video) (e error) {
	tk, e := eth.ConnectToken()
	if e != nil {
		return e
	}
	privateKey, err := crypto.HexToECDSA(eth.key)
	if err != nil {
		logrus.Fatal(err)
		return err
	}

	opt := bind.NewKeyedTransactor(privateKey)
	name := video.Bangumi

	for _, list := range video.VideoGroupList {
		objMax := len(list.Object)
		if objMax > 1 {
			return xerrors.New("multiple group list can only use object with 1 ")
		}

		for i := 0; i < objMax; i++ {
			idxStr := strconv.FormatInt(int64(i+1), 10)
			upperName := strings.ToUpper(name + "@" + idxStr)
			e = eth.CheckExist(upperName)
			if e == nil {
				continue
			}
			transaction, err := tk.InfoInput(opt,
				strings.ToUpper(upperName),
				video.Poster,
				video.Role[0],
				list.Object[i].Link.Hash,
				video.Alias[0],
				list.Sharpness,
				list.Episode,
				list.TotalEpisode,
				list.Season,
				list.Output,
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
			logrus.Info(name + "@" + idxStr + " success")
			logrus.Debugf("receipt is :%x\n", string(receipt.TxHash[:]))
		}
	}
	return nil
}

// InfoInput ...
func infoInput(eth *ETH, video *model.Video) (e error) {
	if video == nil || video.VideoGroupList == nil {
		return
	}

	vgMax := len(video.VideoGroupList)
	fn := singleInput
	if vgMax > 1 {
		fn = multipleInput
	}
	if err := fn(eth, video); err != nil {
		return err
	}

	return nil
}
