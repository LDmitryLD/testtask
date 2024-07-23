package config

import (
	"os"
	"strconv"

	"go.uber.org/zap"
)

type AppConf struct {
	AppName    string
	DB         DB
	GRPCServer GRPCServer
	Logger     Logger
}

type DB struct {
	Driver   string
	Name     string
	User     string
	Password string
	Host     string
	Port     string
	MaxConn  int
	Timeout  int
}

type GRPCServer struct {
	Port string
}

type Logger struct {
	Level string
}

func NewAppConf() AppConf {
	return AppConf{
		AppName: os.Getenv("APP_NAME"),
		DB: DB{
			Driver:   os.Getenv("DB_DRIVER"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		GRPCServer: GRPCServer{
			Port: os.Getenv("GRPC_PORT"),
		},
	}
}

func (a *AppConf) Init(logger *zap.Logger) {
	dbTimeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
	if err != nil {
		logger.Fatal("config: parse db timeout error", zap.Error(err))
	}
	dbMaxConn, err := strconv.Atoi(os.Getenv("MAX_CONN"))
	if err != nil {
		logger.Fatal("config: parse db max connection error", zap.Error(err))
	}

	a.DB.Timeout = dbTimeout
	a.DB.MaxConn = dbMaxConn
}
