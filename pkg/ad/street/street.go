package street

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sku4/ad-parser/pkg/ad/model"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
)

func GetID(ctx context.Context, conn pool.Pooler, name string) (*ID, error) {
	_ = ctx

	call := tarantool.NewCallRequest("street.get_id").Args([]interface{}{name})
	resp, err := conn.Do(call, pool.RW).Get()
	if err != nil {
		return nil, err
	}

	var streetIDTnt []*ID
	err = mapstructure.Decode(resp.Data, &streetIDTnt)
	if err != nil {
		return nil, err
	}

	if len(streetIDTnt) == 0 {
		return nil, model.ErrParseResponse
	}

	return streetIDTnt[0], nil
}

func GetTypes(ctx context.Context, conn pool.Pooler) (map[uint8]*Type, error) {
	_ = ctx

	call := tarantool.NewCallRequest("street.get_types").Args([]interface{}{})
	resp, err := conn.Do(call, pool.PreferRO).Get()
	if err != nil {
		return nil, err
	}

	var streetTypesTnt []*Types
	err = mapstructure.Decode(resp.Data, &streetTypesTnt)
	if err != nil {
		return nil, err
	}

	if len(streetTypesTnt) == 0 {
		return nil, model.ErrParseResponse
	}

	tnt := make(map[uint8]*Type, len(streetTypesTnt[0].Types))
	for _, t := range streetTypesTnt[0].Types {
		tnt[t.ID] = t
	}

	return tnt, nil
}

func GetStreet(ctx context.Context, conn pool.Pooler, id uint64) (*Ext, error) {
	var streetsTnt []*model.StreetTnt
	req := tarantool.NewSelectRequest(model.SpaceStreet).
		Index(model.IndexPrimary).
		Limit(1).
		Iterator(tarantool.IterEq).
		Key(tarantool.UintKey{I: uint(id)}).
		Context(ctx)
	err := conn.Do(req, pool.PreferRW).GetTyped(&streetsTnt)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get street: primary select %d", id))
	}

	if len(streetsTnt) == 0 {
		return nil, fmt.Errorf("get street id %d: %w", id, model.ErrNotFound)
	}

	streetTnt := streetsTnt[0]

	types, err := GetTypes(ctx, conn)
	if err != nil {
		return nil, errors.Wrap(err, "get street: get types")
	}

	return &Ext{
		Street: streetTnt,
		Type:   types[streetTnt.Type],
	}, nil
}
