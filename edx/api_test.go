package edx

import (
	"errors"
	"io"
	"testing"
)

func Test_Assert(t *testing.T) {
	err := Code(2).
		WithOrigin(Code(1).WithOrigin(io.EOF)).WithLabel("label").
		WithMsg("msg").WithTraceID("trace_id")

	yes := errors.Is(err, io.EOF)
	if !yes {
		t.Fatalf("test assert failed yes=%v", yes)
	}

	t.Logf("test assert yes=%v", yes)
	t.Logf("test err to string %s", err.String())
}
