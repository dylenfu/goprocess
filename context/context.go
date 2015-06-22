package goprocessctx

import (
	goprocess "github.com/jbenet/goprocess"
	context "golang.org/x/net/context"
)

// WithContext constructs and returns a Process that respects
// given context. It is the equivalent of:
//
//   func ProcessWithContext(ctx context.Context) goprocess.Process {
//     p := goprocess.WithParent(goprocess.Background())
//     go func() {
//       <-ctx.Done()
//       p.Close()
//     }()
//     return p
//   }
//
func WithContext(ctx context.Context) goprocess.Process {
	if ctx == nil {
		panic("nil Context")
	}

	p := goprocess.WithParent(goprocess.Background())
	go func() {
		<-ctx.Done()
		p.Close()
	}()
	return p
}

// WithContextAndTeardown is a helper function to set teardown at initiation
// of WithContext
func WithContextAndTeardown(ctx context.Context, tf goprocess.TeardownFunc) goprocess.Process {
	if ctx == nil {
		panic("nil Context")
	}
	p := goprocess.WithTeardown(tf)
	go func() {
		<-ctx.Done()
		p.Close()
	}()
	return p
}

// WaitForContext makes p WaitFor ctx. When Closing, p waits for
// ctx.Done(), before being Closed(). It is simply:
//
//   p.WaitFor(goprocess.WithContext(ctx))
//
func WaitForContext(ctx context.Context, p goprocess.Process) {
	p.WaitFor(WithContext(ctx))
}

// WithProcessClosing returns a context.Context derived from ctx that
// is cancelled as p is Closing (after: <-p.Closing()). It is simply:
//
//   func WithProcessClosing(ctx context.Context, p goprocess.Process) context.Context {
//     ctx, cancel := context.WithCancel(ctx)
//     go func() {
//       <-p.Closing()
//       cancel()
//     }()
//     return ctx
//   }
//
func WithProcessClosing(ctx context.Context, p goprocess.Process) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-p.Closing()
		cancel()
	}()
	return ctx
}

// WithProcessClosed returns a context.Context that is cancelled
// after Process p is Closed. It is the equivalent of:
//
//   func WithProcessClosed(ctx context.Context, p goprocess.Process) context.Context {
//     ctx, cancel := context.WithCancel(ctx)
//     go func() {
//       <-p.Closed()
//       cancel()
//     }()
//     return ctx
//   }
//
func WithProcessClosed(ctx context.Context, p goprocess.Process) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-p.Closed()
		cancel()
	}()
	return ctx
}
