package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Anomalies struct {
	SessionID string
	Value     float64
	Timestamp string
}

type DataBase struct {
	DB *gorm.DB
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`
}

func New(p Postgres) (*DataBase, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		p.Host, p.Port, p.User, p.DBName, p.Password, p.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Anomalies{})
	if err != nil {
		return nil, err
	}

	return &DataBase{
		DB: db,
	}, nil
}

func (d *DataBase) WriteAnomaly(sessionID string, frequency float64, timestamp string) {
	d.DB.Create(&Anomalies{
		SessionID: sessionID,
		Value:     frequency,
		Timestamp: timestamp,
	})
}
