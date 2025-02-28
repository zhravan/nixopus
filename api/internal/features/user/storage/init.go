package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type UserStorage struct {
	DB  *bun.DB
	Ctx context.Context
}

func CreateNewUserStorage(db *bun.DB, ctx context.Context) *UserStorage {
	return &UserStorage{
		DB:  db,
		Ctx: ctx,
	}
}
