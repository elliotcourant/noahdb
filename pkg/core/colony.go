package core

type Colony interface {
	Shards()
	Nodes()
	Tenants()
	Network()
	Settings()
	Pool()
	Schema()
	Sequences()
}
