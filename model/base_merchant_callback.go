package model

// CallbackMerchant
// @Description: 回调给商家的结构
type CallbackMerchant struct {
	Type         int         `json:"type"`
	Code         int         `json:"code"`
	Success      bool        `json:"success"`
	Message      string      `json:"message"`
	MerchantNo   string      `json:"merchant_no"`
	SerialNo     string      `json:"serial_no"`
	RelaSerialNo string      `json:"rela_serial_no"`
	Timestamp    int64       `json:"timestamp"`
	Data         interface{} `json:"data"`
}
