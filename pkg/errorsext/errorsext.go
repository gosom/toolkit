package errorsext

import "runtime"

type StackTracer interface {
	Stacktrace() string
}

type stacktraceError struct {
	cause      error
	stacktrace string
}

func (s *stacktraceError) Error() string {
	return s.cause.Error()
}

func (s *stacktraceError) Stacktrace() string {
	return s.stacktrace
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}

	const buffSize = 1024 * 4
	stackBuf := make([]byte, buffSize)
	stackLen := runtime.Stack(stackBuf, false)

	ans := stacktraceError{
		cause:      err,
		stacktrace: string(stackBuf[:stackLen]),
	}

	return &ans
}
