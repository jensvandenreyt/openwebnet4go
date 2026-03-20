package message

// What is an interface for the WHAT part of an OpenWebNet frame.
type What interface {
	Value() int
}

// WhatCommandTranslation is the prefix value for command translations.
const WhatCommandTranslation = 1000

// Dim is an interface for the DIM part of an OpenWebNet frame.
type Dim interface {
	Value() int
}
