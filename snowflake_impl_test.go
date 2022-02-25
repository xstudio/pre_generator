// Copyright 2020 donghuili All rights reserved
// Use of this source code is governed by a BSD-style

// Package service usage
package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var n = New()

func TestSnowFlake_Generate(t *testing.T) {
	type args struct {
		pre int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{pre: 1234},
		},
		{
			name: "test2",
			args: args{pre: 0},
		},
		{
			name: "test3",
			args: args{pre: 1099511627775},
		},
		{
			name: "test4",
			args: args{pre: 8888888},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := n.Generate(tt.args.pre)
			pre, err := n.ParseString(id.String())
			assert.NotEmpty(t, id)
			assert.Nil(t, err)
			assert.Equal(t, tt.args.pre, pre)
		})
	}
}

func TestGenerateDuplicateID(t *testing.T) {
	var x, y ID
	for i := 0; i < 1000000; i++ {
		y = n.Generate(999)
		if x == y {
			t.Errorf("x(%s) & y(%s) are the same", x, y)
		}
		x = y
	}
}

func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = n.Generate(9999999999)
	}
}

func BenchmarkParseString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = n.ParseString("00000003e70003cd949200010000")
	}
}
