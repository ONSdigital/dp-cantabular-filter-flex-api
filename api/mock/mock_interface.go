// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	http "net/http"
	reflect "reflect"
	time "time"

	cantabular "github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	gql "github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	dataset "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	model "github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// Mockresponder is a mock of responder interface.
type Mockresponder struct {
	ctrl     *gomock.Controller
	recorder *MockresponderMockRecorder
}

// MockresponderMockRecorder is the mock recorder for Mockresponder.
type MockresponderMockRecorder struct {
	mock *Mockresponder
}

// NewMockresponder creates a new mock instance.
func NewMockresponder(ctrl *gomock.Controller) *Mockresponder {
	mock := &Mockresponder{ctrl: ctrl}
	mock.recorder = &MockresponderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockresponder) EXPECT() *MockresponderMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *Mockresponder) Error(arg0 context.Context, arg1 http.ResponseWriter, arg2 int, arg3 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", arg0, arg1, arg2, arg3)
}

// Error indicates an expected call of Error.
func (mr *MockresponderMockRecorder) Error(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*Mockresponder)(nil).Error), arg0, arg1, arg2, arg3)
}

// Errors mocks base method.
func (m *Mockresponder) Errors(arg0 context.Context, arg1 http.ResponseWriter, arg2 int, arg3 []error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Errors", arg0, arg1, arg2, arg3)
}

// Errors indicates an expected call of Errors.
func (mr *MockresponderMockRecorder) Errors(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errors", reflect.TypeOf((*Mockresponder)(nil).Errors), arg0, arg1, arg2, arg3)
}

// JSON mocks base method.
func (m *Mockresponder) JSON(arg0 context.Context, arg1 http.ResponseWriter, arg2 int, arg3 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "JSON", arg0, arg1, arg2, arg3)
}

// JSON indicates an expected call of JSON.
func (mr *MockresponderMockRecorder) JSON(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JSON", reflect.TypeOf((*Mockresponder)(nil).JSON), arg0, arg1, arg2, arg3)
}

// StatusCode mocks base method.
func (m *Mockresponder) StatusCode(arg0 http.ResponseWriter, arg1 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StatusCode", arg0, arg1)
}

// StatusCode indicates an expected call of StatusCode.
func (mr *MockresponderMockRecorder) StatusCode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StatusCode", reflect.TypeOf((*Mockresponder)(nil).StatusCode), arg0, arg1)
}

// Mockdatastore is a mock of datastore interface.
type Mockdatastore struct {
	ctrl     *gomock.Controller
	recorder *MockdatastoreMockRecorder
}

// MockdatastoreMockRecorder is the mock recorder for Mockdatastore.
type MockdatastoreMockRecorder struct {
	mock *Mockdatastore
}

// NewMockdatastore creates a new mock instance.
func NewMockdatastore(ctrl *gomock.Controller) *Mockdatastore {
	mock := &Mockdatastore{ctrl: ctrl}
	mock.recorder = &MockdatastoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockdatastore) EXPECT() *MockdatastoreMockRecorder {
	return m.recorder
}

// AddFilterDimension mocks base method.
func (m *Mockdatastore) AddFilterDimension(arg0 context.Context, arg1 string, arg2 model.Dimension) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFilterDimension", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFilterDimension indicates an expected call of AddFilterDimension.
func (mr *MockdatastoreMockRecorder) AddFilterDimension(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFilterDimension", reflect.TypeOf((*Mockdatastore)(nil).AddFilterDimension), arg0, arg1, arg2)
}

// AddFilterOutputEvent mocks base method.
func (m *Mockdatastore) AddFilterOutputEvent(arg0 context.Context, arg1 string, arg2 *model.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFilterOutputEvent", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFilterOutputEvent indicates an expected call of AddFilterOutputEvent.
func (mr *MockdatastoreMockRecorder) AddFilterOutputEvent(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFilterOutputEvent", reflect.TypeOf((*Mockdatastore)(nil).AddFilterOutputEvent), arg0, arg1, arg2)
}

// CreateFilter mocks base method.
func (m *Mockdatastore) CreateFilter(arg0 context.Context, arg1 *model.Filter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFilter", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFilter indicates an expected call of CreateFilter.
func (mr *MockdatastoreMockRecorder) CreateFilter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFilter", reflect.TypeOf((*Mockdatastore)(nil).CreateFilter), arg0, arg1)
}

// CreateFilterOutput mocks base method.
func (m *Mockdatastore) CreateFilterOutput(arg0 context.Context, arg1 *model.FilterOutput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFilterOutput", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFilterOutput indicates an expected call of CreateFilterOutput.
func (mr *MockdatastoreMockRecorder) CreateFilterOutput(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFilterOutput", reflect.TypeOf((*Mockdatastore)(nil).CreateFilterOutput), arg0, arg1)
}

// DeleteFilterDimension mocks base method.
func (m *Mockdatastore) DeleteFilterDimension(arg0 context.Context, arg1, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFilterDimension", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFilterDimension indicates an expected call of DeleteFilterDimension.
func (mr *MockdatastoreMockRecorder) DeleteFilterDimension(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFilterDimension", reflect.TypeOf((*Mockdatastore)(nil).DeleteFilterDimension), arg0, arg1, arg2)
}

// DeleteFilterDimensionOptions mocks base method.
func (m *Mockdatastore) DeleteFilterDimensionOptions(arg0 context.Context, arg1, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFilterDimensionOptions", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFilterDimensionOptions indicates an expected call of DeleteFilterDimensionOptions.
func (mr *MockdatastoreMockRecorder) DeleteFilterDimensionOptions(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFilterDimensionOptions", reflect.TypeOf((*Mockdatastore)(nil).DeleteFilterDimensionOptions), arg0, arg1, arg2)
}

// GetFilter mocks base method.
func (m *Mockdatastore) GetFilter(arg0 context.Context, arg1 string) (*model.Filter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilter", arg0, arg1)
	ret0, _ := ret[0].(*model.Filter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilter indicates an expected call of GetFilter.
func (mr *MockdatastoreMockRecorder) GetFilter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilter", reflect.TypeOf((*Mockdatastore)(nil).GetFilter), arg0, arg1)
}

// GetFilterDimension mocks base method.
func (m *Mockdatastore) GetFilterDimension(ctx context.Context, fID, dimName string) (model.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilterDimension", ctx, fID, dimName)
	ret0, _ := ret[0].(model.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilterDimension indicates an expected call of GetFilterDimension.
func (mr *MockdatastoreMockRecorder) GetFilterDimension(ctx, fID, dimName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilterDimension", reflect.TypeOf((*Mockdatastore)(nil).GetFilterDimension), ctx, fID, dimName)
}

// GetFilterDimensionOptions mocks base method.
func (m *Mockdatastore) GetFilterDimensionOptions(arg0 context.Context, arg1, arg2 string, arg3, arg4 int) ([]string, int, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilterDimensionOptions", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetFilterDimensionOptions indicates an expected call of GetFilterDimensionOptions.
func (mr *MockdatastoreMockRecorder) GetFilterDimensionOptions(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilterDimensionOptions", reflect.TypeOf((*Mockdatastore)(nil).GetFilterDimensionOptions), arg0, arg1, arg2, arg3, arg4)
}

// GetFilterDimensions mocks base method.
func (m *Mockdatastore) GetFilterDimensions(arg0 context.Context, arg1 string, arg2, arg3 int) ([]model.Dimension, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilterDimensions", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]model.Dimension)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFilterDimensions indicates an expected call of GetFilterDimensions.
func (mr *MockdatastoreMockRecorder) GetFilterDimensions(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilterDimensions", reflect.TypeOf((*Mockdatastore)(nil).GetFilterDimensions), arg0, arg1, arg2, arg3)
}

// GetFilterOutput mocks base method.
func (m *Mockdatastore) GetFilterOutput(arg0 context.Context, arg1 string) (*model.FilterOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilterOutput", arg0, arg1)
	ret0, _ := ret[0].(*model.FilterOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilterOutput indicates an expected call of GetFilterOutput.
func (mr *MockdatastoreMockRecorder) GetFilterOutput(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilterOutput", reflect.TypeOf((*Mockdatastore)(nil).GetFilterOutput), arg0, arg1)
}

// RemoveFilterDimensionOption mocks base method.
func (m *Mockdatastore) RemoveFilterDimensionOption(ctx context.Context, filterID, dimension, option, currentETag string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFilterDimensionOption", ctx, filterID, dimension, option, currentETag)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveFilterDimensionOption indicates an expected call of RemoveFilterDimensionOption.
func (mr *MockdatastoreMockRecorder) RemoveFilterDimensionOption(ctx, filterID, dimension, option, currentETag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFilterDimensionOption", reflect.TypeOf((*Mockdatastore)(nil).RemoveFilterDimensionOption), ctx, filterID, dimension, option, currentETag)
}

// UpdateFilterDimension mocks base method.
func (m *Mockdatastore) UpdateFilterDimension(ctx context.Context, filterID, dimensionName string, dimension model.Dimension, currentETag string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFilterDimension", ctx, filterID, dimensionName, dimension, currentETag)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFilterDimension indicates an expected call of UpdateFilterDimension.
func (mr *MockdatastoreMockRecorder) UpdateFilterDimension(ctx, filterID, dimensionName, dimension, currentETag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFilterDimension", reflect.TypeOf((*Mockdatastore)(nil).UpdateFilterDimension), ctx, filterID, dimensionName, dimension, currentETag)
}

// UpdateFilterOutput mocks base method.
func (m *Mockdatastore) UpdateFilterOutput(arg0 context.Context, arg1 *model.FilterOutput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFilterOutput", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFilterOutput indicates an expected call of UpdateFilterOutput.
func (mr *MockdatastoreMockRecorder) UpdateFilterOutput(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFilterOutput", reflect.TypeOf((*Mockdatastore)(nil).UpdateFilterOutput), arg0, arg1)
}

// Mockvalidator is a mock of validator interface.
type Mockvalidator struct {
	ctrl     *gomock.Controller
	recorder *MockvalidatorMockRecorder
}

// MockvalidatorMockRecorder is the mock recorder for Mockvalidator.
type MockvalidatorMockRecorder struct {
	mock *Mockvalidator
}

// NewMockvalidator creates a new mock instance.
func NewMockvalidator(ctrl *gomock.Controller) *Mockvalidator {
	mock := &Mockvalidator{ctrl: ctrl}
	mock.recorder = &MockvalidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockvalidator) EXPECT() *MockvalidatorMockRecorder {
	return m.recorder
}

// Valid mocks base method.
func (m *Mockvalidator) Valid() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Valid")
	ret0, _ := ret[0].(error)
	return ret0
}

// Valid indicates an expected call of Valid.
func (mr *MockvalidatorMockRecorder) Valid() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Valid", reflect.TypeOf((*Mockvalidator)(nil).Valid))
}

// Mockgenerator is a mock of generator interface.
type Mockgenerator struct {
	ctrl     *gomock.Controller
	recorder *MockgeneratorMockRecorder
}

// MockgeneratorMockRecorder is the mock recorder for Mockgenerator.
type MockgeneratorMockRecorder struct {
	mock *Mockgenerator
}

// NewMockgenerator creates a new mock instance.
func NewMockgenerator(ctrl *gomock.Controller) *Mockgenerator {
	mock := &Mockgenerator{ctrl: ctrl}
	mock.recorder = &MockgeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockgenerator) EXPECT() *MockgeneratorMockRecorder {
	return m.recorder
}

// Timestamp mocks base method.
func (m *Mockgenerator) Timestamp() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Timestamp")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Timestamp indicates an expected call of Timestamp.
func (mr *MockgeneratorMockRecorder) Timestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Timestamp", reflect.TypeOf((*Mockgenerator)(nil).Timestamp))
}

// URL mocks base method.
func (m *Mockgenerator) URL(host, path string, args ...interface{}) string {
	m.ctrl.T.Helper()
	varargs := []interface{}{host, path}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "URL", varargs...)
	ret0, _ := ret[0].(string)
	return ret0
}

// URL indicates an expected call of URL.
func (mr *MockgeneratorMockRecorder) URL(host, path interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{host, path}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "URL", reflect.TypeOf((*Mockgenerator)(nil).URL), varargs...)
}

// UniqueTimestamp mocks base method.
func (m *Mockgenerator) UniqueTimestamp() primitive.Timestamp {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UniqueTimestamp")
	ret0, _ := ret[0].(primitive.Timestamp)
	return ret0
}

// UniqueTimestamp indicates an expected call of UniqueTimestamp.
func (mr *MockgeneratorMockRecorder) UniqueTimestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UniqueTimestamp", reflect.TypeOf((*Mockgenerator)(nil).UniqueTimestamp))
}

// MockcantabularClient is a mock of cantabularClient interface.
type MockcantabularClient struct {
	ctrl     *gomock.Controller
	recorder *MockcantabularClientMockRecorder
}

// MockcantabularClientMockRecorder is the mock recorder for MockcantabularClient.
type MockcantabularClientMockRecorder struct {
	mock *MockcantabularClient
}

// NewMockcantabularClient creates a new mock instance.
func NewMockcantabularClient(ctrl *gomock.Controller) *MockcantabularClient {
	mock := &MockcantabularClient{ctrl: ctrl}
	mock.recorder = &MockcantabularClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcantabularClient) EXPECT() *MockcantabularClientMockRecorder {
	return m.recorder
}

// GetArea mocks base method.
func (m *MockcantabularClient) GetArea(arg0 context.Context, arg1 cantabular.GetAreaRequest) (*cantabular.GetAreaResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArea", arg0, arg1)
	ret0, _ := ret[0].(*cantabular.GetAreaResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArea indicates an expected call of GetArea.
func (mr *MockcantabularClientMockRecorder) GetArea(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArea", reflect.TypeOf((*MockcantabularClient)(nil).GetArea), arg0, arg1)
}

// GetDimensionOptions mocks base method.
func (m *MockcantabularClient) GetDimensionOptions(arg0 context.Context, arg1 cantabular.GetDimensionOptionsRequest) (*cantabular.GetDimensionOptionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDimensionOptions", arg0, arg1)
	ret0, _ := ret[0].(*cantabular.GetDimensionOptionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDimensionOptions indicates an expected call of GetDimensionOptions.
func (mr *MockcantabularClientMockRecorder) GetDimensionOptions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDimensionOptions", reflect.TypeOf((*MockcantabularClient)(nil).GetDimensionOptions), arg0, arg1)
}

// GetDimensionsByName mocks base method.
func (m *MockcantabularClient) GetDimensionsByName(arg0 context.Context, arg1 cantabular.GetDimensionsByNameRequest) (*cantabular.GetDimensionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDimensionsByName", arg0, arg1)
	ret0, _ := ret[0].(*cantabular.GetDimensionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDimensionsByName indicates an expected call of GetDimensionsByName.
func (mr *MockcantabularClientMockRecorder) GetDimensionsByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDimensionsByName", reflect.TypeOf((*MockcantabularClient)(nil).GetDimensionsByName), arg0, arg1)
}

// GetGeographyDimensionsInBatches mocks base method.
func (m *MockcantabularClient) GetGeographyDimensionsInBatches(ctx context.Context, datasetID string, batchSize, maxWorkers int) (*gql.Dataset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGeographyDimensionsInBatches", ctx, datasetID, batchSize, maxWorkers)
	ret0, _ := ret[0].(*gql.Dataset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGeographyDimensionsInBatches indicates an expected call of GetGeographyDimensionsInBatches.
func (mr *MockcantabularClientMockRecorder) GetGeographyDimensionsInBatches(ctx, datasetID, batchSize, maxWorkers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGeographyDimensionsInBatches", reflect.TypeOf((*MockcantabularClient)(nil).GetGeographyDimensionsInBatches), ctx, datasetID, batchSize, maxWorkers)
}

// StaticDatasetQuery mocks base method.
func (m *MockcantabularClient) StaticDatasetQuery(arg0 context.Context, arg1 cantabular.StaticDatasetQueryRequest) (*cantabular.StaticDatasetQuery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StaticDatasetQuery", arg0, arg1)
	ret0, _ := ret[0].(*cantabular.StaticDatasetQuery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StaticDatasetQuery indicates an expected call of StaticDatasetQuery.
func (mr *MockcantabularClientMockRecorder) StaticDatasetQuery(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StaticDatasetQuery", reflect.TypeOf((*MockcantabularClient)(nil).StaticDatasetQuery), arg0, arg1)
}

// StatusCode mocks base method.
func (m *MockcantabularClient) StatusCode(arg0 error) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StatusCode", arg0)
	ret0, _ := ret[0].(int)
	return ret0
}

// StatusCode indicates an expected call of StatusCode.
func (mr *MockcantabularClientMockRecorder) StatusCode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StatusCode", reflect.TypeOf((*MockcantabularClient)(nil).StatusCode), arg0)
}

// MockdatasetAPIClient is a mock of datasetAPIClient interface.
type MockdatasetAPIClient struct {
	ctrl     *gomock.Controller
	recorder *MockdatasetAPIClientMockRecorder
}

// MockdatasetAPIClientMockRecorder is the mock recorder for MockdatasetAPIClient.
type MockdatasetAPIClientMockRecorder struct {
	mock *MockdatasetAPIClient
}

// NewMockdatasetAPIClient creates a new mock instance.
func NewMockdatasetAPIClient(ctrl *gomock.Controller) *MockdatasetAPIClient {
	mock := &MockdatasetAPIClient{ctrl: ctrl}
	mock.recorder = &MockdatasetAPIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockdatasetAPIClient) EXPECT() *MockdatasetAPIClientMockRecorder {
	return m.recorder
}

// GetDatasetCurrentAndNext mocks base method.
func (m *MockdatasetAPIClient) GetDatasetCurrentAndNext(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (dataset.Dataset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDatasetCurrentAndNext", ctx, userAuthToken, serviceAuthToken, collectionID, datasetID)
	ret0, _ := ret[0].(dataset.Dataset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDatasetCurrentAndNext indicates an expected call of GetDatasetCurrentAndNext.
func (mr *MockdatasetAPIClientMockRecorder) GetDatasetCurrentAndNext(ctx, userAuthToken, serviceAuthToken, collectionID, datasetID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDatasetCurrentAndNext", reflect.TypeOf((*MockdatasetAPIClient)(nil).GetDatasetCurrentAndNext), ctx, userAuthToken, serviceAuthToken, collectionID, datasetID)
}

// GetMetadataURL mocks base method.
func (m *MockdatasetAPIClient) GetMetadataURL(id, edition, version string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetadataURL", id, edition, version)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetMetadataURL indicates an expected call of GetMetadataURL.
func (mr *MockdatasetAPIClientMockRecorder) GetMetadataURL(id, edition, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetadataURL", reflect.TypeOf((*MockdatasetAPIClient)(nil).GetMetadataURL), id, edition, version)
}

// GetOptionsInBatches mocks base method.
func (m *MockdatasetAPIClient) GetOptionsInBatches(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension string, batchSize, maxWorkers int) (dataset.Options, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOptionsInBatches", ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension, batchSize, maxWorkers)
	ret0, _ := ret[0].(dataset.Options)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOptionsInBatches indicates an expected call of GetOptionsInBatches.
func (mr *MockdatasetAPIClientMockRecorder) GetOptionsInBatches(ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension, batchSize, maxWorkers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOptionsInBatches", reflect.TypeOf((*MockdatasetAPIClient)(nil).GetOptionsInBatches), ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension, batchSize, maxWorkers)
}

// GetVersion mocks base method.
func (m *MockdatasetAPIClient) GetVersion(ctx context.Context, userAuthToken, svcAuthToken, downloadSvcAuthToken, collectionID, datasetID, edition, version string) (dataset.Version, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersion", ctx, userAuthToken, svcAuthToken, downloadSvcAuthToken, collectionID, datasetID, edition, version)
	ret0, _ := ret[0].(dataset.Version)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVersion indicates an expected call of GetVersion.
func (mr *MockdatasetAPIClientMockRecorder) GetVersion(ctx, userAuthToken, svcAuthToken, downloadSvcAuthToken, collectionID, datasetID, edition, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersion", reflect.TypeOf((*MockdatasetAPIClient)(nil).GetVersion), ctx, userAuthToken, svcAuthToken, downloadSvcAuthToken, collectionID, datasetID, edition, version)
}

// Mockcoder is a mock of coder interface.
type Mockcoder struct {
	ctrl     *gomock.Controller
	recorder *MockcoderMockRecorder
}

// MockcoderMockRecorder is the mock recorder for Mockcoder.
type MockcoderMockRecorder struct {
	mock *Mockcoder
}

// NewMockcoder creates a new mock instance.
func NewMockcoder(ctrl *gomock.Controller) *Mockcoder {
	mock := &Mockcoder{ctrl: ctrl}
	mock.recorder = &MockcoderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockcoder) EXPECT() *MockcoderMockRecorder {
	return m.recorder
}

// Code mocks base method.
func (m *Mockcoder) Code() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Code")
	ret0, _ := ret[0].(int)
	return ret0
}

// Code indicates an expected call of Code.
func (mr *MockcoderMockRecorder) Code() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Code", reflect.TypeOf((*Mockcoder)(nil).Code))
}
