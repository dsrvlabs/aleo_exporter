package rpc

import (
	"context"
	"io/ioutil"
	"net/http"
)

type (
	RPCClient struct {
		httpClient http.Client
		rpcAddr    string
	}
)

func NewRPCClient(rpcAddr string) *RPCClient {
	c := &RPCClient{
		httpClient: http.Client{},
		rpcAddr:    rpcAddr,
	}

	return c
}

func (c *RPCClient) rpcRequest(ctx context.Context, endpoint string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.rpcAddr+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
