// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	mongo "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"sync"
)

// Ensure, that DatastoreMock does implement service.Datastore.
// If this is not the case, regenerate this file with moq.
var _ service.Datastore = &DatastoreMock{}

// DatastoreMock is a mock implementation of service.Datastore.
//
// 	func TestSomethingThatUsesDatastore(t *testing.T) {
//
// 		// make and configure a mocked service.Datastore
// 		mockedDatastore := &DatastoreMock{
// 			CheckerFunc: func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error {
// 				panic("mock out the Checker method")
// 			},
// 			ConnFunc: func() *mongo.MongoConnection {
// 				panic("mock out the Conn method")
// 			},
// 			CreateFilterFunc: func(contextMoqParam context.Context, filter *model.Filter) error {
// 				panic("mock out the CreateFilter method")
// 			},
// 			GetFilterFunc: func(contextMoqParam context.Context, s string) (*model.Filter, error) {
// 				panic("mock out the GetFilter method")
// 			CreateFilterOutputFunc: func(contextMoqParam context.Context, filterOutputResponse *model.FilterOutputResponse) error {
// 				panic("mock out the CreateFilterOutput method")
// 			},
// 			GetFilterDimensionsFunc: func(contextMoqParam context.Context, s string) ([]model.Dimension, error) {
// 				panic("mock out the GetFilterDimensions method")
// 			},
// 		}
//
// 		// use mockedDatastore in code that requires service.Datastore
// 		// and then make assertions.
//
// 	}
type DatastoreMock struct {
	// CheckerFunc mocks the Checker method.
	CheckerFunc func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error

	// ConnFunc mocks the Conn method.
	ConnFunc func() *mongo.MongoConnection

	// CreateFilterFunc mocks the CreateFilter method.
	CreateFilterFunc func(contextMoqParam context.Context, filter *model.Filter) error

	// GetFilterFunc mocks the GetFilter method.
	GetFilterFunc func(contextMoqParam context.Context, s string) (*model.Filter, error)
	// CreateFilterOutputFunc mocks the CreateFilterOutput method.
	CreateFilterOutputFunc func(contextMoqParam context.Context, filterOutputResponse *model.FilterOutputResponse) error

	// GetFilterDimensionsFunc mocks the GetFilterDimensions method.
	GetFilterDimensionsFunc func(contextMoqParam context.Context, s string) ([]model.Dimension, error)

	// calls tracks calls to the methods.
	calls struct {
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// CheckState is the checkState argument value.
			CheckState *healthcheck.CheckState
		}
		// Conn holds details about calls to the Conn method.
		Conn []struct {
		}
		// CreateFilter holds details about calls to the CreateFilter method.
		CreateFilter []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Filter is the filter argument value.
			Filter *model.Filter
		}
		// GetFilter holds details about calls to the GetFilter method.
		GetFilter []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// S is the s argument value.
			S string
		// CreateFilterOutput holds details about calls to the CreateFilterOutput method.
		CreateFilterOutput []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// FilterOutputResponse is the filterOutputResponse argument value.
			FilterOutputResponse *model.FilterOutputResponse
		}
		// GetFilterDimensions holds details about calls to the GetFilterDimensions method.
		GetFilterDimensions []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// S is the s argument value.
			S string
		}
	}
	lockChecker             sync.RWMutex
	lockConn                sync.RWMutex
	lockCreateFilter        sync.RWMutex
	lockGetFilter           sync.RWMutex
	lockCreateFilterOutput  sync.RWMutex
	lockGetFilterDimensions sync.RWMutex
}

// Checker calls CheckerFunc.
func (mock *DatastoreMock) Checker(contextMoqParam context.Context, checkState *healthcheck.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("DatastoreMock.CheckerFunc: method is nil but Datastore.Checker was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		CheckState      *healthcheck.CheckState
	}{
		ContextMoqParam: contextMoqParam,
		CheckState:      checkState,
	}
	mock.lockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	mock.lockChecker.Unlock()
	return mock.CheckerFunc(contextMoqParam, checkState)
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//     len(mockedDatastore.CheckerCalls())
func (mock *DatastoreMock) CheckerCalls() []struct {
	ContextMoqParam context.Context
	CheckState      *healthcheck.CheckState
} {
	var calls []struct {
		ContextMoqParam context.Context
		CheckState      *healthcheck.CheckState
	}
	mock.lockChecker.RLock()
	calls = mock.calls.Checker
	mock.lockChecker.RUnlock()
	return calls
}

// Conn calls ConnFunc.
func (mock *DatastoreMock) Conn() *mongo.MongoConnection {
	if mock.ConnFunc == nil {
		panic("DatastoreMock.ConnFunc: method is nil but Datastore.Conn was just called")
	}
	callInfo := struct {
	}{}
	mock.lockConn.Lock()
	mock.calls.Conn = append(mock.calls.Conn, callInfo)
	mock.lockConn.Unlock()
	return mock.ConnFunc()
}

// ConnCalls gets all the calls that were made to Conn.
// Check the length with:
//     len(mockedDatastore.ConnCalls())
func (mock *DatastoreMock) ConnCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockConn.RLock()
	calls = mock.calls.Conn
	mock.lockConn.RUnlock()
	return calls
}

// CreateFilter calls CreateFilterFunc.
func (mock *DatastoreMock) CreateFilter(contextMoqParam context.Context, filter *model.Filter) error {
	if mock.CreateFilterFunc == nil {
		panic("DatastoreMock.CreateFilterFunc: method is nil but Datastore.CreateFilter was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Filter          *model.Filter
	}{
		ContextMoqParam: contextMoqParam,
		Filter:          filter,
	}
	mock.lockCreateFilter.Lock()
	mock.calls.CreateFilter = append(mock.calls.CreateFilter, callInfo)
	mock.lockCreateFilter.Unlock()
	return mock.CreateFilterFunc(contextMoqParam, filter)
}

// CreateFilterCalls gets all the calls that were made to CreateFilter.
// Check the length with:
//     len(mockedDatastore.CreateFilterCalls())
func (mock *DatastoreMock) CreateFilterCalls() []struct {
	ContextMoqParam context.Context
	Filter          *model.Filter
} {
	var calls []struct {
		ContextMoqParam context.Context
		Filter          *model.Filter
	}
	mock.lockCreateFilter.RLock()
	calls = mock.calls.CreateFilter
	mock.lockCreateFilter.RUnlock()
	return calls
}

// GetFilter calls GetFilterFunc.
func (mock *DatastoreMock) GetFilter(contextMoqParam context.Context, s string) (*model.Filter, error) {
	if mock.GetFilterFunc == nil {
		panic("DatastoreMock.GetFilterFunc: method is nil but Datastore.GetFilter was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		S               string
	}{
		ContextMoqParam: contextMoqParam,
		S:               s,
	}
	mock.lockGetFilter.Lock()
	mock.calls.GetFilter = append(mock.calls.GetFilter, callInfo)
	mock.lockGetFilter.Unlock()
	return mock.GetFilterFunc(contextMoqParam, s)
}

// GetFilterCalls gets all the calls that were made to GetFilter.
// Check the length with:
//     len(mockedDatastore.GetFilterCalls())
func (mock *DatastoreMock) GetFilterCalls() []struct {
	ContextMoqParam context.Context
	S               string
} {
	var calls []struct {
		ContextMoqParam context.Context
		S               string
	}
	mock.lockGetFilter.RLock()
	calls = mock.calls.GetFilter
	mock.lockGetFilter.RUnlock()
// CreateFilterOutput calls CreateFilterOutputFunc.
func (mock *DatastoreMock) CreateFilterOutput(contextMoqParam context.Context, filterOutputResponse *model.FilterOutputResponse) error {
	if mock.CreateFilterOutputFunc == nil {
		panic("DatastoreMock.CreateFilterOutputFunc: method is nil but Datastore.CreateFilterOutput was just called")
	}
	callInfo := struct {
		ContextMoqParam      context.Context
		FilterOutputResponse *model.FilterOutputResponse
	}{
		ContextMoqParam:      contextMoqParam,
		FilterOutputResponse: filterOutputResponse,
	}
	mock.lockCreateFilterOutput.Lock()
	mock.calls.CreateFilterOutput = append(mock.calls.CreateFilterOutput, callInfo)
	mock.lockCreateFilterOutput.Unlock()
	return mock.CreateFilterOutputFunc(contextMoqParam, filterOutputResponse)
}

// CreateFilterOutputCalls gets all the calls that were made to CreateFilterOutput.
// Check the length with:
//     len(mockedDatastore.CreateFilterOutputCalls())
func (mock *DatastoreMock) CreateFilterOutputCalls() []struct {
	ContextMoqParam      context.Context
	FilterOutputResponse *model.FilterOutputResponse
} {
	var calls []struct {
		ContextMoqParam      context.Context
		FilterOutputResponse *model.FilterOutputResponse
	}
	mock.lockCreateFilterOutput.RLock()
	calls = mock.calls.CreateFilterOutput
	mock.lockCreateFilterOutput.RUnlock()
	return calls
}

// GetFilterDimensions calls GetFilterDimensionsFunc.
func (mock *DatastoreMock) GetFilterDimensions(contextMoqParam context.Context, s string) ([]model.Dimension, error) {
	if mock.GetFilterDimensionsFunc == nil {
		panic("DatastoreMock.GetFilterDimensionsFunc: method is nil but Datastore.GetFilterDimensions was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		S               string
	}{
		ContextMoqParam: contextMoqParam,
		S:               s,
	}
	mock.lockGetFilterDimensions.Lock()
	mock.calls.GetFilterDimensions = append(mock.calls.GetFilterDimensions, callInfo)
	mock.lockGetFilterDimensions.Unlock()
	return mock.GetFilterDimensionsFunc(contextMoqParam, s)
}

// GetFilterDimensionsCalls gets all the calls that were made to GetFilterDimensions.
// Check the length with:
//     len(mockedDatastore.GetFilterDimensionsCalls())
func (mock *DatastoreMock) GetFilterDimensionsCalls() []struct {
	ContextMoqParam context.Context
	S               string
} {
	var calls []struct {
		ContextMoqParam context.Context
		S               string
	}
	mock.lockGetFilterDimensions.RLock()
	calls = mock.calls.GetFilterDimensions
	mock.lockGetFilterDimensions.RUnlock()
	return calls
}
