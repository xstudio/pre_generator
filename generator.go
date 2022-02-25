// Copyright 2020 donghuili All rights reserved
// Use of this source code is governed by a BSD-style
//go:generate mockgen -source=generator.go -destination=./mocks/generator.go -package=mocks

package generator

// Generator interface
type Generator interface {
	// Generate generate number by predefined value, which has the same trend
	Generate(pre int64) ID
	// ParseString parse a generated number string, returns predefined value
	ParseString(id string) (int64, error)
}
