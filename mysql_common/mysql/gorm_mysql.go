package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// 连接管理器
type RDBManager struct {
	OpenTx bool     // 是否开启事务
	DsName string   // 数据源名称
	Db     *gorm.DB // 非事务实例
	Tx     *gorm.Tx // 事务实例
	Errors []error  // 操作过程中记录的错误
}

// 数据库配置
type DBConfig struct {
	DsName   string // 数据源名称
	Host     string // 地址IP
	Port     string // 数据库端口
	Database string // 数据库名称
	Username string // 账号
	Password string // 密码
}

// db连接
var (
	MASTER = "DB1"                    // 默认主数据源
	RDBs   = map[string]*RDBManager{} // 初始化时加载数据源到集合
)

// 初始化多个数据库配置文件
func BuildByConfig(input ...DBConfig) {
	if len(input) == 0 {
		panic("数据源配置不能为空")
	}
	for _, v := range input {
		db, err := MysqlSetup(v)
		if err != nil {
			log.Printf("数据库链接失败 %s ", err.Error())
			return
		}
		if len(v.DsName) == 0 {
			v.DsName = MASTER
		}
		rdb := &RDBManager{
			Db:     db,
			DsName: v.DsName,
		}
		RDBs[v.DsName] = rdb
	}
}

// Setup 初始化连接
func MysqlSetup(conf DBConfig) (*gorm.DB, error) {
	//启用打印日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level: Silent、Error、Warn、Info
			Colorful:      false,       // 禁用彩色打印
		},
	)
	// repository = newConnection()
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database)
	dialector := mysql.New(mysql.Config{
		DSN:                       dbURI, // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	})
	// conn, err := gorm.Open(dialector, &gorm.Config{
	// 	Logger: newLogger,
	// })
	conn, err := gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Print(err.Error())
		return conn, err
	}
	sqlDB, err := conn.DB()
	if err != nil {
		log.Print("connect repository server failed.")
	}
	sqlDB.SetMaxIdleConns(100) // 设置最大连接数
	sqlDB.SetMaxOpenConns(100) // 设置最大的空闲连接数
	sqlDB.SetConnMaxLifetime(3600)
	return conn, nil
}
