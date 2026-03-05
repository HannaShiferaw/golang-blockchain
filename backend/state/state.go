package state

import "sync"

type Actor struct {
	ID         string
	Name       string
	Role       string
	PublicKey  interface{}
	PrivateKey interface{}
}

type Batch struct {
	BatchID    string
	CoffeeType string
	Quantity   float64
	ExportDate string
	Status     string
}

var (
	actors = make(map[string]*Actor)
	batches = make(map[string]*Batch)
	mutex   = &sync.Mutex{}
)

func AddActor(a *Actor) {
	mutex.Lock()
	defer mutex.Unlock()
	actors[a.Role] = a
}

func GetActorByRole(role string) *Actor {
	mutex.Lock()
	defer mutex.Unlock()
	return actors[role]
}

func GetAllActors() []*Actor {
	mutex.Lock()
	defer mutex.Unlock()
	list := []*Actor{}
	for _, a := range actors {
		list = append(list, a)
	}
	return list
}

func GetBatch(batchID string) *Batch {
	mutex.Lock()
	defer mutex.Unlock()
	return batches[batchID]
}

func UpsertBatch(b *Batch) {
	mutex.Lock()
	defer mutex.Unlock()
	batches[b.BatchID] = b
}

func GetAllBatches() []*Batch {
	mutex.Lock()
	defer mutex.Unlock()
	list := []*Batch{}
	for _, b := range batches {
		list = append(list, b)
	}
	return list
}