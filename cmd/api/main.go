package main

import (
	"database/sql"
	"equiptrack/config"
	"equiptrack/internal/server"
	"flag"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/sirupsen/logrus"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "configPath", "config/config-test.yml", "path to config file")
}

func main() {
	flag.Parse()

	log.Println("Starting api server")
	log.Printf("configPath: %s\n", configPath)

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logrus.New()
	appLogger.SetLevel(logrus.DebugLevel)

	psqlDB, err := newDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	s := server.NewServer(cfg, psqlDB, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}

func newDB(c *config.Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Database.DbHost,
		c.Database.DbPort,
		c.Database.DbUser,
		c.Database.Dbname,
		c.Database.DbPassword,
	)
	db, err := sql.Open(c.Database.DbDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
