package initialisation

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

type Version struct {
	Version, Script string
}

func TestCanInitWhenTableDBVersionDoesntExist(test *testing.T) {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	db.Connect()
	dropAllTables(db)
	db.Query("drop table db_version")
	db.Close()

	Initialize()

	assert.Equal(test, 1, len(getVersions()))
	assert.Equal(test, "0.0.0", getVersions()[0].Version)
	assert.Equal(test, "initialisation", getVersions()[0].Script)
}

func TestDoesntInitWhenTableDBVersionExists(test *testing.T) {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	db.Connect()
	dropAllTables(db)
	db.Query("create table db_version (id INTEGER PRIMARY KEY AUTO_INCREMENT , script VARCHAR(255), version VARCHAR(255), state VARCHAR(255))")
	db.Query("insert into db_version (script, version, state) values ('initialisation', '0.0.0', 'ok')")
	db.Query("drop table db_version")
	db.Close()

	Initialize()

	assert.Equal(test, 1, len(getVersions()))
	assert.Equal(test, "0.0.0", getVersions()[0].Version)
}

func dropAllTables(db mysql.Conn) {
	rows, _, _ := db.Query("show tables")
	var tables []string
	for _, row := range rows {
		tables = append(tables, row.Str(0))
	}
	concatenateTables := strings.Join(tables, ", ")
	db.Query("drop table " + concatenateTables)
}

func getVersions() []Version {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	db.Connect()
	rows, _, _ := db.Query("select version, script from db_version order by id")
	db.Close()
	var versions []Version
	for _, row := range rows {
		versions = append(versions, Version{Version: row.Str(0), Script: row.Str(1)})
	}
	return versions
}