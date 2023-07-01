package libgodone

type Operation int

const (
	HELP     Operation = iota
	NEW      Operation = iota
	INIT     Operation = iota
	RUN      Operation = iota
	BUILD    Operation = iota
	BUILDRUN Operation = iota
	CLEAN    Operation = iota
)
