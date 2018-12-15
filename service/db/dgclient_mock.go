package db

import (
	"github.com/stretchr/testify/mock"
)

// MockDGClient is a mock implementation of data with a Dgraph client for testing purposes
type MockDGClient struct {
	mock.Mock
}

// CreateSchema ...mock
func (s *MockDGClient) CreateSchema(sm Schema) error {
	args := s.Called(sm)
	return args.Error(0)
}

// DropSchema ...mock
func (s *MockDGClient) DropSchema(name string) error {
	args := s.Called(name)
	return args.Error(0)
}

// GetEntity ...mock
func (s *MockDGClient) GetEntity(meta string, uuid string) (map[string]interface{}, error) {
	args := s.Called(meta, uuid)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// GetAllByClusterAndType ...mock
func (s *MockDGClient) GetAllByClusterAndType(meta string, cluster string) (map[string]interface{}, error) {
	args := s.Called(meta, cluster)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// CreateEntity ...mock
func (s *MockDGClient) CreateEntity(meta string, data map[string]interface{}) (map[string]string, error) {
	args := s.Called(meta, data)
	return args.Get(0).(map[string]string), args.Error(1)
}

// DeleteEntity ...mock by uuid
func (s *MockDGClient) DeleteEntity(uuid string) error {
	args := s.Called(uuid)
	return args.Error(0)
}

//CreateOrDeleteEdge ...mock
func (s *MockDGClient) CreateOrDeleteEdge(fromType string, fromUID string, toType string, toUID string, rel string, op Action) error {
	args := s.Called(fromType, fromUID, toType, toUID, rel, op)
	return args.Error(0)
}

// SetFieldToNull ...mock
func (s *MockDGClient) SetFieldToNull(delMap map[string]interface{}) error {
	args := s.Called(delMap)
	return args.Error(0)
}

//UpdateEntity ...mock
func (s *MockDGClient) UpdateEntity(meta string, uuid string, data map[string]interface{}) error {
	args := s.Called(meta, uuid, data)
	return args.Error(0)
}

//GetQueryResult ...mock
func (s *MockDGClient) GetQueryResult(queryMap map[string][]string) (map[string]interface{}, error) {
	args := s.Called(queryMap)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// Close ...mock
func (s *MockDGClient) Close() error {
	args := s.Called()
	return args.Error(0)
}
