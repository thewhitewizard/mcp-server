package mcp

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/thewhitewizard/thegraph-mcp-server/pkg/chain"
	"github.com/thewhitewizard/thegraph-mcp-server/pkg/thegraph"
)

// Handler functions
func handleGetVouchers(_ context.Context, request mcp.CallToolRequest, client *thegraph.Client) (*mcp.CallToolResult, error) {
	vouchers, err := client.GetVouchers()
	if err != nil {
		return nil, err
	}
	result := ""
	owner, _ := request.Params.Arguments["owner"].(string)

	for _, v := range vouchers.Data.Vouchers {
		if owner == "" || strings.EqualFold(owner, v.Owner.ID) {
			result += formatVoucher(v) + "\n"
		}
	}

	return mcp.NewToolResultText(result), nil
}

func handleGetLastBlock(ctx context.Context, _ mcp.CallToolRequest, client *chain.Client) (*mcp.CallToolResult, error) {
	block, err := client.CurrentBlock(ctx)
	if err != nil {
		return nil, err
	}
	result := strconv.FormatUint(block, 10)

	return mcp.NewToolResultText(result), nil
}

func handleWalletInfo(ctx context.Context, request mcp.CallToolRequest, client *chain.Client) (*mcp.CallToolResult, error) {
	wallet, _ := request.Params.Arguments["wallet"].(string)

	balanceXRLC := client.GetBalance(ctx, wallet, chain.DECIMAL_18)
	balanceSRLC := client.GetBalanceForToken(ctx, wallet, chain.BELLECOUR_PROXY_ADDR, chain.DECIMAL_9)
	balanceLRLC := client.GetLockRLCBalance(ctx, wallet, chain.BELLECOUR_PROXY_ADDR, chain.DECIMAL_9)

	result := fmt.Sprintf("xRLC=%s, sRLC:%s, lockRLC=%s", balanceXRLC, balanceSRLC, balanceLRLC)

	return mcp.NewToolResultText(result), nil
}
