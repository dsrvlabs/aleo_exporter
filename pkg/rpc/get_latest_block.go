package rpc

import (
		"context"
		"encoding/json"
		"fmt"

		"k8s.io/klog/v2"

)

type GetLatestBlockResponse struct {
	Header Header `json:"header"`
}

type Header struct {
	PreviousStateRoot string `json:"previous_state_root"`
	TransactionsRoot  string `json:"transactions_root"`
	FinalizeRoot      string `json:"finalize_root"`
	RatificationsRoot string `json:"ratifications_root"`
	SolutionsRoot     string `json:"solutions_root"`
	SubdagRoot        string `json:"subdag_root"`
	Metadata          *Metadata `json:"metadata"`
}

type Metadata struct {
		Network               int   `json:"network"`
		Round                 uint64   `json:"round"`
		Height                uint32   `json:"height"`
		CumulativeWeight      uint64 `json:"cumulative_weight"`
		CumulativeProofTarget uint64 `json:"cumulative_proof_target"`
		CoinbaseTarget        uint64 `json:"coinbase_target"`
		ProofTarget           uint64   `json:"proof_target"`
		LastCoinbaseTarget    uint64 `json:"last_coinbase_target"`
		LastCoinbaseTimestamp uint64   `json:"last_coinbase_timestamp"`
		Timestamp             int64   `json:"timestamp"`
}

// https://developer.aleo.org/testnet/public_endpoints/latest_block
func (c *RPCClient) GetLatestBlock(ctx context.Context) (*Metadata, error) {
	body, err := c.rpcRequest(ctx, "/testnet3/latest/block")

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("latest block: %v", string(body))

	var resp GetLatestBlockResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	metadata := resp.Header.Metadata

	return metadata, nil
}

