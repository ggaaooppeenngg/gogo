# gogo
based on the golang branch dev.go2go

## 修改 ast 支持 gogo 关键字

## 修改 checker(types/stmt.go)

## 内存分配

直接从操作系统分配，暂时没有内存管理器也没有GC。

## 修改标准库

标准库有一些栈边界的检查。但实际上gogoroutine的栈可能已经切换到别的栈上了，过不了`src/runtime/proc.go`里的`reentersyscall`的栈边界检查。
```
if _g_.syscallsp < _g_.stack.lo || _g_.stack.hi < _g_.syscallsp {
        systemstack(func() {
                print("entersyscall inconsistent ", hex(_g_.syscallsp), " [", hex(_g_.stack.lo), ",", hex(_g_.stack.hi), "]\n")
                throw("entersyscall")
        })
}
```

内建函数`println`是用的runtime内的系统调用接口没有这个检查可以过渡一下。

## 上下文切换

通过修改栈上返回地址的指针让gosave退出到其他函数里面。
根据gosave的返回确定是被切换过来的还是正常执行的，如果不是切换过来的会因为返回值是false而进入到if外面的调用，如果是切换过来的会直接进入下一条指令也就是if里面。

## 同步

## Example
```
package main

import (
	runtime "github.com/ggaaooppeenngg/gogo/runtime"
	"time"
)

//go:noline
func f1() {
	for i := 0; i < 100; i++ {
		println("vim-go", "123")
		runtime.Gogoschedule()
	}
}

//go:noline
func f2() {
	println("proc 2 run")
	for i := 0; i < 100; i++ {
		println("vim-go2", "321")
		runtime.Gogoschedule()
	}
}
func main() {
	gogo f1()
	gogo f2()
	time.Sleep(100 * time.Second)
	return
}
```
