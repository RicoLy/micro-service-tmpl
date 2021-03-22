package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"micro-service-tmpl/utils/log"
	"micro-service-tmpl/utils/viper"
)

var (
	MasterDB *gorm.DB //主数据库连接
	SlaveDB  *gorm.DB //从数据库连接
)

func init() {
	var (
		masterCfg DbConfig //主数据库配置
		slaveCfg  DbConfig //从数据库配置
		err       error
	)

	if err = viper.ViperConf.UnmarshalKey("mysql.Master", &masterCfg); err != nil {
		log.GetLogger().Fatal("数据库获取配置文件失败" + err.Error())
	}

	if err = viper.ViperConf.UnmarshalKey("mysql.Slave", &slaveCfg); err != nil {
		log.GetLogger().Fatal("数据库获取配置文件失败" + err.Error())
	}

	// 主
	if MasterDB, err = InitMysql(masterCfg); err != nil {
		log.GetLogger().Fatal("数据库连接失败" + err.Error())
	}
	// 从
	if SlaveDB, err = InitMysql(slaveCfg); err != nil {
		log.GetLogger().Fatal("数据库连接失败" + err.Error())
	}
}

type DbConfig struct {
	Host    string //地址
	Port    int    //端口
	Name    string //用户
	Pass    string //密码
	DBName  string //库名
	Charset string //编码
	MaxIdle int    //最大空闲连接
	MaxOpen int    //最大连接数
}

func InitMysql(cfg DbConfig) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		cfg.Name,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	defer func() { //报错释放资源
		if err := recover(); err != nil {
			if err1 := MasterDBClose(); err1 != nil {
				panic(err1)
			}
			//if err2 := SlaveDBClose(); err2 != nil {
			//	panic(err2)
			//}
			panic(err)
		}
	}()

	if db, err = gorm.Open("mysql", dsn); err != nil {
		panic(err)
	}
	sqlDb := db.DB()
	sqlDb.SetMaxIdleConns(cfg.MaxIdle) //空闲连接数
	sqlDb.SetMaxOpenConns(cfg.MaxOpen) //最大连接数

	db.LogMode(true) //打开sql执行日志
	db = db.Debug()  //debug模式

	//添加钩子函数
	//addCallBackFunc(db)

	//数据迁移生成表
	//err = db.AutoMigrate(
	//	//todo dbModel 指针
	//	&AiDistort{},
	//).Error

	//db.SingularTable(true)	//数据迁移生成表结尾不带s
	return
}

//添加钩子函数
//func addCallBackFunc(db *gorm.DB) {
//
//	// 创建时雪花算法生成ID
//	db.Callback().Create().Before("gorm:create").Register("auto_insert_id", func(scope *gorm.Scope) {
//		// 若主键为空自动填充 ID
//		isTrue := scope.PrimaryKeyZero()
//		if isTrue {
//			_ = scope.SetColumn("ID", NextId())
//		}
//	})
//
//	// 修改时间 使用时间戳
//	db.Callback().Update().Before("gorm:update").Register("update_at_to_stamp", func(scope *gorm.Scope) {
//		if _, ok := scope.Get("gorm:update_column"); !ok {
//			_ = scope.SetColumn("UpdatedAt", time.Now().Unix())
//		}
//	})
//}

// close MasterDb
func MasterDBClose() error {
	if MasterDB != nil {
		return MasterDB.Close()
	}
	return nil
}

// close slaveDb
//func SlaveDBClose() error {
//	if SlaveDB != nil {
//		return SlaveDB.Close()
//	}
//	return nil
//}
