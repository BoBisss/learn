package config

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"personal_blog/global"
	"personal_blog/model"
)

type Config struct {
	Dsn string
}

var AppConfig Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Configuration initialization failed：%v", err.Error())
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Configuration serialization failed:%v", err.Error())
	}
	initDB()
}

func initDB() {
	dsn := AppConfig.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Connecting to the database failed:%v", err.Error())
	}
	global.DB = db

	//迁移模型
	err1 := db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	if err1 != nil {
		log.Fatalf("Model migration failed:%v", err1.Error())
	}
}
