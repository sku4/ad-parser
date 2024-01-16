package realt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sku4/ad-parser/model"
	"github.com/sku4/ad-parser/pkg/logger"
	"github.com/tarantool/go-tarantool/v2/datetime"
	"github.com/tarantool/go-tarantool/v2/decimal"
)

type Realt struct {
	crcTable *crc32.Table
}

func New() *Realt {
	return &Realt{
		crcTable: crc32.MakeTable(crc32.IEEE),
	}
}

const (
	profileID    = 3
	graphQLURL   = "https://bitrixcryptopayment.run/ad/realt/"
	graphQLQuery = "query searchObjects($data: GetObjectsByAddressInput!) {\n  " +
		"searchObjects(data: $data) {\n    body {\n      results {\n        location\n        " +
		"createdAt\n        updatedAt\n        price\n        buildingYear\n        " +
		"pricePerM2\n        storeys\n        storey\n        rooms\n        " +
		"images\n        areaTotal\n        areaLiving\n        areaMax\n        areaKitchen\n        " +
		"areaMin\n        areaLand\n        objectType\n        code\n        " +
		"streetName\n        address\n        houseNumber\n        buildingNumber\n        " +
		"category\n        numberOfBeds\n        toilet\n        }\n      " +
		"pagination {\n        page\n        pageSize\n        totalCount\n      }\n      " +
		"}\n    ...StatusAndErrors\n  }\n}\n\nfragment StatusAndErrors on INullResponse {\n  success\n  errors {\n    " +
		"code\n    title\n    message\n    field\n  }\n}"
	graphQLPageSize = 1000
	topLeftLat      = 53.822171699379794
	topLeftLong     = 27.36090453127423
	bottomRightLat  = 53.97823316350124
	bottomRightLong = 27.731935461117997
	urlFlatMask     = "https://realt.by/sale-flats/object/%d/"
	urlCottagesMask = "https://realt.by/sale-cottages/object/%d/"
	adsCap          = 360
	toiletO         = 0
	toilet1         = 1
	toilet2         = 2
)

var (
	categories = []string{
		"5",
		"11",
	}
)

func (r *Realt) GetCode() string {
	return "realt"
}

func (r *Realt) GetID() uint16 {
	return profileID
}

func (r *Realt) Auth(ctx context.Context) error {
	_ = ctx
	return nil
}

//nolint:gocyclo,funlen
func (r *Realt) SearchArticles(ctx context.Context, page *model.Page) ([]*model.Ad, error) {
	log := logger.Get()

	realtPage := r.getCurrentPage(page)
	categoryID, err := strconv.Atoi(categories[realtPage.CategoryID])
	if err != nil {
		return nil, errors.Wrap(err, "category id convert to int")
	}

	realtReq := &ReqGraphQL{
		OperationName: "searchObjects",
		Query:         graphQLQuery,
		Variables: ReqVariables{
			Data: ReqData{
				Pagination: ReqPagination{
					Page:     page.Num,
					PageSize: graphQLPageSize,
				},
				Where: ReqWhere{
					Category: categoryID,
					Geo: ReqGeo{
						Bbox: [][]float64{
							{topLeftLong, topLeftLat},
							{bottomRightLong, bottomRightLat},
						},
					},
				},
				Sort: []ReqSort{
					{"updatedAt", "DESC"},
				},
			},
		},
	}

	realtReqData, err := json.Marshal(realtReq)
	if err != nil {
		return nil, errors.Wrap(err, "search url marshal")
	}

	resp, err := r.request(ctx, graphQLURL, realtReqData)
	if err != nil {
		return nil, errors.Wrap(err, "search url request")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, model.ErrTooManyRequests
	}

	var realtResp *RespGraphQL
	err = json.NewDecoder(resp.Body).Decode(&realtResp)
	if err != nil {
		return nil, fmt.Errorf("body response decode %s: %w", r.GetCode(), err)
	}

	ads := make([]*model.Ad, 0, adsCap)
	for _, realtAd := range realtResp.Data.SearchObjects.Body.Results {
		var house *string
		var street *string
		h := ""
		if realtAd.HouseNumber != nil && *realtAd.HouseNumber > 0 {
			h = strconv.Itoa(*realtAd.HouseNumber)
		}
		if realtAd.BuildingNumber != nil && *realtAd.BuildingNumber != "" {
			h = fmt.Sprintf("%s/%s", h, *realtAd.BuildingNumber)
		}
		if h != "" {
			house = &h
		}
		if realtAd.StreetName != nil && *realtAd.StreetName != "" {
			sp := *realtAd.StreetName
			street = &sp
		}

		var locLat, locLong *float64
		var rooms, floor, floors *uint8
		var year *uint16
		var m2Main, m2Living, m2Kitchen *float64
		var bathroom *string

		if len(realtAd.Location) > 1 {
			locLatPoint := realtAd.Location[1]
			locLat = &locLatPoint
			locLongPoint := realtAd.Location[0]
			locLong = &locLongPoint
		}

		if realtAd.Rooms != nil && *realtAd.Rooms > 0 {
			rp := *realtAd.Rooms
			rooms = &rp
		}

		if realtAd.Storey != nil && *realtAd.Storey > 0 {
			fp := *realtAd.Storey
			floor = &fp
		}

		if realtAd.Storeys != nil && *realtAd.Storeys > 0 {
			fp := *realtAd.Storeys
			floors = &fp
		}

		if realtAd.Year != nil && *realtAd.Year > 0 {
			yp := *realtAd.Year
			year = &yp
		}

		if realtAd.AreaTotal != nil && *realtAd.AreaTotal > 0 {
			at := *realtAd.AreaTotal
			m2Main = &at
		}

		if realtAd.AreaLiving != nil && *realtAd.AreaLiving > 0 {
			al := *realtAd.AreaLiving
			m2Living = &al
		}

		if realtAd.AreaKitchen != nil && *realtAd.AreaKitchen > 0 {
			ak := *realtAd.AreaKitchen
			m2Kitchen = &ak
		}

		if realtAd.Toilet != nil {
			switch *realtAd.Toilet {
			case toiletO:
				bp := "Раздельный"
				bathroom = &bp
			case toilet1:
				bp := "Совмещенный"
				bathroom = &bp
			case toilet2:
				bp := "2 и более"
				bathroom = &bp
			}
		}

		var price *decimal.Decimal
		var errPrice error
		if realtAd.Price != nil && *realtAd.Price > 0 {
			fs := strconv.FormatFloat(*realtAd.Price, 'f', 2, 64)
			price, errPrice = decimal.NewDecimalFromString(fs)
			if errPrice != nil {
				log.Warnf("error price convert %.f to decimal: %s", *realtAd.Price, errPrice)
			}
		}

		photos := make([]string, 0, len(realtAd.Images))
		photos = append(photos, realtAd.Images...)
		if len(photos) > 0 {
			photos = photos[0:1]
		}

		link := ""
		switch categories[realtPage.CategoryID] {
		case "11":
			link = fmt.Sprintf(urlCottagesMask, realtAd.Code)
		default:
			link = fmt.Sprintf(urlFlatMask, realtAd.Code)
		}

		extID := crc32.Checksum([]byte(link), r.crcTable)

		created, errTime := datetime.NewDatetime(realtAd.CreatedAt.UTC())
		if errTime != nil {
			log.Warnf("error time convert %v to datetime: %s", realtAd.CreatedAt.UTC(), errTime)
		}

		modelAd := &model.Ad{
			ExtID:     extID,
			Created:   created,
			URL:       link,
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

	pagination := realtResp.Data.SearchObjects.Body.Pagination
	pageCount := 1
	if pagination.PageSize > 0 {
		pageCount = (pagination.TotalCount / pagination.PageSize) + 1
	}
	if pagination.PageSize == 0 || page.Num == pageCount || len(realtResp.Data.SearchObjects.Body.Results) == 0 {
		if realtPage.CategoryID == len(categories)-1 {
			return ads, model.ErrLastPage
		}
		realtPage.CategoryID++
		page.Num = 0
	}

	page.Next = realtPage

	return ads, nil
}

func (r *Realt) DownloadArticle(ctx context.Context, modelAd *model.Ad) (*model.Ad, error) {
	_ = ctx
	return modelAd, nil
}

func (r *Realt) request(ctx context.Context, url string, jsonBody []byte) (*http.Response, error) {
	var client http.Client
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error create request: %w", err)
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error request body page: %w", err)
	}

	return resp, nil
}

func (r *Realt) getCurrentPage(page *model.Page) *Page {
	realtPage := &Page{}
	if kp, ok := page.Next.(*Page); ok {
		realtPage = kp
	}

	return realtPage
}
