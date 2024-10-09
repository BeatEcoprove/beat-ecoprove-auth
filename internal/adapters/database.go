package adapters

import (
	"errors"
	"fmt"

	"github.com/BeatEcoprove/identityService/config"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormDatabase *GormDatabase

type GormDatabase struct {
	conn *gorm.DB
}

func (gdc *GormDatabase) Close() error {

	return errors.New("")
}

func (gdc *GormDatabase) GetConnectionString() string {
	return getConnectionString()
}

func (gdc *GormDatabase) GetOrm() interfaces.Orm {
	return gdc.conn
}

func newDatabaseGorm() (*GormDatabase, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: getConnectionString(),
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	return &GormDatabase{
		conn: db,
	}, nil
}

func getConnectionString() string {
	env := config.GetCofig()

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		env.POSTGRES_USER,
		env.POSTGRES_PASSWORD,
		env.POSTGRES_HOST,
		env.POSTGRES_PORT,
		env.POSTGRES_DB,
	)
}

func GetDatabase() interfaces.Database {
	var err error

	if gormDatabase == nil {
		gormDatabase, err = newDatabaseGorm()

		if err != nil {
			panic(err)
		}
	}

	return gormDatabase
}
