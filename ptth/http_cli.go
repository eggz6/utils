package ptth

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	defHTTPClient *HTTPClient
)

func init() {
	defHTTPClient = NewHTTPClient(&http.Client{}, WithUnmarshal(json.Unmarshal))
}

type HTTPClientOption func(cli *HTTPClient)

func WithUnmarshal(f func(data []byte, v interface{}) error) HTTPClientOption {
	return func(cli *HTTPClient) {
		cli.Unmarshal = f
	}
}

type HTTPClient struct {
	http.Client
	Unmarshal func(data []byte, v interface{}) error
}

func NewHTTPClient(cli *http.Client, opts ...HTTPClientOption) *HTTPClient {
	res := &HTTPClient{Client: *cli, Unmarshal: json.Unmarshal}
	for _, opt := range opts {
		opt(res)
	}

	return res
}

// DoCtx do the req and json unmarshal the resp into ref
// ref must be ptr
func DoCtx(ctx context.Context, req *http.Request, resp interface{}) (*http.Response, error) {
	r := req.Clone(ctx)

	return Do(r, resp)
}

// Do do the req and json unmarshal the resp into ref
// ref must be ptr
func Do(req *http.Request, ref interface{}) (*http.Response, error) {
	return defHTTPClient.Do(req, ref)
}

func (h *HTTPClient) Do(req *http.Request, ref interface{}) (*http.Response, error) {
	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, e(err, HTTPErrCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, e(err, BodyReadCode)
	}

	defer func() {
		resp.Body.Close()
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}()

	if defHTTPClient.Unmarshal == nil {
		return resp, nil
	}

	err = defHTTPClient.Unmarshal(body, ref)
	if err != nil {
		return resp, e(err, UnmarshalCode)
	}

	return resp, nil
}

func e(err error, code int) *Err {
	return &Err{code: code, origin: err}
}

type Err struct {
	code   int
	origin error
}

func (e *Err) Code() int {
	return e.code
}

func (e *Err) Error() string {
	return e.origin.Error()
}

func IsReadBodyErr(e error) (*Err, bool) {
	return assertErr(e, func(err *Err) bool {
		return err.code == BodyReadCode
	})
}

func IsUnmarshalErr(e error) (*Err, bool) {
	return assertErr(e, func(err *Err) bool {
		return err.code == UnmarshalCode
	})
}

func IsHTTPErr(e error) (*Err, bool) {
	return assertErr(e, func(err *Err) bool {
		return err.code == HTTPErrCode
	})
}

func assertErr(e error, condition func(err *Err) bool) (*Err, bool) {
	res, ok := e.(*Err)
	if !ok {
		return nil, false
	}

	ok = condition(res)
	if !ok {
		return res, false
	}

	return res, true
}

const (
	HTTPErrCode   = 100000
	BodyReadCode  = 100001
	UnmarshalCode = 100002
	NotPtrCode    = 100003
)
