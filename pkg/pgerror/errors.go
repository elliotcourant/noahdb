package pgerror

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/pkg/errors"
)

var _ error = &Error{}

func (pg *Error) Error() string {
	return pg.Message
}

// FullError can be used when the hint and/or detail are to be tested.
func FullError(err error) string {
	var errString string
	if pqErr, ok := err.(*pq.Error); ok {
		errString = formatMsgHintDetail("pq: ", pqErr.Message, pqErr.Hint, pqErr.Detail)
	} else if pg, ok := GetPGCause(err); ok {
		errString = formatMsgHintDetail("", err.Error(), pg.Hint, pg.Detail)
	} else {
		errString = err.Error()
	}
	return errString
}

func formatMsgHintDetail(prefix, msg, hint, detail string) string {
	b := make([]string, 2)
	b[0] = prefix
	b[1] = msg
	if hint != "" {
		b = append(b, "\nHINT: ")
		b = append(b, hint)
	}
	if detail != "" {
		b = append(b, "\nDETAIL: ")
		b = append(b, detail)
	}
	return strings.Join(b, " ")

	// var b strings.Builder
	// b.WriteString(prefix)
	// b.WriteString(msg)
	// if hint != "" {
	//     b.WriteString("\nHINT: ")
	//     b.WriteString(hint)
	// }
	// if detail != "" {
	//     b.WriteString("\nDETAIL: ")
	//     b.WriteString(detail)
	// }
	// return b.String()
}

// NewErrorWithDepthf creates an Error and extracts the context
// information at the specified depth level.
func NewErrorWithDepthf(depth int, code string, format string, args ...interface{}) *Error {
	// srcCtx := makeSrcCtx(depth + 1)
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Code:    code,
		// Source:  &srcCtx,
	}
}

// NewError creates an Error.
func NewError(code string, msg string) *Error {
	return NewErrorWithDepthf(1, code, "%s", msg)
}

// NewErrorWithDepth creates an Error with context extracted from the
// specified depth.
func NewErrorWithDepth(depth int, code string, msg string) *Error {
	return NewErrorWithDepthf(depth+1, code, "%s", msg)
}

// NewErrorf creates an Error with a format string.
func NewErrorf(code string, format string, args ...interface{}) *Error {
	return NewErrorWithDepthf(1, code, format, args...)
}

// NewDangerousStatementErrorf creates a new Error for "rejected dangerous statements".
func NewDangerousStatementErrorf(format string, args ...interface{}) *Error {
	var buf bytes.Buffer
	buf.WriteString("rejected: ")
	fmt.Fprintf(&buf, format, args...)
	buf.WriteString(" (sql_safe_updates = true)")
	return NewErrorWithDepthf(1, CodeWarningError, "%s", buf.String())
}

// NewWrongNumberOfPreparedStatements creates new an Error for trying to prepare
// a query string containing more than one statement.
func NewWrongNumberOfPreparedStatements(n int) *Error {
	return NewErrorWithDepthf(1, CodeInvalidPreparedStatementDefinitionError,
		"prepared statement had %d statements, expected 1", n)
}

// SetHintf annotates an Error object with a hint.
func (pg *Error) SetHintf(f string, args ...interface{}) *Error {
	pg.Hint = fmt.Sprintf(f, args...)
	return pg
}

// SetDetailf annotates an Error object with details.
func (pg *Error) SetDetailf(f string, args ...interface{}) *Error {
	pg.Detail = fmt.Sprintf(f, args...)
	return pg
}

// makeSrcCtx creates a Error_Source value with contextual information
// about the caller at the requested depth.
// func makeSrcCtx(depth int) Error_Source {
//     //panic("test")
//     //f, l, fun := caller.Lookup(depth + 1)
//     //return Error_Source{File: f, Line: int32(l), Function: fun}
// }

// GetPGCause returns an unwrapped Error.
func GetPGCause(err error) (*Error, bool) {
	switch pgErr := errors.Cause(err).(type) {
	case *Error:
		return pgErr, true

	default:
		return nil, false
	}
}

// UnimplementedWithIssueErrorf constructs an error with the formatted message
// and a link to the passed issue. Recorded as "#<issue>" in tracking.
func UnimplementedWithIssueErrorf(issue int, format string, args ...interface{}) error {
	err := NewErrorWithDepthf(1, CodeFeatureNotSupportedError, "unimplemented: "+format, args...)
	err.InternalCommand = fmt.Sprintf("#%d", issue)
	return err.SetHintf("See: https://github.com/readystock/noah/issues/%d", issue)
}

// UnimplementedWithIssueError constructs an error with the given message
// and a link to the passed issue. Recorded as "#<issue>" in tracking.
func UnimplementedWithIssueError(issue int, msg string) error {
	err := NewErrorWithDepthf(1, CodeFeatureNotSupportedError, "unimplemented: %s", msg)
	err.InternalCommand = fmt.Sprintf("#%d", issue)
	return err.SetHintf("See: https://github.com/readystock/noah/issues/%d", issue)
}

// Unimplemented constructs an unimplemented feature error.
//
// `feature` is used for tracking, and is not included when the error printed.
func Unimplemented(feature, msg string) *Error {
	return UnimplementedWithDepth(1, feature, msg)
}

// UnimplementedWithDepth constructs an implemented feature error,
// tracking the context at the specified depth.
func UnimplementedWithDepth(depth int, feature, msg string) *Error {
	err := NewErrorWithDepthf(depth+1, CodeFeatureNotSupportedError, "%s", msg)
	err.InternalCommand = feature
	return err
}
