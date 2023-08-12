package ad

import (
	"context"
	"time"

	"github.com/sku4/ad-parser/pkg/ad/ad"
	"github.com/sku4/ad-parser/pkg/ad/model"
	"github.com/sku4/ad-parser/pkg/ad/profile"
	"github.com/sku4/ad-parser/pkg/ad/street"
	"github.com/sku4/ad-parser/pkg/ad/subscription"
	"github.com/tarantool/go-tarantool/v2/pool"
)

type Client struct {
	conn pool.Pooler
}

func NewClient(conn pool.Pooler) *Client {
	return &Client{
		conn: conn,
	}
}

func (c *Client) StreetGetID(ctx context.Context, name string) (*street.ID, error) {
	return street.GetID(ctx, c.conn, name)
}

func (c *Client) StreetGetTypes(ctx context.Context) (map[uint8]*street.Type, error) {
	return street.GetTypes(ctx, c.conn)
}

func (c *Client) StreetGet(ctx context.Context, id uint64) (*street.Ext, error) {
	return street.GetStreet(ctx, c.conn, id)
}

func (c *Client) AdsClean(ctx context.Context, timeTo time.Time, profileID uint16) (uint64, error) {
	return ad.Clean(ctx, c.conn, timeTo, profileID)
}

func (c *Client) AdFilter(ctx context.Context, fields map[string]any) ([]*model.AdLocationTnt, error) {
	return ad.Filter(ctx, c.conn, fields)
}

func (c *Client) ProfileGetByCode(ctx context.Context, code string) uint16 {
	return profile.GetByCode(ctx, code)
}

func (c *Client) ProfileGetByID(ctx context.Context, id uint16) string {
	return profile.GetByID(ctx, id)
}

func (c *Client) SubscriptionFilter(ctx context.Context, fields map[string]any) ([]int64, error) {
	return subscription.Filter(ctx, c.conn, fields)
}

func (c *Client) SubscriptionGetByTgID(ctx context.Context, tgID int64, limit int, after string) (
	*subscription.GetByTgIDTnt, error) {
	return subscription.GetByTgID(ctx, c.conn, tgID, limit, after)
}
