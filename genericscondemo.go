//This uses Go generics, concurrency (goroutines, channels, Context package), and anonymous functions

//we will:

/*
Run functions asychronously
Wait for multiple functions to complete
Add cancellation support to a function
Take the output of one concurrent function and use it as the input for another concurrent function
*/

package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

// Func represents any function that returns a Promise when passed to Run or Then.
//function calls could have a context.Context for the first parameter and an error for the last return value.

type Func[T, V any] func(context.Context, T) (V, error)

var (
	// ErrIncomplete is returned when GetNow is invoked and the Func associated with the Promise hasn't completed.
	ErrIncomplete = errors.New("incomplete")
)

// Promise represents a potential or actual result from running a Func.
/*
The Promise has the return value and the error, along with a done channel. The val and err fields are unexported; we want to block access to these until the values are populated.
*/
type Promise[V any] struct {
	val  V
	err  error
	done <-chan struct{}
}

// Get returns the value and the error (if any) for the Promise. Get waits until the Func associated with this
// Promise has completed. If the Func has completed, Get returns immediately.
/*
 Get the function output from the Promise. We do this using the Get method on *Promise:
*/
func (p *Promise[V]) Get() (V, error) {
	<-p.done
	return p.val, p.err
}

// GetNow returns the value and the error (if any) for the Promise. If the Func associated with this Promise has
// not completed, GetNow returns the zero value for the return type and ErrIncomplete.
func (p *Promise[V]) GetNow() (V, error) {
	select {
	case <-p.done:
		return p.val, p.err
	default:
		var zero V
		return zero, ErrIncomplete
	}
}

// Run produces a Promise for the supplied Func, evaluating the supplied context.Context and data. The Promise is returned immediately, no matter how long it takes for the Func to complete processing.

/*
Pass a context.Context, the data to process, and a function that meets the Func type into Run, and it immediately returns back a *Promise that will (eventually) contain the result of runing the supplied function with the supplied data. Unlike languages that have async built into them, there’s no special keyword that indicates that a function always runs concurrently. That means we can take any function that meets our Func type and launch it with Run. We are using anonymous functions.
*/
func Run[T, V any](ctx context.Context, t T, f Func[T, V]) *Promise[V] {
	done := make(chan struct{})
	p := Promise[V]{
		done: done,
	}
	go func() {
		defer close(done)
		p.val, p.err = f(ctx, t)
	}()
	return &p
}

/*
To build a Wait function that works when the Promise instances don't all return values of the same type, define a new, non-generic interface, Waiting, make Promise implement this interface:
*/
// Waiting defines an interface for the parameters to the Wait function.
type Waiting interface {
	Wait() error
}

// Wait allows a Promise to implement the Waiting interface. It is similar to Get, but only returns the error.
func (p *Promise[V]) Wait() error {
	<-p.done
	return p.err
}

//write Wait in terms of Waiting
// Wait takes in zero or more Waiting instances and pauses until one returns an error or all of them complete successfully.
// It returns the first error from a Waiting or nil, if no Waiting returns an error.
func Wait(ws ...Waiting) error {
	var wg sync.WaitGroup
	wg.Add(len(ws))
	errChan := make(chan error, len(ws))
	done := make(chan struct{})
	for _, w := range ws {
		go func(w Waiting) {
			defer wg.Done()
			err := w.Wait()
			if err != nil {
				errChan <- err
			}
		}(w)
	}
	go func() {
		defer close(done)
		wg.Wait()
	}()
	select {
	case err := <-errChan:
		return err
	case <-done:
	}
	return nil
}

// WithCancellation takes in a Func and returns a Func that implements the passed-in Func's behavior, but adds context cancellation.

/*
WithCancellation takes advantage of our old friend, the done channel pattern. In fact, we use two done channels. Every context has a Done method that returns a channel. This channel is closed when a context is cancelled either by a timeout or by a call to a cancel function associated with the context. We also create our own done channel. The function returned by WithCancellation invokes the passed-in Func in a goroutine. When that goroutine completes, our done channel is closed. We use a select statement to wait for either the done channel to close or for the channel returned by the context’sDone method to be closed. If ours closes first, we return the results of the passed-in Func. If the context’s done channel closes first, we return a zero value and the error from the context.
Returning that zero value relies on an interesting Go generics trick. In order to get a zero value for a generic type, you use the code var zero V . Remember, if you don’t assign a value to the variable in a var statement, the variable is set to the type’s zero value. Any time you need to get the zero value for a generic type, take advantage of this pattern.
*/
/*
 let’s do it for them by writing a function that takes in a Func and returns a Func that implements context cancellation:

generics and closures
*/
func WithCancellation[T, V any](f Func[T, V]) Func[T, V] {
	return func(ctx context.Context, t T) (V, error) {
		done := make(chan struct{})
		var val V
		var err error
		go func() {
			defer close(done)
			val, err = f(ctx, t)
		}()
		select {
		case <-ctx.Done():
			var zero V
			return zero, ctx.Err()
		case <-done:
		}
		return val, err
	}
}

// Then produces a Promise for the supplied Func, evaluating the supplied context.Context and Promise. The returned Promise is
// returned immediately, no matter how long it takes for the Func to complete processing. If the supplied Promise returns a
// non-nil error, the error is propagated to the returned Promise and the passed-in Func is not run.
//an operation that people like to do with promises: chain them together.

func Then[T, V any](ctx context.Context, p *Promise[T], f Func[T, V]) *Promise[V] {
	done := make(chan struct{})
	out := Promise[V]{
		done: done,
	}
	go func() {
		defer close(done)
		val, err := p.Get()
		if err != nil {
			out.err = err
			return
		}
		val2, err := f(ctx, val)
		out.val = val2
		out.err = err
	}()
	return &out
}

/*
Then looks a lot like Run, with a few minor differences. The first difference is that we pass in a *Promise[T] instead of a T. Second, rather than use the passed-in value directly, we call p.Get() to retrieve the value from the Promise. If the Promise contains a non-nil error, then we assign the error to the new Promise and return immediately. Otherwise, we call our new function with the value from the passed-in Promise.
*/

func main() {
	ctx := context.Background()
	p1 := Run(ctx, 10, func(ctx context.Context, i int) (int, error) {
		return i * 2, nil
	})
	p2 := Then(ctx, p1, func(ctx context.Context, i int) (string, error) {
		return strconv.Itoa(i), nil
	})
	val, err := p2.Get()
	fmt.Println(val, err)
}

/*
two new types( Func and Promise )and three functions (Run, WithCancellation, and Then
*/
