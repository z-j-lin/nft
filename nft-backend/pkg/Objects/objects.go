package objects

type State struct {
	HighestProcessedBlock int64
	InSync                bool
}
