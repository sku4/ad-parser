package kufar

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	dec "github.com/shopspring/decimal"
	"github.com/sku4/ad-parser/model"
	"github.com/sku4/ad-parser/pkg/logger"
	"github.com/tarantool/go-tarantool/v2/datetime"
	"github.com/tarantool/go-tarantool/v2/decimal"
)

type Kufar struct {
	crcTable *crc32.Table
}

func New() *Kufar {
	return &Kufar{
		crcTable: crc32.MakeTable(crc32.IEEE),
	}
}

const (
	profileID = 1
	searchURL = "" +
		"https://api.kufar.by/search-api/v1/search/rendered-paginated" +
		"?cat=%s&cur=USD&cursor=%s" +
		"&gtsy=country-belarus~province-minsk~locality-minsk&lang=ru&size=200&typ=sell"
	yamsURL     = "https://yams.kufar.by/api/v1/kufar-ads/images/%s/%s.jpg?rule=list_thumbs_2x"
	rmsURL      = "https://rms.kufar.by/v1/list_thumbs_2x/%s"
	roundPlaces = 2
	roundNumber = 100
)

var (
	categories = []string{
		"1010",
		"1020",
	}
)

func (k *Kufar) GetCode() string {
	return "kufar"
}

func (k *Kufar) GetID() uint16 {
	return profileID
}

func (k *Kufar) Auth(ctx context.Context) error {
	_ = ctx
	return nil
}

//nolint:gocyclo,funlen
func (k *Kufar) SearchArticles(ctx context.Context, page *model.Page) ([]*model.Ad, error) {
	log := logger.Get()

	kufarPage := k.getCurrentPage(page)

	url := fmt.Sprintf(searchURL, categories[kufarPage.CategoryID], kufarPage.Cursor)
	resp, err := k.request(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "search url request")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, model.ErrTooManyRequests
	}

	var kufarResp Resp
	err = json.NewDecoder(resp.Body).Decode(&kufarResp)
	if err != nil {
		return nil, fmt.Errorf("body response decode %s: %w", k.GetCode(), err)
	}

	ads := make([]*model.Ad, 0)
	for _, kufarAd := range kufarResp.Ads {
		var house *string
		address := ""
		for _, param := range kufarAd.AccountParameters {
			if param.P == "address" {
				address = param.V
			}
		}
		var street *string
		streetSplit, houseSplit := k.addressSplit(address)
		if streetSplit != "" {
			street = &streetSplit
		}
		if houseSplit != "" {
			house = &houseSplit
		}

		var locLat, locLong *float64
		var rooms, floor, floors *uint8
		var year *uint16
		var m2Main, m2Living, m2Kitchen *float64
		var bathroom *string
		for _, param := range kufarAd.AdParameters {
			switch param.P {
			case "coordinates":
				if loc, ok := param.V.([]interface{}); ok && len(loc) == 2 {
					locLatPoint := loc[1].(float64)
					locLat = &locLatPoint
					locLongPoint := loc[0].(float64)
					locLong = &locLongPoint
				}
			case "rooms":
				if rs, ok := param.V.(string); ok {
					rc, errConv := strconv.ParseUint(rs, 10, 0)
					if errConv != nil {
						log.Warnf("error rooms convert %s to uint: %s", rs, errConv)
					} else {
						roomsPoint := uint8(rc)
						rooms = &roomsPoint
					}
				}
			case "floor":
				if f, ok := param.V.([]interface{}); ok && len(f) > 0 {
					floorPoint := uint8(f[0].(float64))
					floor = &floorPoint
				}
			case "re_number_floors":
				if fs, ok := param.V.(string); ok {
					fc, errConv := strconv.ParseUint(fs, 10, 0)
					if errConv != nil {
						log.Warnf("error floors convert %s to uint: %s", fs, errConv)
					} else {
						floorsPoint := uint8(fc)
						floors = &floorsPoint
					}
				}
			case "year_built":
				if y, ok := param.V.(float64); ok {
					yearPoint := uint16(y)
					year = &yearPoint
				}
			case "size":
				if m, ok := param.V.(float64); ok {
					m2Main = &m
				}
			case "size_living_space":
				if m, ok := param.V.(float64); ok {
					m2Living = &m
				}
			case "size_kitchen":
				if m, ok := param.V.(float64); ok {
					m2Kitchen = &m
				}
			case "bathroom":
				if b, ok := param.Vl.(string); ok {
					bathroom = &b
				}
			}
		}

		var price *decimal.Decimal
		var errPrice error
		if kufarAd.PriceUsd != nil && *kufarAd.PriceUsd != "" {
			price, errPrice = decimal.NewDecimalFromString(*kufarAd.PriceUsd)
			if errPrice != nil {
				log.Warnf("error price convert %s to decimal: %s", kufarAd.PriceUsd, errPrice)
			} else {
				price = decimal.NewDecimal(price.Div(dec.New(roundNumber, 0)).Round(roundPlaces))
			}
		}

		photos := make([]string, 0, len(kufarAd.Images))
		for _, i := range kufarAd.Images {
			if i.MediaStorage == "rms" {
				photos = append(photos, fmt.Sprintf(rmsURL, i.Path))
			} else if i.MediaStorage == "yams" && i.ID != "" {
				yams := []rune(i.ID)
				if len(yams) > 1 {
					photos = append(photos, fmt.Sprintf(yamsURL, string(yams[:2]), i.ID))
				}
			}
		}
		if len(photos) > 0 {
			photos = photos[0:1]
		}

		extID := crc32.Checksum([]byte(kufarAd.AdLink), k.crcTable)
		created, errTime := datetime.NewDatetime(kufarAd.ListTime.UTC())
		if errTime != nil {
			log.Warnf("error time convert %v to datetime: %s", kufarAd.ListTime, errTime)
		}

		modelAd := &model.Ad{
			ExtID:     extID,
			Created:   created,
			URL:       kufarAd.AdLink,
			Street:    street,
			House:     house,
			LocLat:    locLat,
			LocLong:   locLong,
			Price:     price,
			Rooms:     rooms,
			Floor:     floor,
			Floors:    floors,
			Year:      year,
			Photos:    photos,
			M2Main:    m2Main,
			M2Living:  m2Living,
			M2Kitchen: m2Kitchen,
			Bathroom:  bathroom,
		}
		ads = append(ads, modelAd)
	}

	next := ""
	for _, p := range kufarResp.Pagination.Pages {
		if p.Label == "next" {
			next = *p.Token
			break
		}
	}

	if next == "" {
		if kufarPage.CategoryID == len(categories)-1 {
			return ads, model.ErrLastPage
		}
		kufarPage.CategoryID++
	}

	kufarPage.Cursor = next
	page.Next = kufarPage

	return ads, nil
}

func (k *Kufar) DownloadArticle(ctx context.Context, modelAd *model.Ad) (*model.Ad, error) {
	_ = ctx
	return modelAd, nil
}

func (k *Kufar) addressSplit(address string) (street, house string) {
	street = address
	addressSplit := strings.Split(address, ",")
	if len(addressSplit) > 1 {
		street = addressSplit[0]
		house = addressSplit[1]
	} else if len(addressSplit) > 0 {
		street = addressSplit[0]
	}

	return strings.TrimSpace(street), strings.TrimSpace(house)
}

func (k *Kufar) request(ctx context.Context, url string) (*http.Response, error) {
	var client http.Client
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error request body page: %w", err)
	}

	return resp, nil
}

func (k *Kufar) getCurrentPage(page *model.Page) *Page {
	kufarPage := &Page{}
	if kp, ok := page.Next.(*Page); ok {
		kufarPage = kp
	}

	return kufarPage
}
