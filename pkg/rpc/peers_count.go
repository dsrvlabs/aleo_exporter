package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

// https://developer.aleo.org/testnet/public_endpoints/get_peers_count
func (c *RPCClient) GetPeerCount(ctx context.Context) (uint32, error) {
	body, err := c.rpcRequest(ctx, "/testnet3/peers/count")

	if body == nil {
		return 0, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return 0, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("peer count: %v", string(body))

	var resp uint32
	if err = json.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("failed to decode response body: %w", err)
	}

	return resp, nil
}
