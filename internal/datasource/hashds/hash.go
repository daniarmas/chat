package hashds

import (
	"golang.org/x/crypto/bcrypt"
)

type HashDatasource interface {
	Hash(value string) (string, error)
	CheckHash(value, hash string) bool
}

type hashDatasource struct {
}

func NewBcryptHash() HashDatasource {
	return &hashDatasource{}
}

func (ds hashDatasource) Hash(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), 10)
	return string(bytes), err
}

func (ds hashDatasource) CheckHash(value, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	return err == nil
}
