package minappDataDecrypt

type Minapp struct {
	AppId         string `json:"appId" binding:"required"`
	Code          string `json:"code" binding:"required"`
	Iv            string `json:"iv" binding:"required"`
	EncryptedData string `json:"encryptedData" binding:"required"`
}

type PhoneNumber struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}
