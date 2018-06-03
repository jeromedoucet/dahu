package persistence

import bolt "github.com/coreos/bbolt"

type bucketCreationTransaction interface {
	CreateBucketIfNotExists(name []byte) (*bolt.Bucket, error)
}
