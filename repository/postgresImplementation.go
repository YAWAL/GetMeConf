// Package repository contains repository interfaces as well as their implementations for given databases.
package repository

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/YAWAL/GetMeConf/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const okResult = "OK"

// PostgresStorage wraps the database connection.
type PostgresStorage struct {
	DB *gorm.DB
}

// PostgresConfig represents a configuration for the postgres database connection
type PostgresConfig struct {
	Schema                   string
	DSN                      string
	MaxOpenedConnectionsToDb int
	MaxIdleConnectionsToDb   int
	ConnMaxLifetimeMinutes   int
}

// initPostgresDB initiates database connection using environmental variables.
func initPostgresDB(conf *PostgresConfig) (db *gorm.DB, err error) {
	db, err = gorm.Open(conf.Schema, conf.DSN)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxOpenConns(conf.MaxOpenedConnectionsToDb)
	db.DB().SetMaxIdleConns(conf.MaxIdleConnectionsToDb)
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(conf.ConnMaxLifetimeMinutes))
	return db, nil
}

// NewPostgresStorage returns a pointer to new PostgresStorage structure.
func NewPostgresStorage(conf *PostgresConfig) (*PostgresStorage, error) {
	db, err := initPostgresDB(conf)
	return &PostgresStorage{
		DB: db,
	}, err
}

//FindMongoDBConfig returns a config record from database using the unique name
func (s *PostgresStorage) FindMongoDBConfig(configName string) (*entity.Mongodb, error) {
	result := entity.Mongodb{}
	err := s.DB.Where("domain = ?", configName).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//FindAllMongoDBConfig returns all config record of one type from database
func (s *PostgresStorage) FindAllMongoDBConfig() ([]entity.Mongodb, error) {
	var mongoConfigs []entity.Mongodb
	err := s.DB.Find(&mongoConfigs).Error
	if err != nil {
		return nil, err
	}
	return mongoConfigs, nil
}

//SaveMongoDBConfig saves new config record to the database
func (s *PostgresStorage) SaveMongoDBConfig(config *entity.Mongodb) (string, error) {
	s.DB.LogMode(true)
	err := s.DB.Create(config).Error
	if err != nil {
		return "", err
	}
	return okResult, nil
}

//DeleteMongoDBConfig removes config record from database
func (s *PostgresStorage) DeleteMongoDBConfig(configName string) (string, error) {
	rowsAffected := s.DB.Delete(entity.Mongodb{}, "domain = ?", configName).RowsAffected
	if rowsAffected < 1 {
		return "", errors.New("could not delete from database")
	}
	return fmt.Sprintf("deleted %d row(s)", rowsAffected), nil
}

//UpdateMongoDBConfig updates a record in database, rewriting the fields if string fields are not empty
func (s *PostgresStorage) UpdateMongoDBConfig(newConfig *entity.Mongodb) (string, error) {
	var persistedConfig entity.Mongodb
	err := s.DB.Where("domain = ?", newConfig.Domain).Find(&persistedConfig).Error
	if err != nil {
		return "", err
	}
	if newConfig.Host != "" && newConfig.Port != "" {
		err = s.DB.Exec("UPDATE mongodbs SET mongodb = ?, port = ?, host = ? WHERE domain = ?",
			strconv.FormatBool(newConfig.Mongodb), newConfig.Port, newConfig.Host, persistedConfig.Domain).Error
		if err != nil {
			return "", err
		}
		return okResult, nil
	}
	return "", errors.New("fields are empty")
}

//FindTempConfig returns a config record from database using the unique name.
func (s *PostgresStorage) FindTempConfig(configName string) (*entity.Tempconfig, error) {
	result := entity.Tempconfig{}
	err := s.DB.Where("rest_api_root = ?", configName).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//FindAllTempConfig returns all config record of one type from database.
func (s *PostgresStorage) FindAllTempConfig() ([]entity.Tempconfig, error) {
	var tempConfigs []entity.Tempconfig
	err := s.DB.Find(&tempConfigs).Error
	if err != nil {
		return nil, err
	}
	return tempConfigs, nil
}

//SaveTempConfig saves new config record to the database.
func (s *PostgresStorage) SaveTempConfig(config *entity.Tempconfig) (string, error) {
	err := s.DB.Create(config).Error
	if err != nil {
		return "", err
	}
	return okResult, nil
}

//DeleteTempConfig removes config record from database.
func (s *PostgresStorage) DeleteTempConfig(configName string) (string, error) {
	rowsAffected := s.DB.Delete(entity.Tempconfig{}, "rest_api_root = ?", configName).RowsAffected
	if rowsAffected < 1 {
		return "", errors.New("could not delete from database")
	}
	return fmt.Sprintf("deleted %d row(s)", rowsAffected), nil
}

//UpdateTempConfig updates a record in database, rewriting the fields if string fields are not empty.
func (s *PostgresStorage) UpdateTempConfig(newConfig *entity.Tempconfig) (string, error) {
	var persistedConfig entity.Tempconfig
	err := s.DB.Where("rest_api_root = ?", newConfig.RestApiRoot).Find(&persistedConfig).Error
	if err != nil {
		return "", err
	}
	if newConfig.Host != "" && newConfig.Port != "" && newConfig.Remoting != "" {
		err = s.DB.Exec("UPDATE tempconfigs SET remoting = ?, port = ?, host = ?, legasy_explorer = ? WHERE rest_api_root = ?",
			newConfig.Remoting, newConfig.Port, newConfig.Host, strconv.FormatBool(newConfig.LegasyExplorer), persistedConfig.RestApiRoot).Error
		if err != nil {
			return "", err
		}
		return okResult, nil
	}
	return "", errors.New("fields are empty")
}

//FindTsConfig returns a config record from database using the unique name.
func (s *PostgresStorage) FindTsConfig(configName string) (*entity.Tsconfig, error) {
	result := entity.Tsconfig{}
	err := s.DB.Where("module = ?", configName).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//FindAllTsConfig returns all config record of one type from database.
func (s *PostgresStorage) FindAllTsConfig() ([]entity.Tsconfig, error) {
	var tsConfigs []entity.Tsconfig
	err := s.DB.Find(&tsConfigs).Error
	if err != nil {
		return nil, err
	}
	return tsConfigs, nil
}

//SaveTsConfig saves new config record to the database.
func (s *PostgresStorage) SaveTsConfig(config *entity.Tsconfig) (string, error) {
	_, err := s.DB.Create(config).Rows()
	if err != nil {
		return "", err
	}
	return okResult, nil
}

//DeleteTsConfig removes config record from database.
func (s *PostgresStorage) DeleteTsConfig(configName string) (string, error) {
	rowsAffected := s.DB.Delete(entity.Tsconfig{}, "module = ?", configName).RowsAffected
	if rowsAffected < 1 {
		return "", errors.New("could not delete from database")
	}
	return fmt.Sprintf("deleted %d row(s)", rowsAffected), nil
}

//UpdateTsConfig updates a record in database, rewriting the fields if string fields are not empty.
func (s *PostgresStorage) UpdateTsConfig(newConfig *entity.Tsconfig) (string, error) {
	var persistedConfig entity.Tsconfig
	err := s.DB.Where("module = ?", newConfig.Module).Find(&persistedConfig).Error
	if err != nil {
		return "", err
	}
	if newConfig.Target != "" {
		err = s.DB.Exec("UPDATE tsconfigs SET target = ?, source_map = ?, excluding = ? WHERE module = ?",
			newConfig.Target, strconv.FormatBool(newConfig.SourceMap), strconv.Itoa(newConfig.Excluding), persistedConfig.Module).Error
		if err != nil {
			return "", err
		}
		return okResult, nil
	}
	return "", errors.New("fields are empty")
}
