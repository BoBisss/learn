package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"personal_blog/global"
	"personal_blog/model"
)

type Config struct {
	Dsn  string
	Port string
}

var AppConfig Config

func InitConfig() {
	InitLogger()

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		global.Logger.Error("Configuration initialization failed", zap.Error(err))
		log.Fatalf("Configuration initialization failed：%v", err.Error())
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		global.Logger.Error("Configuration serialization failed", zap.Error(err))
		log.Fatalf("Configuration serialization failed:%v", err.Error())
	}
	initDB()
}

func initDB() {
	dsn := AppConfig.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		global.Logger.Error("Connecting to the database failed", zap.Error(err))
		log.Fatalf("Connecting to the database failed:%v", err.Error())
	}
	global.DB = db

	//迁移模型
	err1 := db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	if err1 != nil {
		global.Logger.Error("Model migration failed", zap.Error(err1))
		log.Fatalf("Model migration failed:%v", err1.Error())
	}

	global.Logger.Info("Database connected successfully")
}
