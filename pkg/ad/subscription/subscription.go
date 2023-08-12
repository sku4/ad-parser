package subscription

import (
	"context"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sku4/ad-parser/pkg/ad/model"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
)

const (
	batchLimit = 100000
)

func Filter(ctx context.Context, conn pool.Pooler, fields map[string]any) ([]int64, error) {
	var after string
	tgIds := make([]int64, 0)
	for {
		select {
		case <-ctx.Done():
			return tgIds, nil
		default:
		}

		var fnBody FilterTntBody
		call := tarantool.NewCallRequest("subscription.filter").
			Args([]interface{}{fields, batchLimit, after})
		err := conn.Do(call, pool.PreferRO).GetTyped(&fnBody)
		if err != nil {
			return nil, err
		}

		subFilterTnt, errParse := fnBody.Parse()
		if errParse != nil {
			return nil, errParse
		}

		if subFilterTnt.Status != http.StatusOK {
			return nil, errors.Wrap(model.ErrInternalServerError, subFilterTnt.Code)
		}

		after = subFilterTnt.After
		tgIds = append(tgIds, subFilterTnt.TgIds...)
		if after == "" {
			break
		}
	}

	return tgIds, nil
}

func GetByTgID(ctx context.Context, conn pool.Pooler, tgID int64, limit int, after string) (*GetByTgIDTnt, error) {
	_ = ctx

	call := tarantool.NewCallRequest("subscription.get_by_tg_id").Args([]interface{}{tgID, limit, after})
	resp, err := conn.Do(call, pool.PreferRO).Get()
	if err != nil {
		return nil, err
	}

	var subsTgIDTnt []*GetByTgIDTnt
	err = mapstructure.Decode(resp.Data, &subsTgIDTnt)
	if err != nil {
		return nil, err
	}

	if len(subsTgIDTnt) == 0 {
		return nil, model.ErrParseResponse
	}

	subTgIDTnt := subsTgIDTnt[0]

	if subTgIDTnt.Status != http.StatusOK {
		return nil, errors.Wrap(model.ErrInternalServerError, subTgIDTnt.Code)
	}

	return subTgIDTnt, nil
}
