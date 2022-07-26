package index

type TermIndexBuilder struct {
	indexRoot TermIndexEntry
}

func NewTermIndexBuilder(persister AddressPersister, ingoringRegionNames []string) TermIndexBuilder {
	return TermIndexBuilder{}
}
