// Package repository contains repository interfaces as well as their implementations for given databases.
package repository

import (
	"github.com/YAWAL/GetMeConf/entity"
)

// Storage interface collects all the methods to interact with a database.
type Storage interface {
	Migrate() error

	FindMongoDBConfig(configName string) (*entity.Mongodb, error)
	FindAllMongoDBConfig() ([]entity.Mongodb, error)
	UpdateMongoDBConfig(config *entity.Mongodb) (string, error)
	SaveMongoDBConfig(config *entity.Mongodb) (string, error)
	DeleteMongoDBConfig(configName string) (string, error)

	FindTempConfig(configName string) (*entity.Tempconfig, error)
	FindAllTempConfig() ([]entity.Tempconfig, error)
	UpdateTempConfig(config *entity.Tempconfig) (string, error)
	SaveTempConfig(config *entity.Tempconfig) (string, error)
	DeleteTempConfig(configName string) (string, error)

	FindTsConfig(configName string) (*entity.Tsconfig, error)
	FindAllTsConfig() ([]entity.Tsconfig, error)
	UpdateTsConfig(config *entity.Tsconfig) (string, error)
	SaveTsConfig(config *entity.Tsconfig) (string, error)
	DeleteTsConfig(configName string) (string, error)
}
