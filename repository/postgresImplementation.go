// Package repository contains repository interfaces as well as their implementations for given databases.
package repository

import (
	"fmt"
	"strconv"

	"time"

	"os"

	"errors"

	"net/url"

	"github.com/YAWAL/GetMeConf/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gopkg.in/gormigrate.v1"
)

const (
	pdbScheme           = "PDB_SCHEME"
	pdbHost             = "PDB_HOST"
	pdbPort             = "PDB_PORT"
	pdbUser             = "PDB_USER"
	pdbPassword         = "PDB_PASSWORD"
	pdbName             = "PDB_NAME"
	maxOpCon            = "MAX_OPENED_CONNECTIONS_TO_DB"
	maxIdleCon          = "MAX_IDLE_CONNECTIONS_TO_DB"
	vConnMaxLifetimeMin = "MB_CONN_MAX_LIFETIME_MINUTES"
)

var (
	defaultDbScheme                 = "postgres"
	defaultDbHost                   = "horton.elephantsql.com"
	defaultDbPort                   = "5432"
	defaultDbUser                   = "dlxifkbx"
	defaultDbPassword               = "L7Cey-ucPY4L3T6VFlFdNykNE4jO0VjV"
	defaultDbName                   = "dlxifkbx"
	defaultMaxOpenedConnectionsToDb = 5
	defaultMaxIdleConnectionsToDb   = 0
	defaultmbConnMaxLifetimeMinutes = 30
)

var logger *zap.Logger

// ServiceConfig structure contains the configuration information for the database.
type postgresConfig struct {
	dbSchema                 string
	dbHost                   string `yaml:"dbhost"`
	dbPort                   string `yaml:"dbport"`
	dbUser                   string `yaml:"dbUser"`
	dbPassword               string `yaml:"dbPassword"`
	dbName                   string `yaml:"dbName"`
	maxOpenedConnectionsToDb int    `yaml:"maxOpenedConnectionsToDb"`
	maxIdleConnectionsToDb   int    `yaml:"maxIdleConnectionsToDb"`
	mbConnMaxLifetimeMinutes int    `yaml:"mbConnMaxLifetimeMinutes"`
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

// InitZapLogger is used to set the logger for the package.
func InitZapLogger(zlog *zap.Logger) {
	logger = zlog
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

func (c *postgresConfig) validate() {
	if c.dbSchema == "" {
		logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbScheme))
		c.dbSchema = defaultDbScheme
	}
	if c.dbHost == "" {
		logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbHost))
		c.dbHost = defaultDbHost
	}
	if c.dbPort == "" {
		logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbPort))
		c.dbPort = defaultDbPort
	}
	if c.dbUser == "" {
		logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbUser))
		c.dbUser = defaultDbUser
	}
	if c.dbPassword == "" {
		logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbPassword))
		c.dbPassword = defaultDbPassword
	}
	if c.dbName == "" {
		logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbName))
		c.dbName = defaultDbName
	}
	if c.maxOpenedConnectionsToDb == 0 {
		logger.Info("maxOpenedConnectionsToDb = 0", zap.Int("default value is used ", defaultMaxOpenedConnectionsToDb))
		c.maxOpenedConnectionsToDb = defaultMaxOpenedConnectionsToDb
	}
	if c.maxIdleConnectionsToDb == 0 {
		logger.Info("maxIdleConnectionsToDb = 0", zap.Int("default value is used ", defaultMaxIdleConnectionsToDb))
		c.maxIdleConnectionsToDb = defaultMaxIdleConnectionsToDb
	}
	if c.mbConnMaxLifetimeMinutes == 0 {
		logger.Info("mbConnMaxLifetimeMinutes = 0", zap.Int("default value is used ", defaultmbConnMaxLifetimeMinutes))
		c.mbConnMaxLifetimeMinutes = defaultmbConnMaxLifetimeMinutes
	}
}

func initPostgresConfig() *postgresConfig {
	c := new(postgresConfig)
	c.dbSchema = os.Getenv(pdbScheme)
	c.dbHost = os.Getenv(pdbHost)
	c.dbPort = os.Getenv(pdbPort)
	c.dbUser = os.Getenv(pdbUser)
	c.dbPassword = os.Getenv(pdbPassword)
	c.dbName = os.Getenv(pdbName)
	var err error
	c.maxOpenedConnectionsToDb, err = strconv.Atoi(os.Getenv(maxOpCon))
	if err != nil {
		logger.Info("error during reading env. variable. Could not convert from string to int", zap.Error(err))
	}
	c.maxIdleConnectionsToDb, err = strconv.Atoi(os.Getenv(maxIdleCon))
	if err != nil {
		logger.Info("error during reading env. variable. Could not convert from string to int", zap.Error(err))
	}
	c.mbConnMaxLifetimeMinutes, err = strconv.Atoi(os.Getenv(vConnMaxLifetimeMin))
	if err != nil {
		logger.Info("error during reading env. variable. Could not convert from string to int", zap.Error(err))
	}
	return c
}

// InitPostgresDB initiates database connection using environmental variables.
func InitPostgresDB() (db *gorm.DB, err error) {
	conf := initPostgresConfig()
	conf.validate()
	dbInf := url.URL{Scheme: conf.dbSchema, User: url.UserPassword(conf.dbUser, conf.dbPassword), Host: conf.dbHost + ":" + conf.dbPort, Path: conf.dbName}
	db, err = gorm.Open("postgres", dbInf.String()+"?sslmode=disable")

	if err != nil {
		logger.Info("error during connection to postgres database has occurred", zap.Error(err))
		return nil, err
	}

	db.DB().SetMaxOpenConns(conf.maxOpenedConnectionsToDb)
	db.DB().SetMaxIdleConns(conf.maxIdleConnectionsToDb)
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(conf.mbConnMaxLifetimeMinutes))
	logger.Info("connection to postgres database has been established")

	if err = migrate(db); err != nil {
		logger.Info("error during migration", zap.Error(err))
		return nil, err
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
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
		logger.Info("could not migrate", zap.Error(err))
	}
	logger.Info("Migration did run successfully")
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
		logger.Info("error during saving to database", zap.Error(err))
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
			logger.Info("error during updating", zap.Error(err))
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
		logger.Info("error during saving to database", zap.Error(err))
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
			logger.Info("error during updating", zap.Error(err))
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
		logger.Info("error during saving to database", zap.Error(err))
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
			logger.Info("error during updating", zap.Error(err))
			return "", err
		}
		return "OK", nil
	}
	return "", errors.New("fields are empty")
}
