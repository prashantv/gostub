package gostub

import (
	"reflect"
	"testing"
)

// Stub replaces the value stored at varPtrToStub with stubVal.
// varToStub must be a pointer to the variable.
// The variable will be reset to the original value at the end of the test
// but it can be reset earlier using the returned function.
func Stub[T any](t testing.TB, varPtrToStub *T, stubVal T) (reset func()) {
	orig := *varPtrToStub
	*varPtrToStub = stubVal
	cleanup := func() {
		*varPtrToStub = orig
	}
	t.Cleanup(cleanup)
	return cleanup
}

type Func[T any] interface {
	func(any) T | func(any, any) T
}

// FuncReturning creates a new function with type funcType that returns results.
func FuncReturning(funcType reflect.Type, results ...interface{}) reflect.Value {
	var resultValues []reflect.Value
	for i, r := range results {
		var retValue reflect.Value
		if r == nil {
			// We can't use reflect.ValueOf(nil), so we need to create the zero value.
			retValue = reflect.Zero(funcType.Out(i))
		} else {
			// We cannot simply use reflect.ValueOf(r) as that does not work for
			// interface types, as reflect.ValueOf receives the dynamic type, which
			// is the underlying type. e.g. for an error, it may *errors.errorString.
			// Instead, we make the return type's expected interface value using
			// reflect.New, and set the data to the passed in value.
			tempV := reflect.New(funcType.Out(i))
			tempV.Elem().Set(reflect.ValueOf(r))
			retValue = tempV.Elem()
		}
		resultValues = append(resultValues, retValue)
	}
	return reflect.MakeFunc(funcType, func(_ []reflect.Value) []reflect.Value {
		return resultValues
	})
}

// ResetSingle resets a single stubbed variable back to its original value.
func (s *Stubs) ResetSingle(varToStub interface{}) {
	v := reflect.ValueOf(varToStub)
	originalVal, ok := s.stubs[v]
	if !ok {
		panic("cannot reset variable as it has not been stubbed yet")
	}

	v.Elem().Set(originalVal)
}
