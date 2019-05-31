package seed

import (
	"context"
	"math/big"
	"strconv"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
	"gopkg.in/urfave/cli.v2"
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
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
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

// CheckNameExists ...
func (eth *ETH) CheckNameExists(ban string, idx ...int) (e error) {
	idxStr := ""
	nb := ""
	if idx == nil {
		idx = append(idx, 1)
	}

	for _, i := range idx {
		idxStr = strconv.FormatInt(int64(i), 10)
		nb = strings.ToUpper(ban + "@" + idxStr)
		e = eth.CheckExist(nb)
		if e != nil {
			return e
		}
	}
	return nil
}

// GetLastVersionCode ...
func (eth *ETH) GetLastVersionCode() (code *big.Int, e error) {
	tk, e := eth.ConnectDHash()
	if e != nil {
		return nil, e
	}
	code, _, e = tk.GetLatest(&bind.CallOpts{Pending: true})
	return
}

// GetLastHash ...
func (eth *ETH) GetLastHash() (ver, hash string, e error) {
	tk, e := eth.ConnectDHash()
	if e != nil {
		return "", "", e
	}
	_, version, e := tk.GetLatest(&bind.CallOpts{Pending: true})
	if e != nil {
		return "", "", e
	}
	s, e := tk.GetHash(&bind.CallOpts{Pending: true}, version)
	if e != nil {
		return "", "", e
	}
	return version, s, nil
}

// CheckExist ...
func (eth *ETH) CheckExist(ban string) (e error) {
	tk, e := eth.ConnectBangumi()
	if e != nil {
		return e
	}
	hash, e := tk.QueryHash(&bind.CallOpts{Pending: true}, ban)
	if e != nil {
		return e
	}
	log.With("size", len(hash), "hash", hash, "name", ban).Info("checked")
	if hash == "" || hash == "," || len(hash) != 46 {
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
func Contract(key string, contract string) (e error) {
	var videos = new([]*model.Video)

	if e = model.DB().Find(videos); e != nil {
		return e
	}
	if key == "" {
		key = getSeedKey()
	}
	log.Debug("key is: ", key)
	eth := NewETH(key)
	eth.ContractAddress = contract
	if eth == nil {
		return xerrors.New("nil eth")
	}
	for _, v := range *videos {

		e = eth.InfoInput(v)
		if e != nil {
			log.Error("contract err:", v.Bangumi, e)
			return e
		}
	}
	return
}

// CmdContract ...
func CmdContract(app *cli.App) *cli.Command {
	flags := append(app.Flags,
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"a"},
			Value:   "",
			Usage:   "contract address",
		},
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Value:   "video",
			Usage:   "contract process type",
		},
		&cli.StringFlag{
			Name:  "ban",
			Usage: "ban no to check",
		},
		&cli.StringFlag{
			Name:    "appver",
			Value:   "v0.0.1",
			Aliases: []string{"av"},
			Usage:   "set the application version",
		},
		&cli.StringFlag{
			Name:  "hash",
			Usage: "set the app ipfs hash",
		},
	)
	return &cli.Command{
		Name:    "contract",
		Aliases: []string{"C"},
		Usage:   "contract share the video info to contract.",
		Action: func(context *cli.Context) error {
			log.Info("contract call")
			key := ""
			if context.NArg() > 0 {
				key = context.Args().Get(0)
			}
			address := context.String("a")
			if address == "" {
				panic("address must set use -address,-a")
			}
			version := context.String("av")
			path := context.String("p")
			hash := context.String("hash")
			eth := NewETH(key)
			eth.ContractAddress = address
			switch context.String("t") {
			case "video":
				log.Info("video:", context.String("ban"))
				return Contract(key, address)
			case "check":
				ban := context.String("ban")
				if ban == "" {
					return nil
				}
				log.Info("check:", context.String("ban"))
				e := eth.CheckNameExists(context.String("ban"), 1, 2, 3, 4, 5, 6, 7, 8)
				if e != nil {
					log.Error(e)
					return e
				}
			case "hot":

			case "app":
				if path != "" {
					if err := eth.UpdateAppWithPath(version, path); err != nil {
						log.Error(err)
						return err
					}
					return nil
				} else if hash != "" {
					e := eth.UpdateApp(version, hash)
					if e != nil {
						log.Error(e)
						return e
					}
				} else {
					ver, lastHash, e := eth.GetLastHash()
					if e != nil {
						log.Error(e)
						return e
					}
					log.With("version", ver, "hash", lastHash).Info("last")
				}

			}

			return nil
		},
		Subcommands: nil,
		Flags:       flags,
	}
}

// ConnectBangumi ...
func (eth *ETH) ConnectBangumi() (*BangumiData, error) {
	tk, err := NewBangumiData(common.HexToAddress(eth.ContractAddress), eth.conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
		return &BangumiData{}, nil
	}
	return tk, nil
}

// ConnectDHash ...
func (eth *ETH) ConnectDHash() (*Dhash, error) {
	tk, err := NewDhash(common.HexToAddress(eth.ContractAddress), eth.conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
		return &Dhash{}, nil
	}
	return tk, nil
}

func singleInput(eth *ETH, video *model.Video) (e error) {
	tk, e := eth.ConnectBangumi()
	if e != nil {
		return e
	}
	privateKey, err := crypto.HexToECDSA(eth.key)
	if err != nil {
		log.Fatal(err)
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
			//log.Fatalf("tx mining error:%v\n", err)
			return err
		}
		log.Info(name + "@" + idxStr + " success")
		log.Debugf("receipt is :%x\n", string(receipt.TxHash[:]))
	}
	return nil
}

func multipleInput(eth *ETH, video *model.Video) (e error) {
	tk, e := eth.ConnectBangumi()
	if e != nil {
		return e
	}
	privateKey, err := crypto.HexToECDSA(eth.key)
	if err != nil {
		log.Fatal(err)
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
				//log.Fatalf("tx mining error:%v\n", err)
				return err
			}
			log.Info(name + "@" + idxStr + " success")
			log.Debugf("receipt is :%x\n", string(receipt.TxHash[:]))
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
	return fn(eth, video)
}

// UpdateAppWithPath ...
func (eth *ETH) UpdateAppWithPath(version, path string) (e error) {
	obj, e := rest.AddFile(path)
	if e != nil {
		return e
	}
	return eth.UpdateApp(version, obj.Hash)
}

// UpdateApp ...
func (eth *ETH) UpdateApp(version string, hash string) (e error) {
	return update(eth, version, hash)
}

func update(eth *ETH, version string, hash string) (e error) {
	tk, e := eth.ConnectDHash()
	if e != nil {
		return e
	}
	privateKey, err := crypto.HexToECDSA(eth.key)
	if err != nil {
		log.Fatal(err)
		return err
	}
	code, e := eth.GetLastVersionCode()
	if e != nil {
		log.Fatal(err)
		return e
	}
	one := big.NewInt(1)
	code = code.Add(code, one)

	opt := bind.NewKeyedTransactor(privateKey)
	transaction, err := tk.UpdateVersion(opt, version, hash, code)
	if err != nil {
		return err
	}

	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, eth.conn, transaction)
	if err != nil {
		//log.Fatalf("tx mining error:%v\n", err)
		return err
	}
	log.Info(version + "@" + hash + " success")
	log.Debugf("receipt is :%x\n", string(receipt.TxHash[:]))
	return nil
}
