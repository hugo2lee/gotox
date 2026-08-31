package ormx_test

import (
	"testing"

	"github.com/hugo2lee/gotox/ormx"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestOrmxDatabaseWiring(t *testing.T) {
	defaultDB := &gorm.DB{}
	namedDB := &gorm.DB{}

	orm := ormx.NewWithDBs(nil, nil, map[string]*gorm.DB{
		ormx.DefaultProjectName: defaultDB,
		"named":                namedDB,
	})

	assert.Same(t, defaultDB, orm.GetDB())
	assert.Same(t, namedDB, orm.GetDB("named"))
}
