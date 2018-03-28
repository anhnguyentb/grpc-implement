package global

import (
	"testing"
	"github.com/anhnguyentb/grpc-implement/models"
)

func init() {
	LoadConfig()
	LoadLogger(true)
}

func TestLoadDatabase(t *testing.T) {

	err := LoadDatabase()
	if err != nil {
		t.Fatal("Fatal error to load database connection")
	}
	defer Db.Close()
}

func TestCreateSchema(t *testing.T) {

	err := LoadDatabase()
	if err != nil {
		t.Fatal("Fatal error to load database connection")
	}
	defer Db.Close()

	err = CreateSchema()
	if err != nil {
		t.Errorf("Fail to create schema %s \n", err)
	}

	_, err = Db.Query(&models.Audit{}, "SELECT * FROM audits")
	if err != nil {
		t.Errorf("Fail to create schema %s \n", err)
	}
}


