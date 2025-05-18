package handlers

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethpandaops/dora/clients/execution"
	"github.com/ethpandaops/dora/clients/execution/rpc"
	"github.com/ethpandaops/dora/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock RPC client
type mockRPCClient struct {
	mock.Mock
	execution.Client
}

func (m *mockRPCClient) GetRPCClient() *rpc.ExecutionClient {
	args := m.Called()
	return args.Get(0).(*rpc.ExecutionClient)
}

func (m *mockRPCClient) GetLastHead() (uint64, common.Hash) {
	return 0, common.Hash{}
}

func (m *mockRPCClient) GetName() string {
	return "mock"
}

func (m *mockRPCClient) GetEndpoint() string {
	return "mock"
}

func (m *mockRPCClient) GetType() execution.ClientType {
	return execution.ClientType(0)
}

// Mock RPC
type mockRPC struct {
	mock.Mock
}

func (m *mockRPC) GetEthClient() *ethclient.Client {
	args := m.Called()
	return args.Get(0).(*ethclient.Client)
}

// Mock eth client
type mockEthClient struct {
	mock.Mock
	ethclient.Client
}

func (m *mockEthClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	args := m.Called(ctx, account, blockNumber)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *mockEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	args := m.Called(ctx, number)
	return args.Get(0).(*types.Block), args.Error(1)
}

func (m *mockEthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	args := m.Called(ctx, txHash)
	return args.Get(0).(*types.Receipt), args.Error(1)
}

func (m *mockEthClient) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	args := m.Called(ctx, tx, block, index)
	return args.Get(0).(common.Address), args.Error(1)
}

func TestHandleAccountBalance(t *testing.T) {
	// Create mock clients
	mockRPC := new(mockRPCClient)
	mockRPCClient := new(mockRPC)
	mockEth := new(mockEthClient)

	// Setup test address and balance
	testAddr := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
	testBalance := big.NewInt(1000000000000000000) // 1 ETH

	// Setup mock expectations
	mockRPC.On("GetRPCClient").Return(mockRPCClient)
	mockRPCClient.On("GetEthClient").Return(mockEth)
	mockEth.On("BalanceAt", mock.Anything, testAddr, (*big.Int)(nil)).Return(testBalance, nil)

	// Create test request
	req := httptest.NewRequest("GET", "/api?module=account&action=balance&address="+testAddr.Hex(), nil)
	w := httptest.NewRecorder()

	// Setup global service with mock client
	services.GlobalBeaconService = &services.ChainService{
		ExecutionPool: &execution.Pool{},
	}
	services.GlobalBeaconService.ExecutionPool.AddEndpoint(mockRPC)

	// Call handler
	EtherscanAPI(w, req)

	// Parse response
	var resp EtherscanResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)

	// Verify response
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Message)
	assert.Equal(t, testBalance.String(), resp.Result)

	// Verify mock expectations
	mockRPC.AssertExpectations(t)
	mockRPCClient.AssertExpectations(t)
	mockEth.AssertExpectations(t)
}

func TestHandleTxList(t *testing.T) {
	// Create mock clients
	mockRPC := new(mockRPCClient)
	mockRPCClient := new(mockRPC)
	mockEth := new(mockEthClient)

	// Setup test data
	testAddr := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
	latestBlock := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(1000),
		Time:   1234567890,
	})
	testTx := types.NewTransaction(0, testAddr, big.NewInt(1000000000000000000), 21000, big.NewInt(1000000000), nil)
	testReceipt := &types.Receipt{
		TxHash:            testTx.Hash(),
		TransactionIndex: 0,
		BlockHash:        latestBlock.Hash(),
		BlockNumber:      latestBlock.Number(),
	}
	senderAddr := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44f")

	// Setup mock expectations
	mockRPC.On("GetRPCClient").Return(mockRPCClient)
	mockRPCClient.On("GetEthClient").Return(mockEth)
	mockEth.On("BlockByNumber", mock.Anything, (*big.Int)(nil)).Return(latestBlock, nil)
	mockEth.On("BlockByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(latestBlock, nil)
	mockEth.On("TransactionReceipt", mock.Anything, mock.AnythingOfType("common.Hash")).Return(testReceipt, nil)
	mockEth.On("TransactionSender", mock.Anything, mock.AnythingOfType("*types.Transaction"), mock.AnythingOfType("common.Hash"), mock.AnythingOfType("uint")).Return(senderAddr, nil)

	// Create test request
	req := httptest.NewRequest("GET", "/api?module=account&action=txlist&address="+testAddr.Hex(), nil)
	w := httptest.NewRecorder()

	// Setup global service with mock client
	services.GlobalBeaconService = &services.ChainService{
		ExecutionPool: &execution.Pool{},
	}
	services.GlobalBeaconService.ExecutionPool.AddEndpoint(mockRPC)

	// Call handler
	EtherscanAPI(w, req)

	// Parse response
	var resp EtherscanResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)

	// Verify response
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Message)
	txs, ok := resp.Result.([]map[string]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, txs)

	// Verify mock expectations
	mockRPC.AssertExpectations(t)
	mockRPCClient.AssertExpectations(t)
	mockEth.AssertExpectations(t)
}
