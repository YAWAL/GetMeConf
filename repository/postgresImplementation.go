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
	"gopkg.in/gormigrate.v1"
)

type Storage interface {
	Migrate(db *gorm.DB) error
}

type PostgresStorage struct {
	MongoDBRepo *MongoDBConfigRepoImpl
	TsRepo      *TsConfigRepoImpl
	TempRepo    *TempConfigRepoImpl
}

// MongoDBConfigRepoImpl represents an implementation of a MongoDB configs repository.
type MongoDBConfigRepoImpl struct {
	DB *gorm.DB
}

// TsConfigRepoImpl represents an implementation of a Tsconfigs repository.
type TsConfigRepoImpl struct {
	DB *gorm.DB
}

// TempConfigRepoImpl represents an implementation of a Tempconfigs repository.
type TempConfigRepoImpl struct {
	DB *gorm.DB
}

type PostgresConfig struct {
	Shema                    string
	DSN                      string
	MaxOpenedConnectionsToDb int
	MaxIdleConnectionsToDb   int
	MbConnMaxLifetimeMinutes int
}

// InitPostgresDB initiates database connection using environmental variables.
func initPostgresDB(conf *PostgresConfig) (db *gorm.DB, err error) {
	db, err = gorm.Open(conf.Shema, conf.DSN)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxOpenConns(conf.MaxOpenedConnectionsToDb)
	db.DB().SetMaxIdleConns(conf.MaxIdleConnectionsToDb)
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(conf.MbConnMaxLifetimeMinutes))
	return db, nil
}

// NewPostgresStorage returns a pointer to new PostgresStorage structure.
func NewPostgresStorage(conf *PostgresConfig) (*PostgresStorage, error) {
	db, err := initPostgresDB(conf)
	return &PostgresStorage{
		MongoDBRepo: &MongoDBConfigRepoImpl{DB: db},
		TsRepo:      &TsConfigRepoImpl{DB: db},
		TempRepo:    &TempConfigRepoImpl{DB: db},
	}, err
}

// NewMongoDBConfigRepo returns a new MongoDB configs repository.
func NewMongoDBConfigRepo(db *gorm.DB) MongoDBConfigRepo {
	return &MongoDBConfigRepoImpl{
		DB: db,
	}
}

// NewTempConfigRepo returns a new Tempconfigs repository.
func NewTempConfigRepo(db *gorm.DB) TempConfigRepo {
	return &TempConfigRepoImpl{
		DB: db,
	}
}

// NewTsConfigRepo returns a new TsConfig repository.
func NewTsConfigRepo(db *gorm.DB) TsConfigRepo {
	return &TsConfigRepoImpl{
		DB: db,
	}
}

func (s *PostgresStorage) Migrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "Initial",
			Migrate: func(tx *gorm.DB) error {
				type Mongodb struct {
					//gorm.Model
					Domain  string `gorm:"primary_key"`
					Mongodb bool
					Host    string
					Port    string
				}
				type Tsconfig struct {
					//gorm.Model
					Module    string `gorm:"primary_key"`
					Target    string
					SourceMap bool
					Excluding int
				}
				type Tempconfig struct {
					//gorm.Model
					RestApiRoot    string `gorm:"primary_key"`
					Host           string
					Port           string
					Remoting       string
					LegasyExplorer bool
				}
				return tx.AutoMigrate(&Mongodb{}, &Tsconfig{}, &Tempconfig{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("mongodbs", "tsconfigs", "tempconfigs").Error
			},
		},
	})

	err := m.Migrate()
	if err != nil {
		return err
	}
	return err
}

//Find returns a config record from database using the unique name
func (r *MongoDBConfigRepoImpl) Find(configName string) (*entity.Mongodb, error) {
	result := entity.Mongodb{}
	err := r.DB.Where("domain = ?", configName).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//FindAll returns all config record of one type from database
func (r *MongoDBConfigRepoImpl) FindAll() ([]entity.Mongodb, error) {
	var confSlice []entity.Mongodb
	err := r.DB.Find(&confSlice).Error
	if err != nil {
		return nil, err
	}
	return confSlice, nil
}

//Save saves new config record to the database
func (r *MongoDBConfigRepoImpl) Save(config *entity.Mongodb) (string, error) {
	err := r.DB.Create(config).Error
	if err != nil {
		return "", err
	}
	return "OK", nil
}

//Delete removes config record from database
func (r *MongoDBConfigRepoImpl) Delete(configName string) (string, error) {
	rowsAffected := r.DB.Delete(entity.Mongodb{}, "domain = ?", configName).RowsAffected
	if rowsAffected < 1 {
		return "", errors.New("could not delete from database")
	}
	return fmt.Sprintf("deleted %d row(s)", rowsAffected), nil
}

//Update updates a record in database, rewriting the fields if string fields are not empty
func (r *MongoDBConfigRepoImpl) Update(newConfig *entity.Mongodb) (string, error) {
	var persistedConfig entity.Mongodb
	err := r.DB.Where("domain = ?", newConfig.Domain).Find(&persistedConfig).Error
	if err != nil {
		return "", err
	}
	if newConfig.Host != "" && newConfig.Port != "" {
		err = r.DB.Exec("UPDATE mongodbs SET mongodb = ?, port = ?, host = ? WHERE domain = ?", strconv.FormatBool(newConfig.Mongodb), newConfig.Port, newConfig.Host, persistedConfig.Domain).Error
		if err != nil {
			return "", err
		}
		return "OK", nil
	}
	return "", errors.New("fields are empty")
}

//Find returns a config record from database using the unique name
func (r *TempConfigRepoImpl) Find(configName string) (*entity.Tempconfig, error) {
	result := entity.Tempconfig{}
	err := r.DB.Where("rest_api_root = ?", configName).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//FindAll returns all config record of one type from database
func (r *TempConfigRepoImpl) FindAll() ([]entity.Tempconfig, error) {
	var confSlice []entity.Tempconfig
	err := r.DB.Find(&confSlice).Error
	if err != nil {
		return nil, err
	}
	return confSlice, nil
}

//Save saves new config record to the database
func (r *TempConfigRepoImpl) Save(config *entity.Tempconfig) (string, error) {
	err := r.DB.Create(config).Error
	if err != nil {
		return "", err
	}
	return "OK", nil
}

//Delete removes config record from database
func (r *TempConfigRepoImpl) Delete(configName string) (string, error) {
	rowsAffected := r.DB.Delete(entity.Tempconfig{}, "rest_api_root = ?", configName).RowsAffected
	if rowsAffected < 1 {
		return "", errors.New("could not delete from database")
	}
	return fmt.Sprintf("deleted %d row(s)", rowsAffected), nil
}

//Update updates a record in database, rewriting the fields if string fields are not empty
func (r *TempConfigRepoImpl) Update(newConfig *entity.Tempconfig) (string, error) {
	var persistedConfig entity.Tempconfig
	err := r.DB.Where("rest_api_root = ?", newConfig.RestApiRoot).Find(&persistedConfig).Error
	if err != nil {
		return "", err
	}
	if newConfig.Host != "" && newConfig.Port != "" && newConfig.Remoting != "" {
		err = r.DB.Exec("UPDATE tempconfigs SET remoting = ?, port = ?, host = ?, legasy_explorer = ? WHERE rest_api_root = ?", newConfig.Remoting, newConfig.Port, newConfig.Host, strconv.FormatBool(newConfig.LegasyExplorer), persistedConfig.RestApiRoot).Error
		if err != nil {
			return "", err
		}
		return "OK", nil
	}
	return "", errors.New("fields are empty")
}

//Find returns a config record from database using the unique name
func (r *TsConfigRepoImpl) Find(configName string) (*entity.Tsconfig, error) {
	result := entity.Tsconfig{}
	err := r.DB.Where("module = ?", configName).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//FindAll returns all config record of one type from database
func (r *TsConfigRepoImpl) FindAll() ([]entity.Tsconfig, error) {
	var confSlice []entity.Tsconfig
	err := r.DB.Find(&confSlice).Error
	if err != nil {
		return nil, err
	}
	return confSlice, nil
}

//Save saves new config record to the database
func (r *TsConfigRepoImpl) Save(config *entity.Tsconfig) (string, error) {
	err := r.DB.Create(config).Error
	if err != nil {
		return "", err
	}
	return "OK", nil
}

//Delete removes config record from database
func (r *TsConfigRepoImpl) Delete(configName string) (string, error) {
	rowsAffected := r.DB.Delete(entity.Tsconfig{}, "module = ?", configName).RowsAffected
	if rowsAffected < 1 {
		return "", errors.New("could not delete from database")
	}
	return fmt.Sprintf("deleted %d row(s)", rowsAffected), nil
}

//Update updates a record in database, rewriting the fields if string fields are not empty
func (r *TsConfigRepoImpl) Update(newConfig *entity.Tsconfig) (string, error) {
	var persistedConfig entity.Tsconfig
	err := r.DB.Where("module = ?", newConfig.Module).Find(&persistedConfig).Error
	if err != nil {
		return "", err
	}
	if newConfig.Target != "" {
		err = r.DB.Exec("UPDATE tsconfigs SET target = ?, source_map = ?, excluding = ? WHERE module = ?", newConfig.Target, strconv.FormatBool(newConfig.SourceMap), strconv.Itoa(newConfig.Excluding), persistedConfig.Module).Error
		if err != nil {
			return "", err
		}
		return "OK", nil
	}
	return "", errors.New("fields are empty")
}
