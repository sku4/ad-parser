package street

import "github.com/sku4/ad-parser/pkg/ad/model"

type ID struct {
	Status int    `mapstructure:"status"`
	Code   string `mapstructure:"code"`
	ID     uint64 `mapstructure:"id"`
}

type Types struct {
	Status int     `mapstructure:"status"`
	Code   string  `mapstructure:"code"`
	Types  []*Type `mapstructure:"types"`
}

type Type struct {
	ID      uint8    `mapstructure:"id"`
	Short   string   `mapstructure:"short"`
	Any     []string `mapstructure:"any"`
	InStart bool     `mapstructure:"in_start"`
}

type Ext struct {
	Street *model.StreetTnt
	Type   *Type
}
