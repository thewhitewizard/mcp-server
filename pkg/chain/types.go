package chain

import (
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	emptyBalance    = "-"
	heartbeat_retry = 10 * time.Second
	wallet_regex    = `^0x[a-fA-F0-9]{40}$`
	ten             = 10.0
	DEFAULT_URL     = "https://bellecour.iex.ec/"
)

type Client struct {
	conn *ethclient.Client
	url  string
	down bool
}
