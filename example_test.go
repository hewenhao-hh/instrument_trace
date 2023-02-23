package instrument_trace_test

import (
	trace "instrument_trace"
	"sync"
)

func a() {
	defer trace.Trace()()
	b()
}

func b() {
	defer trace.Trace()()
	c()
}

func c() {
	defer trace.Trace()()
	d()
}

func d() {
	defer trace.Trace()()
}

func ExampleTrace() {
	a()
	// Output:
	// g[00001]:    ->instrument_trace_test.a
	// g[00001]:        ->instrument_trace_test.b
	// g[00001]:            ->instrument_trace_test.c
	// g[00001]:                ->instrument_trace_test.d
	// g[00001]:                <-instrument_trace_test.d
	// g[00001]:            <-instrument_trace_test.c
	// g[00001]:        <-instrument_trace_test.b
	// g[00001]:    <-instrument_trace_test.a
}

// trace2/trace.go
func A1() {
	defer trace.Trace()()
	B1()
}

func B1() {
	defer trace.Trace()()
	C1()
}

func C1() {
	defer trace.Trace()()
	D()
}

func D() {
	defer trace.Trace()()
}

func A2() {
	defer trace.Trace()()
	B2()
}
func B2() {
	defer trace.Trace()()
	C2()
}
func C2() {
	defer trace.Trace()()
	D()
}

func ExampleTrace1() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		A2()
		wg.Done()
	}()

	A1()
	wg.Wait()
	// Output:
	// g[00001]:    ->instrument_trace_test.A1
	// g[00001]:        ->instrument_trace_test.B1
	// g[00001]:            ->instrument_trace_test.C1
	// g[00001]:                ->instrument_trace_test.D
	// g[00001]:                <-instrument_trace_test.D
	// g[00001]:            <-instrument_trace_test.C1
	// g[00001]:        <-instrument_trace_test.B1
	// g[00001]:    <-instrument_trace_test.A1
	// g[00019]:    ->instrument_trace_test.A2
	// g[00019]:        ->instrument_trace_test.B2
	// g[00019]:            ->instrument_trace_test.C2
	// g[00019]:                ->instrument_trace_test.D
	// g[00019]:                <-instrument_trace_test.D
	// g[00019]:            <-instrument_trace_test.C2
	// g[00019]:        <-instrument_trace_test.B2
	// g[00019]:    <-instrument_trace_test.A2
}
