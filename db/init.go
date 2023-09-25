package db

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"consul-ext/pkg/config"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/proxy"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Pool  *sqlx.DB
	Gpool *gorm.DB
)

func Init() {
	dbStr := config.Data.Database.Username + ":" + config.Data.Database.Password + "@tcp(" + config.Data.Database.Host + ":" + fmt.Sprint(config.Data.Database.Port) + ")/" + config.Data.Database.DBname + "?parseTime=true&loc=Local"
	maxIdleConns := 20
	dbMaxIdleConns := os.Getenv("DB_MAX_IDLE_CONNS")
	if dbMaxIdleConns != "" {
		count, err := strconv.Atoi(dbMaxIdleConns)
		if err != nil {
			panic(err)
		}
		maxIdleConns = count
	}
	maxLifeTime := 10
	dbMaxLifeTime := os.Getenv("DB_LIFE_TIME")
	if dbMaxLifeTime != "" {
		count, err := strconv.Atoi(dbMaxLifeTime)
		if err != nil {
			panic(err)
		}
		maxLifeTime = count
	}
	maxOpenConns := 100
	dbMaxOpenConns := os.Getenv("DB_MAX_OPEN_CONNS")
	if dbMaxOpenConns != "" {
		count, err := strconv.Atoi(dbMaxOpenConns)
		if err != nil {
			panic(err)
		}
		maxOpenConns = count
	}
	dialer := proxy.FromEnvironment()
	mysql.RegisterDial("tcp", func(network string) (net.Conn, error) {
		return dialer.Dial("tcp", network)
	})
	var err error
	Pool, err = sqlx.Open("mysql", dbStr)
	if err != nil {
		panic(err.Error())
	}
	Pool.SetMaxIdleConns(maxIdleConns)
	Pool.SetMaxOpenConns(maxOpenConns)
	Pool.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Minute)
	err = Pool.Ping()
	if err != nil {
		panic(err.Error())
	}
	Gpool, err = gorm.Open(gormMysql.New(gormMysql.Config{Conn: Pool}), &gorm.Config{})
	if err != nil {
		log.Println("gorm connect failed ", err.Error())
		panic(err.Error())
	}
	log.Println("mysql connect suc")
}
