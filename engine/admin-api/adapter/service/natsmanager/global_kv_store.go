package natsmanager

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
)

func (n *Client) CreateGlobalKeyValueStore(ctx context.Context, product string) (string, error) {
	req := natspb.CreateGlobalKeyValueStoreRequest{
		ProductId: product,
	}

	res, err := n.client.CreateGlobalKeyValueStore(ctx, &req)
	if err != nil {
		return "", fmt.Errorf("creating global key-value store: %w", err)
	}

	return res.GlobalKeyValueStore, err
}

func (n *Client) DeleteGlobalKeyValueStore(ctx context.Context, product string) error {
	req := natspb.DeleteGlobalKeyValueStoreRequest{
		ProductId: product,
	}

	_, err := n.client.DeleteGlobalKeyValueStore(ctx, &req)
	if err != nil {
		return fmt.Errorf("creating global key-value store: %w", err)
	}

	return err
}
