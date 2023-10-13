package main

import (
	"github.com/tidwall/buntdb"
)

func openDatabase() (*buntdb.DB, error) {
	// db, err := buntdb.Open("data.db")
	db, err := buntdb.Open(":memory:")
	if err != nil {
		return nil, err
	}

	indexes, err := db.Indexes()
	if err != nil {
		return nil, err
	}

	if len(indexes) == 0 {
		err = db.CreateSpatialIndex("subscription_map", "subscription:map:*", buntdb.IndexRect)
		if err != nil {
			return nil, err
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

	return db, nil
}
