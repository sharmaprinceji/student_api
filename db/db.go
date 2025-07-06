package db

import (
	"github.com/sharmaprinceji/student-api/internal/config"
	"github.com/sharmaprinceji/student-api/internal/storage/sqlite"
	"github.com/sharmaprinceji/student-api/internal/storage"
)

func Mydb(cfg *config.Config) (storage.Storage, error) {
	st, err := sqlite.New(cfg) // Make sure sqlite.New accepts *config.Config too
	if err != nil {
		return nil, err
	}
	return st, nil
}