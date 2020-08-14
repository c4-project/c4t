// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package quantity contains the quantity sets for various parts of the tester.
//
// These sets appear grouped into this package primarily for reasons of dependency cycle breaking; various bits of the
// tester at various different levels need access to them.
package quantity

import (
	"log"
	"reflect"

	"github.com/MattWindsor91/act-tester/internal/helper/stringhelp"
)

// GenericOverride substitutes any quantities in new that are non-zero for those in *old (which must be a pointer).
func GenericOverride(old, new interface{}) {
	qv := reflect.ValueOf(old).Elem()
	nv := reflect.ValueOf(new)

	nf := nv.NumField()
	for i := 0; i < nf; i++ {
		k := nv.Field(i)
		if !k.IsZero() {
			qv.Field(i).Set(k)
		}
	}
}

// LogWorkers dumps the number of workers configured by nworkers to the logger l.
func LogWorkers(l *log.Logger, nworkers int) {
	l.Println("running across", stringhelp.PluralQuantity(nworkers, "worker", "", "s"))
}
