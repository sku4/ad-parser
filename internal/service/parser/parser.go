package parser

import (
	"context"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/sku4/ad-parser/configs"
	"github.com/sku4/ad-parser/internal/repository"
	"github.com/sku4/ad-parser/pkg/logger"
)

//go:generate mockgen -source=parser.go -destination=mocks/parser.go

type Service struct {
	repos      *repository.Repository
	wg         *sync.WaitGroup
	mu         sync.RWMutex
	cacheClean *lru.Cache[uint16, time.Time]
}

func NewService(repos *repository.Repository) *Service {
	log := logger.Get()
	cacheClean, err := lru.New[uint16, time.Time](len(codeProfiles))
	if err != nil {
		log.Fatalf("error init lru cache: %s", err)
	}

	return &Service{
		repos:      repos,
		wg:         &sync.WaitGroup{},
		cacheClean: cacheClean,
	}
}

func (s *Service) Run(ctx context.Context) (err error) {
	log := logger.Get()
	cfg := configs.Get(ctx)

	for _, code := range cfg.Profiles {
		if _, ok := codeProfiles[code]; !ok {
			log.Errorf("Parser '%s' not found", code)
			continue
		}

		s.wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, code string) {
			defer wg.Done()

			for {
				needClean := false
				s.mu.RLock()
				codeProfile := codeProfiles[code]
				var cleanTime *time.Time
				now := time.Now()
				if t, ok := s.cacheClean.Get(codeProfile.GetID()); ok {
					cleanTime = &t
				}
				if cleanTime == nil || cleanTime.Before(now.Add(-cfg.Parser.CleanTime)) {
					s.cacheClean.Add(codeProfile.GetID(), now)
					needClean = true
				}
				s.mu.RUnlock()

				if needClean {
					log.Infof("Parser '%s' is running with clean", code)
				} else {
					log.Infof("Parser '%s' is running", code)
				}

				profile := NewProfile(s.repos, codeProfiles[code], needClean)

				if err = profile.Parse(ctx); err != nil {
					log.Errorf("Parser '%s' not might parse: %s", code, err)
				}
				log.Infof("Parser '%s' was ends of work", code)

				timer := time.NewTimer(cfg.Parser.CheckTime)
				select {
				case <-ctx.Done():
					return
				case <-timer.C:
				}
			}
		}(ctx, s.wg, code)
	}

	return nil
}

func (s *Service) Shutdown() error {
	s.wg.Wait()

	return nil
}
