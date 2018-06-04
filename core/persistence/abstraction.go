package persistence

/*
* Group of interface that is an isolation layer to
* Bolt. Allow better testing and better control.
 */

import bolt "github.com/coreos/bbolt"

// transaction for bucket creation
type bucketCreationTransaction interface {
	CreateBucketIfNotExists(name []byte) (*bolt.Bucket, error)
}
