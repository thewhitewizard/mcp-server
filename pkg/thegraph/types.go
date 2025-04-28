package thegraph

// #region Common struct
type Owner struct {
	ID string `json:"id,omitempty"`
}

// #endregion

// #region Voucher struct
type VoucherResponse struct {
	Data VoucherData `json:"data,omitempty"`
}
type VoucherType struct {
	ID   string `json:"id,omitempty"`
	Desc string `json:"description,omitempty"`
}

type Voucher struct {
	VoucherType VoucherType `json:"voucherType,omitempty"`
	ID          string      `json:"id,omitempty"`
	Owner       Owner       `json:"owner,omitempty"`
	Expiration  string      `json:"expiration,omitempty"`
	Value       string      `json:"value,omitempty"`
	Balance     string      `json:"balance,omitempty"`
}
type VoucherData struct {
	Vouchers []Voucher `json:"vouchers,omitempty"`
}

// #endregion
