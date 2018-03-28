package global

import (
	"github.com/go-pg/pg"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"fmt"
	"github.com/anhnguyentb/grpc-implement/models"
	"github.com/go-pg/pg/orm"
)

var Db *pg.DB

func LoadDatabase() error {

	var err error
	Db, err = GetConnection()
	if err != nil {
		return err
	}

	return nil
}

func GetConnection() (*pg.DB, error) {

	for _, val := range []string{"host", "username", "password", "name"} {
		if !viper.IsSet("database." + val) {

			Log.Errorw(
				"Database configuration is not exists",
				zap.String("key", "database." + val),
			)
			return nil, fmt.Errorf("Database configuration is not exists\n")
		}
	}

	db := pg.Connect(&pg.Options{
		Addr: viper.GetString("database.host"),
		User: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		Database: viper.GetString("database.name"),
	})

	_, err := db.Exec("SELECT 1")
	if err != nil {
		Log.Infow(
			"Cannot connect to database",
			zap.String("host", viper.GetString("database.host")),
			zap.String("user", viper.GetString("database.username")),
			zap.Error(err),
		)
		return nil, err
	}

	Log.Infow(
		"Database is connected",
		zap.String("host", viper.GetString("database.host")),
		zap.String("user", viper.GetString("database.username")),
	)

	return db, nil
}

func CreateSchema() error {

	Log.Info("Creating schema for audits table")
	return Db.CreateTable(&models.Audit{}, &orm.CreateTableOptions{IfNotExists: true})
}
