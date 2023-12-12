package ptth

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/suite"
)

type HttpClientTestSuite struct {
	suite.Suite
	mockErr error
}

func (s *HttpClientTestSuite) SetupTest() {
	s.mockErr = errors.New("mock error")
}

func TestUtilTestSuite(t *testing.T) {
	suite.Run(t, new(HttpClientTestSuite))
}

func (s *HttpClientTestSuite) TestDo_Nil_Ref() {
	req, _ := http.NewRequest("GET", "www.google.com", nil)
	tests := []struct {
		name      string
		err       error
		want      *http.Response
		req       *http.Request
		patch     func() *gomonkey.Patches
		assertErr func(e error) (*Err, bool)
	}{
		{
			name: "success",
			req:  req,
			want: &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"hello":"jack"}`)))},
			patch: func() *gomonkey.Patches {
				res := gomonkey.NewPatches()
				res.ApplyMethod(reflect.TypeOf(&http.Client{}), "Do", func(_ *http.Client,
					_ *http.Request) (*http.Response, error) {
					return &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"hello":"jack"}`)))}, nil
				})

				return res
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.patch != nil {
				res := tt.patch()
				defer res.Reset()
			}

			resp, err := DoCtx(context.Background(), tt.req, nil)

			s.Equal(tt.err, err)
			s.Equal(tt.want, resp)

			if err != nil {
				_, ok := tt.assertErr(err)
				s.Equal(true, ok)
			}
		})
	}
}
func (s *HttpClientTestSuite) TestDoCtx() {
	req, _ := http.NewRequest("GET", "www.google.com", nil)
	tests := []struct {
		name      string
		err       error
		req       *http.Request
		want      map[string]interface{}
		patch     func() *gomonkey.Patches
		assertErr func(e error) (*Err, bool)
	}{
		{
			name: "success",
			want: map[string]interface{}{"hello": "jack"},
			req:  req,
			patch: func() *gomonkey.Patches {
				res := gomonkey.NewPatches()
				res.ApplyMethod(reflect.TypeOf(&http.Client{}), "Do", func(_ *http.Client,
					_ *http.Request) (*http.Response, error) {
					return &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"hello":"jack"}`)))}, nil
				})

				return res
			},
		},
		{
			name: "bad_http",
			req:  req,
			patch: func() *gomonkey.Patches {
				res := gomonkey.NewPatches()
				res.ApplyMethod(reflect.TypeOf(&http.Client{}), "Do", func(_ *http.Client,
					_ *http.Request) (*http.Response, error) {
					return nil, s.mockErr
				})

				return res
			},
			err:       e(s.mockErr, HTTPErrCode),
			want:      map[string]interface{}{},
			assertErr: IsHTTPErr,
		},
		{
			name: "bad_read",
			want: map[string]interface{}{},
			req:  req,
			patch: func() *gomonkey.Patches {
				res := gomonkey.NewPatches()
				res.ApplyMethod(reflect.TypeOf(&http.Client{}), "Do", func(_ *http.Client,
					_ *http.Request) (*http.Response, error) {
					return &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer([]byte(``)))}, nil
				})
				res.ApplyFunc(ioutil.ReadAll, func(r io.Reader) ([]byte, error) { return nil, s.mockErr })

				return res
			},
			err:       e(s.mockErr, BodyReadCode),
			assertErr: IsReadBodyErr,
		},
		{
			name: "bad_unmarshal",
			want: map[string]interface{}{},
			req:  req,
			patch: func() *gomonkey.Patches {
				res := gomonkey.NewPatches()
				res.ApplyMethod(reflect.TypeOf(&http.Client{}), "Do", func(_ *http.Client,
					_ *http.Request) (*http.Response, error) {
					return &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"hello":"jack"}`)))}, nil
				})

				res.ApplyFunc(defHTTPClient.Unmarshal, func(data []byte, v interface{}) error { return s.mockErr })

				return res
			},
			err:       e(s.mockErr, UnmarshalCode),
			assertErr: IsUnmarshalErr,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.patch != nil {
				res := tt.patch()
				defer res.Reset()
			}

			ref := map[string]interface{}{}
			_, err := DoCtx(context.Background(), tt.req, &ref)

			s.Equal(tt.err, err)
			s.Equal(tt.want, ref)

			if err != nil {
				_, ok := tt.assertErr(err)
				s.Equal(true, ok)
			}
		})
	}
}
