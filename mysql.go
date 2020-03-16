package connection

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	dbs                = make(map[string]*Repository)
	configs            = make(map[string]*MySQLConfig)
	defaultDatabaseKey = ""
	defaultMySQLConfig = MySQLConfig{
		host:     "127.0.0.1",
		port:     "3306",
		user:     "root",
		password: "",
		database: "",
	}
)

type MySQLOption func(config *MySQLConfig)

func MySQLHost(host string) MySQLOption {
	return func(config *MySQLConfig) {
		config.host = host
	}
}

func MySQLPort(port string) MySQLOption {
	return func(config *MySQLConfig) {
		config.port = port
	}
}

func MySQLUsername(username string) MySQLOption {
	return func(config *MySQLConfig) {
		config.user = username
	}
}

func MySQLPassword(password string) MySQLOption {
	return func(config *MySQLConfig) {
		config.password = password
	}
}

func MySQLDatabase(database string) MySQLOption {
	return func(config *MySQLConfig) {
		config.database = database
	}
}

type MySQLConfig struct {
	host     string
	port     string
	user     string
	password string
	database string
}

func NewMySQLConfig(options ...MySQLOption) *MySQLConfig {

	config := defaultMySQLConfig
	for _, option := range options {
		option(&config)
	}

	if conf, ok := configs[config.database]; ok {
		return conf
	}

	if len(configs) == 0 {
		defaultDatabaseKey = config.database
	}

	configs[config.database] = &config

	return &config
}

func (m *MySQLConfig) Connect() (err error) {

	var repo = &Repository{
		config: m,
	}
	return repo.Connect()
}

func (m *MySQLConfig) Close() error {
	if database, ok := dbs[m.database]; ok {
		return database.Close()
	}
	return nil
}

func (m MySQLConfig) getConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&loc=%s&parseTime=true",
		m.user, m.password, m.host, m.port, m.database, "Local")
}

type Repository struct {
	*gorm.DB
	config *MySQLConfig
}

func (repo *Repository) Connect() (err error) {

	if repo.isConnected() {
		return
	}

	return repo.reconnect()
}

func (repo *Repository) reconnect() (err error) {
	logger.WithField("connectionString", repo.config.getConnectionString()).
		Info("mysql connect")

	repo.DB, err = gorm.Open("mysql", repo.config.getConnectionString())
	if err != nil {
		logger.Errorf("mysql connect: %v", err)
		return
	}

	repo.DB.SingularTable(true)

	logger.Info("mysql connect succeed")

	// 将连接放入缓存
	dbs[repo.config.database] = repo

	return
}

func (repo *Repository) isConnected() bool {
	if repo.DB == nil {
		return false
	}

	db := repo.DB.DB()
	if db == nil {
		return false
	}

	if err := db.Ping(); err != nil {
		return false
	}

	return true
}

func (repo *Repository) Specify(database string) *Repository {

	if repo, ok := dbs[database]; ok {
		_ = repo.Connect()
		return repo
	}

	var config MySQLConfig
	config = *repo.config
	config.database = database

	newRepo := &Repository{
		config: &config,
	}
	if err := newRepo.Connect(); err != nil {
		logger.WithField("config", config).
			WithField("error", err).
			Error("connect error")
	}

	return newRepo
}

func (repo *Repository) Begin() *gorm.DB {

	return repo.DB.Begin()
}

func GetMySQLSpecifyDatabase(database string) *Repository {

	if len(dbs) == 0 {
		return nil
	}

	if repo, ok := dbs[database]; ok {
		_ = repo.Connect()
		return repo
	}

	return dbs[defaultDatabaseKey].Specify(database)
}

func GetMySQL() *Repository {

	return dbs[defaultDatabaseKey]
}

func SetMySQL(database string, repo *Repository) {

	if len(dbs) == 0 {
		defaultDatabaseKey = database
	}
	dbs[database] = repo
}
