package gopromise

import "sync"

type ResolveFn[T any] func(T)

type RejectFn func(error)

type Op[T any] func(resolve ResolveFn[T], reject RejectFn)

type Promise[T any] struct {
	op   Op[T]
	done chan struct{}

	once sync.Once
	data T
	err  error
}

func New[T any](op Op[T]) *Promise[T] {
	p := &Promise[T]{
		op:   op,
		done: make(chan struct{}),
	}
	p.run()
	return p
}

func (p *Promise[T]) run() {
	go func() {
		p.op(p.resolveFn, p.rejectFn)
	}()
}

func (p *Promise[T]) resolveFn(data T) {
	p.once.Do(func() {
		p.data = data
		close(p.done)
	})
}

func (p *Promise[T]) rejectFn(err error) {
	p.once.Do(func() {
		p.err = err
		close(p.done)
	})
}

func (p *Promise[T]) Await() (T, error) {
	<-p.done
	return p.data, p.err
}

func (p *Promise[T]) Then(fn FollowupFn[T, T]) *Promise[T] {
	return Then[T, T](p, fn)
}

type FollowupFn[A, B any] func(A, error) (B, error)

func Then[A, B any](p *Promise[A], fn FollowupFn[A, B]) *Promise[B] {
	return New[B](func(resolveB ResolveFn[B], rejectFn RejectFn) {
		data, err := p.Await()
		if err != nil {
			rejectFn(err)
			return
		}
		b, err := fn(data, err)
		if err != nil {
			rejectFn(err)
			return
		}
		resolveB(b)
	})
}
