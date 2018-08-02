package model

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type PublicModel interface {
	ToPublicModel()
}

func generateId(id []byte) ([]byte, error) {
	if id != nil && string(id) != "" {
		return nil, errors.New(fmt.Sprintf("the id %+v already defined", string(id)))
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return []byte(strconv.Itoa(r.Int())), nil
}
