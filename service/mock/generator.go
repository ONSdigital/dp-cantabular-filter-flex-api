// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	"github.com/google/uuid"
	"sync"
	"time"
)

// Ensure, that GeneratorMock does implement service.Generator.
// If this is not the case, regenerate this file with moq.
var _ service.Generator = &GeneratorMock{}

// GeneratorMock is a mock implementation of service.Generator.
//
// 	func TestSomethingThatUsesGenerator(t *testing.T) {
//
// 		// make and configure a mocked service.Generator
// 		mockedGenerator := &GeneratorMock{
// 			PSKFunc: func() ([]byte, error) {
// 				panic("mock out the PSK method")
// 			},
// 			TimestampFunc: func() time.Time {
// 				panic("mock out the Timestamp method")
// 			},
// 			URLFunc: func(host string, path string, args ...interface{}) string {
// 				panic("mock out the URL method")
// 			},
// 			UUIDFunc: func() (uuid.UUID, error) {
// 				panic("mock out the UUID method")
// 			},
// 		}
//
// 		// use mockedGenerator in code that requires service.Generator
// 		// and then make assertions.
//
// 	}
type GeneratorMock struct {
	// PSKFunc mocks the PSK method.
	PSKFunc func() ([]byte, error)

	// TimestampFunc mocks the Timestamp method.
	TimestampFunc func() time.Time

	// URLFunc mocks the URL method.
	URLFunc func(host string, path string, args ...interface{}) string

	// UUIDFunc mocks the UUID method.
	UUIDFunc func() (uuid.UUID, error)

	// calls tracks calls to the methods.
	calls struct {
		// PSK holds details about calls to the PSK method.
		PSK []struct {
		}
		// Timestamp holds details about calls to the Timestamp method.
		Timestamp []struct {
		}
		// URL holds details about calls to the URL method.
		URL []struct {
			// Host is the host argument value.
			Host string
			// Path is the path argument value.
			Path string
			// Args is the args argument value.
			Args []interface{}
		}
		// UUID holds details about calls to the UUID method.
		UUID []struct {
		}
	}
	lockPSK       sync.RWMutex
	lockTimestamp sync.RWMutex
	lockURL       sync.RWMutex
	lockUUID      sync.RWMutex
}

// PSK calls PSKFunc.
func (mock *GeneratorMock) PSK() ([]byte, error) {
	if mock.PSKFunc == nil {
		panic("GeneratorMock.PSKFunc: method is nil but Generator.PSK was just called")
	}
	callInfo := struct {
	}{}
	mock.lockPSK.Lock()
	mock.calls.PSK = append(mock.calls.PSK, callInfo)
	mock.lockPSK.Unlock()
	return mock.PSKFunc()
}

// PSKCalls gets all the calls that were made to PSK.
// Check the length with:
//     len(mockedGenerator.PSKCalls())
func (mock *GeneratorMock) PSKCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockPSK.RLock()
	calls = mock.calls.PSK
	mock.lockPSK.RUnlock()
	return calls
}

// Timestamp calls TimestampFunc.
func (mock *GeneratorMock) Timestamp() time.Time {
	if mock.TimestampFunc == nil {
		panic("GeneratorMock.TimestampFunc: method is nil but Generator.Timestamp was just called")
	}
	callInfo := struct {
	}{}
	mock.lockTimestamp.Lock()
	mock.calls.Timestamp = append(mock.calls.Timestamp, callInfo)
	mock.lockTimestamp.Unlock()
	return mock.TimestampFunc()
}

// TimestampCalls gets all the calls that were made to Timestamp.
// Check the length with:
//     len(mockedGenerator.TimestampCalls())
func (mock *GeneratorMock) TimestampCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockTimestamp.RLock()
	calls = mock.calls.Timestamp
	mock.lockTimestamp.RUnlock()
	return calls
}

// URL calls URLFunc.
func (mock *GeneratorMock) URL(host string, path string, args ...interface{}) string {
	if mock.URLFunc == nil {
		panic("GeneratorMock.URLFunc: method is nil but Generator.URL was just called")
	}
	callInfo := struct {
		Host string
		Path string
		Args []interface{}
	}{
		Host: host,
		Path: path,
		Args: args,
	}
	mock.lockURL.Lock()
	mock.calls.URL = append(mock.calls.URL, callInfo)
	mock.lockURL.Unlock()
	return mock.URLFunc(host, path, args...)
}

// URLCalls gets all the calls that were made to URL.
// Check the length with:
//     len(mockedGenerator.URLCalls())
func (mock *GeneratorMock) URLCalls() []struct {
	Host string
	Path string
	Args []interface{}
} {
	var calls []struct {
		Host string
		Path string
		Args []interface{}
	}
	mock.lockURL.RLock()
	calls = mock.calls.URL
	mock.lockURL.RUnlock()
	return calls
}

// UUID calls UUIDFunc.
func (mock *GeneratorMock) UUID() (uuid.UUID, error) {
	if mock.UUIDFunc == nil {
		panic("GeneratorMock.UUIDFunc: method is nil but Generator.UUID was just called")
	}
	callInfo := struct {
	}{}
	mock.lockUUID.Lock()
	mock.calls.UUID = append(mock.calls.UUID, callInfo)
	mock.lockUUID.Unlock()
	return mock.UUIDFunc()
}

// UUIDCalls gets all the calls that were made to UUID.
// Check the length with:
//     len(mockedGenerator.UUIDCalls())
func (mock *GeneratorMock) UUIDCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockUUID.RLock()
	calls = mock.calls.UUID
	mock.lockUUID.RUnlock()
	return calls
}
