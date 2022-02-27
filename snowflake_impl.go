package generator

import (
	"encoding/hex"
	"errors"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

const (
	numBytes = 14
)

var (
	// Epoch is set to the twitter snowflake epoch of Feb 24 2022 20:04:31 UTC+8 in milliseconds
	// You may customize this to set a different epoch for your application.
	Epoch int64 = 1645704271000

	// NodeBytes holds the number of bytes to use for Node
	NodeBytes = 2 // 16bit

	// StepBytes holds the number of bytes to use for Step
	StepBytes = 2 // 16bit

	// TimeBytes timestamp, can use until 2056-12-28 15:58:18 UTC+8
	TimeBytes = 5 // 40bit

	// PreBits predifine high value, max value is 1099511627775
	PreBytes = 5 // 40bit

	once      sync.Once
	generator *SnowFlake
)

var (
	_ Generator = generator // implement check hint to compiler

	pool = sync.Pool{ // bytes pool for generate
		New: func() interface{} {
			b := make([]byte, numBytes*2)
			return &b
		},
	}

	New = func() *SnowFlake {
		once.Do(func() {
			var err error
			generator, err = NewSnowFlake(1) // default node id, TODO: support distributed
			if err != nil {
				panic(err)
			}
		})

		return generator
	}
)

// ID An ID is a custom type used for a snowflake ID.  This is used so we can
// attach methods onto the ID.
// nocopy can be embedded in a struct to help prevent shallow copies.
// This does not rely on a Go language feature, but rather a special case
// within the vet checker.
type ID struct {
	noCopy [0]sync.Mutex
	id     string
}

func (i *ID) String() string {
	return i.id
}

// SnowFlake generator implement
// node
type SnowFlake struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	step  int64
	node  int64
	pre   int64

	stepMax int64

	preShift  int
	stepShift int
	nodeShift int
	timeShift int

	num [numBytes]byte
}

// NewSnowFlake returns a new snowflake node that can be used to generate snowflake
func NewSnowFlake(node int64) (*SnowFlake, error) {
	n := SnowFlake{}
	n.node = node
	n.stepMax = (1 << (StepBytes * 8)) - 1

	nodeMax := (1 << NodeBytes * 8) - 1
	if n.node < 0 || n.node > int64(nodeMax) {
		return nil, errors.New("Node number must be between 0 and " + strconv.Itoa(nodeMax))
	}

	n.timeShift = PreBytes
	n.nodeShift = PreBytes + TimeBytes
	n.stepShift = PreBytes + TimeBytes + NodeBytes

	curTime := time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(Epoch/1000, (Epoch%1000)*1000000).Sub(curTime))

	return &n, nil
}

func (n *SnowFlake) Generate(pre int64) ID {
	n.mu.Lock()
	n.num = [numBytes]byte{} // reset array

	now := time.Since(n.epoch).Nanoseconds() / 1000000

	if now == n.time {
		n.step = (n.step + 1) & n.stepMax

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.time = now
	n.pre = pre

	// bigEndian int64 to custom size bytes
	n.toBytes(n.preShift, PreBytes, n.pre)
	n.toBytes(n.timeShift, TimeBytes, n.time)
	n.toBytes(n.nodeShift, NodeBytes, n.node)
	n.toBytes(n.stepShift, StepBytes, n.step)

	// to hex string, id is the result variable, use bytes pool for 0 allocs/op
	id := pool.Get().(*[]byte)
	hex.Encode(*id, n.num[:])
	pool.Put(id)

	n.mu.Unlock()

	return ID{id: *(*string)(unsafe.Pointer(id))}
}

func (n *SnowFlake) toBytes(index, length int, v int64) {
	for i := 0; i < length; i++ {
		n.num[index+i] = byte(v >> ((length - i - 1) * 8))
	}
}

func (n *SnowFlake) toInt64(num []byte, index, length int) int64 {
	var v int64
	for i := 0; i < length; i++ {
		v |= int64(num[index+i]) << ((length - i - 1) * 8)
	}

	return v
}

func (n *SnowFlake) ParseString(id string) (int64, error) {
	num := make([]byte, numBytes) // num is tmp variable, 0 heap allocs/op, unnecessary use bytes pool
	if _, err := hex.Decode(num, []byte(id)); err != nil {
		return 0, err
	}

	return n.toInt64(num, n.preShift, PreBytes), nil
}
