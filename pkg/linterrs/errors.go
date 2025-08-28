// Package linterrs contains common errors encountered
package linterrs

import "errors"

var (
	// ErrRequirementsFailed is returned when required analyzer failed.
	ErrRequirementsFailed = errors.New("prerequisite analyzer failed")
)
