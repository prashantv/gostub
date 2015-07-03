package gostub

import "reflect"

// Stub replaces the value stored at varToStub with stubVal.
// varToStub must be a pointer to the variable. stubVal should have a type
// that is assignable to the variable.
func Stub(varToStub interface{}, stubVal interface{}) *Stubs {
	return newStub().Stub(varToStub, stubVal)
}

// Stubs represents a set of stubbed variables that can be reset.
type Stubs struct {
	// stubs is a map from the variable pointer (being stubbed) to the original value.
	stubs map[reflect.Value]reflect.Value
}

func newStub() *Stubs {
	return &Stubs{make(map[reflect.Value]reflect.Value)}
}

// Stub replaces the value stored at varToStub with stubVal.
// varToStub must be a pointer to the variable. stubVal should have a type
// that is assignable to the variable.
func (s *Stubs) Stub(varToStub interface{}, stubVal interface{}) *Stubs {
	v := reflect.ValueOf(varToStub)
	stub := reflect.ValueOf(stubVal)

	// Ensure varToStub is a pointer to the variable.
	if v.Type().Kind() != reflect.Ptr {
		panic("variable to stub is expected to be a pointer")
	}

	if _, ok := s.stubs[v]; !ok {
		// Store the original value if this is the first time varPtr is being stubbed.
		s.stubs[v] = reflect.ValueOf(v.Elem().Interface())
	}

	// *varToStub = stubVal
	v.Elem().Set(stub)
	return s
}

// Reset resets all stubbed variables back to their original values.
func (s *Stubs) Reset() {
	for v, originalVal := range s.stubs {
		v.Elem().Set(originalVal)
	}
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
