# instrument_trace
一个支持多 Goroutine 分析函数调用链的工具。
## 一. 基本方法使用 `defer trace.Trace()()` 分析函数调用链
```go
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
```
## 二. 安装及 demo 使用
1. clone：`git clone https://github.com/hewenhao-hh/instrument_trace.git` 
2. 运行 makefile 编译二进制包
2. 运行 demo：`instrument -w examples/demo/demo.go` 
## 三. 核心 code
```go
// 1. 获取 Goroutine ID
// 参考源码 $GOROOT/src/net/http/h2_bundle.go http2curGoroutineID函数
func curGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}
	return n
}
```
```go
// 2. Trace 函数
func printTrace(id uint64, name, arrow string, indent int) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "    "
	}
	fmt.Printf("g[%05d]:%s%s%s\n", id, indents, arrow, name)
}

var mu sync.Mutex
var m = make(map[uint64]int)

func Trace() func() {
	// runtime.Caller 的参数标识的是要获取的是哪一个栈帧的信息。当参数为 0 时，返回的是 Caller 函数的调用者的函数信息，在这里就是 Trace 函数。但我们需要的是 Trace 函数的调用者的信息，于是我们传入 1。
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}

	fn := runtime.FuncForPC(pc)
	name := fn.Name()
	gid := curGoroutineID()

	mu.Lock()
	indents := m[gid]
	m[gid] = indents + 1
	mu.Unlock()
	printTrace(gid, name, "->", indents+1)
	return func() {
		mu.Lock()
		indents := m[gid]
		m[gid] = indents - 1
		mu.Unlock()
		printTrace(gid, name, "<-", indents)
	}
}
```
四. 目录结构
```
$ tree -F
.
|-- LICENSE
|-- Makefile
|-- README.md
|-- cmd/
|   `-- instrument/
|       `-- main.go # instrument命令行工具的main包
|-- example_test.go # test
|-- examples/
|   `-- demo/
|       |-- demo.go.orig
|       `-- go.mod
|-- go.mod
|-- go.sum
|-- instrument.exe*
|-- instrumenter/   # 自动注入逻辑的相关结构
|   |-- ast/
|   |   `-- ast.go
|   `-- instrumenter.go
`-- trace.go        # trace 核心代码
```