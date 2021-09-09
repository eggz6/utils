package edx

import "fmt"

type E uint

func Code(code E) Err {
	res := &entity{code: uint(code)}

	return res
}

type errEntity interface {
	error

	Origin() error
	Code() uint
	Label() string
	Msg() string
	TraceID() string
	String() string
}

type Err interface {
	errEntity

	WithLabel(label string) Err
	WithMsg(msg string) Err
	WithTraceID(traceID string) Err
	WithOrigin(err error) Err

	Unwrap() error
}

type entity struct {
	error

	code    uint
	label   string
	msg     string
	traceID string
}

func (err *entity) WithLabel(label string) Err {
	err.label = label

	return err
}

func (err *entity) WithMsg(msg string) Err {
	err.msg = msg

	return err
}

func (err *entity) WithTraceID(traceID string) Err {
	err.traceID = traceID

	return err
}

func (err *entity) WithOrigin(e error) Err {
	err.error = e

	return err
}

func (err *entity) Unwrap() error {
	return err.Origin()
}

func (err *entity) Origin() error {
	val, ok := (err.error).(*entity)
	if !ok {
		return err.error
	}

	return val.Origin()
}

func (err *entity) Error() string {
	return err.error.Error()
}

func (err *entity) Code() uint {
	return err.code
}

func (err *entity) Label() string {
	return err.label
}

func (err *entity) Msg() string {
	return err.msg
}

func (err *entity) TraceID() string {
	return err.traceID
}

func (err *entity) String() string {
	return fmt.Sprintf("code=%v, err=%v, label=%v, msg=%v, edx_trace_id=%v",
		err.Code(), err.Error(), err.Label(), err.Msg(), err.TraceID())
}
