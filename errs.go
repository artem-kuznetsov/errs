package errs

import (
	"encoding/json"
	"fmt"
	"runtime"
)

const (
	NonSerializableValue = "non-serializable"
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

func NewFrameFunc(args FuncArgs) FrameFunc {
	return FrameFunc{
		Name: callerName(),
		Args: args,
	}
}

func NewCauseFunc(name string, args FuncArgs) *CauseFunc {
	return &CauseFunc{
		Name: name,
		Args: args,
	}
}

func (f *FrameFunc) AddData(d FuncData) {
	if f.Data == nil {
		f.Data = d
	}
	for k, v := range d {
		f.Data[k] = v
	}
}

func (err *errorWrapper) Error() string {
	b, _ := json.Marshal(err)
	return string(b)
}

func (err *errorWrapper) addFrame() {
	err.CallStack = append(err.CallStack, frame{})
}

func (err *errorWrapper) withWrapMessage(msg string) {
	err.CallStack[len(err.CallStack)-1].wrapContext.Message = msg
}

func (err *errorWrapper) withWrapPlace() {
	err.CallStack[len(err.CallStack)-1].wrapContext.Place = callerLocation()
}

func (err *errorWrapper) withFrameFunc(f FrameFunc) {
	err.CallStack[len(err.CallStack)-1].FrameFunc = f
}

func (err *errorWrapper) withCauseFunc(f *CauseFunc) {
	err.CallStack[len(err.CallStack)-1].CauseFunc = f
}

func Wrap(err error, frameFunc FrameFunc, causeFunc *CauseFunc, msg string) error {
	ew := wrap(err)

	ew.addFrame()
	ew.withWrapMessage(msg)
	ew.withWrapPlace()
	ew.withFrameFunc(frameFunc)

	if ew.CallStack[0].CauseFunc == nil {
		ew.withCauseFunc(causeFunc)
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
