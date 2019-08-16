package seed

import (
	"context"
	"crypto/sha1"
	"fmt"
	"math"
	"sync"

	"github.com/glvd/seed/model"
	shell "github.com/godcong/go-ipfs-restapi"
	api "github.com/ipfs/go-ipfs-http-client"
	jsoniter "github.com/json-iterator/go"
	ma "github.com/multiformats/go-multiaddr"
)

// Options ...
type Options func(*Seed)

// AfterInitOptions ...
type AfterInitOptions func(*Seed)

// Thread ...
type Thread struct {
	wg sync.WaitGroup
}

// Threader ...
type Threader interface {
	Runnable
	BeforeRun(seed *Seed)
	AfterRun(seed *Seed)
}

// Runnable ...
type Runnable interface {
	Run(context.Context)
}

// Stepper ...
type Stepper int

// StepperNone ...
const (
	// StepperNone ...
	StepperNone Stepper = iota
	// StepperInformation ...
	StepperInformation
	// StepperMoveInfo ...
	StepperMoveInfo
	// StepperProcess ...
	StepperProcess
	// StepperMoveproc ...
	StepperMoveproc
	// StepperTransfer ...
	StepperTransfer
	// StepperPin ...
	StepperPin
	// StepperCheck ...
	StepperCheck
	// StepperUpdate ...
	StepperUpdate
	// StepperMax ...
	StepperMax
)

// Seed ...
type Seed struct {
	Shell       *shell.Shell
	API         *API
	maindb      *Database
	Workspace   string
	Scale       int64
	NoCheck     bool
	Unfinished  map[string]*model.Unfinished
	Videos      map[string]*model.Video
	Moves       map[string]string
	MaxLimit    int
	From        string
	wg          *sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	skipConvert bool
	preAdd      bool
	noSlice     bool
	upScale     bool
	threads     int
	thread      []Threader
	ignores     map[string][]byte
	err         error
	skipExist   bool
}

// Stop ...
func (seed *Seed) Stop() {
	if seed.cancel != nil {
		seed.cancel()
	}
}

// Err ...
func (seed *Seed) Err() error {
	return seed.err
}

// Start ...
func (seed *Seed) Start() {
	seed.wg.Add(1)
	go func() {
		log.Info("first running")
		defer seed.wg.Done()
		for i, thread := range seed.thread {
			log.With("index", i)
			if thread != nil {
				thread.BeforeRun(seed)
				thread.Run(seed.ctx)
				thread.AfterRun(seed)
			}
		}
	}()
}

// Wait ...
func (seed *Seed) Wait() {
	seed.wg.Wait()
}

// NewSeed ...
func NewSeed(ops ...Optioner) *Seed {
	ctx, cancel := context.WithCancel(context.Background())
	seed := &Seed{
		Unfinished: make(map[string]*model.Unfinished),
		Videos:     make(map[string]*model.Video),
		Moves:      make(map[string]string),
		MaxLimit:   math.MaxUint16,
		wg:         &sync.WaitGroup{},
		ctx:        ctx,
		cancel:     cancel,
		threads:    0,
		thread:     make([]Threader, StepperMax),
		ignores:    make(map[string][]byte),
	}

	seed.Register(ops...)

	if seed.Shell == nil {
		seed.Shell = shell.NewShell("localhost:5001")
	}

	if seed.API == nil {
		addr, e := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/5001")
		if e != nil {
			log.Error(e)
			return nil
		}
		addrAPI, e := api.NewApi(addr)
		if e != nil {
			log.Error(e)
			return nil
		}
		seed.API = addrAPI
	}

	return seed
}

// Register ...
func (seed *Seed) Register(ops ...Optioner) {
	for _, op := range ops {
		op.Option(seed)
	}
}

// SyncDatabase ...
func SyncDatabase() AfterInitOptions {
	return func(seed *Seed) {
		if seed.maindb == nil {
			panic("nil database")
		}
		e := model.Sync(seed.maindb.eng)
		if e != nil {
			panic(e)
		}
	}
}

// SkipConvertOption ...
func SkipConvertOption() Options {
	return func(seed *Seed) {
		seed.skipConvert = true
	}
}

// SkipExistOption ...
func SkipExistOption() Options {
	return func(seed *Seed) {
		seed.skipExist = true
	}
}

// NoSliceOption ...
func NoSliceOption() Options {
	return func(seed *Seed) {
		seed.noSlice = true
	}
}

// MaxLimitOption ...
func MaxLimitOption(max int) Options {
	return func(seed *Seed) {
		seed.MaxLimit = max
	}
}

// PreAddOption ...
func PreAddOption() Options {
	return func(seed *Seed) {
		seed.preAdd = true
	}
}

// databaseOption ...
func databaseOption(db *Database) Options {
	return func(seed *Seed) {
		seed.maindb = db
	}
}

// InformationOption ...
func informationOption(info *Information) Options {
	return func(seed *Seed) {
		seed.thread[StepperInformation] = info
	}
}

// UnfinishedOption ...
func UnfinishedOption(unfins ...*model.Unfinished) Options {
	return func(seed *Seed) {
		if seed.Unfinished == nil {
			seed.Unfinished = make(map[string]*model.Unfinished)
		}
		for _, u := range unfins {
			if u == nil {
				continue
			}

			if u.Hash != "" {
				seed.Unfinished[u.Hash] = u
			}
		}
	}
}

// ShellOption ...
func ShellOption(s string) Options {
	return func(seed *Seed) {
		log.Info("ipfs: ", s)
		seed.Shell = shell.NewShell(s)
	}
}

// APIOption ...
func APIOption(s string) Options {
	var e error
	return func(seed *Seed) {
		var addr ma.Multiaddr
		addr, e = ma.NewMultiaddr(s)
		if e != nil {
			log.Error(e)
			return
		}
		seed.API, e = api.NewApi(addr)
		if e != nil {
			log.Error(e)
		}
	}
}

// processOption ...
func processOption(process *process) Options {
	return func(seed *Seed) {
		seed.thread[StepperProcess] = process
	}
}

// pinOption ...
func pinOption(pin *pin) Options {
	return func(seed *Seed) {
		seed.thread[StepperPin] = pin
	}
}

// IgnoreOption ...
func IgnoreOption(ignores ...string) Options {
	return func(seed *Seed) {
		for _, i := range ignores {
			seed.ignores[i] = nil
		}
	}
}

// ThreadOption ...
func ThreadOption(t int) Options {
	return func(seed *Seed) {
		seed.threads = t
	}
}

// Extend ...
type Extend struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// VideoSource ...
type VideoSource struct {
	Bangumi      string    `json:"bangumi"`       //番号 no
	SourceHash   string    `json:"source_hash"`   //原片hash
	Type         string    `json:"type"`          //类型：film，FanDrama
	Format       string    `json:"format"`        //输出：3D，2D
	VR           string    `json:"vr"`            //VR格式：左右，上下，平面
	Thumb        string    `json:"thumb"`         //缩略图
	Intro        string    `json:"intro"`         //简介 title
	Alias        []string  `json:"alias"`         //别名，片名
	VideoEncode  string    `json:"video_encode"`  //视频编码
	AudioEncode  string    `json:"audio_encode"`  //音频编码
	Files        []string  `json:"files"`         //存放路径
	HashFiles    []string  `json:"hash_files"`    //已上传Hash
	CheckFiles   []string  `json:"check_files"`   //Unfinished checksum
	Slice        bool      `json:"sliceAdd"`      //是否HLS切片
	Encrypt      bool      `json:"encrypt"`       //加密
	Key          string    `json:"key"`           //秘钥
	M3U8         string    `json:"m3u8"`          //M3U8名
	SegmentFile  string    `json:"segment_file"`  //ts切片名
	PosterPath   string    `json:"poster_path"`   //海报路径
	Poster       string    `json:"poster"`        //海报HASH
	ExtendList   []*Extend `json:"extend_list"`   //扩展信息
	Role         []string  `json:"role"`          //角色列表 stars
	Director     string    `json:"director"`      //导演
	Systematics  string    `json:"systematics"`   //分级
	Season       string    `json:"season"`        //季
	Episode      string    `json:"episode"`       //集数
	TotalEpisode string    `json:"total_episode"` //总集数
	Sharpness    string    `json:"sharpness"`     //清晰度
	Publish      string    `json:"publish"`       //发行日
	Date         string    `json:"date"`          //发行日
	Length       string    `json:"length"`        //片长
	Producer     string    `json:"producer"`      //制片商
	Series       string    `json:"series"`        //系列
	Tags         []string  `json:"tags"`          //标签
	Publisher    string    `json:"publisher"`     //发行商
	Language     string    `json:"language"`      //语言
	Caption      string    `json:"caption"`       //字幕
	MagnetLinks  []string  `json:"magnet_links"`  //磁链
	Uncensored   bool      `json:"uncensored"`    //有码,无码
}

// Hash ...
func Hash(v interface{}) string {
	bytes, e := jsoniter.Marshal(v)
	if e != nil {
		return ""
	}
	return fmt.Sprintf("%x", sha1.Sum([]byte(bytes)))
}

// SkipTypeVerify ...
func SkipTypeVerify(tp string, v ...interface{}) bool {
	for i := range v {
		if v1, b := (v[i]).(string); b {
			if v1 == tp {
				return true
			}
		}
	}
	return false
}
