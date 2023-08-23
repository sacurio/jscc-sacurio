package db

import (
	"fmt"

	"github.com/sacurio/jb-challenge/internal/app/user"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type dataBase struct {
	user   string
	pwd    string
	port   string
	host   string
	dbName string
	DB     *gorm.DB
	logger *logrus.Logger
}

func NewDatabase(user, pwd, port, host, dbName string, logger *logrus.Logger) *dataBase {
	return &dataBase{
		user:   user,
		pwd:    pwd,
		port:   port,
		host:   host,
		dbName: dbName,
		logger: logger,
	}
}

func (d *dataBase) buildDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.user, d.pwd, d.host, d.port, d.dbName)
}

// ConnectMySQL configure and sets the MySQL database connection.
func (d *dataBase) connect() (*gorm.DB, error) {
	d.logger.Info("Attempting to establish a database connection...")
	db, err := gorm.Open(mysql.Open(d.buildDSN()), &gorm.Config{})
	if err != nil {
		d.logger.Errorf("Database connection was not acomplished: %s", err.Error())
		return nil, err
	}

	d.logger.Info("Database connections successfully")
	return db, nil
}

func (d *dataBase) SetupDB() error {
	dbConn, err := d.connect()
	if err != nil {
		d.logger.Panicf("Database connection was not acomplished, %s", err.Error())
	}

	d.DB = dbConn
	dbConn.AutoMigrate(&user.User{})

	return nil
}
