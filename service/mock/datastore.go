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
// 			AddFilterDimensionFunc: func(ctx context.Context, s string, dimension model.Dimension) error {
// 				panic("mock out the AddFilterDimension method")
// 			},
// 			CheckerFunc: func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error {
// 				panic("mock out the Checker method")
// 			},
// 			ConnFunc: func() *mongo.MongoConnection {
// 				panic("mock out the Conn method")
// 			},
// 			CreateFilterFunc: func(contextMoqParam context.Context, filter *model.Filter) error {
// 				panic("mock out the CreateFilter method")
// 			},
// 			CreateFilterOutputFunc: func(contextMoqParam context.Context, filterOutput *model.FilterOutput) error {
// 				panic("mock out the CreateFilterOutput method")
// 			},
// 			GetFilterFunc: func(contextMoqParam context.Context, s string) (*model.Filter, error) {
// 				panic("mock out the GetFilter method")
// 			},
// 			GetFilterDimensionsFunc: func(contextMoqParam context.Context, s string, n1 int, n2 int) ([]model.Dimension, int, error) {
// 				panic("mock out the GetFilterDimensions method")
// 			},
// 		}
//
// 		// use mockedDatastore in code that requires service.Datastore
// 		// and then make assertions.
//
// 	}
type DatastoreMock struct {
	// AddFilterDimensionFunc mocks the AddFilterDimension method.
	AddFilterDimensionFunc func(ctx context.Context, s string, dimension model.Dimension) error

	// CheckerFunc mocks the Checker method.
	CheckerFunc func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error

	// ConnFunc mocks the Conn method.
	ConnFunc func() *mongo.MongoConnection

	// CreateFilterFunc mocks the CreateFilter method.
	CreateFilterFunc func(contextMoqParam context.Context, filter *model.Filter) error

	// CreateFilterOutputFunc mocks the CreateFilterOutput method.
	CreateFilterOutputFunc func(contextMoqParam context.Context, filterOutput *model.FilterOutput) error

	// GetFilterFunc mocks the GetFilter method.
	GetFilterFunc func(contextMoqParam context.Context, s string) (*model.Filter, error)

	// GetFilterDimensionsFunc mocks the GetFilterDimensions method.
	GetFilterDimensionsFunc func(contextMoqParam context.Context, s string, n1 int, n2 int) ([]model.Dimension, int, error)

	// calls tracks calls to the methods.
	calls struct {
		// AddFilterDimension holds details about calls to the AddFilterDimension method.
		AddFilterDimension []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// S is the s argument value.
			S string
			// Dimension is the dimension argument value.
			Dimension model.Dimension
		}
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
		// CreateFilterOutput holds details about calls to the CreateFilterOutput method.
		CreateFilterOutput []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// FilterOutput is the filterOutput argument value.
			FilterOutput *model.FilterOutput
		}
		// GetFilter holds details about calls to the GetFilter method.
		GetFilter []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// S is the s argument value.
			S string
		}
		// GetFilterDimensions holds details about calls to the GetFilterDimensions method.
		GetFilterDimensions []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// S is the s argument value.
			S string
			// N1 is the n1 argument value.
			N1 int
			// N2 is the n2 argument value.
			N2 int
		}
	}
	lockAddFilterDimension  sync.RWMutex
	lockChecker             sync.RWMutex
	lockConn                sync.RWMutex
	lockCreateFilter        sync.RWMutex
	lockCreateFilterOutput  sync.RWMutex
	lockGetFilter           sync.RWMutex
	lockGetFilterDimensions sync.RWMutex
}

// AddFilterDimension calls AddFilterDimensionFunc.
func (mock *DatastoreMock) AddFilterDimension(ctx context.Context, s string, dimension model.Dimension) error {
	if mock.AddFilterDimensionFunc == nil {
		panic("DatastoreMock.AddFilterDimensionFunc: method is nil but Datastore.AddFilterDimension was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		S         string
		Dimension model.Dimension
	}{
		Ctx:       ctx,
		S:         s,
		Dimension: dimension,
	}
	mock.lockAddFilterDimension.Lock()
	mock.calls.AddFilterDimension = append(mock.calls.AddFilterDimension, callInfo)
	mock.lockAddFilterDimension.Unlock()
	return mock.AddFilterDimensionFunc(ctx, s, dimension)
}

// AddFilterDimensionCalls gets all the calls that were made to AddFilterDimension.
// Check the length with:
//     len(mockedDatastore.AddFilterDimensionCalls())
func (mock *DatastoreMock) AddFilterDimensionCalls() []struct {
	Ctx       context.Context
	S         string
	Dimension model.Dimension
} {
	var calls []struct {
		Ctx       context.Context
		S         string
		Dimension model.Dimension
	}
	mock.lockAddFilterDimension.RLock()
	calls = mock.calls.AddFilterDimension
	mock.lockAddFilterDimension.RUnlock()
	return calls
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

// CreateFilterOutput calls CreateFilterOutputFunc.
func (mock *DatastoreMock) CreateFilterOutput(contextMoqParam context.Context, filterOutput *model.FilterOutput) error {
	if mock.CreateFilterOutputFunc == nil {
		panic("DatastoreMock.CreateFilterOutputFunc: method is nil but Datastore.CreateFilterOutput was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		FilterOutput    *model.FilterOutput
	}{
		ContextMoqParam: contextMoqParam,
		FilterOutput:    filterOutput,
	}
	mock.lockCreateFilterOutput.Lock()
	mock.calls.CreateFilterOutput = append(mock.calls.CreateFilterOutput, callInfo)
	mock.lockCreateFilterOutput.Unlock()
	return mock.CreateFilterOutputFunc(contextMoqParam, filterOutput)
}

// CreateFilterOutputCalls gets all the calls that were made to CreateFilterOutput.
// Check the length with:
//     len(mockedDatastore.CreateFilterOutputCalls())
func (mock *DatastoreMock) CreateFilterOutputCalls() []struct {
	ContextMoqParam context.Context
	FilterOutput    *model.FilterOutput
} {
	var calls []struct {
		ContextMoqParam context.Context
		FilterOutput    *model.FilterOutput
	}
	mock.lockCreateFilterOutput.RLock()
	calls = mock.calls.CreateFilterOutput
	mock.lockCreateFilterOutput.RUnlock()
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
	return calls
}

// GetFilterDimensions calls GetFilterDimensionsFunc.
func (mock *DatastoreMock) GetFilterDimensions(contextMoqParam context.Context, s string, n1 int, n2 int) ([]model.Dimension, int, error) {
	if mock.GetFilterDimensionsFunc == nil {
		panic("DatastoreMock.GetFilterDimensionsFunc: method is nil but Datastore.GetFilterDimensions was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		S               string
		N1              int
		N2              int
	}{
		ContextMoqParam: contextMoqParam,
		S:               s,
		N1:              n1,
		N2:              n2,
	}
	mock.lockGetFilterDimensions.Lock()
	mock.calls.GetFilterDimensions = append(mock.calls.GetFilterDimensions, callInfo)
	mock.lockGetFilterDimensions.Unlock()
	return mock.GetFilterDimensionsFunc(contextMoqParam, s, n1, n2)
}

// GetFilterDimensionsCalls gets all the calls that were made to GetFilterDimensions.
// Check the length with:
//     len(mockedDatastore.GetFilterDimensionsCalls())
func (mock *DatastoreMock) GetFilterDimensionsCalls() []struct {
	ContextMoqParam context.Context
	S               string
	N1              int
	N2              int
} {
	var calls []struct {
		ContextMoqParam context.Context
		S               string
		N1              int
		N2              int
	}
	mock.lockGetFilterDimensions.RLock()
	calls = mock.calls.GetFilterDimensions
	mock.lockGetFilterDimensions.RUnlock()
	return calls
}
