package storage

import (
	"time"

	"github.com/boltdb/bolt"
)

func NewBoltDB(path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// package storage

// import (
// 	"github.com/go-pg/pg"
// 	"github.com/go-pg/pg/orm"
// )

// //NewPostgresDB returns a new postgres database with proper connections
// func NewPostgresDB(username, password, address, name string) *pg.DB {
// 	db := pg.Connect(&pg.Options{
// 		User:     username,
// 		Password: password,
// 		Addr:     address,
// 		Database: name,
// 	})

// 	return db
// }

// //InitDB creates new postgres tables following the struct model passed to it
// func InitDB(database *pg.DB, models []interface{}) error {
// 	for _, model := range models {
// 		err := database.CreateTable(model, &orm.CreateTableOptions{
// 			IfNotExists: true,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }