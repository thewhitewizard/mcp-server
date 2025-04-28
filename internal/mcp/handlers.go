package mcp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
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

	return mcp.NewToolResultText(result + "\n"), nil
}

func formatVoucher(v thegraph.Voucher) string {
	timestamp, _ := strconv.ParseInt(v.Expiration, 10, 64)
	date := time.Unix(timestamp, 0)

	return fmt.Sprintf("ID=%s Type=%s Owner=%s Value=%s Balance=%s Expiration=%s",
		v.ID, v.VoucherType.Desc, v.Owner.ID, v.Value, v.Balance, date.Format("2006-01-02 15:04:05 MST"))
}
