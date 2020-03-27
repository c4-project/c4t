// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package builder

// Observer is the interface for things that observe a builder.
type Observer interface {
	// OnBuildStart executes when the builder starts processing.
	OnBuildStart(Manifest)

	// OnBuildRequest executes when a corpus builder request is processed.
	OnBuildRequest(Request)

	// OnBuildFinish executes when the builder stops processing.
	OnBuildFinish()
}

// OnBuildStart sends an OnBuildStart message to each observer in obs.
func OnBuildStart(m Manifest, obs ...Observer) {
	for _, o := range obs {
		o.OnBuildStart(m)
	}
}

// OnBuildRequest sends an OnBuildRequest message to each observer in obs.
func OnBuildRequest(r Request, obs ...Observer) {
	for _, o := range obs {
		o.OnBuildRequest(r)
	}
}

// OnBuildStart sends an OnBuildFinish message to each observer in obs.
func OnBuildFinish(obs ...Observer) {
	for _, o := range obs {
		o.OnBuildFinish()
	}
}
