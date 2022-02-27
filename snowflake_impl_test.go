// Copyright 2020 donghuili All rights reserved
// Use of this source code is governed by a BSD-style

// Package service usage
package generator

import (
	"testing"
	"unsafe"

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
			id1 := n.Generate(tt.args.pre)
			id := id1.String()
			pre, err := n.ParseString(id)
			assert.NotEmpty(t, id)
			assert.Nil(t, err)
			assert.Equal(t, tt.args.pre, pre)
		})
	}
}

func TestGenerateDuplicateID_Nocopy(t *testing.T) {
	var x, y ID
	for i := 0; i < 2; i++ {
		y = n.Generate(99)
		if sx, sy := x.String(), y.String(); i > 0 && sx != sy {
			t.Errorf("x(%s) & y(%s) are the same", sx, sy)
		}
		x = y
	}
}

func clone(s ID) ID {
	b := make([]byte, len(s.String()))
	copy(b, s.String())
	return ID{id: *(*string)(unsafe.Pointer(&b))}
}

func TestGenerateDuplicateID(t *testing.T) {
	var x, y ID
	for i := 0; i < 2; i++ {
		y = n.Generate(99)
		if sx, sy := x.String(), y.String(); sx == sy {
			t.Errorf("x(%s) & y(%s) are the same", sx, sy)
		}
		x = clone(y)
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
