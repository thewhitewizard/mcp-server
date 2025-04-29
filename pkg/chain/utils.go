package chain

import "regexp"

const (
	BELLECOUR_PROXY_ADDR = "0x3eca1b216a7df1c7689aeb259ffb83adfb894e7f"
	DECIMAL_18           = 18
	DECIMAL_9            = 9
)

func IsValidEthereumAddressWithChecksum(address string) bool {
	return regexp.MustCompile(wallet_regex).MatchString(address)
}
