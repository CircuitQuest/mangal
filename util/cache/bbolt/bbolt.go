package bbolt

import (
	"time"

	"github.com/luevano/mangal/config"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/util"
	bolt "go.etcd.io/bbolt"
)

const TTLBucketName = "ttl"

var _ gokv.Store = (*Store)(nil)

// Store is a gokv.Store implementation for bbolt (formerly known as Bolt / Bolt DB).
type Store struct {
	db         *bolt.DB
	bucketName string
	codec      encoding.Codec
	ttl        time.Duration
}

// Set stores the given value for the given key.
// Values are automatically marshalled to JSON or gob (depending on the configuration).
// The key must not be "" and the value must not be nil.
func (s Store) Set(k string, v interface{}) error {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return err
	}

	// First turn the passed object into something that bbolt can handle
	data, err := s.codec.Marshal(v)
	if err != nil {
		return err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		if err := b.Put([]byte(k), data); err != nil {
			return err
		}

		if s.ttl == 0 {
			return nil
		}
		bTTL := tx.Bucket([]byte(TTLBucketName))
		return bTTL.Put([]byte(s.bucketName+"-"+k), []byte(time.Now().UTC().Format(time.RFC3339Nano)))
	})
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves the stored value for the given key.
// You need to pass a pointer to the value, so in case of a struct
// the automatic unmarshalling can populate the fields of the object
// that v points to with the values of the retrieved object's values.
// If no value is found it returns (false, nil).
// The key must not be "" and the pointer must not be nil.
func (s Store) Get(k string, v interface{}) (found bool, err error) {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return false, err
	}

	var data []byte
	err = s.db.View(func(tx *bolt.Tx) error {
		if s.ttl != 0 {
			// If the data has expired (surpassed the TTL) then
			// don't even look for the data
			bTTL := tx.Bucket([]byte(TTLBucketName))
			ttlData := bTTL.Get([]byte(s.bucketName + "-" + k))
			if ttlData != nil {
				ttl, err := time.Parse(time.RFC3339Nano, string(ttlData))
				if err != nil {
					return err
				}
				if time.Now().UTC().After(ttl.Add(s.ttl)) {
					return nil
				}
			}
		}

		b := tx.Bucket([]byte(s.bucketName))
		txData := b.Get([]byte(k))
		// txData is only valid during the transaction.
		// Its value must be copied to make it valid outside of the tx.
		// TODO: Benchmark if it's faster to copy + close tx,
		// or to keep the tx open until unmarshalling is done.
		if txData != nil {
			// `data = append([]byte{}, txData...)` would also work, but the following is more explicit
			data = make([]byte, len(txData))
			copy(data, txData)
		}
		return nil
	})
	if err != nil {
		return false, nil
	}

	// If no value was found return false
	if data == nil {
		return false, nil
	}

	return true, s.codec.Unmarshal(data, v)
}

// Delete deletes the stored value for the given key.
// Deleting a non-existing key-value pair does NOT lead to an error.
// The key must not be "".
func (s Store) Delete(k string) error {
	if err := util.CheckKey(k); err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		if err := b.Delete([]byte(k)); err != nil {
			return err
		}

		if s.ttl == 0 {
			return nil
		}
		bTTL := tx.Bucket([]byte(TTLBucketName))
		return bTTL.Delete([]byte(s.bucketName + "-" + k))
	})
}

// Close closes the store.
// It must be called to make sure that all open transactions finish and to release all DB resources.
func (s Store) Close() error {
	return s.db.Close()
}

// Options are the options for the bbolt store.
type Options struct {
	// TTL Time-To-Live (expire time) of the data.
	// A TTL of 0 (or any negative number) means no expire time.
	// Optional (time.Hour * 24 * 7 by default).
	TTL time.Duration
	// Bucket name for storing the key-value pairs.
	// Optional ("default" by default).
	BucketName string
	// Path of the DB file.
	// Optional ("bbolt.db" by default).
	Path string
	// Encoding format.
	// Optional (encoding.JSON by default).
	Codec encoding.Codec
}

// DefaultOptions is an Options object with default values.
// BucketName: "default", Path: "bbolt.db", Codec: encoding.JSON
var DefaultOptions = Options{
	TTL:        time.Hour * 24 * 7,
	BucketName: "default",
	Path:       "bbolt.db",
	Codec:      encoding.JSON,
}

// NewStore creates a new bbolt store.
// Note: bbolt uses an exclusive write lock on the database file so it cannot be shared by multiple processes.
// So when creating multiple clients you should always use a new database file (by setting a different Path in the options).
//
// You must call the Close() method on the store when you're done working with it.
func NewStore(options Options) (Store, error) {
	result := Store{}

	// Set default values
	if options.BucketName == "" {
		options.BucketName = DefaultOptions.BucketName
	}
	if options.Path == "" {
		options.Path = DefaultOptions.Path
	}
	if options.Codec == nil {
		options.Codec = DefaultOptions.Codec
	}
	if options.TTL < 0 {
		options.TTL = 0
	}

	// Open DB
	db, err := bolt.Open(options.Path, config.Download.ModeDB.Get(), nil)
	if err != nil {
		return result, err
	}

	// Create a bucket if it doesn't exist yet.
	// In bbolt key/value pairs are stored to and read from buckets.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(options.BucketName))
		if err != nil {
			return err
		}

		if options.TTL == 0 {
			return nil
		}
		_, err = tx.CreateBucketIfNotExists([]byte(TTLBucketName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return result, err
	}

	result.db = db
	result.bucketName = options.BucketName
	result.codec = options.Codec
	result.ttl = options.TTL

	return result, nil
}
