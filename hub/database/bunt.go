package database

import (
	"fmt"

	"github.com/kiel-live/kiel-live/hub/graph/model"
	"github.com/tidwall/buntdb"
)

type BuntDatabase struct {
	db *buntdb.DB
}

func NewBunt() Database {
	return &BuntDatabase{}
}

func (b *BuntDatabase) Open() error {
	// db, err := buntdb.Open("data.db")
	db, err := buntdb.Open(":memory:")
	if err != nil {
		return err
	}

	b.db = db

	indexes, err := db.Indexes()
	if err != nil {
		return err
	}

	if len(indexes) == 0 {
		err = db.CreateSpatialIndex("subscription_map", "subscription:map:*", buntdb.IndexRect)
		if err != nil {
			return err
		}

		// err = db.CreateSpatialIndex("pois", "poi:*:pos", buntdb.IndexRect)
		// if err != nil {
		// 	return nil, err
		// }

		// err = db.CreateIndex("stops", "stop:*:details", buntdb.IndexString)
		// if err != nil {
		// 	return nil, err
		// }

		// err = db.CreateIndex("stop_arrivals", "stop:*:arrivals", buntdb.IndexString)
		// if err != nil {
		// 	return nil, err
		// }
	}

	return nil
}

func (b *BuntDatabase) Close() error {
	return b.db.Close()
}

func (b *BuntDatabase) GetStops(opts *ListOptions) ([]*model.Stop, error) {
	return nil, nil
}

func (b *BuntDatabase) GetStop(id string) (*model.Stop, error) {
	return nil, nil
}

func (b *BuntDatabase) SetStop(stop *model.Stop) error {
	return b.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("stop:"+stop.ID, fmt.Sprintf("[%f %f]", stop.Location.Longitude, stop.Location.Latitude), nil)
		return err
	})
}

func (b *BuntDatabase) DeleteStop(id string) error {
	return b.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete("stop:" + id)
		return err
	})
}

func (b *BuntDatabase) GetVehicles(opts *ListOptions) ([]*model.Vehicle, error) {
	return nil, nil
}

func (b *BuntDatabase) GetVehicle(id string) (*model.Vehicle, error) {
	var vehicle *model.Vehicle
	err := b.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("vehicle:" + id)
		if err != nil {
			return err
		}

		vehicle = &model.Vehicle{
			ID: id,
			Location: &model.Location{
				Longitude: 0,
				Latitude:  0,
			},
		}

		fmt.Sscanf(val, "[%f %f]", &vehicle.Location.Longitude, &vehicle.Location.Latitude)
		return nil
	})
	return vehicle, err
}

func (b *BuntDatabase) SetVehicle(vehicle *model.Vehicle) error {
	return b.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("vehicle:"+vehicle.ID, fmt.Sprintf("[%f %f]", vehicle.Location.Longitude, vehicle.Location.Latitude), nil)
		return err
	})
}

func (b *BuntDatabase) DeleteVehicle(id string) error {
	return b.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete("vehicle:" + id)
		return err
	})
}
