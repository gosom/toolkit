package rollbar

import (
	"context"
	"net/http"

	"github.com/rollbar/rollbar-go"

	"github.com/gosom/toolkit/pkg/errorsext"
)

type Config struct {
	TOKEN       string
	ENVIRONMENT string
}

type ErrorReporter struct {
}

func NewErrorReporter(params Config) *ErrorReporter {
	rollbar.SetToken(params.TOKEN)
	rollbar.SetEnvironment(params.ENVIRONMENT)

	return &ErrorReporter{}
}

func (r *ErrorReporter) ReportError(_ context.Context, args ...any) {
	var reportArgs []any

	custom := map[string]any{}

	var (
		key string
		ok  bool
		val any
	)

	for i := 0; i < len(args); i += 2 {
		key, ok = args[i].(string)
		if !ok {
			key = "invalid_key"
		}

		if i+1 < len(args) {
			v := args[i+1]

			switch v := v.(type) {
			case error:
				reportArgs = append(reportArgs, v)

				if ste, ok := v.(errorsext.StackTracer); ok {
					custom["stacktrace"] = ste.Stacktrace()
				}
			case *http.Request:
				reportArgs = append(reportArgs, v)
			default:
				val = v
			}
		} else {
			val = nil
		}

		custom[key] = val
	}

	reportArgs = append(reportArgs, custom)

	rollbar.Error(reportArgs...)
}

func (r ErrorReporter) Close() {
	rollbar.Close()
}

func (r *ErrorReporter) ReportPanic(_ context.Context, args ...any) {
	rollbar.Critical(args...)
	rollbar.Wait()
}
