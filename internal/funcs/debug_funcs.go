// Package funcs contains types used to identify log functions
package funcs

import (
	"github.com/hashicorp/go-set/v3"
)

// Opt is the option type for creating LogFunctions.
type Opt func(*DebugFuncs)

// WithNormalDebugFns adds fns set to normal and all
// debug functions.
func WithNormalDebugFns(fns *set.Set[Description]) Opt {
	return func(lf *DebugFuncs) {
		for fn := range fns.Items() {
			lf.AddNormal(fn)
		}
	}
}

// DebugFuncs describes all known debug logging functions.
type DebugFuncs struct {
	// requiresGuard contains found log functions that always require a guard.
	// Examples include dump functions that iterate over some collection printing debug logs
	// in a loop
	requiresGuard *set.Set[Description]
	// normal contains all "normal" debug functions found. These are often
	// simple wrappers around debug calls, functions that only contain debug logging
	normal *set.Set[Description]
	// all is the union of requiresGuard and normal
	all *set.Set[Description]
}

// NewDebugFuncs returns creates a DebugFuncs.
func NewDebugFuncs(opts ...Opt) *DebugFuncs {
	lf := &DebugFuncs{
		requiresGuard: set.New[Description](1),
		normal:        set.New[Description](1),
		all:           set.New[Description](1),
	}

	for _, opt := range opts {
		opt(lf)
	}

	return lf
}

// AddRequiresGuard adds a log function description that always requires a guard.
func (lf *DebugFuncs) AddRequiresGuard(f Description) {
	lf.all.Insert(f)
	lf.requiresGuard.Insert(f)
}

// AddNormal adds a log function description that does not alway require a guard.
func (lf *DebugFuncs) AddNormal(f Description) {
	lf.all.Insert(f)
	lf.normal.Insert(f)
}

// RequiresGuard returns the set of all known log functions that must always be guarded.
func (lf *DebugFuncs) RequiresGuard() *set.Set[Description] {
	return lf.requiresGuard.Copy()
}

// Normal returns the set of all known log functions that do not always require a guard.
func (lf *DebugFuncs) Normal() *set.Set[Description] {
	return lf.normal.Copy()
}

// All returns a set of all log functions, both normal and those that always require a guard.
func (lf *DebugFuncs) All() *set.Set[Description] {
	return lf.all.Copy()
}
