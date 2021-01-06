// Copyright 2020 donghuili All rights reserved
// Use of this source code is governed by a BSD-style

// Package service usage
package service

import (
	"testing"
)

func TestSnowFlake(t *testing.T) {
	type args struct {
		pre int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "equal2",
			args: args{pre: 1001},
		},
		{
			name: "equal2",
			args: args{pre: 20032032},
		},
		{
			name: "equal3",
			args: args{pre: 20032032},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := New()
			got := n.Generate(tt.args.pre)
			t.Logf("got: %v", got)

			gotPre, gotTime, gotNode, gotStep, err := n.ParseString(got.String())
			t.Logf("parsed: %v, %v, %v, %v", gotPre, gotTime, gotNode, gotStep)
			if (err != nil) != false {
				t.Errorf("ParseString() error = %v, wantErr %v", err, false)
				return
			}
			if gotPre != tt.args.pre {
				t.Errorf("ParseString() gotPre = %v, want %v", gotPre, tt.args.pre)
			}
			if gotTime <= 0 {
				t.Errorf("ParseString() gotTime = %v, want > 0", gotTime)
			}
			if gotNode != 1 {
				t.Errorf("ParseString() gotNode = %v, want 1", gotNode)
			}
		})
	}
}
