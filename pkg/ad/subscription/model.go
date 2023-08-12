package subscription

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sku4/ad-parser/pkg/ad/model"
)

type FilterTnt struct {
	Status int     `json:"status"`
	Code   string  `json:"code"`
	After  string  `json:"after"`
	TgIds  []int64 `json:"tg_ids"`
}

type FilterTntBody struct {
	Data map[string]any
}

func (f FilterTntBody) Parse() (*FilterTnt, error) {
	data, err := json.Marshal(f.Data)
	if err != nil {
		return nil, errors.Wrap(err, "subscription filter marshal")
	}

	var tnt FilterTnt
	err = json.Unmarshal(data, &tnt)
	if err != nil {
		return nil, errors.Wrap(err, "subscription filter unmarshal")
	}

	return &tnt, nil
}

type GetByTgIDTnt struct {
	Status        int                      `mapstructure:"status"`
	Code          string                   `mapstructure:"code"`
	After         string                   `mapstructure:"after"`
	All           int64                    `mapstructure:"all"`
	Subscriptions []*model.SubscriptionTnt `mapstructure:"subscriptions"`
}
