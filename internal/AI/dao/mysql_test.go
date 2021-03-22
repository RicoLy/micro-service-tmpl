package dao

import (
	"fmt"
	"testing"
)

func TestInitMysql(t *testing.T) {
	cfg := DbConfig{
		Host:    "192.168.142.128",
		Port:    3306,
		Name:    "root",
		Pass:    "mysqlly",
		DBName:  "ai",
		Charset: "utf8mb4",
		MaxIdle: 50,
		MaxOpen: 50,
	}
	db, err := InitMysql(cfg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(db)
}
