package main

import (
	"github.com/artonge/go-gtfs"
	"github.com/hashicorp/go-memdb"
)

func importGTFS(db *memdb.MemDB, g *gtfs.GTFS) error {
	// create stop times table
	txn := db.Txn(true)

	txn.DeleteAll("stop_times", "")

	for _, stopTime := range g.StopsTimes {
		err := txn.Insert("stop_times", stopTime)
		if err != nil {
			return err
		}
	}

	for _, trip := range g.Trips {
		err := txn.Insert("trips", trip)
		if err != nil {
			return err
		}

		// delete the last stop time for each trip (highest stop_sequence),
		// because the vehicle does not depart in this trip anymore, the trip is finished
		stopTimesIt, err := txn.Get("stop_times", "trip_id", trip.ID)
		if err != nil {
			return err
		}

		var lastStopTime *gtfs.StopTime
		for obj := stopTimesIt.Next(); obj != nil; obj = stopTimesIt.Next() {
			stopTime := obj.(gtfs.StopTime)

			if lastStopTime == nil {
				lastStopTime = &stopTime
			} else if stopTime.StopSeq > lastStopTime.StopSeq {
				lastStopTime = &stopTime
			}
		}

		if lastStopTime != nil {
			err = txn.Delete("stop_times", lastStopTime)
			if err != nil {
				return err
			}
		}
	}

	for _, route := range g.Routes {
		err := txn.Insert("routes", route)
		if err != nil {
			return err
		}
	}

	for _, calendar := range g.Calendars {
		err := txn.Insert("calendars", calendar)
		if err != nil {
			return err
		}
	}

	txn.Commit()

	return nil
}
