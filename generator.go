// Copyright 2020 donghuili All rights reserved
// Use of this source code is governed by a BSD-style
//go:generate mockgen -source=generator.go -destination=./mocks/generator.go -package=mocks

// Package service 发号逻辑
package service

// Generator 发号器
type Generator interface {
	// Generate 使用pre值生成发号 发号值随pre值递增/减
	Generate(pre int64) ID
	// ParseString 解析发号的pre值，生成时间，节点id，毫秒内自增步长
	ParseString(id string) (pre, time, node, step int64, err error)
}
