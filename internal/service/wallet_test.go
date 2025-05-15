package service_test

import (
	"context"
	"errors"
	"java-code/internal/domain"
	"java-code/internal/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockWalletRepository struct {
	mock.Mock
}

type MockTxManager struct {
	mock.Mock
}

func (m *MockWalletRepository) DepositBalance(ctx context.Context, op domain.WalletOperation) error {
	args := m.Called(ctx, op)
	return args.Error(0)
}

func (m *MockWalletRepository) WithdrawBalance(ctx context.Context, op domain.WalletOperation) error {
	args := m.Called(ctx, op)
	return args.Error(0)
}

func (m *MockWalletRepository) GetWalletBalance(ctx context.Context, id uuid.UUID) (domain.Wallet, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Wallet), args.Error(1)
}

func (m *MockTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {

	return fn(ctx)
}

func TestOperationWithWallet(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name        string
		op          domain.WalletOperation
		setupMock   func(*MockWalletRepository)
		expectedErr error
	}{
		{
			name: "succes deposit",
			op: domain.WalletOperation{
				WalletID:      id,
				OperationType: domain.Deposit,
				Amount:        500,
			},
			setupMock: func(m *MockWalletRepository) {
				m.On("DepositBalance", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "deposit fails",
			op: domain.WalletOperation{
				WalletID:      id,
				OperationType: domain.Deposit,
				Amount:        1000,
			},
			setupMock: func(m *MockWalletRepository) {
				m.On("DepositBalance", mock.Anything, mock.Anything).Return(domain.Errdb)
			},
			expectedErr: domain.Errdb,
		},
		{
			name: "invalid type operation",
			op: domain.WalletOperation{
				WalletID:      id,
				OperationType: "INVALID_TYPE",
				Amount:        500,
			},
			setupMock: func(m *MockWalletRepository) {},

			expectedErr: domain.ErrInvalidOperation,
		},
		{
			name: "withdraw succes",
			op: domain.WalletOperation{
				WalletID:      id,
				OperationType: domain.Withdraw,
				Amount:        300,
			},
			setupMock: func(m *MockWalletRepository) {
				m.On("GetWalletBalance", mock.Anything, id).Return(domain.Wallet{Balance: 1000}, nil)
				m.On("WithdrawBalance", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "withdraw fail",
			op: domain.WalletOperation{
				WalletID:      id,
				OperationType: domain.Withdraw,
				Amount:        300,
			},
			setupMock: func(m *MockWalletRepository) {
				m.On("GetWalletBalance", mock.Anything, id).Return(domain.Wallet{Balance: 1000}, nil)
				m.On("WithdrawBalance", mock.Anything, mock.Anything).Return(domain.Errdb)
			},
			expectedErr: domain.Errdb,
		},
		{
			name: "check amount =0",
			op: domain.WalletOperation{
				WalletID:      id,
				OperationType: domain.Deposit,
				Amount:        0,
			},
			setupMock: func(m *MockWalletRepository) {
				m.AssertNotCalled(t, "DepositBalance", mock.Anything, mock.Anything)
			},
			expectedErr: domain.ErrAmout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockWalletRepository)
			mockTxManager := new(MockTxManager)

			tt.setupMock(mockRepo)
			svc := service.NewWallet(mockRepo, mockTxManager)

			err := svc.OperationWithWallet(context.Background(), tt.op)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

