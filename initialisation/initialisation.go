package initialisation

import (
	. "db_versioning/log"
	"fmt"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

func Initialize(schema string) {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", schema)
	err := db.Connect()
	if err != nil {
		Fail("Error while connecting to database : %s \n", err.Error())
	}

	row, _, err := db.QueryFirst("show tables like 'db_version'")
	if err != nil {
		Fail("Error while fetching db_version : %s \n", err.Error())
	}
	if row == nil {
		db.Query("create table db_version (id INTEGER PRIMARY KEY AUTO_INCREMENT , script VARCHAR(255), version VARCHAR(255), state VARCHAR(255))")
		db.Query("insert into db_version (script, version, state) values ('initialisation', '0.0.0', 'ok')")
		fmt.Printf("Database schema version initialized for '%s' \n", schema)
	} else {
		fmt.Printf("Database schema version is already initialized for '%s' \n", schema)
	}
	db.Close()
}
