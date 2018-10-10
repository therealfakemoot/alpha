package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

// Conf is for deserializing the JSON config file.
type Conf struct {
	Token  string `json:"token"`
	Prefix string `json:"prefix"`
	Status string `json:"status"`
	State  map[string]interface{}
}

// LoadConfig parses the JSON config file and initializes or updates the provided boltdb instance.
func LoadConfig(db *bolt.DB) {
	var v Conf

	user, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to load config file: %s", err)
	}

	home := user.HomeDir

	f, err := os.Open(filepath.Join(home, "alpha.json"))
	if err != nil {
		log.Fatalf("Unable to load config file: %s", err)
	}

	rawConf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Unable to read config file: %s", err)
	}

	err = json.Unmarshal(rawConf, &v)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}

	// This is the initializer. It will create a bucket for the bot's core info and populate that bucket with info from the config.json, or it will update existing keys to match what's currently in the config.
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("core"))

		if err != nil {
			return err
		}

		err = b.Put([]byte("token"), []byte(v.Token))
		if err != nil {
			return err
		}

		err = b.Put([]byte("prefix"), []byte(v.Prefix))
		if err != nil {
			return err
		}

		err = b.Put([]byte("status"), []byte(v.Status))
		if err != nil {
			return err
		}

		return nil
	})
}
