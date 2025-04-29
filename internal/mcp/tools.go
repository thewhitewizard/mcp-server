package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/thewhitewizard/thegraph-mcp-server/pkg/chain"
	"github.com/thewhitewizard/thegraph-mcp-server/pkg/thegraph"
)

// RegisterTools registers all the TheGraph API tools with the MCP server
func RegisterTools(s *server.MCPServer, thegraphClient *thegraph.Client, chainClient *chain.Client) {
	// 1. GetVouchers
	getVouchers := mcp.NewTool("getVouchers",
		mcp.WithDescription("Get Vouchers"),
		mcp.WithString("owner",
			mcp.Description("The Owner of the voucher (optionnal, return all if empty)"),
		),
	)
	s.AddTool(getVouchers, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetVouchers(ctx, request, thegraphClient)
	})

	// 2. GetLastBlock
	getLastBlockTool := mcp.NewTool("getLastBlock",
		mcp.WithDescription("Get Last Block"),
	)
	s.AddTool(getLastBlockTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetLastBlock(ctx, request, chainClient)
	})

	// 3. getWalletInfo
	getWalletInfo := mcp.NewTool("getWalletInfo",
		mcp.WithDescription("Get info for this wallet"),
		mcp.WithString("wallet",
			mcp.Required(),
			mcp.Description("wallet to fetch info"),
		),
	)
	s.AddTool(getWalletInfo, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleWalletInfo(ctx, request, chainClient)
	})
}
