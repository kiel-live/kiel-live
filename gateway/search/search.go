package search

import (
	"context"

	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/kiel-live/kiel-live/pkg/search"
)

type MemorySearch struct {
	vehicleSearch *search.SearchIndex[*models.Vehicle]
	stopSearch    *search.SearchIndex[*models.Stop]
	multiSearch   *search.MultiIndex
}

func NewMemorySearch() Search {
	vehicleSearch := search.NewSearchIndex(
		func(v *models.Vehicle) string { return v.ID },
		func(v *models.Vehicle) []search.SearchableEntry {
			return []search.SearchableEntry{
				{Text: v.Name, Score: 1.0},
				{Text: v.Provider, Score: 0.8},
				{Text: v.Description, Score: 0.5},
			}
		},
	)

	stopSearch := search.NewSearchIndex(
		func(s *models.Stop) string { return s.ID },
		func(s *models.Stop) []search.SearchableEntry {
			entries := []search.SearchableEntry{
				{Text: s.Name, Score: 1.0},
				{Text: s.Provider, Score: 0.8},
			}

			// TODO: does this make sense => searching a route name would return multiple stops
			for _, route := range s.Routes {
				entries = append(entries, search.SearchableEntry{
					Text:  route.Name,
					Score: 0.7,
				})
			}
			for _, vehicle := range s.Vehicles {
				entries = append(entries, search.SearchableEntry{
					Text:  vehicle.Name,
					Score: 0.7,
				})
			}
			return entries
		},
	)

	return &MemorySearch{
		vehicleSearch: vehicleSearch,
		stopSearch:    stopSearch,
		multiSearch: search.NewMultiIndex().
			AddWeightedIndex("stops", stopSearch, 2.0).
			AddWeightedIndex("vehicles", vehicleSearch, 1.0),
	}
}

func (m *MemorySearch) SetStop(_ context.Context, stop *models.Stop) error {
	m.stopSearch.Set(stop)
	return nil
}

func (m *MemorySearch) DeleteStop(_ context.Context, id string) error {
	m.stopSearch.Delete(id)
	return nil
}

func (m *MemorySearch) SetVehicle(_ context.Context, vehicle *models.Vehicle) error {
	m.vehicleSearch.Set(vehicle)
	return nil
}

func (m *MemorySearch) DeleteVehicle(_ context.Context, id string) error {
	m.vehicleSearch.Delete(id)
	return nil
}

func (m *MemorySearch) Search(_ context.Context, query string, limit int) ([]search.SearchResult, error) {
	return m.multiSearch.Search(query, limit), nil
}
