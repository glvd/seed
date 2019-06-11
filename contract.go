package seed

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	jsoniter "github.com/json-iterator/go"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
	"math/big"
	"strconv"
	"strings"
)

// GatwayAddress ...
const defaultGatewayAddress = "https://ropsten.infura.io/QVsqBu3yopMu2svcHqRj"

var eth *ETH

func bangumi() ProcessorFunc {
	return func(eth *ETH) interface{} {
		tk, err := NewBangumiData(common.HexToAddress(eth.contract), eth.conn)
		if err != nil {
			return nil
		}
		return tk
	}
}

func dhash() ProcessorFunc {
	return func(eth *ETH) interface{} {
		tk, err := NewDhash(common.HexToAddress(eth.contract), eth.conn)
		if err != nil {
			return nil
		}
		return tk
	}
}

// InitGlobalETH ...
func InitGlobalETH(key, contract string, processorFunc ...ProcessorFunc) (err error) {
	eth = NewETH(key, contract)
	eth.conn, err = ethclient.Dial(defaultGatewayAddress)
	if err != nil {
		return err
	}
	eth.processor = append(eth.processor, processorFunc...)
	eth.processor = append(eth.processor, bangumi(), dhash())
	return nil
}

// ProcessorFunc ...
type ProcessorFunc func(eth *ETH) interface{}

//ETH ...
type ETH struct {
	conn      *ethclient.Client
	key       string
	contract  string
	gateway   string
	processor []ProcessorFunc
}

// NewETH ...
func NewETH(key, contract string) *ETH {
	// Create an IPC based RPC connection to a remote node and instantiate a contract binding
	return &ETH{
		conn:      nil,
		key:       key,
		gateway:   defaultGatewayAddress,
		contract:  contract,
		processor: nil,
	}
}

// PrivateKey ...
func (eth *ETH) PrivateKey() (key *ecdsa.PrivateKey) {
	privateKey, err := crypto.HexToECDSA(eth.key)
	if err != nil {
		return nil
	}
	return privateKey
}

// RegisterContract ...
func (eth *ETH) RegisterContract(fn func(eth *ETH) interface{}) {
	eth.processor = append(eth.processor, fn)
}

// ProcContract ...
func (eth *ETH) ProcContract(fn func(v interface{}) (bool, error)) (e error) {
	for _, proc := range eth.processor {
		if b, e := fn(proc(eth)); b {
			return e
		}
	}
	return nil
}

// InfoInput ...
func InfoInput(video *model.Video) (e error) {
	e = infoInput(eth, video)
	if e != nil {
		return e
	}
	return nil
}

// CheckNameExists ...
func CheckNameExists(ban string, idx ...int) (e error) {
	idxStr := ""
	nb := ""
	if idx == nil {
		idx = append(idx, 1)
	}

	for _, i := range idx {
		idxStr = strconv.FormatInt(int64(i), 10)
		nb = strings.ToUpper(ban + "@" + idxStr)
		e = CheckExist(nb)
		if e != nil {
			return e
		}
	}
	return nil
}

// GetLastVersionCode ...
func GetLastVersionCode() (code *big.Int, e error) {
	err := eth.ProcContract(func(v interface{}) (b bool, e error) {
		data, b := v.(*Dhash)
		if !b {
			return false, nil
		}
		code, _, e = data.GetLatest(&bind.CallOpts{Pending: true})
		return true, e
	})
	if err != nil {
		return nil, err
	}
	return
}

// GetLastVersionHash ...
func GetLastVersionHash() (ver, hash string, e error) {
	err := eth.ProcContract(func(v interface{}) (b bool, e error) {
		data, b := v.(*Dhash)
		if !b {
			return false, nil
		}
		_, ver, e = data.GetLatest(&bind.CallOpts{Pending: true})
		if e != nil {
			return false, e
		}
		hash, e = data.GetHash(&bind.CallOpts{Pending: true}, ver)
		if e != nil {
			return true, e
		}
		return true, e
	})
	if err != nil {
		return "", "", err
	}
	return
}

// UpdateHotList ...
func UpdateHotList(list ...string) (e error) {
	err := eth.ProcContract(func(v interface{}) (b bool, e error) {
		data, b := v.(*Dhash)
		if !b {
			return false, nil
		}
		opt := bind.NewKeyedTransactor(eth.PrivateKey())
		bytes, e := jsoniter.Marshal(list)
		if e != nil {
			return false, e
		}
		transaction, e := data.UpdateHotList(opt, string(bytes))
		if e != nil {
			return true, e
		}
		ctx := context.Background()
		receipt, err := bind.WaitMined(ctx, eth.conn, transaction)
		if err != nil {
			return true, err
		}
		log.Debugf("receipt is :%x\n", string(receipt.TxHash[:]))
		return true, nil
	})
	if err != nil {
		return err
	}
	return
}

// GetHostList ...
func GetHostList() (list []string) {

	err := eth.ProcContract(func(v interface{}) (b bool, e error) {
		data, b := v.(*Dhash)
		if !b {
			return false, nil
		}
		ll, e := data.GetHotList(&bind.CallOpts{Pending: true})
		if e != nil {
			return true, e
		}
		e = jsoniter.Unmarshal([]byte(ll), &list)
		if e != nil {
			return false, e
		}
		return true, nil
	})
	if err != nil {
		return nil
	}
	return list
}

// CheckExist ...
func CheckExist(ban string) (e error) {
	return eth.ProcContract(func(v interface{}) (b bool, e error) {
		data, b := v.(*BangumiData)
		if !b {
			return false, nil
		}
		hash, e := data.QueryHash(&bind.CallOpts{Pending: true}, "ban")
		log.With("size", len(hash), "hash", hash, "name", ban).Info("checked")
		if hash == "" || hash == "," || len(hash) != 46 {
			return true, xerrors.New(ban + " hash is not found!")
		}
		return true, e
	})

}

// Close ...
func (eth *ETH) Close() {
	if eth.conn == nil {
		return
	}
	eth.conn.Close()
}

// ContractProcessor ...
//type ContractProcessor interface {
//	ContractProc(eth *ETH) error
//}

func singleInput(eth *ETH, video *model.Video) (e error) {
	return eth.ProcContract(func(v interface{}) (b bool, e error) {
		data, b := (v).(*BangumiData)
		if !b {
			return false, nil
		}
		opt := bind.NewKeyedTransactor(eth.PrivateKey())
		name := video.Bangumi
		//list := video.VideoGroupList[0]
		//objMax := len(list.Object)
		//objMaxStr := strconv.FormatInt(int64(objMax), 10)
		//for i := 0; i < objMax; i++ {
		//	idxStr := strconv.FormatInt(int64(i+1), 10)
		upperName := strings.ToUpper(name + "@" + video.Episode)
		e = CheckExist(upperName)
		if e == nil {
			return
		}
		transaction, err := data.InfoInput(opt,
			strings.ToUpper(upperName),
			video.PosterHash,
			video.Role[0],
			video.M3U8Hash,
			video.Alias[0],
			video.Sharpness,
			video.Episode,
			video.TotalEpisode,
			video.Season,
			video.Format,
			"",
			"")
		if err != nil {
			return true, err
		}
		ctx := context.Background()
		receipt, err := bind.WaitMined(ctx, eth.conn, transaction)
		if err != nil {
			return true, err
		}
		log.Info(name + "@" + video.Episode + " success")
		log.Debugf("receipt is :%x\n", string(receipt.TxHash[:]))
		return true, nil
	})

}

func multipleInput(eth *ETH, video *model.Video) (e error) {
	return eth.ProcContract(func(v interface{}) (b bool, e error) {
		data, b := (v).(*BangumiData)
		if !b {
			return false, nil
		}
		opt := bind.NewKeyedTransactor(eth.PrivateKey())
		name := video.Bangumi

		upperName := strings.ToUpper(name + "@" + video.Episode)
		e = CheckExist(upperName)
		if e == nil {
			return
		}
		transaction, err := data.InfoInput(opt,
			strings.ToUpper(upperName),
			video.PosterHash,
			video.Role[0],
			video.M3U8Hash,
			video.Alias[0],
			video.Sharpness,
			video.Episode,
			video.TotalEpisode,
			video.Season,
			video.Format,
			"",
			"")
		if err != nil {
			return true, err
		}
		ctx := context.Background()
		receipt, err := bind.WaitMined(ctx, eth.conn, transaction)
		if err != nil {
			//log.Fatalf("tx mining error:%v\n", err)
			return true, err
		}
		log.Info(name + "@" + video.Episode + " success")
		log.Debugf("receipt is :%x\n", string(receipt.TxHash[:]))

		return true, nil
	})

}

// InfoInput ...
func infoInput(eth *ETH, video *model.Video) (e error) {
	//if video == nil || video.VideoGroupList == nil {
	//	return
	//}
	//
	//vgMax := len(video.VideoGroupList)
	fn := singleInput
	//if vgMax > 1 {
	//	fn = multipleInput
	//}
	return fn(eth, video)
}

// UpdateApp ...
func UpdateApp(version string, hash string) (e error) {
	return updateAppHash(version, hash)
}

func updateAppHash(version string, hash string) (e error) {
	return eth.ProcContract(func(v interface{}) (b bool, e error) {
		data, b := v.(*Dhash)
		if !b {
			return false, nil
		}

		code, e := GetLastVersionCode()
		if e != nil {
			return true, e
		}
		one := big.NewInt(1)
		code = code.Add(code, one)
		key := eth.PrivateKey()
		opt := bind.NewKeyedTransactor(key)
		transaction, err := data.UpdateVersion(opt, version, hash, code)
		if err != nil {
			return true, e
		}

		ctx := context.Background()
		receipt, err := bind.WaitMined(ctx, eth.conn, transaction)
		if err != nil {
			return true, e
		}
		log.Info(version + "@" + hash + " success")
		log.Debugf("receipt is :%x\n", string(receipt.TxHash[:]))
		return true, nil
	})
}
