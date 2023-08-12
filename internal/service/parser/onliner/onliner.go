package onliner

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"net/http"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/sku4/ad-parser/model"
	"github.com/sku4/ad-parser/pkg/logger"
	"github.com/tarantool/go-tarantool/v2/datetime"
	"github.com/tarantool/go-tarantool/v2/decimal"
)

type Onliner struct {
	crcTable *crc32.Table
}

func New() *Onliner {
	return &Onliner{
		crcTable: crc32.MakeTable(crc32.IEEE),
	}
}

const (
	profileID = 2
	searchURL = "" +
		"https://r.onliner.by/sdapi/pk.api/search/apartments" +
		"?bounds[lb][lat]=53.822171699379794" +
		"&bounds[lb][long]=27.36090453127423" +
		"&bounds[rt][lat]=53.97823316350124" +
		"&bounds[rt][long]=27.73193546111799" +
		"&page=%d&limit=750"
)

var (
	delCities = []string{
		"Минск", "Беларусь", "Minsk",
	}
)

func (o *Onliner) GetCode() string {
	return "onliner"
}

func (o *Onliner) GetID() uint16 {
	return profileID
}

func (o *Onliner) Auth(ctx context.Context) error {
	_ = ctx
	return nil
}

func (o *Onliner) SearchArticles(ctx context.Context, page *model.Page) ([]*model.Ad, error) {
	log := logger.Get()

	url := fmt.Sprintf(searchURL, page.Num)
	resp, err := o.request(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "search url request")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, model.ErrTooManyRequests
	}

	var onlinerResp Resp
	err = json.NewDecoder(resp.Body).Decode(&onlinerResp)
	if err != nil {
		return nil, fmt.Errorf("body response decode %s: %w", o.GetCode(), err)
	}

	ads := make([]*model.Ad, 0)
	for _, onlinerAd := range onlinerResp.Apartments {
		var house *string
		address := onlinerAd.Location.Address
		if address == "" {
			address = onlinerAd.Location.UserAddress
		}
		var street *string
		streetSplit, houseSplit := o.addressSplit(address)
		if streetSplit != "" {
			street = &streetSplit
		}
		if houseSplit != "" {
			house = &houseSplit
		}

		var locLat, locLong *float64
		var rooms, floor, floors *uint8
		var m2Main, m2Living, m2Kitchen *float64

		lat := onlinerAd.Location.Latitude
		locLat = &lat
		long := onlinerAd.Location.Longitude
		locLong = &long
		roomsPoint := uint8(onlinerAd.NumberOfRooms)
		rooms = &roomsPoint
		floorPoint := uint8(onlinerAd.Floor)
		floor = &floorPoint
		floorsPoint := uint8(onlinerAd.NumberOfFloors)
		floors = &floorsPoint
		if onlinerAd.Area.Total != nil {
			m2Main = onlinerAd.Area.Total
		}
		if onlinerAd.Area.Living != nil {
			m2Living = onlinerAd.Area.Living
		}
		if onlinerAd.Area.Kitchen != nil {
			m2Kitchen = onlinerAd.Area.Kitchen
		}

		var price *decimal.Decimal
		var errPrice error
		if onlinerAd.Price.Converted.USD.Amount != nil && *onlinerAd.Price.Converted.USD.Amount != "" {
			price, errPrice = decimal.NewDecimalFromString(*onlinerAd.Price.Converted.USD.Amount)
			if errPrice != nil {
				log.Warnf("error price convert %s to decimal: %s",
					onlinerAd.Price.Converted.USD.Amount, errPrice)
			}
		}

		photos := make([]string, 0, 1)
		if onlinerAd.Photo != "" {
			photos = append(photos, onlinerAd.Photo)
		}

		extID := crc32.Checksum([]byte(onlinerAd.URL), o.crcTable)
		created, errTime := datetime.NewDatetime(onlinerAd.CreatedAt.UTC())
		if errTime != nil {
			log.Warnf("error time convert %v to datetime: %s", onlinerAd.CreatedAt, errTime)
		}

		modelAd := &model.Ad{
			ExtID:     extID,
			Created:   created,
			URL:       onlinerAd.URL,
			Street:    street,
			House:     house,
			LocLat:    locLat,
			LocLong:   locLong,
			Price:     price,
			Rooms:     rooms,
			Floor:     floor,
			Floors:    floors,
			Photos:    photos,
			M2Main:    m2Main,
			M2Living:  m2Living,
			M2Kitchen: m2Kitchen,
		}
		ads = append(ads, modelAd)
	}

	if page.Num == onlinerResp.Page.Last {
		return ads, model.ErrLastPage
	}

	return ads, nil
}

func (o *Onliner) DownloadArticle(ctx context.Context, modelAd *model.Ad) (*model.Ad, error) {
	_ = ctx
	return modelAd, nil
}

func (o *Onliner) addressSplit(address string) (street, house string) {
	street = address
	addressSplit := strings.Split(address, ",")
	addressParts := make([]string, 0, len(addressSplit))
	for _, as := range addressSplit {
		f := false
		for _, del := range delCities {
			if strings.Contains(as, del) {
				f = true
				break
			}
		}
		if !f && as != "" {
			addressParts = append(addressParts, strings.TrimSpace(as))
		}
	}

	if len(addressParts) > 1 {
		last := addressParts[len(addressParts)-1]
		if o.checkHasDigit(last) {
			street = strings.Join(addressParts[0:len(addressParts)-1], ", ")
			house = last
		} else {
			street = strings.Join(addressParts, ", ")
		}
	} else if len(addressParts) > 0 {
		streetSplit := strings.Split(addressParts[0], ".")
		if len(streetSplit) > 1 {
			last := streetSplit[len(streetSplit)-1]
			if o.checkHasDigit(last) {
				street = strings.Join(streetSplit[0:len(streetSplit)-1], ". ")
				house = last
			} else {
				street = addressParts[0]
			}
		} else {
			streetSplit = strings.Split(addressParts[0], " ")
			if len(streetSplit) > 1 {
				last := streetSplit[len(streetSplit)-1]
				if o.checkHasDigit(last) {
					street = strings.Join(streetSplit[0:len(streetSplit)-1], " ")
					house = last
				} else {
					street = addressParts[0]
				}
			} else if len(streetSplit) > 0 {
				street = addressParts[0]
			}
		}
	}

	return strings.TrimSpace(street), strings.TrimSpace(house)
}

func (o *Onliner) checkHasDigit(str string) bool {
	reg, _ := regexp.Compile(`\d+`)
	numbers := reg.FindAllString(str, -1)

	return len(numbers) > 0
}

func (o *Onliner) request(ctx context.Context, url string) (*http.Response, error) {
	var client http.Client
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error create request: %w", err)
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error request body page: %w", err)
	}

	return resp, nil
}
