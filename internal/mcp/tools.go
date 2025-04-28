package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/thewhitewizard/thegraph-mcp-server/pkg/thegraph"
)

// RegisterTools registers all the TheGraph API tools with the MCP server
func RegisterTools(s *server.MCPServer, client *thegraph.Client) {
	// 1. GetVouchers
	voucherTool := mcp.NewTool("getVouchers",
		mcp.WithDescription("Get Vouchers"),
		mcp.WithString("owner",
			mcp.Description("The Owner of the voucher (optionnal, return all if empty)"),
		),
	)
	s.AddTool(voucherTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetVouchers(ctx, request, client)
	})
}
