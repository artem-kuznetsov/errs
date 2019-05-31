package errs

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type errorWrapper struct {
	Err       string  `json:"error"`
	CallStack []frame `json:"callstack"`
}

type frame struct {
	wrapContext
	FrameFunc
	CauseFunc *CauseFunc `json:"error_cause,omitempty"`
}

type wrapContext struct {
	Message string `json:"wrap_message,omitempty"`
	Place   string `json:"wrap_place"`
}

type FrameFunc struct {
	Name string   `json:"func_name"`
	Args FuncArgs `json:"func_args,omitempty"`
	Data FuncData `json:"func_data,omitempty"`
}

type CauseFunc struct {
	Name string   `json:"func_name"`
	Args FuncArgs `json:"func_args,omitempty"`
}

type FuncArgs map[string]interface{}

type FuncData map[string]interface{}

//--------------------------------------------------

func NewFrameFunc() FrameFunc {
	return FrameFunc{
		Name: callerName(),
	}
}

func (f FrameFunc) AddArgs(a FuncArgs) FrameFunc {
	f.Args = a
	return f
}

func (f FrameFunc) AddData(d FuncData) FrameFunc {
	if f.Data == nil {
		f.Data = d
		return f
	}
	for k, v := range d {
		f.Data[k] = v
	}
	return f
}

//--------------------------------------------------

func NewCauseFunc() CauseFunc {
	return CauseFunc{}
}

func (f CauseFunc) AddName(n string) CauseFunc {
	f.Name = n
	return f
}

func (f CauseFunc) AddArgs(a FuncArgs) CauseFunc {
	f.Args = a
	return f
}

//--------------------------------------------------

func (err *errorWrapper) Error() string {
	b, _ := json.Marshal(err)
	return fmt.Sprint(string(b))
}

func (err *errorWrapper) addFrame() *errorWrapper {
	err.CallStack = append(err.CallStack, frame{})
	return err
}

func (err *errorWrapper) withWrapMessage(msg string) *errorWrapper {
	err.CallStack[len(err.CallStack)-1].wrapContext.Message = msg
	return err
}

func (err *errorWrapper) withWrapPlace() *errorWrapper {
	err.CallStack[len(err.CallStack)-1].wrapContext.Place = callerLocation()
	return err
}

func (err *errorWrapper) withFrameFunc(f FrameFunc) *errorWrapper {
	err.CallStack[len(err.CallStack)-1].FrameFunc = f
	return err
}

func (err *errorWrapper) withCauseFunc(f *CauseFunc) *errorWrapper {
	err.CallStack[len(err.CallStack)-1].CauseFunc = f
	return err
}

//--------------------------------------------------

func Wrap(err error, frameFunc FrameFunc, causeFunc *CauseFunc, msg string) error {
	ew := wrap(err).addFrame().
		withWrapMessage(msg).
		withWrapPlace().
		withFrameFunc(frameFunc)

	if ew.CallStack[0].CauseFunc == nil {
		ew = ew.withCauseFunc(causeFunc)
	}
	return ew
}

func wrap(err error) *errorWrapper {
	if ew, ok := err.(*errorWrapper); ok {
		return ew
	}
	return &errorWrapper{
		Err: err.Error(),
	}
}

func callerName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

func callerLocation() string {
	_, file, line, _ := runtime.Caller(3)
	return fmt.Sprintf("%s:%d", file, line)
}
