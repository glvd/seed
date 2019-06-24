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
	"time"

	"github.com/go-xorm/xorm"
	"github.com/godcong/go-trait"
	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
)

var db *xorm.Engine
var syncTable = map[string]interface{}{}
var log = trait.NewZapSugar(zap.String("package", "model"))

// Database ...
type Database struct {
	ShowSQL  bool   `toml:"show_sql"`
	UseCache bool   `json:"use_cache"`
	Type     string `toml:"type"`
	Addr     string `toml:"addr"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Schema   string `toml:"schema"`
	Charset  string `toml:"charset"`
	Prefix   string `toml:"prefix"`
	Loc      string `toml:"loc"`
	location string
}

// DefaultDB ...
func DefaultDB() *Database {
	return &Database{
		ShowSQL:  true,
		UseCache: true,
		Type:     "mysql",
		Addr:     "localhost",
		Port:     "3306",
		Username: "root",
		Password: "111111",
		Schema:   "yinhe",
		Loc:      url.QueryEscape("Asia/Shanghai"),
		Charset:  "utf8mb4",
		Prefix:   "",
	}
}

// SetLocation ...
func (d *Database) SetLocation(loc string) {
	d.location = url.QueryEscape(loc)
}

// Location ...
func (d *Database) Location() string {
	if d.location != "" {
		return d.location
	}
	return url.QueryEscape(d.Loc)
}

// Source ...
func (d *Database) Source() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?loc=%s&charset=%s&parseTime=true",
		d.Username, d.Password, d.Addr, d.Port, d.Schema, url.QueryEscape(d.location), d.Charset)
}

// RegisterTable ...
func RegisterTable(v interface{}) {
	tof := reflect.TypeOf(v).Name()
	syncTable[tof] = v
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

// DB ...
func DB() *xorm.Engine {
	if db == nil {
		if err := InitSQLite3(); err != nil {
			panic(err)
		}
	}
	return db
}

// InitSQLite3 ...
func InitSQLite3() (e error) {
	eng, e := xorm.NewEngine("sqlite3", "seed.db")
	if e != nil {
		return e
	}
	eng.ShowSQL(true)
	eng.ShowExecTime(true)
	result, e := eng.Exec("PRAGMA journal_mode = OFF;")
	if e != nil {
		return e
	}
	log.Info("result:", result)
	for idx, val := range syncTable {
		log.Info("syncing ", idx)
		e := eng.Sync2(val)
		if e != nil {
			return e
		}
	}

	db = eng
	return nil
}

// InitDB ...
func InitDB(db, source string) (eng *xorm.Engine, e error) {
	eng, e = xorm.NewEngine(db, source)
	if e != nil {
		return
	}
	//eng.ShowSQL(true)
	//eng.ShowExecTime(true)
	//for idx, val := range syncTable {
	//	log.Info("syncing ", idx)
	//	e = eng.Sync2(val)
	//	if e != nil {
	//		return
	//	}
	//}
	return eng, nil
}

// LoadToml ...
func LoadToml(path string) (db *Database) {
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
