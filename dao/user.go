package dao

import (
	"gin_demo/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var GlobalDb1 *gorm.DB

func DataLoad() {
	db, _ := gorm.Open(mysql.New(mysql.Config{ //配置
		DSN: "root:123456@tcp(127.0.0.1:3306)/gindemo?charset=utf8mb4&parseTime=True&loc=Local",
	}), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "gindemo_",
			SingularTable: false,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxLifetime(10) //数据池
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	GlobalDb1 = db  //全局变量GlobalDb赋值
	TestUserCreat() //若无表单自动创建表单
}

func TestUserCreat() {
	GlobalDb1.AutoMigrate(&model.User{})
}

// 若没有这个用户返回 false，反之返回 true
func SelectUser(username string) bool {
	var u struct {
		Username string
	}
	GlobalDb1.Model(&model.User{}).Where("username = ?", username).Find(&u)
	if u.Username == "" {
		return false
	} else {
		return true
	}
}