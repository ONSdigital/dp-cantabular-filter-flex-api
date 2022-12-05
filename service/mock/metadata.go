// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	"sync"
)

// Ensure, that MetadataAPIClientMock does implement service.MetadataAPIClient.
// If this is not the case, regenerate this file with moq.
var _ service.MetadataClient = &MetadataAPIClientMock{}

// MetadataAPIClientMock is a mock implementation of service.MetadataAPIClient.
//
//	func TestSomethingThatUsesMetadataAPIClient(t *testing.T) {
//
//		// make and configure a mocked service.MetadataAPIClient
//		mockedMetadataAPIClient := &MetadataAPIClientMock{
//			GetDefaultClassificationFunc: func(ctx context.Context, req cantabularmetadata.GetDefaultClassificationRequest) (*cantabularmetadata.GetDefaultClassificationResponse, error) {
//				panic("mock out the GetDefaultClassification method")
//			},
//		}
//
//		// use mockedMetadataAPIClient in code that requires service.MetadataAPIClient
//		// and then make assertions.
//
//	}
type MetadataAPIClientMock struct {
	// GetDefaultClassificationFunc mocks the GetDefaultClassification method.
	GetDefaultClassificationFunc func(ctx context.Context, req cantabularmetadata.GetDefaultClassificationRequest) (*cantabularmetadata.GetDefaultClassificationResponse, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetDefaultClassification holds details about calls to the GetDefaultClassification method.
		GetDefaultClassification []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Req is the req argument value.
			Req cantabularmetadata.GetDefaultClassificationRequest
		}
	}
	lockGetDefaultClassification sync.RWMutex
}

// GetDefaultClassification calls GetDefaultClassificationFunc.
func (mock *MetadataAPIClientMock) GetDefaultClassification(ctx context.Context, req cantabularmetadata.GetDefaultClassificationRequest) (*cantabularmetadata.GetDefaultClassificationResponse, error) {
	if mock.GetDefaultClassificationFunc == nil {
		panic("MetadataAPIClientMock.GetDefaultClassificationFunc: method is nil but MetadataAPIClient.GetDefaultClassification was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Req cantabularmetadata.GetDefaultClassificationRequest
	}{
		Ctx: ctx,
		Req: req,
	}
	mock.lockGetDefaultClassification.Lock()
	mock.calls.GetDefaultClassification = append(mock.calls.GetDefaultClassification, callInfo)
	mock.lockGetDefaultClassification.Unlock()
	return mock.GetDefaultClassificationFunc(ctx, req)
}

// GetDefaultClassificationCalls gets all the calls that were made to GetDefaultClassification.
// Check the length with:
//
//	len(mockedMetadataAPIClient.GetDefaultClassificationCalls())
func (mock *MetadataAPIClientMock) GetDefaultClassificationCalls() []struct {
	Ctx context.Context
	Req cantabularmetadata.GetDefaultClassificationRequest
} {
	var calls []struct {
		Ctx context.Context
		Req cantabularmetadata.GetDefaultClassificationRequest
	}
	mock.lockGetDefaultClassification.RLock()
	calls = mock.calls.GetDefaultClassification
	mock.lockGetDefaultClassification.RUnlock()
	return calls
}