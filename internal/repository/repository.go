package repository

import (
	"context"
	"time"

	"github.com/sku4/ad-parser/internal/repository/tarantool/ad"
	"github.com/sku4/ad-parser/model"
	"github.com/tarantool/go-tarantool/v2/pool"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

type Ad interface {
	Put(ctx context.Context, ad *model.Ad, profileID uint16) error
	Clean(ctx context.Context, timeTo time.Time, profileID uint16) error
}

type Repository struct {
	Ad
}

func NewRepository(conn pool.Pooler) *Repository {
	return &Repository{
		Ad: ad.NewAd(conn),
	}
}
