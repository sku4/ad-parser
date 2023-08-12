package ad

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	dec "github.com/shopspring/decimal"
	"github.com/sku4/ad-parser/model"
	client "github.com/sku4/ad-parser/pkg/ad"
	clientModel "github.com/sku4/ad-parser/pkg/ad/model"
	"github.com/sku4/ad-parser/pkg/logger"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/datetime"
	"github.com/tarantool/go-tarantool/v2/decimal"
	"github.com/tarantool/go-tarantool/v2/pool"
)

const (
	roundPlaces = 2
)

type Ad struct {
	conn   pool.Pooler
	client *client.Client
}

func NewAd(conn pool.Pooler) *Ad {
	return &Ad{
		conn:   conn,
		client: client.NewClient(conn),
	}
}

func (ad *Ad) Put(ctx context.Context, modelAd *model.Ad, profileID uint16) error {
	var adsTnt []clientModel.AdTnt
	extIDSelect := tarantool.NewSelectRequest(clientModel.SpaceAd).
		Index(clientModel.IndexExt).
		Limit(1).
		Iterator(tarantool.IterEq).
		Key(tarantool.UintKey{I: uint(modelAd.ExtID)})
	err := ad.conn.Do(extIDSelect, pool.PreferRW).GetTyped(&adsTnt)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("put: ext_id select %d", modelAd.ExtID))
	}

	updated, err := datetime.NewDatetime(time.Now().UTC())
	if err != nil {
		return errors.Wrap(err, "put: time convert to datetime")
	}

	// get street id
	if modelAd.Street != nil && *modelAd.Street != "" {
		streetGetIDTnt, errStreet := ad.client.StreetGetID(ctx, *modelAd.Street)
		if errStreet != nil {
			return errors.Wrap(errStreet, "put: street.get_id")
		}
		if streetGetIDTnt.Status == http.StatusOK {
			modelAd.StreetID = &streetGetIDTnt.ID
		}
	}

	modelAd.Updated = updated
	modelAd.Profile = profileID

	if modelAd.Price != nil && modelAd.M2Main != nil && *modelAd.M2Main > 0 {
		m2Main := dec.NewFromFloat(*modelAd.M2Main)
		modelAd.PriceM2 = decimal.NewDecimal(modelAd.Price.Div(m2Main).Round(roundPlaces))
	}

	if len(adsTnt) > 0 {
		// if ad exists - update u_time
		operations := tarantool.NewOperations().
			Assign(clientModel.SpaceAdFieldUTime, updated).
			Assign(clientModel.SpaceAdFieldLocLat, modelAd.LocLat).
			Assign(clientModel.SpaceAdFieldLocLong, modelAd.LocLong).
			Assign(clientModel.SpaceAdFieldHouse, modelAd.House).
			Assign(clientModel.SpaceAdFieldPrice, modelAd.Price).
			Assign(clientModel.SpaceAdFieldPriceM2, modelAd.PriceM2).
			Assign(clientModel.SpaceAdFieldRooms, modelAd.Rooms).
			Assign(clientModel.SpaceAdFieldFloor, modelAd.Floor).
			Assign(clientModel.SpaceAdFieldFloors, modelAd.Floors).
			Assign(clientModel.SpaceAdFieldYear, modelAd.Year).
			Assign(clientModel.SpaceAdFieldPhotos, modelAd.Photos).
			Assign(clientModel.SpaceAdFieldM2Main, modelAd.M2Main).
			Assign(clientModel.SpaceAdFieldM2Living, modelAd.M2Living).
			Assign(clientModel.SpaceAdFieldM2Kitchen, modelAd.M2Kitchen).
			Assign(clientModel.SpaceAdFieldBathroom, modelAd.Bathroom)
		if modelAd.StreetID == nil {
			operations.Assign(clientModel.SpaceAdFieldStreetID, nil)
		} else {
			operations.Assign(clientModel.SpaceAdFieldStreetID, uint(*modelAd.StreetID))
		}
		timeUpdate := tarantool.NewUpdateRequest(clientModel.SpaceAd).
			Index(clientModel.IndexExt).
			Key(tarantool.UintKey{I: uint(modelAd.ExtID)}).
			Operations(operations)
		_, errUpd := ad.conn.Do(timeUpdate, pool.RW).Get()
		if errUpd != nil {
			return errors.Wrap(errUpd, "put: update")
		}

		return nil
	}

	// put ad
	adTuple, err := modelAd.ConvertToTuple()
	if err != nil {
		return errors.Wrap(err, "put")
	}

	callPut := tarantool.NewCallRequest("box.space.ad:put").Args([]interface{}{adTuple})
	_, err = ad.conn.Do(callPut, pool.RW).Get()
	if err != nil {
		return errors.Wrap(err, "put: call")
	}

	// send broadcast event as put new ad
	ad.conn.Do(tarantool.NewBroadcastRequest(clientModel.EventNewAd).Value(true), pool.RO)

	return nil
}

func (ad *Ad) Clean(ctx context.Context, timeTo time.Time, profileID uint16) error {
	log := logger.Get()

	// clean ads
	cntClean, err := ad.client.AdsClean(ctx, timeTo, profileID)
	if err != nil {
		return errors.Wrap(err, "clean: call")
	}

	profileCode := ad.client.ProfileGetByID(ctx, profileID)

	log.Infof("Clean %d tuples before time %s for '%s' profile",
		cntClean, timeTo.Format(time.DateTime), profileCode)

	return nil
}
