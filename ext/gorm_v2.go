package ext

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vdongchina/ratgo/utils/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

// Db配置
type Config struct {
	Identification            string // 数据库配置标识
	DriverName                string `json:"driverName"`     // 驱动名称，目前支持: mysql
	DSN                       string `json:"dataSourceName"` // DSN data source name
	SkipInitializeWithVersion bool   // 根据当前 MySQL 版本自动配置
	DefaultStringSize         uint   // string 类型字段的默认长度
	DefaultDatetimePrecision  *int
	DisableDatetimePrecision  bool // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	DontSupportRenameIndex    bool // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	DontSupportRenameColumn   bool // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	DontSupportForShareClause bool
	// 连接池配置
	MaxOpenConn   int           `json:"maxOpenConn"` // SetMaxOpenConns 设置打开数据库连接的最大数量
	MaxIdleConn   int           `json:"maxIdleConn"` // SetMaxIdleConns 设置空闲连接池中连接的最大数量
	MaxLifetime   time.Duration `json:"maxLifetime"` // SetConnMaxLifetime 设置了连接可复用的最大时间
	SingularTable string        `json:"SingularTable"`
}

// *gorm.DB 存储容器
type DbStorage struct {
	Abstract
	dbMap   map[string]*gorm.DB    // gorm存储
	Plugins map[string]gorm.Plugin // gorm插件
	lock    sync.RWMutex           // 排它锁
}

var GormV2 *DbStorage

func init() {
	GormV2 = &DbStorage{
		dbMap:   map[string]*gorm.DB{},
		Plugins: map[string]gorm.Plugin{},
		//lock:    new(sync.RWMutex),
	}
}

// 根据数据库标识获取对应 *gorm.DB
func (ds *DbStorage) GormDB(identification string) *gorm.DB {
	if db, ok := ds.dbMap[identification]; ok {
		return db
	}
	ds.lock.Lock()
	defer ds.lock.Unlock()
	if db, ok := ds.dbMap[identification]; ok {
		return db
	}
	fmt.Println("获取数据库类....")
	db, err := ds.gormDB(identification)
	if err != nil {
		logrus.Panic(fmt.Sprintf("get db failed. error:%v", err))
	} else {
		ds.dbMap[identification] = db
	}
	return db
}

// 根据 identification 读取配置并获取 *gorm.DB
func (ds *DbStorage) gormDB(identification string) (db *gorm.DB, err error) {
	// 读取配置
	config := types.AnyMap(ds.config.Get(identification).ToAnyMap())
	if len(config) == 0 {
		return db, errors.New(fmt.Sprintf("get config failed by identification '%s'", identification))
	}

	// gorm 配置
	gormConfig := &gorm.Config{}

	// 日志输出
	logConfig := types.AnyMap(config.Get("Log").ToAnyMap())
	if logConfig.Get("Turn").ToString() == "on" {
		gormConfig.Logger = logger.New(log.New(os.Stdout, "\r\n"+logConfig.Get("Prefix").ToString(), log.LstdFlags), logger.Config{
			SlowThreshold: time.Duration(logConfig.Get("SlowThreshold").ToInt64()) * time.Millisecond,
			LogLevel:      logger.LogLevel(logConfig.Get("LogLevel").ToInt()),
			Colorful:      logConfig.Get("Colorful").ToBool(),
		})
	}

	// 选择数据库驱动
	driverName := config.Get("DriverName").ToString()
	switch driverName {
	case "mysql":
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       config.Get("DSN").ToString(),
			SkipInitializeWithVersion: config.Get("SkipInitializeWithVersion").ToBool(),
			DefaultStringSize:         uint(config.Get("DefaultStringSize").ToInt()),
			//DefaultDatetimePrecision:  &config.Get("DefaultDatetimePrecision").ToInt(),
			DisableDatetimePrecision:  config.Get("DisableDatetimePrecision").ToBool(),
			DontSupportRenameIndex:    config.Get("DontSupportRenameIndex").ToBool(),
			DontSupportRenameColumn:   config.Get("DontSupportRenameColumn").ToBool(),
			DontSupportForShareClause: config.Get("DontSupportForShareClause").ToBool(),
		}), gormConfig)
	default:
		err = errors.New(fmt.Sprintf("get db driver '%s' failed.", driverName))
	}
	if err != nil {
		return
	}

	// 注册插件
	if plugin, ok := ds.Plugins[identification]; ok {
		err = db.Use(plugin)
		if err != nil {
			return db, err
		}
	}

	// 连接池设置
	if sqlDB, err := db.DB(); err != nil {
		return db, err
	} else {
		connPoolConfig := types.AnyMap(config.Get("ConnPool").ToAnyMap())
		sqlDB.SetMaxOpenConns(connPoolConfig.Get("MaxOpenConn").ToInt())                                   // SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDB.SetMaxIdleConns(connPoolConfig.Get("MaxIdleConn").ToInt())                                   // SetMaxIdleConns 用于设置连接池中空闲连接的最大数
		sqlDB.SetConnMaxLifetime(time.Duration(connPoolConfig.Get("maxLifetime").ToInt64()) * time.Second) // SetConnMaxLifetime 设置了连接可复用的最大时间
	}
	return db, nil
}

// 注册插件
func (ds *DbStorage) RegisterPlugins(identification string, plugin gorm.Plugin) {
	ds.Plugins[identification] = plugin
}
