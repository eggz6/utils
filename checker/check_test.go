package checker

import (
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/suite"
)

type UtilSuite struct {
	suite.Suite
	mockErr error
}

func (s *UtilSuite) SetupTest() {
	s.mockErr = errors.New("mock error")
}

func TestUtilTestSuite(t *testing.T) {
	suite.Run(t, new(UtilSuite))
}

func (s *UtilSuite) TestCheckYes() {
	cases := []struct {
		name  string
		err   error
		patch func() *gomonkey.Patches
		check func() error
	}{
		{
			name: "success",
			check: func() error {
				return Check().NoEmptyString("app", "val").NoZero("age", 1).Invalid("line", func() bool { return true }).Yes()
			},
		},
		{
			name: "first_error",
			check: func() error {
				return Check().NoEmptyString("app", "").NoZero("age", 1).Invalid("line", func() bool { return true }).Yes()
			},
			err: errors.New("empty string app"),
		},
		{
			name: "empty_string",
			check: func() error {
				return Check().NoEmptyString("app", "").Yes()
			},
			err: errors.New("empty string app"),
		},
		{
			name: "zero",
			check: func() error {
				return Check().NoZero("app", 0).Yes()
			},
			err: errors.New("app is zero"),
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			if c.patch != nil {
				p := c.patch()

				defer p.Reset()
			}

			err := c.check()
			s.Equal(c.err, err)
		})
	}
}
