package service

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

var (
	// Epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC in milliseconds
	// You may customize this to set a different epoch for your application.
	Epoch int64 = 1288834974657

	// NodeBits holds the number of bits to use for Node
	// Remember, you have a total 22 bits to share between Node/Step
	NodeBits uint8 = 10

	// StepBits holds the number of bits to use for Step
	// Remember, you have a total 22 bits to share between Node/Step
	StepBits uint8 = 12

	// TimeBits timestamp
	TimeBits uint8 = 41

	// PreBits predifine high value
	PreBits uint8 = 37

	// TextLength text
	TextLength int = 32
)

var once sync.Once
var generator Generator

// New init
var New = func() Generator {
	once.Do(func() {
		var err error
		generator, err = NewSnowFlake(1)
		if err != nil {
			panic(err)
		}
	})

	return generator
}

// SnowFlake A Node struct holds the basic information needed for a snowflake generator
// node
type SnowFlake struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	step  int64
	node  int64
	pre   int64

	stepMax int64

	nodeShift uint8
	timeShift uint8
	preShift  uint8

	preMask  float64
	timeMask float64
	nodeMask float64
}

// ID An ID is a custom type used for a snowflake ID.  This is used so we can
// attach methods onto the ID.
type ID string

// NewSnowFlake returns a new snowflake node that can be used to generate snowflake
// IDs
func NewSnowFlake(node int64) (*SnowFlake, error) {
	n := SnowFlake{}
	n.node = node
	n.stepMax = (1 << StepBits) - 1
	n.nodeShift = StepBits
	n.timeShift = NodeBits + StepBits
	n.preShift = NodeBits + StepBits + TimeBits

	n.preMask = math.Pow(2, float64(n.preShift))
	n.timeMask = math.Pow(2, float64(n.timeShift))
	n.nodeMask = math.Pow(2, float64(n.nodeShift))

	nodeMax := (1 << NodeBits) - 1
	if n.node < 0 || n.node > int64(nodeMax) {
		return nil, errors.New("Node number must be between 0 and " + strconv.Itoa(nodeMax))
	}

	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(Epoch/1000, (Epoch%1000)*1000000).Sub(curTime))

	return &n, nil
}

// Generate 使用pre值生成发号（32位长度字符串 不足高位补0） 发号值随pre值递增/减
func (n *SnowFlake) Generate(pre int64) ID {

	n.mu.Lock()

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

	r := ID(decimal.NewFromInt(n.pre).Mul(decimal.NewFromFloat(n.preMask)).
		Add(decimal.NewFromInt(n.time).Mul(decimal.NewFromFloat(n.timeMask))).
		Add(decimal.NewFromInt(n.node).Mul(decimal.NewFromFloat(n.nodeMask))).
		Add(decimal.NewFromInt(n.step)).String())

	n.mu.Unlock()

	return r
}

// ParseString 解析发号的pre值，生成时间，节点id，毫秒内自增步长
func (n *SnowFlake) ParseString(id string) (pre, time, node, step int64, err error) {
	var idDecimal decimal.Decimal

	idDecimal, err = decimal.NewFromString(id)
	if err != nil {
		return
	}

	preDecimal := idDecimal.Div(decimal.NewFromFloat(n.preMask))
	pre = preDecimal.IntPart()

	levelDecimal := idDecimal.Sub(decimal.NewFromInt(pre).Mul(decimal.NewFromFloat(n.preMask)))
	timeDecimal := levelDecimal.Div(decimal.NewFromFloat(n.timeMask))
	time = timeDecimal.IntPart()

	levelDecimal = levelDecimal.Sub(decimal.NewFromInt(time).Mul(decimal.NewFromFloat(n.timeMask)))
	nodeDecimal := levelDecimal.Div(decimal.NewFromFloat(n.nodeMask))
	node = nodeDecimal.IntPart()

	levelDecimal = levelDecimal.Sub(decimal.NewFromInt(node).Mul(decimal.NewFromFloat(n.nodeMask)))

	step = levelDecimal.IntPart()

	return
}

// String format
func (f ID) String() string {
	l, s := len(f), string(f)
	if l < TextLength {
		s = strings.Repeat("0", TextLength-l) + s
	}
	return s
}
