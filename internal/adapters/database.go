package adapters

import (
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

func (gdc *GormDatabase) GetConnectionString() string {
	return getConnectionString()
}

func (gdc *GormDatabase) GetOrm() interfaces.Orm {
	return gdc.conn
}

func (gdc *GormDatabase) Close() {
}

func newDatabaseGorm(connectionString string) (*GormDatabase, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connectionString,
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
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.POSTGRES_USER,
		env.POSTGRES_PASSWORD,
		env.POSTGRES_HOST,
		env.POSTGRES_PORT,
		env.POSTGRES_DB,
	)
}

func GetDatabaseWithConnectionString(connectionString string) interfaces.Database {
	var err error

	if gormDatabase == nil {
		gormDatabase, err = newDatabaseGorm(connectionString)

		if err != nil {
			panic(err)
		}
	}

	return gormDatabase
}

func GetDatabase() interfaces.Database {
	connectionString := getConnectionString()

	if err := config.Migrate(connectionString, false); err != nil {
		panic(err)
	}

	return GetDatabaseWithConnectionString(connectionString)
}
