package hive

type prefix byte

const (
	dataNodesPrefix prefix = 'd'
)

func (p prefix) Bytes() []byte {
	return []byte{byte(p)}
}
