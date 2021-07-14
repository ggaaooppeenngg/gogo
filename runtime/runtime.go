package runtime

import (
	"time"
	"unsafe"
)

type GoGoBuf struct {
	SP uintptr // stack pointer
	PC uintptr // program counter
}

const (
	_Grunnable = iota
	_Grunning
	_Gdead
)

type GoGoRoutine struct {
	Sched GoGoBuf
	Stack uintptr

	status int
	goid   int

	schedlink *GoGoRoutine // next gogoroutine on queue
}

func Malloc(n int) uintptr {
	return mal(n)
}

// mmap syscall
func mmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) (unsafe.Pointer, int)

func mal(n int) uintptr {
	ptr, eno := mmap(unsafe.Pointer(nil), uintptr(n), _PROT_READ|_PROT_WRITE, _MAP_ANON|_MAP_PRIVATE, -1, 0)
	if eno != 0 {
		panic("can not allocate memory")
	}
	return uintptr(ptr)
}

var gogoid int

//go:noinline
func Gogoschedule() {
	if !gosave(&sched.curgg.Sched) {
		// 主动调度，不是 jump 过来的
		// println("in gosave at schedule")
		gogogo(&sched.gg0.Sched)
	}
	println("out of gosave at schedule")
	// jump 过来会返回到原来的函数
	// 这里会被编译器优化掉导致gosave的下一行不会跑到这里
	// 导致函数不会返回回原来的函数
	// 加上 noinline
}

func goexit() {
	sched.curgg.status = _Gdead
	//TODO: free stack memory
	Gogoschedule()
}

var globalgoid int

func NewProc(f interface{}, args ...interface{}) {
	pc := FuncPC(f)
	stack := Malloc(1024)
	sp := stack + 1024 - 4*8
	*(*uintptr)(unsafe.Pointer(sp - 8)) = FuncPC(goexit) + 1
	gogoRoutine := GoGoRoutine{}
	gogoRoutine.Sched.PC = pc
	gogoRoutine.Sched.SP = sp - 8 - 8
	gogoRoutine.Stack = stack
	globalgoid++
	gogoRoutine.goid = globalgoid
	gogoRoutine.status = _Grunnable
	ggput(&gogoRoutine)
}

func loopghead() {
	g := sched.gghead
	for g != nil {
		println("goid", g.goid)
		g = g.schedlink
	}
}

var count int

func schedule() {
	if gosave(&sched.gg0.Sched) {
		println("gogoschedule")
		// gosave 永远是 false，这个地方只能被 gogo 进来。
		curgg := sched.curgg
		switch curgg.status {
		case _Grunnable:
			panic("invalid status")
		case _Grunning:
			curgg.status = _Grunnable
			ggput(curgg)
			break
		case _Gdead:
			break
		}
	}
	// println("schedule")
	for count < 300 {
		count++
		// println("find g")
		gg := ggget()
		if gg == nil {
			println("find no g sleep")
			loopghead()
			time.Sleep(time.Second)
			continue
		}
		gg.status = _Grunning
		sched.curgg = gg
		gogogo(&gg.Sched)
	}

}

//
// gogoroutine

var GlobalGoGoRoutines map[int]*GoGoRoutine = make(map[int]*GoGoRoutine)
var GlobalGoGoID int

var NextGoGoroutine *GoGoRoutine

func gosave(gogobuf *GoGoBuf) bool

func gogogo(gogobuf *GoGoBuf)

// get a gogoroutine from the queue
func ggget() (gg *GoGoRoutine) {
	gg = sched.gghead
	if gg != nil {
		sched.gghead = gg.schedlink
		if sched.gghead == nil {
			sched.ggtail = nil
		}
		sched.ggwait--
	}
	return gg
}

// put a gogoroutine on the queue
func ggput(gg *GoGoRoutine) {
	if sched.gghead == nil {
		sched.gghead = gg
	} else {
		sched.ggtail.schedlink = gg
	}
	sched.ggtail = gg
	sched.ggwait++
	return
}

func GogoSchedule() {
}

const (
	_ENOMEM = 0xc

	_PROT_NONE  = 0x0
	_PROT_READ  = 0x1
	_PROT_WRITE = 0x2
	_PROT_EXEC  = 0x4

	_MAP_ANON    = 0x20
	_MAP_PRIVATE = 0x2
	_MAP_FIXED   = 0x10
)

var sched struct {
	gg0   *GoGoRoutine
	curgg *GoGoRoutine

	gghead *GoGoRoutine
	ggtail *GoGoRoutine
	ggwait int
}

var g struct {
}

func init() {
	sched.gg0 = &GoGoRoutine{}
	sched.curgg = &GoGoRoutine{}
	go schedule()
}
