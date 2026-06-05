package errext_test

import (
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/neumachen/errext"
)

var errBench = errors.New("sentinel")

func BenchmarkNewError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = errext.NewError(errBench)
	}
}

func BenchmarkWrap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = errext.Wrap(errBench, 0)
	}
}

func BenchmarkWrapPrefix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = errext.WrapPrefix(errBench, "ctx", 0)
	}
}

func BenchmarkError(b *testing.B) {
	err := errext.WrapPrefix(errBench, "ctx", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkErrorsIs(b *testing.B) {
	err := errext.Wrap(errBench, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = errors.Is(err, errBench)
	}
}

func BenchmarkStackFramesFirstCall(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := errext.NewError(errBench)
		_ = err.StackFrames()
	}
}

func BenchmarkStackFramesRepeated(b *testing.B) {
	err := errext.NewError(errBench)
	_ = err.StackFrames()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.StackFrames()
	}
}

func BenchmarkRuntimeStack(b *testing.B) {
	err := errext.NewError(errBench)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.RuntimeStack()
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	err := errext.NewError(errBench)
	md := json.RawMessage(`{"k":"v"}`)
	if e := err.SetMetadata(&md); e != nil {
		b.Fatal(e)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, e := json.Marshal(err); e != nil {
			b.Fatal(e)
		}
	}
}

func BenchmarkLogValue(b *testing.B) {
	err := errext.NewError(errBench)
	te := err.(*errext.TraceError)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = te.LogValue()
	}
}

func BenchmarkSetMetadata(b *testing.B) {
	err := errext.NewError(errBench)
	md := json.RawMessage(`{"k":"v"}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if e := err.SetMetadata(&md); e != nil {
			b.Fatal(e)
		}
	}
}

func BenchmarkUnmarshalMetadata(b *testing.B) {
	err := errext.NewError(errBench)
	md := json.RawMessage(`{"key":"value","n":42}`)
	if e := err.SetMetadata(&md); e != nil {
		b.Fatal(e)
	}
	type out struct {
		Key string `json:"key"`
		N   int    `json:"n"`
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var o out
		if e := err.UnmarshalMetadata(&o); e != nil {
			b.Fatal(e)
		}
	}
}

// BenchmarkLogValueWithSlogHandler measures end-to-end slog throughput with
// a TraceError attribute.
func BenchmarkLogValueWithSlogHandler(b *testing.B) {
	err := errext.NewError(errBench)
	logger := slog.New(slog.NewJSONHandler(devNullWriter{}, nil))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("event", slog.Any("err", err))
	}
}

type devNullWriter struct{}

func (devNullWriter) Write(p []byte) (int, error) { return len(p), nil }
