package repository

import (
	"fmt"

	"github.com/YAWAL/GetMeConf/entity"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/gormigrate.v1"
)

var serviceMigrations = []*gormigrate.Migration{
	{ID: "20181906",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&entity.Mongodb{},
			).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable(
				&entity.Mongodb{},
			).Error
		},
	},
	{ID: "20181907",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&entity.Tempconfig{},
			).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable(
				&entity.Tempconfig{},
			).Error
		},
	},
	{ID: "20181908",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&entity.Tsconfig{},
			).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable(
				&entity.Tsconfig{},
			).Error
		},
	},
}

// RunMigrations - run migration of data.
func RunMigrations(db *gorm.DB) error {
	//	db.LogMode(c.Config.Debug())
	return gormigrate.New(db, gormigrate.DefaultOptions, serviceMigrations).Migrate()
}

// RollbackMigrations rollbacks all migrations.
func RollbackMigrations(db *gorm.DB) (err error) {
	//	db.LogMode(c.Config.Debug())
	var errs gorm.Errors

	migrate := gormigrate.New(db, gormigrate.DefaultOptions, serviceMigrations)
	for _, migration := range serviceMigrations {
		if err := migrate.RollbackLast(); err != nil {
			errs.Add(errors.Wrap(err, fmt.Sprintf("rollback migration error, ID %s", migration.ID)))
		}
	}

	if len(errs.GetErrors()) > 0 {
		return errors.New(errs.Error())
	}
	return
}
