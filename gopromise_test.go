package gopromise

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAwait(t *testing.T) {
	p := New(func(resolve ResolveFn[int], reject RejectFn) {
		resolve(2)
	})

	d, err := p.Await()
	assert.Equal(t, d, 2)
	assert.Nil(t, err)
}

func TestAwaitAsync(t *testing.T) {
	p := New(func(resolve ResolveFn[int], reject RejectFn) {
		go func() {
			time.Sleep(100 * time.Millisecond)
			resolve(2)
		}()
	})

	d, err := p.Await()
	assert.Equal(t, d, 2)
	assert.Nil(t, err)
}

func TestThen(t *testing.T) {
	p := New(func(resolve ResolveFn[int], reject RejectFn) {
		resolve(2)
	})

	ch := make(chan int)
	Then(p, func(d int, err error) (any, error) {
		ch <- d * 2
		return nil, nil
	})

	d := <-ch
	assert.Equal(t, d, 4)
}

func TestThenMultiple(t *testing.T) {
	p := New(func(resolve ResolveFn[int], reject RejectFn) {
		resolve(2)
	})

	ch := make(chan int)
	n := 3
	for i := 0; i < n; i++ {
		Then(p, func(d int, err error) (any, error) {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			ch <- d * 2
			return nil, nil
		})
	}

	for i := 0; i < n; i++ {
		d := <-ch
		assert.Equal(t, d, 4)
	}
}

func TestThenChain(t *testing.T) {
	p1 := New(func(resolve ResolveFn[int], reject RejectFn) {
		resolve(2)
	})

	p2 := Then(p1, func(d int, err error) (int, error) {
		return d * 2, err
	})

	p3 := Then(p2, func(d int, err error) (string, error) {
		return strconv.Itoa(d), err
	})

	d1, err1 := p1.Await()
	assert.Equal(t, d1, 2)
	assert.Nil(t, err1)

	d2, err2 := p2.Await()
	assert.Equal(t, d2, 4)
	assert.Nil(t, err2)

	d3, err3 := p3.Await()
	assert.Equal(t, d3, "4")
	assert.Nil(t, err3)
}

func TestThenChainMethod(t *testing.T) {
	p := New(func(resolve ResolveFn[int], reject RejectFn) {
		resolve(2)
	}).Then(func(d int, err error) (int, error) {
		return d * 2, err
	}).Then(func(d int, err error) (int, error) {
		return d * 3, err
	})

	d, err := p.Await()
	assert.Equal(t, d, 12)
	assert.Nil(t, err)
}
