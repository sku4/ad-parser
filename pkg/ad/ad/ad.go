package ad

import (
	"context"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sku4/ad-parser/pkg/ad/model"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/datetime"
	"github.com/tarantool/go-tarantool/v2/pool"
)

const (
	batchLimitClean  = 10000
	batchLimitFilter = 100000
)

func Clean(ctx context.Context, conn pool.Pooler, timeTo time.Time, profileID uint16) (uint64, error) {
	timeToTnt, err := datetime.NewDatetime(timeTo.UTC())
	if err != nil {
		return 0, err
	}

	var after string
	var cnt uint64
	for {
		select {
		case <-ctx.Done():
			return cnt, nil
		default:
		}

		var fnBody CleanTntBody
		call := tarantool.NewCallRequest("ad.clean").
			Args([]interface{}{profileID, timeToTnt, batchLimitClean, after})
		err = conn.Do(call, pool.RW).GetTyped(&fnBody)
		if err != nil {
			return 0, err
		}

		adCleanTnt, errParse := fnBody.Parse()
		if errParse != nil {
			return 0, errParse
		}

		after = adCleanTnt.After
		cnt += adCleanTnt.Cnt
		if after == "" {
			break
		}
	}

	return cnt, nil
}

func Filter(ctx context.Context, conn pool.Pooler, fields map[string]any) ([]*model.AdLocationTnt, error) {
	var after string
	locs := make([]*model.AdLocationTnt, 0)
	for {
		select {
		case <-ctx.Done():
			return locs, nil
		default:
		}

		call := tarantool.NewCallRequest("ad.filter").Args([]interface{}{fields, batchLimitFilter, after})
		resp, err := conn.Do(call, pool.PreferRO).Get()
		if err != nil {
			return nil, err
		}

		var adsFilterTnt []*FilterTnt
		err = mapstructure.Decode(resp.Data, &adsFilterTnt)
		if err != nil {
			return nil, err
		}

		if len(adsFilterTnt) == 0 {
			return nil, model.ErrParseResponse
		}
		adFilterTnt := adsFilterTnt[0]

		if adFilterTnt.Status != http.StatusOK {
			return nil, errors.Wrap(model.ErrInternalServerError, adFilterTnt.Code)
		}

		after = adFilterTnt.After
		locs = append(locs, adFilterTnt.Locations...)
		if after == "" {
			break
		}
	}

	return locs, nil
}
