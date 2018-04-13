// Package repository contains repository interfaces as well as their implementations for given databases.
package repository

import (
	"github.com/YAWAL/GetMeConf/entity"
)

// MongoDBConfigRepo is a repository interface for MongoDB configs.
type MongoDBConfigRepo interface {
	Find(configName string) (*entity.Mongodb, error)
	FindAll() ([]entity.Mongodb, error)
	Update(config *entity.Mongodb) (string, error)
	Save(config *entity.Mongodb) (string, error)
	Delete(configName string) (string, error)
}

// TempConfigRepo is a repository interface for Tempconfigs.
type TempConfigRepo interface {
	Find(configName string) (*entity.Tempconfig, error)
	FindAll() ([]entity.Tempconfig, error)
	Update(config *entity.Tempconfig) (string, error)
	Save(config *entity.Tempconfig) (string, error)
	Delete(configName string) (string, error)
}

// TsConfigRepo is a repository interface for Tsconfigs.
type TsConfigRepo interface {
	Find(configName string) (*entity.Tsconfig, error)
	FindAll() ([]entity.Tsconfig, error)
	Update(config *entity.Tsconfig) (string, error)
	Save(config *entity.Tsconfig) (string, error)
	Delete(configName string) (string, error)
}
