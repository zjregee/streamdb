package storage

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
)

const (
	defaultDBPath = "./data"
)

func NewDB() *DB {
	return &DB{path: defaultDBPath}
}

func NewDBWithPath(path string) *DB {
	return &DB{path: path}
}

type DB struct {
	db   *bolt.DB
	path string
}

func (db *DB) Open() error {
	var err error
	db.db, err = bolt.Open(db.path, 0600, nil)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) ReadData(bucketName, key string) (string, error) {
	var result string
	err := db.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		value := bucket.Get([]byte(key))
		result = string(value)
		return nil
	})
	return result, err
}

func (db *DB) WriteData(bucketName, key, value string) error {
	err := db.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		if err := bucket.Put([]byte(key), []byte(value)); err != nil {
			return err
		}
		return nil
	})
	return err
}
