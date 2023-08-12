package service

import (
	"context"

	"github.com/sku4/ad-parser/internal/repository"
	"github.com/sku4/ad-parser/internal/service/parser"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go

type Runner interface {
	Run(context.Context) error
	Shutdown() error
}

type Service struct {
	Parser Runner
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Parser: parser.NewService(repos),
	}
}
