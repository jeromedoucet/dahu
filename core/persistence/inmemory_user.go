package persistence

import (
	"context"
	"encoding/json"
	"errors"

	bolt "github.com/coreos/bbolt"
	"github.com/jeromedoucet/dahu/core/model"
)

func (i *inMemory) GetUser(id string, ctx context.Context) (*model.User, error) {
	var user model.User
	err := i.doViewAction(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing users. The database may be corrupted !")
		}
		data := b.Get([]byte(id))
		mErr := json.Unmarshal(data, &user)
		return mErr
	})
	if err == nil {
		return &user, nil
	} else {
		return nil, err
	}
}
