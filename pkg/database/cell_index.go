package database

import (
	"sync"

	"github.com/golang/geo/s2"
)

type CellIndex struct {
	sync.RWMutex
	index map[s2.CellID]map[string]struct{}
}

func NewCellIndex() *CellIndex {
	return &CellIndex{
		index: make(map[s2.CellID]map[string]struct{}),
	}
}

func (c *CellIndex) UpdateItem(itemID string, newIDs []s2.CellID, oldIDs []s2.CellID) {
	c.Lock()
	defer c.Unlock()

	toDelete := make(map[s2.CellID]struct{})
	for _, oldID := range oldIDs {
		toDelete[oldID] = struct{}{}
	}

	toAdd := make(map[s2.CellID]struct{})
	for _, newID := range newIDs {
		if _, ok := toDelete[newID]; ok {
			// keep item and therefore don't delete it
			delete(toDelete, newID)
		} else {
			toAdd[newID] = struct{}{}
		}
	}

	for id := range toDelete {
		if _, ok := c.index[id]; ok {
			delete(c.index[id], itemID)

			// delete empty cell
			if len(c.index[id]) == 0 {
				delete(c.index, id)
			}
		}
	}

	for id := range toAdd {
		// create new cell if it doesn't exist
		if _, ok := c.index[id]; !ok {
			c.index[id] = make(map[string]struct{})
		}
		c.index[id][itemID] = struct{}{}
	}
}

func (c *CellIndex) AddItem(itemID string, cellIDs []s2.CellID) {
	c.UpdateItem(itemID, cellIDs, nil)
}

func (c *CellIndex) RemoveItem(itemID string, cellIDs []s2.CellID) {
	c.UpdateItem(itemID, nil, cellIDs)
}

func (c *CellIndex) GetItemIDs(cellID s2.CellID) []string {
	c.RLock()
	defer c.RUnlock()

	itemIDs, ok := c.index[cellID]
	if !ok {
		return nil
	}

	ids := make([]string, 0, len(itemIDs))
	for id := range itemIDs {
		ids = append(ids, id)
	}

	return ids
}
