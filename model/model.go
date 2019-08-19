package model

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/godcong/go-trait"
	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"golang.org/x/xerrors"
)

var globalDB *xorm.Engine
var syncTable = map[string]interface{}{}
var log = trait.NewZapFileSugar()

//SetGlobalDB set db
func SetGlobalDB(eng *xorm.Engine) {
	globalDB = eng
}

//DB get global db
func DB() *xorm.Engine {
	if globalDB != nil {
		return globalDB
	}
	panic(xerrors.New("nil db"))
}

// DatabaseConfig ...
type DatabaseConfig struct {
	ShowSQL      bool   `toml:"show_sql"`
	ShowExecTime bool   `toml:"show_exec_time"`
	UseCache     bool   `json:"use_cache"`
	Type         string `toml:"type"`
	Addr         string `toml:"addr"`
	Port         string `toml:"port"`
	Username     string `toml:"username"`
	Password     string `toml:"password"`
	Schema       string `toml:"schema"`
	Charset      string `toml:"charset"`
	Prefix       string `toml:"prefix"`
	Loc          string `toml:"loc"`
	location     string
}

// DefaultDB ...
func DefaultDB() *DatabaseConfig {
	return &DatabaseConfig{
		ShowSQL:  true,
		UseCache: true,
		Type:     "mysql",
		Addr:     "localhost",
		Port:     "3306",
		Username: "root",
		Password: "111111",
		Schema:   "glvd",
		Loc:      "Asia/Shanghai",
		Charset:  "utf8mb4",
		Prefix:   "",
	}
}

// SetLocation ...
func (d *DatabaseConfig) SetLocation(loc string) {
	d.location = url.QueryEscape(loc)
}

// Location ...
func (d *DatabaseConfig) Location() string {
	if d.Loc != "" {
		d.location = url.QueryEscape(d.Loc)
		d.Loc = "" //clear the loc buf
	}
	return d.location
}

// Source ...
func (d *DatabaseConfig) Source() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?loc=%s&charset=%s&parseTime=true",
		d.Username, d.Password, d.Addr, d.Port, d.Schema, url.QueryEscape(d.location), d.Charset)
}

var (
	tableMu       sync.RWMutex
	tableRegister = make(map[string]interface{})
)

// Register ...
func Register(name string, table interface{}) {
	tableMu.Lock()
	defer tableMu.Unlock()
	if table == nil {
		panic("table: Register table is nil")
	}
	if _, dup := tableRegister[name]; dup {
		panic("table: Register called twice for table " + name)
	}
	tableRegister[name] = table
}

// RegisterTable ...
func RegisterTable(v interface{}) {
	Register(reflect.TypeOf(v).Name(), v)
}

// Tables ...
func Tables() []interface{} {
	var r []interface{}
	for _, tb := range syncTable {
		r = append(r, tb)
	}
	return r
}

// Sync ...
func Sync(db *xorm.Engine) (e error) {
	for idx, val := range syncTable {
		log.Info("syncing ", idx)
		e = db.Sync2(val)
		if e != nil {
			return
		}
	}
	return nil
}

// SQLite3DB ...
func SQLite3DB(name string) string {
	return fmt.Sprintf("file:%s?cache=shared&mode=rwc&_journal_mode=WAL", name)
}

// InitSQLite3 ...
func InitSQLite3(name string) (eng *xorm.Engine, e error) {
	eng, e = xorm.NewEngine("sqlite3", SQLite3DB(name))
	if e != nil {
		return nil, e
	}

	return eng, nil
}

// MustDatabase ...
func MustDatabase(engine *xorm.Engine, err error) *xorm.Engine {
	if err != nil {
		panic(err)
	}
	return engine
}

// InitDB ...
func InitDB(db *DatabaseConfig) (eng *xorm.Engine, e error) {
	eng, e = xorm.NewEngine(db.Type, db.Source())
	if e != nil {
		return
	}
	return eng, nil
}

//LoadDatabaseConfig ...
func LoadDatabaseConfig(path string) (db *DatabaseConfig) {
	db = DefaultDB()
	tree, err := toml.LoadFile(path)
	if err != nil {
		return db
	}
	err = tree.Unmarshal(db)
	if err != nil {
		return db
	}
	return db
}

// Model ...
type Model struct {
	ID        string     `xorm:"id pk"`
	CreatedAt time.Time  `xorm:"created_at created"`
	UpdatedAt time.Time  `xorm:"updated_at updated"`
	DeletedAt *time.Time `xorm:"deleted_at deleted"`
	Version   int        `xorm:"version"`
}

// Modeler ...
type Modeler interface {
	GetID() string
	SetID(string)
	GetVersion() int
	SetVersion(int)
}

// BeforeInsert ...
func (m *Model) BeforeInsert() {
	if m.ID == "" {
		m.ID = uuid.Must(uuid.NewRandom()).String()
	}
}

// MustSession ...
func MustSession(session *xorm.Session) *xorm.Session {
	if session == nil {
		return DB().NewSession()
	}
	return session
}

// Checksum ...
func Checksum(filepath string) string {
	hash := sha1.New()
	file, e := os.Open(filepath)
	if e != nil {
		return ""
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	_, e = io.Copy(hash, reader)
	if e != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}
