package models

// Blockchain model
type Blockchain struct {
	GenesisBlock Block
	Blocks       []Block
}
