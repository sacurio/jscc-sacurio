package db

import (
	"fmt"

	"github.com/sacurio/jb-challenge/internal/app/model"
	"github.com/sacurio/jb-challenge/internal/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	user   string
	pwd    string
	port   string
	host   string
	dbName string
	DB     *gorm.DB
	logger *logrus.Logger
}

func NewDB(config config.DBConfig, logger *logrus.Logger) *DB {
	db := &DB{
		user:   config.User,
		pwd:    config.Pwd,
		port:   config.Port,
		host:   config.Host,
		dbName: config.DBName,
		logger: logger,
	}

	db.setupDB()
	return db
}

func (d *DB) connectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.user, d.pwd, d.host, d.port, d.dbName)
}

// ConnectMySQL configure and sets the MySQL database connection.
func (d *DB) connect() (*gorm.DB, error) {
	d.logger.Info("Attempting to establish a database connection...")
	db, err := gorm.Open(mysql.Open(d.connectionString()), &gorm.Config{})
	if err != nil {
		d.logger.Errorf("Database connection was not acomplished: %s", err.Error())
		return nil, err
	}

	d.logger.Info("Database connections successfully")
	return db, nil
}

func (d *DB) setupDB() error {
	dbConn, err := d.connect()
	if err != nil {
		d.logger.Panicf("Database connection was not acomplished, %s", err.Error())
	}

	d.DB = dbConn
	dbConn.AutoMigrate(&model.User{}, &model.Message{})

	return nil
}
