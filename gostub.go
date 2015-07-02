package gostub

import (
	"fmt"
	"reflect"
)

// Stub stores the original value at varPtr and replaces it with stub.
func Stub(varPtr interface{}, stub interface{}) *Stubs {
	return newStub().Stub(varPtr, stub)
}

// Stubs represents a set of stubbed variables that can be reset.
type Stubs struct {
	// stubs is a map from the variable pointer (being stubbed) to the original value.
	stubs map[reflect.Value]reflect.Value
}

func newStub() *Stubs {
	return &Stubs{make(map[reflect.Value]reflect.Value)}
}

func validateTypes(varType reflect.Type, stubType reflect.Type) error {
	if varType.Kind() != reflect.Ptr {
		return fmt.Errorf("variable to stub is expected to be a pointer (*%v), got %v", varType, stubType)
	}
	return nil
}

// Stub stores the original value at varPtr and replaces it with stub.
func (s *Stubs) Stub(varPtr interface{}, stub interface{}) *Stubs {
	v := reflect.ValueOf(varPtr)
	stubVal := reflect.ValueOf(stub)

	// Ensure varPtr is a pointer to something of type stub.
	if err := validateTypes(v.Type(), stubVal.Type()); err != nil {
		panic(err)
	}

	if _, ok := s.stubs[v]; !ok {
		// Store the original value if this is the first time varPtr is being stubbed.
		s.stubs[v] = reflect.ValueOf(v.Elem().Interface())
	}

	v.Elem().Set(stubVal)
	return s
}

// Reset resets all stubbed variables back to their original values.
func (s *Stubs) Reset() {
	for k, v := range s.stubs {
		k.Elem().Set(v)
	}
}

// ResetSingle resets a single stubbed variable back to the original value.
func (s *Stubs) ResetSingle(varPtr interface{}) {
	varPtrV := reflect.ValueOf(varPtr)
	v, ok := s.stubs[varPtrV]
	if !ok {
		panic("cannot reset variable as it has not been stubbed yet")
	}

	varPtrV.Elem().Set(v)
}
