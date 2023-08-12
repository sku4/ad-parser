package model

import (
	"github.com/tarantool/go-tarantool/v2/datetime"
	"github.com/tarantool/go-tarantool/v2/decimal"
)

type AdTnt struct {
	ID          uint64             `mapstructure:"id" json:"id"`
	ExtID       uint32             `mapstructure:"ext_id" json:"ext_id"`
	Created     *datetime.Datetime `mapstructure:"c_time" json:"c_time"`
	Updated     *datetime.Datetime `mapstructure:"u_time" json:"u_time"`
	QueueStatus string             `mapstructure:"nq_status" json:"nq_status"`
	URL         string             `mapstructure:"url" json:"url"`
	StreetID    *uint64            `mapstructure:"street_id" json:"street_id"`
	House       *string            `mapstructure:"house" json:"house"`
	LocLat      *float64           `mapstructure:"loc_lat" json:"loc_lat"`
	LocLong     *float64           `mapstructure:"loc_long" json:"loc_long"`
	Price       *decimal.Decimal   `mapstructure:"price" json:"price"`
	PriceM2     *decimal.Decimal   `mapstructure:"price_m2" json:"price_m2"`
	Rooms       *uint8             `mapstructure:"rooms" json:"rooms"`
	Floor       *uint8             `mapstructure:"floor" json:"floor"`
	Floors      *uint8             `mapstructure:"floors" json:"floors"`
	Year        *uint16            `mapstructure:"year" json:"year"`
	Photos      []string           `mapstructure:"photos" json:"photos"`
	M2Main      *float64           `mapstructure:"m2_main" json:"m2_main"`
	M2Living    *float64           `mapstructure:"m2_living" json:"m2_living"`
	M2Kitchen   *float64           `mapstructure:"m2_kitchen" json:"m2_kitchen"`
	Bathroom    *string            `mapstructure:"bathroom" json:"bathroom"`
	Profile     uint16             `mapstructure:"profile" json:"profile"`
}

type AdLocationTnt struct {
	ID      *uint64  `mapstructure:"id" json:"id"`
	IDs     []uint64 `mapstructure:"ids" json:"ids"`
	LocLat  float64  `mapstructure:"la" json:"la"`
	LocLong float64  `mapstructure:"lo" json:"lo"`
}
