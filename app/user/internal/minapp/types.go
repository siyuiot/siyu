package minapp

type Minapp struct {
	TokenMp       string `json:"tokenMp"`       // binding:"required"
	AppId         string `json:"appId"`         // binding:"required"
	Code          string `json:"code"`          // binding:"required"
	Iv            string `json:"iv"`            // binding:"required"
	EncryptedData string `json:"encryptedData"` // binding:"required"
}
