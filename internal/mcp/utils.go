package mcp

import (
	"fmt"
	"strconv"
	"time"

	"github.com/thewhitewizard/thegraph-mcp-server/pkg/thegraph"
)

func formatVoucher(v thegraph.Voucher) string {
	timestamp, _ := strconv.ParseInt(v.Expiration, 10, 64)
	date := time.Unix(timestamp, 0)

	return fmt.Sprintf("ID=%s Type=%s Owner=%s Value=%s Balance=%s Expiration=%s",
		v.ID, v.VoucherType.Desc, v.Owner.ID, v.Value, v.Balance, date.Format("2006-01-02 15:04:05 MST"))
}
