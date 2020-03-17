// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package forward

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"
)

// Replayer coordinates reading forwarded builder-status messages from a JSON decoder and replaying them to an observer.
type Replayer struct {
	// Decoder is the decoder on which we are listening for messages to replay.
	Decoder *json.Decoder

	// Observer is the observer to which we are forwarding observations.
	Observer builder.Observer
}

// Run runs the replayer.
func (r *Replayer) Run(ctx context.Context) error {
	for {
		if err := checkClose(ctx); err != nil {
			return err
		}

		var f Forward
		if err := r.Decoder.Decode(&f); err != nil {
			// EOF is entirely expected at some point.
			if errors.Is(err, io.EOF) {
				return ctx.Err()
			}
			return fmt.Errorf("while decoding updates: %w", err)
		}

		if err := r.forwardToObs(f); err != nil {
			return fmt.Errorf("while forwarding updates: %w", err)
		}
	}
}

func (r *Replayer) forwardToObs(f Forward) error {
	switch {
	case f.BuildEnd:
		r.Observer.OnFinish()
		return nil
	case f.Error != nil:
		return f.Error
	case f.BuildStart != nil:
		r.Observer.OnStart(*f.BuildStart)
		return nil
	case f.BuildUpdate != nil:
		r.Observer.OnRequest(*f.BuildUpdate)
		return nil
	default:
		return errors.New("received forward with nothing present")
	}
}

func checkClose(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
