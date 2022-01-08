package datastore

type Entity interface {
}

type DataStore interface {
	Create(Entity)
}
