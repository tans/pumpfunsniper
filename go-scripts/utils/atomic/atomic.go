package atomic_

import (
	"math"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
)

type Float32 struct {
	Value uint32 // Changed to uint32 for 32-bit representation
}

func (af *Float32) Load() float32 {
	return math.Float32frombits(atomic.LoadUint32(&af.Value)) // Changed to Float32frombits and LoadUint32
}

func (af *Float32) Store(value float32) {
	atomic.StoreUint32(&af.Value, math.Float32bits(value)) // Changed to StoreUint32 and Float32bits
}

type Float64 struct {
	Value uint64
}

func (af *Float64) Load() float64 {
	return math.Float64frombits(atomic.LoadUint64(&af.Value))
}

func (af *Float64) Store(value float64) {
	atomic.StoreUint64(&af.Value, math.Float64bits(value))
}

type Bool struct {
	Value int32
}

func (ab *Bool) Load() bool {
	return atomic.LoadInt32(&ab.Value) != 0
}

func (ab *Bool) Store(newValue bool) {
	var newValueInt32 int32
	if newValue {
		newValueInt32 = 1
	}
	atomic.StoreInt32(&ab.Value, newValueInt32)
}

type Uint64 struct {
	Value uint64
}

func (au *Uint64) Load() uint64 {
	return atomic.LoadUint64(&au.Value)
}

func (au *Uint64) Store(value uint64) {
	atomic.StoreUint64(&au.Value, value)
}

type Int64 struct {
	Value int64
}

func (ai *Int64) Load() int64 {
	return atomic.LoadInt64(&ai.Value)
}

func (ai *Int64) Store(value int64) {
	atomic.StoreInt64(&ai.Value, value)
}

type Int32 struct {
	Value int32
}

func (ai *Int32) Load() int32 {
	return atomic.LoadInt32(&ai.Value)
}

func (ai *Int32) Store(value int32) {
	atomic.StoreInt32(&ai.Value, value)
}

func (ai *Int32) CompareAndSwap(oldValue, newValue int32) bool {
	return atomic.CompareAndSwapInt32(&ai.Value, oldValue, newValue)
}

type BigFloat struct {
	mu    sync.Mutex
	Value *big.Float
}

// NewBigFloat returns a new BigFloat with the given initial value.
func NewBigFloat(initialValue *big.Float) *BigFloat {
	return &BigFloat{
		Value: new(big.Float).Set(initialValue),
	}
}

// Load returns a copy of the BigFloat's value.
func (ab *BigFloat) Load() *big.Float {
	ab.mu.Lock()
	defer ab.mu.Unlock()
	if ab.Value == nil {
		return nil
	}
	return new(big.Float).Set(ab.Value)
}

// Store sets the BigFloat's value to newValue.
func (ab *BigFloat) Store(newValue *big.Float) {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	if ab.Value == nil {
		ab.Value = new(big.Float)
	}

	ab.Value.Set(newValue)
}

type BigInt struct {
	mu    sync.Mutex
	Value *big.Int
}

func NewBigInt(initialValue *big.Int) *BigInt {
	return &BigInt{
		Value: new(big.Int).Set(initialValue),
	}
}

func (ab *BigInt) Add(delta *big.Int) *big.Int {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	if ab.Value == nil {
		ab.Value = new(big.Int)
	}

	ab.Value.Add(ab.Value, delta)
	return new(big.Int).Set(ab.Value)
}

func (ab *BigInt) Load() *big.Int {
	ab.mu.Lock()
	defer ab.mu.Unlock()
	if ab.Value == nil {
		return nil
	}
	return new(big.Int).Set(ab.Value)
}

func (ab *BigInt) Store(newValue *big.Int) {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	if ab.Value == nil {
		ab.Value = new(big.Int)
	}

	ab.Value.Set(newValue)
}

type Address struct {
	mu    sync.Mutex
	Value common.Address
}

func NewAddress(initialValue common.Address) *Address {
	return &Address{
		Value: initialValue,
	}
}

func (aa *Address) Load() common.Address {
	aa.mu.Lock()
	defer aa.mu.Unlock()

	return aa.Value
}

func (aa *Address) Store(newValue common.Address) {
	aa.mu.Lock()
	defer aa.mu.Unlock()

	aa.Value = newValue
}

type String struct {
	mu    sync.Mutex
	Value string
}

func NewString(initialValue string) *String {
	return &String{
		Value: initialValue,
	}
}

func (as *String) Load() string {
	as.mu.Lock()
	defer as.mu.Unlock()

	return as.Value
}

func (as *String) Store(newValue string) {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.Value = newValue
}
