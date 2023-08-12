package model

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tarantool/go-tarantool/v2/datetime"
	"github.com/tarantool/go-tarantool/v2/decimal"
)

type Ad struct {
	ExtID     uint32             `json:"ext_id"`
	Created   *datetime.Datetime `json:"c_time"`
	Updated   *datetime.Datetime `json:"u_time"`
	URL       string             `json:"url"`
	StreetID  *uint64            `json:"street_id"`
	House     *string            `json:"house"`
	LocLat    *float64           `json:"loc_lat"`
	LocLong   *float64           `json:"loc_long"`
	Price     *decimal.Decimal   `json:"price"`
	PriceM2   *decimal.Decimal   `json:"price_m2"`
	Rooms     *uint8             `json:"rooms"`
	Floor     *uint8             `json:"floor"`
	Floors    *uint8             `json:"floors"`
	Year      *uint16            `json:"year"`
	Photos    []string           `json:"photos"`
	M2Main    *float64           `json:"m2_main"`
	M2Living  *float64           `json:"m2_living"`
	M2Kitchen *float64           `json:"m2_kitchen"`
	Bathroom  *string            `json:"bathroom"`
	Profile   uint16             `json:"profile"`
	Street    *string            `json:"-"`
}

func (ad Ad) ConvertToTuple() (map[string]interface{}, error) {
	adJSON, errM := json.Marshal(ad)
	if errM != nil {
		return nil, errors.Wrap(errM, "convert to tuple")
	}

	var adTuple map[string]interface{}
	errM = json.Unmarshal(adJSON, &adTuple)
	if errM != nil {
		return nil, errors.Wrap(errM, "convert to tuple")
	}

	adTuple["c_time"] = ad.Created
	adTuple["u_time"] = ad.Updated
	adTuple["price"] = ad.Price
	adTuple["price_m2"] = ad.PriceM2

	return adTuple, nil
}
