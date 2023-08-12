package model

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tarantool/go-tarantool/v2/datetime"
	"github.com/tarantool/go-tarantool/v2/decimal"
)

type SubscriptionTnt struct {
	ID          uint64             `mapstructure:"id" json:"id"`
	TelegramID  int64              `mapstructure:"tg_id" json:"tg_id"`
	Created     *datetime.Datetime `mapstructure:"c_time" json:"c_time"`
	StreetID    *uint64            `mapstructure:"street_id" json:"street_id"`
	House       *string            `mapstructure:"house" json:"house"`
	PriceFrom   *decimal.Decimal   `mapstructure:"price_from" json:"price_from"`
	PriceTo     *decimal.Decimal   `mapstructure:"price_to" json:"price_to"`
	PriceM2From *decimal.Decimal   `mapstructure:"price_m2_from" json:"price_m2_from"`
	PriceM2To   *decimal.Decimal   `mapstructure:"price_m2_to" json:"price_m2_to"`
	RoomsFrom   *uint8             `mapstructure:"rooms_from" json:"rooms_from"`
	RoomsTo     *uint8             `mapstructure:"rooms_to" json:"rooms_to"`
	FloorFrom   *uint8             `mapstructure:"floor_from" json:"floor_from"`
	FloorTo     *uint8             `mapstructure:"floor_to" json:"floor_to"`
	YearFrom    *uint16            `mapstructure:"year_from" json:"year_from"`
	YearTo      *uint16            `mapstructure:"year_to" json:"year_to"`
	M2MainFrom  *float64           `mapstructure:"m2_main_from" json:"m2_main_from"`
	M2MainTo    *float64           `mapstructure:"m2_main_to" json:"m2_main_to"`
}

func (s SubscriptionTnt) ConvertToPutTuple() (map[string]interface{}, error) {
	adJSON, errM := json.Marshal(s)
	if errM != nil {
		return nil, errors.Wrap(errM, "convert to tuple")
	}

	var adTuple map[string]interface{}
	errM = json.Unmarshal(adJSON, &adTuple)
	if errM != nil {
		return nil, errors.Wrap(errM, "convert to tuple")
	}

	adTuple["c_time"] = s.Created
	adTuple["price_from"] = s.PriceFrom
	adTuple["price_to"] = s.PriceTo
	adTuple["price_m2_from"] = s.PriceM2From
	adTuple["price_m2_to"] = s.PriceM2To
	adTuple["id"] = nil

	return adTuple, nil
}

func (s SubscriptionTnt) ConvertToInsertTuple() []interface{} {
	var streetID *uint
	if s.StreetID != nil {
		streetUint := uint(*s.StreetID)
		streetID = &streetUint
	}

	return []interface{}{
		nil, // ID
		s.TelegramID,
		s.Created,
		streetID,
		s.House,
		s.PriceFrom,
		s.PriceTo,
		s.PriceM2From,
		s.PriceM2To,
		s.RoomsFrom,
		s.RoomsTo,
		s.FloorFrom,
		s.FloorTo,
		s.YearFrom,
		s.YearTo,
		s.M2MainFrom,
		s.M2MainTo,
	}
}

func (s SubscriptionTnt) Valid() bool {
	return !(s.StreetID == nil && s.House == nil &&
		s.PriceFrom == nil && s.PriceTo == nil &&
		s.PriceM2From == nil && s.PriceM2To == nil &&
		s.RoomsFrom == nil && s.RoomsTo == nil &&
		s.FloorFrom == nil && s.FloorTo == nil &&
		s.YearFrom == nil && s.YearTo == nil &&
		s.M2MainFrom == nil && s.M2MainTo == nil)
}
