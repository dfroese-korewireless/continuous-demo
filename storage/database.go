package storage

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/dfroese-korewireless/continuous-demo/messages"
	bolt "go.etcd.io/bbolt"
)

const (
	messagesBucket = "messages"
)

// Database is an interface that all storage providers must satisfy
type Database interface {
	StoreMessage(messages.Message) (uint64, error)
	GetAllMessages() ([]messages.Message, error)
	GetMessage(uint64) (messages.Message, error)
}

type boltDB struct {
	filePath string
}

func setupDatabase(s *boltDB) error {
	db, err := bolt.Open(s.filePath, 0666, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = createBucket(db, messagesBucket)
	if err != nil {
		return err
	}

	return nil
}

func createBucket(db *bolt.DB, name string) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

// New returns a new bolt instance that satisfies the Database interface
func New(path string) (Database, error) {
	b := boltDB{
		filePath: path,
	}

	err := setupDatabase(&b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (b *boltDB) StoreMessage(msg messages.Message) (uint64, error) {
	db, err := bolt.Open(b.filePath, 0666, nil)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var msgID uint64
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(messagesBucket))
		id, err := bucket.NextSequence()
		msg.ID = id
		msgID = id

		buf, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		return bucket.Put(itob(id), buf)
	})

	return msgID, nil
}

func (b *boltDB) GetAllMessages() ([]messages.Message, error) {
	db, err := bolt.Open(b.filePath, 0666, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var msgs []messages.Message
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(messagesBucket))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var msg messages.Message
			err := json.Unmarshal(v, &msg)
			if err != nil {
				fmt.Printf("marshalling data into struct: %s\n", err)
			}

			msgs = append(msgs, msg)
		}
		return nil
	})

	return msgs, nil
}

func (b *boltDB) GetMessage(id uint64) (messages.Message, error) {
	db, err := bolt.Open(b.filePath, 0666, nil)
	if err != nil {
		return messages.Message{}, err
	}
	defer db.Close()

	msg := messages.Message{}
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(messagesBucket))
		v := bucket.Get(itob(id))

		if v == nil {
			return fmt.Errorf("looking up value for %d return nil", id)
		}

		err = json.Unmarshal(v, &msg)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return messages.Message{}, err
	}

	return msg, nil
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
