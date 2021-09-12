package objects

type State struct {
	HighestFinalizedBlock uint64
	HighestProcessedBlock uint64
	InSync                bool
}
