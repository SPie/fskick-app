package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbConfig struct {
	database string
	withDebug bool
	withLog bool
}

func CreateDbConfig(database string, withDebug bool, withLog bool) DbConfig {
	return DbConfig{
		database: database,
		withDebug: withDebug,
		withLog: withLog,
	}
}

func OpenDbConnection(cfg DbConfig) (*sql.DB, error) {
	return sql.Open("sqlite3", cfg.database)
}

type ConnectionHandler struct {
	connection *gorm.DB
}

func NewConnectionHandler(cfg DbConfig) (*ConnectionHandler, error) {
	conn, err := createSqliteConnection(getLogger(cfg.withLog), cfg.database)
	if err != nil {
		return nil, err
	}

	if cfg.withDebug {
		conn.SetDebug()
	}

	return conn, nil
}

func createConnectionHandler(connection *gorm.DB) *ConnectionHandler {
	return &ConnectionHandler{connection: connection}
}

func createSqliteConnection(logger logger.Interface, database string) (*ConnectionHandler, error) {
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{Logger: logger})
	if err != nil {
		return &ConnectionHandler{}, err
	}

	return createConnectionHandler(db), nil
}

func getLogger(shouldLog bool) logger.Interface {
	logLevel := logger.Silent
	if shouldLog {
		logLevel = logger.Error
	}

	l := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logLevel,
		Colorful:                  true,
		IgnoreRecordNotFoundError: true,
	})

	return l
}

func (connectionHandler *ConnectionHandler) Create(model interface{}) error {
	result := connectionHandler.connection.Create(model)

	return result.Error
}

func (connectionHandler *ConnectionHandler) Save(model interface{}) error {
	result := connectionHandler.connection.Save(model)

	return result.Error
}

func (connectionHandler *ConnectionHandler) Find(models interface{}, condition ...interface{}) error {
	result := connectionHandler.connection.Find(models, condition...)

	return result.Error
}

func (connectionHandler *ConnectionHandler) FindOne(model interface{}, condition interface{}) error {
	result := connectionHandler.connection.First(model, condition)

	return result.Error
}

func (connectionHandler *ConnectionHandler) GetAll(models interface{}) error {
	result := connectionHandler.connection.Find(models)

	return result.Error
}

func (connectionHandler *ConnectionHandler) Preload(modelName string) *ConnectionHandler {
	chainConnectionHandler := connectionHandler.getInstane()
	chainConnectionHandler.connection = chainConnectionHandler.connection.Preload(modelName)

	return chainConnectionHandler
}

func (connectionHandler *ConnectionHandler) Select(query interface{}, args ...interface{}) *ConnectionHandler {
	chainConnectionHandler := connectionHandler.getInstane()
	chainConnectionHandler.connection = chainConnectionHandler.connection.Select(query, args...)

	return chainConnectionHandler
}

func (connectionHandler *ConnectionHandler) Joins(modelName string) *ConnectionHandler {
	chainConnectionHandler := connectionHandler.getInstane()
	chainConnectionHandler.connection = chainConnectionHandler.connection.Joins(modelName)

	return chainConnectionHandler
}

func (connectionHandler *ConnectionHandler) Count(model interface{}) (int, error) {
	var count int64

	result := connectionHandler.connection.Model(model).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

func (connectionHandler *ConnectionHandler) Where(query interface{}, args ...interface{}) *ConnectionHandler {
	chainConnectionHandler := connectionHandler.getInstane()
	chainConnectionHandler.connection = chainConnectionHandler.connection.Where(query, args...)

	return chainConnectionHandler
}

func (connectionHandler *ConnectionHandler) AutoMigrate(model interface{}) {
	connectionHandler.connection.AutoMigrate(model)
}

func (connectionHandler *ConnectionHandler) Exec(statement string, parameters ...interface{}) error {
	connectionHandler.connection.Exec(statement, parameters...)

	return nil
}

func (connectionHandler *ConnectionHandler) Close() error {
	return nil
}

func (connectionHandler *ConnectionHandler) SetDebug() {
	connectionHandler.connection = connectionHandler.connection.Debug()
}

func (connectionHandler *ConnectionHandler) getInstane() *ConnectionHandler {
	return createConnectionHandler(connectionHandler.connection)
}
