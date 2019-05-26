package errs

import (
	"errors"
	"log"
)

func main() {
	err := f1(1)
	log.Println(err)
}

func f1(arg int) error {
	frame := NewFrameFunc().AddArgs(FuncArgs{"arg": arg})
	arg = arg + 1

	err := f2(arg)
	if err != nil {
		return Wrap(err, frame, nil, "msg")
	}
	return nil
}

func f2(arg int) error {
	frame := NewFrameFunc().AddArgs(FuncArgs{"arg": arg})
	arg = arg + 1

	err := f3(arg)
	if err != nil {
		return Wrap(err, frame, NewCauseFunc().AddName("f3").AddArgs(FuncArgs{"arg": arg}), "msg")
	}
	return nil
}

func f3(arg int) error {
	return errors.New("error")
}
