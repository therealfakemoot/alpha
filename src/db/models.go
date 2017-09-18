package db

import (
	"fmt"
	bdb "github.com/tidwall/buntdb"
)

// Model encapsulates a value that can be serialized to BuntDB.
type Model interface {
	Key() string
}

// MovingAverages blah blah blah
func MovingAverages(tx *bdb.Tx) {}

func initDB() (*bdb.DB, error) {
	db, err := bdb.Open(":memory:")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return db, nil
}

func init() {
	db, err := initDB()

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

}
