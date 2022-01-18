package blot

import (
	"github.com/boltdb/bolt"
	"log"
)

func Setup() *bolt.DB {
	blot, err := bolt.Open("./db/blot/blot.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	blot.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("UrlBucket"))
		return err
	})
	return blot
}
