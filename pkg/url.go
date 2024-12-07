package pkg

func GetRequestUrlForHandShake(baseUrl, clientId string) string {
	return baseUrl + "/oauth2/clients/handshake/" + clientId
}

func GetRequestUrlForOtp(baseUrl, phoneNumber string) string {
	return baseUrl + "/oauth2/otp/authorize/" + phoneNumber
}

func GetRequestUrlForVerifyOtp(baseUrl, phoneNumber string) string {
	return baseUrl + "/oauth2/otp/verify/" + phoneNumber
}

func GetRequestUrlForAccessToken(baseUrl string) string {
	return baseUrl + "/oauth2/token"
}

func GetUrlForTokenValidation(baseUrl string) string {
	return baseUrl + "/oauth2/token/info"
}

func GetUrlForDeactivating(baseUrl string) string {
	return baseUrl + "/oauth2/token/revoke"
}

func GetUrlForUserInfo(baseUrl string) string {
	return baseUrl + "/users"
}

func GetUrlForListOfUserInfo(baseUrl string) string {
	return baseUrl + "/users/info/list"
}
