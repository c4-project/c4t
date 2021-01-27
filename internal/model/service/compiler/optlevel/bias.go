// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package optlevel

import (
	"errors"
	"fmt"
	"strings"
)

// Bias is an enumeration of biases for optimisation levels.
type Bias uint8

const (
	// BiasUnknown signifies that we don't know what the bias of this optimisation level is.
	BiasUnknown Bias = iota
	// BiasDebug signifies that an optimisation level biases towards debuggability.
	// Examples include '/Od' in MSVC (which disables optimisations), and '-Og' in GCC (which doesn't).
	BiasDebug
	// BiasSize signifies that an optimisation level biases towards size.
	// Examples include '-Os' in GCC, and '/O1' in MSVC.
	BiasSize
	// BiasSpeed signifies that an optimisation level biases towards speed.
	// Examples include '-O3' in GCC, and '/O2' in MSVC.
	BiasSpeed
	// NumBias marks the number of bias members.
	NumBias
)

var (
	// ErrBadBias occurs when we try to marshal/unmarshal a bias that doesn't exist.
	ErrBadBias = errors.New("no such bias type")

	biasStrings = [NumBias]string{
		"unknown",
		"debug",
		"size",
		"speed",
	}
)

// BiasOfString tries to get the bias corresponding to s.
func BiasOfString(s string) (Bias, error) {
	for i := BiasUnknown; i < NumBias; i++ {
		if strings.EqualFold(biasStrings[i], s) {
			return i, nil
		}
	}
	return BiasUnknown, fmt.Errorf("%w: %s", ErrBadBias, s)
}

// String converts this bias into a human-readable string.
func (b Bias) String() string {
	ts, err := b.tryString()
	if err != nil {
		return "(ERROR)"
	}
	return ts
}

// MarshalText tries to marshal this bias into text.
func (b Bias) MarshalText() ([]byte, error) {
	ts, err := b.tryString()
	if err != nil {
		return []byte{}, err
	}
	return []byte(ts), nil
}

func (b Bias) tryString() (string, error) {
	if NumBias <= b {
		return "", fmt.Errorf("%w: #%d", ErrBadBias, b)
	}
	return biasStrings[b], nil
}

// UnmarshalText tries to unmarshal text into a Bias.
func (b *Bias) UnmarshalText(text []byte) error {
	var err error
	*b, err = BiasOfString(string(text))
	return err
}
