package handlers

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethpandaops/dora/clients/execution"
	"github.com/ethpandaops/dora/services"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRPCClient is a mock implementation of the RPC client
type MockRPCClient struct {
	mock.Mock
}

// CallContext mocks the RPC CallContext method
func (m *MockRPCClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	callArgs := m.Called(append([]interface{}{ctx, result, method}, args...)...)
	return callArgs.Error(0)
}

// MockExecutionClient is a mock implementation of the execution client
type MockExecutionClient struct {
	mock.Mock
	rpcClient *MockRPCClient
	ethClient *ethclient.Client
}

// NewMockExecutionClient creates a new mock execution client
func NewMockExecutionClient() *MockExecutionClient {
	return &MockExecutionClient{
		rpcClient: &MockRPCClient{},
	}
}

// GetRPCClient returns the mock RPC client
func (m *MockExecutionClient) GetRPCClient() *rpc.Client {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*rpc.Client)
}

// GetEthClient returns the mock ETH client
func (m *MockExecutionClient) GetEthClient() *ethclient.Client {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*ethclient.Client)
}

// MockChainService is a mock implementation of the ChainService
type MockChainService struct {
	executionClients []*execution.Client
}

// NewMockChainService creates a new mock chain service
func NewMockChainService() *MockChainService {
	return &MockChainService{
		executionClients: make([]*execution.Client, 0),
	}
}

// GetExecutionClients returns the list of execution clients
func (m *MockChainService) GetExecutionClients() []*execution.Client {
	return m.executionClients
}

// AddExecutionClient adds a new execution client to the mock service
func (m *MockChainService) AddExecutionClient(client *execution.Client) {
	m.executionClients = append(m.executionClients, client)
}

func TestGetTokenBalance(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name           string
		address        string
		contract       string
		setupMocks     func(*MockExecutionClient, *MockChainService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "successful balance check",
			address:  "0x1234567890123456789012345678901234567890",
			contract: "0x0987654321098765432109876543210987654321",
			setupMocks: func(mockClient *MockExecutionClient, mockChain *MockChainService) {
				// Setup mock expectations
				mockClient.On("GetEthClient").Return(&ethclient.Client{})
				mockClient.On("GetRPCClient").Return(&rpc.Client{})
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"1","message":"OK","result":"1000000000000000000"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockClient := NewMockExecutionClient()
			mockChain := NewMockChainService()

			// Setup test request
			req, err := http.NewRequest("GET", "/api?module=account&action=tokenbalance&address="+tt.address+"&contractaddress="+tt.contract, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetTokenBalance)

			// Setup mocks
			tt.setupMocks(mockClient, mockChain)

			// Create a context with the mock chain service
			ctx := context.WithValue(req.Context(), "chain", mockChain)
			req = req.WithContext(ctx)

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Assert the status code is what we expect.
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Assert the response body is what we expect.
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestGetTokenTx(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name           string
		address        string
		contract       string
		setupMocks     func(*MockExecutionClient, *MockChainService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "successful transaction check",
			address:  "0x1234567890123456789012345678901234567890",
			contract: "0x0987654321098765432109876543210987654321",
			setupMocks: func(mockClient *MockExecutionClient, mockChain *MockChainService) {
				// Setup mock expectations
				mockClient.On("GetEthClient").Return(&ethclient.Client{})
				mockClient.On("GetRPCClient").Return(&rpc.Client{})
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"1","message":"OK","result":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockClient := NewMockExecutionClient()
			mockChain := NewMockChainService()

			// Setup test request
			req, err := http.NewRequest("GET", "/api?module=account&action=tokentx&address="+tt.address+"&contractaddress="+tt.contract, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetTokenTx)

			// Setup mocks
			tt.setupMocks(mockClient, mockChain)

			// Create a context with the mock chain service
			ctx := context.WithValue(req.Context(), "chain", mockChain)
			req = req.WithContext(ctx)

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Assert the status code is what we expect.
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Assert the response body is what we expect.
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}
