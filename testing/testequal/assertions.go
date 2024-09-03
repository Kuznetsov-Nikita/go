//go:build !solution

package testequal

import (
	"fmt"
	"reflect"
)

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if reflect.ValueOf(expected).Kind() == reflect.Struct || reflect.ValueOf(actual).Kind() == reflect.Struct || !reflect.DeepEqual(expected, actual) {
		t.Helper()
		if len(msgAndArgs) > 0 {
			msg := msgAndArgs[0].(string)
			args := msgAndArgs[1:]
			t.Errorf("not equal:\n\texpected: %v\n\tactual  : %v\n\tmessage : %v", expected, actual, fmt.Sprintf(msg, args...))
		} else {
			t.Errorf("not equal:\n\texpected: %v\n\tactual  : %v", expected, actual)
		}
		return false
	}
	return true
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	_, ok1 := expected.(struct{})
	if ok1 {
		return true
	}
	_, ok2 := actual.(struct{})
	if ok2 {
		return true
	}
	if reflect.DeepEqual(expected, actual) {
		t.Helper()
		if len(msgAndArgs) > 0 {
			msg := msgAndArgs[0].(string)
			args := msgAndArgs[1:]
			t.Errorf("equal:\n\texpected: %v\n\tactual  : %v\n\tmessage : %v", expected, actual, fmt.Sprintf(msg, args...))
		} else {
			t.Errorf("equal:\n\texpected: %v\n\tactual  : %v", expected, actual)
		}
		return false
	}
	return true
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	if reflect.ValueOf(expected).Kind() == reflect.Struct || reflect.ValueOf(actual).Kind() == reflect.Struct || !reflect.DeepEqual(expected, actual) {
		t.Helper()
		if len(msgAndArgs) > 0 {
			msg := msgAndArgs[0].(string)
			args := msgAndArgs[1:]
			t.Errorf("not equal:\n\texpected: %v\n\tactual  : %v\n\tmessage : %v", expected, actual, fmt.Sprintf(msg, args...))
		} else {
			t.Errorf("not equal:\n\texpected: %v\n\tactual  : %v", expected, actual)
		}
		t.FailNow()
	}
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	_, ok1 := expected.(struct{})
	if ok1 {
		return
	}
	_, ok2 := actual.(struct{})
	if ok2 {
		return
	}
	if reflect.DeepEqual(expected, actual) {
		t.Helper()
		if len(msgAndArgs) > 0 {
			msg := msgAndArgs[0].(string)
			args := msgAndArgs[1:]
			t.Errorf("equal:\n\texpected: %v\n\tactual  : %v\n\tmessage : %v", expected, actual, fmt.Sprintf(msg, args...))
		} else {
			t.Errorf("equal:\n\texpected: %v\n\tactual  : %v", expected, actual)
		}
		t.FailNow()
	}
}
