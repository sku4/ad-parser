package parser

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sku4/ad-parser/configs"
	"github.com/sku4/ad-parser/internal/repository"
	"github.com/sku4/ad-parser/model"
	"github.com/sku4/ad-parser/pkg/logger"
)

type iProfile interface {
	Auth(context.Context) error
	SearchArticles(ctx context.Context, page *model.Page) (ads []*model.Ad, err error)
	DownloadArticle(ctx context.Context, ad *model.Ad) (*model.Ad, error)
	GetCode() string
	GetID() uint16
}

const (
	chanBufferLen = 10000
	timeSleep     = time.Second * 10
)

type Profile struct {
	iProfile
	repos           *repository.Repository
	urlsChan        chan *model.Ad
	adChan          chan *model.Ad
	rwMutex         *sync.RWMutex
	tooManyReqLimit int
	searchCount     int
	saveCount       int
	checkLastPage   bool
	needClean       bool
}

func NewProfile(repos *repository.Repository, profile iProfile, needClean bool) *Profile {
	return &Profile{
		repos:     repos,
		iProfile:  profile,
		urlsChan:  make(chan *model.Ad, chanBufferLen),
		adChan:    make(chan *model.Ad, chanBufferLen),
		rwMutex:   &sync.RWMutex{},
		needClean: needClean,
	}
}

func (p *Profile) Parse(ctx context.Context) (err error) {
	wg := &sync.WaitGroup{}
	cfg := configs.Get(ctx)

	p.tooManyReqLimit = cfg.Parser.TooManyReqLimit
	start := time.Now()

	// auth
	if err = p.iProfile.Auth(ctx); err != nil {
		return model.ErrProfileNotMightAuth
	}

	// search new articles
	wg.Add(1)
	go func() {
		p.searchArticles(ctx, wg)
	}()

	// download articles
	wg.Add(cfg.Parser.DownloadWorkerCount)
	for i := 0; i < cfg.Parser.DownloadWorkerCount; i++ {
		go func() {
			p.downloadArticles(ctx, wg)
		}()
	}

	// save to db
	wgs := &sync.WaitGroup{}
	wgs.Add(1)
	go func() {
		p.saveArticles(ctx, wgs)
	}()

	wg.Wait()
	close(p.adChan)
	wgs.Wait()

	if p.needClean && p.checkLastPage && p.searchCount == p.saveCount && p.searchCount > 0 {
		p.cleanArticles(ctx, start)
	}

	return nil
}

func (p *Profile) searchArticles(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(p.urlsChan)

	log := logger.Get()

	page := &model.Page{
		Num: 1,
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		p.rwMutex.RLock()
		if p.tooManyReqLimit <= 0 {
			p.rwMutex.RUnlock()
			return
		}
		p.rwMutex.RUnlock()

		urls, err := p.iProfile.SearchArticles(ctx, page)
		if err != nil && !errors.Is(err, model.ErrLastPage) && !errors.Is(err, model.ErrTooManyRequests) {
			log.Errorf("Search articles page num %d error: %s", page.Num, err)
			time.Sleep(timeSleep)
		}

		for _, url := range urls {
			p.urlsChan <- url
		}

		p.rwMutex.Lock()
		p.searchCount += len(urls)
		p.rwMutex.Unlock()

		if errors.Is(err, model.ErrTooManyRequests) {
			p.rwMutex.Lock()
			p.tooManyReqLimit--
			p.rwMutex.Unlock()
			log.Warnf("Search articles too many requests: page num %d", page.Num)
			time.Sleep(timeSleep)
		}

		if errors.Is(err, model.ErrLastPage) {
			p.rwMutex.Lock()
			p.checkLastPage = true
			p.rwMutex.Unlock()
			break
		}

		if err == nil {
			page.Num++
		}
	}
}

func (p *Profile) downloadArticles(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	log := logger.Get()

	for ad := range p.urlsChan {
		select {
		case <-ctx.Done():
			return
		default:
		}

		p.rwMutex.RLock()
		if p.tooManyReqLimit <= 0 {
			p.rwMutex.RUnlock()
			return
		}
		p.rwMutex.RUnlock()

		modelAd, err := p.iProfile.DownloadArticle(ctx, ad)
		if err != nil && !errors.Is(err, model.ErrTooManyRequests) {
			log.Errorf("Download article (%s) error: %s", ad.URL, err)
		}

		if errors.Is(err, model.ErrTooManyRequests) {
			p.rwMutex.Lock()
			p.tooManyReqLimit--
			p.rwMutex.Unlock()
			log.Warnf("Download article too many requests (%s)", ad.URL)
		}

		if modelAd != nil {
			p.adChan <- modelAd
		}
	}
}

func (p *Profile) saveArticles(ctx context.Context, wgs *sync.WaitGroup) {
	defer wgs.Done()
	log := logger.Get()

	successCnt := 0
	profileID := p.iProfile.GetID()
	for ad := range p.adChan {
		err := p.repos.Ad.Put(ctx, ad, profileID)
		if err != nil {
			log.Errorf("Save articles error: %s", err)
		} else {
			successCnt++
		}
	}

	p.rwMutex.Lock()
	p.saveCount += successCnt
	p.rwMutex.Unlock()
}

func (p *Profile) cleanArticles(ctx context.Context, timeStart time.Time) {
	log := logger.Get()

	err := p.repos.Ad.Clean(ctx, timeStart, p.iProfile.GetID())
	if err != nil {
		log.Errorf("Clean articles error: %s", err)
	}
}
