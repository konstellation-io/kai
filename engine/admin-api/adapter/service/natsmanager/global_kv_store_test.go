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
