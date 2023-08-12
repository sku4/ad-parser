package ad

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sku4/ad-parser/pkg/ad/model"
)

type CleanTnt struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Cnt    uint64 `json:"cnt"`
	After  string `json:"after"`
}

type CleanTntBody struct {
	Data map[string]any
}

func (s CleanTntBody) Parse() (*CleanTnt, error) {
	data, err := json.Marshal(s.Data)
	if err != nil {
		return nil, errors.Wrap(err, "ad clean marshal")
	}

	var tnt CleanTnt
	err = json.Unmarshal(data, &tnt)
	if err != nil {
		return nil, errors.Wrap(err, "ad clean unmarshal")
	}

	return &tnt, nil
}

type FilterTnt struct {
	Status    int                    `mapstructure:"status"`
	Code      string                 `mapstructure:"code"`
	After     string                 `mapstructure:"after"`
	Locations []*model.AdLocationTnt `mapstructure:"ads"`
}
