package runtime

import (
	"fmt"
	"reflect"
)

type GoGoRoutine struct {
	F    interface{}
	Args []interface{}
}

func NewProc(f interface{}, args ...interface{}) {
	GlobalGoGoRoutines = append(GlobalGoGoRoutines, &GoGoRoutine{F: f, Args: args})
}

var GlobalGoGoRoutines []*GoGoRoutine

func schedule() {
	for _, gg := range GlobalGoGoRoutines {
		args := []reflect.Value{}
		for _, arg := range gg.Args {
			args = append(args, reflect.ValueOf(arg))
		}
		fmt.Println("call gg")
		reflect.ValueOf(gg.F).Call(args)
	}
}

func init() {
	fmt.Println("init")
	go schedule()
}
