package connection

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetMySQL(t *testing.T) {

	repo1 := &Repository{config: NewMySQLConfig()}
	tests := []struct {
		inputDatabase string
		inputRepo     *Repository
		wantDatabase  *Repository
	}{
		{
			inputDatabase: "test1",
			inputRepo:     repo1,
			wantDatabase:  repo1,
		},
	}

	for _, test := range tests {

		SetMySQL(test.inputDatabase, test.inputRepo)
		assert.Equal(t, test.wantDatabase, dbs[test.inputDatabase])
	}
}

func TestGetMySQL(t *testing.T) {

	repo1 := &Repository{config: NewMySQLConfig()}
	tests := []struct {
		inputDatabase string
		inputRepo     *Repository
		wantDatabase  *Repository
	}{
		{
			inputDatabase: "test1",
			inputRepo:     repo1,
			wantDatabase:  repo1,
		},
	}

	for _, test := range tests {

		defaultDatabaseKey = test.inputDatabase
		dbs[defaultDatabaseKey] = test.inputRepo
		assert.Equal(t, test.wantDatabase, GetMySQL())
	}
}

func TestGetMySQLSpecifyDatabase(t *testing.T) {

	repo1 := &Repository{config: NewMySQLConfig()}
	tests := []struct {
		inputDatabase string
		inputRepo     *Repository
		wantDatabase  *Repository
	}{
		{
			inputDatabase: "test1",
			inputRepo:     repo1,
			wantDatabase:  repo1,
		},
	}

	for _, test := range tests {

		dbs[test.inputDatabase] = test.inputRepo
		assert.Equal(t, test.wantDatabase, GetMySQLSpecifyDatabase(test.inputDatabase))
	}
}

func TestRepository_Transaction(t *testing.T) {

	rand.Seed(time.Now().UnixNano())
	repo1 := &Repository{config: NewMySQLConfig(MySQLDatabase("test"))}
	err := repo1.Connect()
	assert.Nil(t, err)

	type testModel struct {
		Id   uint64 `gorm:"primary_key;auto_increment:false"`
		Name string
	}

	repo1.AutoMigrate(&testModel{})

	var m = testModel{Id: rand.Uint64(), Name: "test"}
	transaction := repo1.Begin()
	transaction.Model(&testModel{}).Create(m)
	assert.False(t, repo1.Model(&testModel{}).NewRecord(m))
	transaction.Commit()

	var readModel testModel
	repo1.Where(&testModel{Id: m.Id}).First(&readModel)
	assert.EqualValues(t, m, readModel)
}
