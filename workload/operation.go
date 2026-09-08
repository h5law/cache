package workload

type OpType uint8

const (
	OpGet OpType = iota
	OpPut
	OpDelete
)

type Operation struct {
	Type OpType
	Key  uint64
}
