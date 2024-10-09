package adapters

import "gorm.io/gorm"

type Orm *gorm.DB

type Database interface {
	GetConnectionString() string
	GetOrm() Orm
}
