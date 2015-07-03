package gostub

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestStubTime(t *testing.T) {
	var timeNow = time.Now

	var fakeTime = time.Date(2015, 7, 1, 0, 0, 0, 0, time.UTC)
	StubFunc(&timeNow, fakeTime)
	expectVal(t, fakeTime, timeNow())
}

func TestReturnErr(t *testing.T) {
	var osRemove = os.Remove

	StubFunc(&osRemove, nil)
	expectVal(t, nil, osRemove("test"))

	e := errors.New("err")
	StubFunc(&osRemove, e)
	expectVal(t, e, osRemove("test"))
}

func TestStubHostname(t *testing.T) {
	var osHostname = os.Hostname

	StubFunc(&osHostname, "fakehost", nil)
	hostname, err := osHostname()
	expectVal(t, "fakehost", hostname)
	expectVal(t, nil, err)

	var errNoHost = errors.New("no hostname")
	StubFunc(&osHostname, "", errNoHost)
	hostname, err = osHostname()
	expectVal(t, "", hostname)
	expectVal(t, errNoHost, err)
}

func TestStubReturnFunc(t *testing.T) {
	var retFunc = func() func() error {
		return func() error {
			return errors.New("err")
		}
	}

	var errInception = errors.New("in limbo")
	StubFunc(&retFunc, func() error {
		return errInception
	})
	expectVal(t, errInception, retFunc()())
}
