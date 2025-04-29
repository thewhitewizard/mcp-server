package chain

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

func NewClient(rpcAddr string) *Client {
	client := &Client{url: rpcAddr, down: true}
	client.Connect()

	go client.heartbeat()

	return client
}

func (c *Client) heartbeat() {
	for {
		block, err := c.CurrentBlock(context.Background())

		if err != nil || block == 0 {
			c.down = true
			c.Connect()
		}

		time.Sleep(heartbeat_retry)
	}
}

func (c *Client) Connect() {
	isOK := false

	for {
		if c.conn != nil {
			c.conn.Close()
			c.conn = nil
		}
		client, err := ethclient.Dial(c.url)

		if err == nil {
			c.conn = client
			isOK = true
		}

		if isOK {
			c.down = false
			break
		} else {
			time.Sleep(heartbeat_retry)
		}
	}
}

func (c *Client) Close() {
	c.conn.Close()
}

// CurrentBlock return the current block
func (c *Client) CurrentBlock(ctx context.Context) (uint64, error) {
	return c.conn.BlockNumber(ctx)
}

// GetBalance return the wallet balance
func (c *Client) GetBalance(ctx context.Context, wallet string, decimals int) string {
	if c.down {
		return emptyBalance
	}
	balance, err := c.conn.BalanceAt(ctx, common.HexToAddress(wallet), nil)

	if err == nil {
		return formatBalance(balance, decimals)
	}

	return emptyBalance
}

// GetBalanceForToken return balance for wallet for specific token
func (c *Client) GetBalanceForToken(ctx context.Context, wallet string, tokenAddress string, decimals int) string {
	if c.down {
		return emptyBalance
	}

	caller, err := NewTokenCaller(common.HexToAddress(tokenAddress), c.conn)

	if err != nil {
		return emptyBalance
	}
	balance, err := caller.BalanceOf(nil, common.HexToAddress(wallet))

	if err == nil {
		return formatBalance(balance, decimals)
	}

	return emptyBalance
}

// GetLockRLCBalance return balance for wallet for Lock RLC, very specific for iExec
func (c *Client) GetLockRLCBalance(ctx context.Context, wallet string, tokenAddress string, decimals int) string {
	if c.down {
		return emptyBalance
	}

	caller, err := NewTokenCaller(common.HexToAddress(tokenAddress), c.conn)

	if err != nil {
		return emptyBalance
	}
	balance, err := caller.FrozenOf(nil, common.HexToAddress(wallet))

	if err == nil {
		return formatBalance(balance, decimals)
	}

	return emptyBalance
}

func formatBalance(balance *big.Int, decimals int) string {
	mul := decimal.NewFromFloat(ten).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(balance.String())
	result := num.Div(mul)

	return result.String()
}
