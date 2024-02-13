//go:build unit

package natsmanager_test

import (
	"context"
	"errors"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
)

func (s *NatsManagerTestSuite) TestCreateGlobalKeyValueStore() {
	var (
		ctx       = context.Background()
		expected  = "global-kv-store"
		clientReq = &natspb.CreateGlobalKeyValueStoreRequest{
			ProductId: productID,
		}
	)

	s.mockService.EXPECT().CreateGlobalKeyValueStore(ctx, clientReq).
		Return(&natspb.CreateGlobalKeyValueStoreResponse{GlobalKeyValueStore: expected}, nil)

	actual, err := s.natsManagerClient.CreateGlobalKeyValueStore(ctx, productID)
	s.Require().NoError(err)
	s.Assert().Equal(expected, actual)
}

func (s *NatsManagerTestSuite) TestCreateGlobalKeyValueStore_ServiceError() {
	var (
		ctx           = context.Background()
		expectedError = errors.New("service error")
		clientReq     = &natspb.CreateGlobalKeyValueStoreRequest{
			ProductId: productID,
		}
	)

	s.mockService.EXPECT().CreateGlobalKeyValueStore(ctx, clientReq).
		Return(nil, expectedError)

	_, err := s.natsManagerClient.CreateGlobalKeyValueStore(ctx, productID)
	s.Assert().ErrorIs(err, expectedError)
}

func (s *NatsManagerTestSuite) TestDeleteGlobalKeyValueStore() {
	var (
		ctx       = context.Background()
		clientReq = &natspb.DeleteGlobalKeyValueStoreRequest{
			ProductId: productID,
		}
	)

	s.mockService.EXPECT().DeleteGlobalKeyValueStore(ctx, clientReq).
		Return(&natspb.DeleteResponse{Message: "Test message"}, nil)

	err := s.natsManagerClient.DeleteGlobalKeyValueStore(ctx, productID)
	s.Require().NoError(err)
}

func (s *NatsManagerTestSuite) TestDeleteGlobalKeyValueStore_ServiceError() {
	var (
		ctx       = context.Background()
		clientReq = &natspb.DeleteGlobalKeyValueStoreRequest{
			ProductId: productID,
		}
		expectedErr = errors.New("service error")
	)

	s.mockService.EXPECT().DeleteGlobalKeyValueStore(ctx, clientReq).
		Return(nil, expectedErr)

	err := s.natsManagerClient.DeleteGlobalKeyValueStore(ctx, productID)
	s.Require().ErrorIs(err, expectedErr)
}
