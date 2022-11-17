package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ConnectionHandler interface {
	Create(model interface{}) error
	Save(model interface{}) error
	FindOne(model interface{}, condition interface{}) error
	Find(models interface{}, condition ...interface{}) error
	GetAll(models interface{}) error
	Select(query interface{}, args ...interface{}) ConnectionHandler
	Preload(modelName string) ConnectionHandler
	Joins(modelName string) ConnectionHandler
	Count(model interface{}) (int, error)
	Where(query interface{}, args ...interface{}) ConnectionHandler
	AutoMigrate(model interface{})
	Close() error
	SetDebug()
}

type connectionHandler struct {
	connection *gorm.DB
}

func NewConnectionHandler(database string, username string, password string, host string, port string, driver string, log bool) (ConnectionHandler, error) {
	var connectionHandler ConnectionHandler
	var err error

	switch driver {
	case "mysql":
		connectionHandler, err = createMysqlConnection(getLogger(log), database, username, password, host, port)
	case "sqlite":
		fallthrough
	default:
		connectionHandler, err = createSqliteConnection(getLogger(log), database)
	}

	return connectionHandler, err
}

func createConnectionHandler(connection *gorm.DB) *connectionHandler {
	return &connectionHandler{connection: connection}
}

func createSqliteConnection(logger logger.Interface, database string) (ConnectionHandler, error) {
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{Logger: logger})
	if err != nil {
		return &connectionHandler{}, err
	}

	return createConnectionHandler(db), nil
}

func createMysqlConnection(
	logger logger.Interface,
	database string,
	username string,
	password string,
	host string,
	port string) (ConnectionHandler, error) {
	db, err := gorm.Open(
		mysql.Open(
			fmt.Sprintf(
				"%s:%s@(%s:%s)/%s?parseTime=true",
				username,
				password,
				host,
				port,
				database,
			),
		),
		&gorm.Config{Logger: logger},
	)
	if err != nil {
		return &connectionHandler{}, err
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

func (connectionHandler *connectionHandler) Create(model interface{}) error {
	result := connectionHandler.connection.Create(model)

	return result.Error
}

func (connectionHandler *connectionHandler) Save(model interface{}) error {
	result := connectionHandler.connection.Save(model)

	return result.Error
}

func (connectionHandler *connectionHandler) Find(models interface{}, condition ...interface{}) error {
	result := connectionHandler.connection.Find(models, condition...)

	return result.Error
}

func (connectionHandler *connectionHandler) FindOne(model interface{}, condition interface{}) error {
	result := connectionHandler.connection.First(model, condition)

	return result.Error
}

func (connectionHandler *connectionHandler) GetAll(models interface{}) error {
	result := connectionHandler.connection.Find(models)

	return result.Error
}

func (connectionHandler *connectionHandler) Preload(modelName string) ConnectionHandler {
	chainConnectionHandler := connectionHandler.getInstane()
	chainConnectionHandler.connection = chainConnectionHandler.connection.Preload(modelName)

	return chainConnectionHandler
}

func (connectionHandler *connectionHandler) Select(query interface{}, args ...interface{}) ConnectionHandler {
	chainConnectionHandler := connectionHandler.getInstane()
	chainConnectionHandler.connection = chainConnectionHandler.connection.Select(query, args...)

	return chainConnectionHandler
}

func (connectionHandler *connectionHandler) Joins(modelName string) ConnectionHandler {
	chainConnectionHandler := connectionHandler.getInstane()
	chainConnectionHandler.connection = chainConnectionHandler.connection.Joins(modelName)

	return chainConnectionHandler
}

func (connectionHandler *connectionHandler) Count(model interface{}) (int, error) {
	var count int64

	result := connectionHandler.connection.Model(model).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

func (connectionHandler *connectionHandler) Where(query interface{}, args ...interface{}) ConnectionHandler {
	chainConnectionHandler := connectionHandler.getInstane()
	chainConnectionHandler.connection = chainConnectionHandler.connection.Where(query, args...)

	return chainConnectionHandler
}

func (connectionHandler *connectionHandler) AutoMigrate(model interface{}) {
	connectionHandler.connection.AutoMigrate(model)
}

func (connectionHandler *connectionHandler) Close() error {
	return nil
}

func (connectionHandler *connectionHandler) SetDebug() {
	connectionHandler.connection = connectionHandler.connection.Debug()
}

func (connectionHandler *connectionHandler) getInstane() *connectionHandler {
	return createConnectionHandler(connectionHandler.connection)
}
