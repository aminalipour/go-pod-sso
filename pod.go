package pod

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/aminalipour/go-pod-sso/constanse"
	"github.com/aminalipour/go-pod-sso/errors"
	"github.com/aminalipour/go-pod-sso/pkg"
	"github.com/aminalipour/go-pod-sso/types"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

// making request for handshake
// if uuid already exist for this user ,  pass it through function
// if not pass it like uuid.Nil
func (cfg *Config) SendHandshakeRequest(requestBody types.HandShakeApiAdditionalDataFromClient, deviceUid uuid.UUID) (types.HandShakeResponse, error) {
	//validate the requestbody
	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		return types.HandShakeResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// var uuidForHandshake uuid.UUID

	// if deviceUid != uuid.Nil {
	// 	uuidForHandshake = deviceUid
	// } else {
	// 	uuidForHandshake = uuid.Nil
	// }

	// create the url data
	urlData, err := pkg.GetUrlDataForHandShakeRequest(requestBody, uuid.Nil)
	if err != nil {
		return types.HandShakeResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// create the request base info's
	requestUrl := pkg.GetRequestUrlForHandShake(cfg.BaseUrl, cfg.ClientId)
	headers := map[string]string{
		"Content-Type":  constanse.ContentTypeForUrlDataPod,
		"Authorization": "Bearer " + cfg.PodAccessToken,
	}

	// make the request
	var handShakeResponse types.HandShakeResponse
	err = pkg.MakeRequestWithUrlData(requestUrl, "POST", urlData, headers, &handShakeResponse)
	if err != nil {
		return types.HandShakeResponse{}, err
	}
	return handShakeResponse, nil
}

// making request for otp
// example request type
//
//	requestBodyForOtp := dto.OTPRequestToPodRequestBody{
//		ResponseType:     "code",
//		IdentityType:     "phone_number",
//		ReferrerType:     "username",
//		LinkDeliveryType: "SMS",
//	}
//
// "host: accounts.pod.ir"
func (cfg *Config) SendOtpRequest(requestBody types.OTPRequestToPodRequestBody, keyId, phoneNumber string) (types.OtpResponse, string, error) {

	//validate the requestbody
	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		return types.OtpResponse{}, "", errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// create the url data
	urlData, err := pkg.GetUrlDataFromGivenStruct(requestBody)
	if err != nil {
		return types.OtpResponse{}, "", errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}
	if cfg.Signature == "" && cfg.PrivateKeyFile == "" {
		return types.OtpResponse{}, "", errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrSignatureKeyOrFileIsMissing,
				"errorDescription": "signature key or file missing",
			},
		)
	}

	signature, sigErr := cfg.GetSignature()
	if sigErr != nil {
		return types.OtpResponse{}, "", errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidSignature,
				"errorDescription": errors.Signature,
			},
		)
	}

	// create the request base info's
	requestUrl := pkg.GetRequestUrlForOtp(cfg.BaseUrl, phoneNumber)
	authorizationHeader := fmt.Sprintf("Signature keyId=\"%s\",signature=\"%s\",headers=\"%s\"", keyId, signature, "host")
	headers := map[string]string{
		"Content-Type":  constanse.ContentTypeForUrlDataPod,
		"host":          pkg.GetHostFromURL(cfg.BaseUrl),
		"Authorization": authorizationHeader,
	}

	// make the request
	var otpSendResponse types.OtpResponse
	err = pkg.MakeRequestWithUrlData(requestUrl, "POST", urlData, headers, &otpSendResponse)
	if err != nil {
		return types.OtpResponse{}, "", err
	}
	return otpSendResponse, signature, nil
}

// making request for otp verification
// pass the signiture created from the SendOtpRequest return , dont pass the file path here
func (cfg *Config) MakeRequestForOtpVerify(requestBody types.VerifyOtpRequestBody, keyId string) (types.OtpVerifyResponse, error) {

	//validate the requestbody
	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		return types.OtpVerifyResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// generating url data for the request to pod
	urlDataForOtp, err := pkg.GetUrlDataFromGivenStruct(requestBody)
	if err != nil {
		return types.OtpVerifyResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	signature, sigErr := cfg.GetSignature()
	if sigErr != nil {
		return types.OtpVerifyResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidSignature,
				"errorDescription": errors.Signature,
			},
		)
	}

	// create the request base info's
	phoneNumber := urlDataForOtp.Get("phone_number")
	requestUrl := pkg.GetRequestUrlForVerifyOtp(cfg.BaseUrl, phoneNumber)
	authorizationHeader := fmt.Sprintf("Signature keyId=\"%s\",signature=\"%s\",headers=\"%s\"", keyId, signature, "host")
	headers := map[string]string{
		"Content-Type":  constanse.ContentTypeForUrlDataPod,
		"host":          pkg.GetHostFromURL(cfg.BaseUrl),
		"Authorization": authorizationHeader,
	}

	// make the request
	var otpSendResponse types.OtpVerifyResponse
	err = pkg.MakeRequestWithUrlData(requestUrl, "POST", urlDataForOtp, headers, &otpSendResponse)
	if err != nil {
		return types.OtpVerifyResponse{}, err
	}
	return otpSendResponse, nil
}

// request for getting access token
// eg of rquest body requirments
//
//	requestToGetAccessToken := dto.GetTokenRequestBody{
//		Code:         code from otp verify ,
//		GrantType:    "authorization_code",
//		IdentityType: "phone_number",
//		Idenify:      phone number,
//	}
func (cfg *Config) MakeRequestForGetAccessToken(requestBody types.GetTokenRequestBody) (types.GetAccessTokenResponse, error) {

	//validate the requestbody
	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		return types.GetAccessTokenResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// generating url data for the request to pod
	urlData, err := pkg.GetUrlDataFromGivenStruct(requestBody)
	if err != nil {
		return types.GetAccessTokenResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// create the request base info's
	requestUrl := pkg.GetRequestUrlForAccessToken(cfg.BaseUrl)
	authorizationHeader := "Basic " + pkg.BasicAuth(cfg.ClientId, cfg.ClientSecret)
	headers := map[string]string{
		"Content-Type":  constanse.ContentTypeForUrlDataPod,
		"host":          pkg.GetHostFromURL(cfg.BaseUrl),
		"Authorization": authorizationHeader,
	}

	// make the request
	var tokenResponse types.GetAccessTokenResponse
	err = pkg.MakeRequestWithUrlData(requestUrl, "POST", urlData, headers, &tokenResponse)
	if err != nil {
		return types.GetAccessTokenResponse{}, err
	}
	return tokenResponse, nil
}

// request for refresh token
// example of request body requirments
//
//	requestToGetRefreshToken := dto.GetTokenRequestBody{
//		GrantType:    "refresh_token",
//		RefreshToken: refreshToken,
//	}
func (cfg *Config) MakeRequestForRefreshToken(requestBody types.GetTokenRequestBody) (types.GetAccessTokenResponse, error) {
	//validate the requestbody
	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		return types.GetAccessTokenResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// generating url data for the token grant
	urlDataForToken, err := pkg.GetUrlDataFromGivenStruct(requestBody)
	if err != nil {
		return types.GetAccessTokenResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// create the request base info's
	requestUrl := pkg.GetRequestUrlForAccessToken(cfg.BaseUrl)
	authorizationHeader := "Basic " + pkg.BasicAuth(cfg.ClientId, cfg.ClientSecret)
	headers := map[string]string{
		"Content-Type":  constanse.ContentTypeForUrlDataPod,
		"host":          pkg.GetHostFromURL(cfg.BaseUrl),
		"Authorization": authorizationHeader,
	}

	// make the request
	var tokenResponse types.GetAccessTokenResponse
	err = pkg.MakeRequestWithUrlData(requestUrl, "POST", urlDataForToken, headers, &tokenResponse)
	if err != nil {
		return types.GetAccessTokenResponse{}, err
	}
	return tokenResponse, nil
}

//	request for deactivating token
//
// requires level 4 access to pod
// Not Fully tested cause of the lack of the access , so it the
// implemenation may have bug
// request body for deactivating access token :
//
//	requestToGetRefreshToken := dto.DenyPermissionRequestBody{
//		Token:         requestBody.AccessToken,
//		TokenTypeHint: "access_token",
//	}
//
// request body for the refresh token :
//
//	requestToGetRefreshToken := dto.DenyPermissionRequestBody{
//		Token:         refreshToken,
//		TokenTypeHint: "refresh_token",
//	}
func (cfg *Config) MakeRequestForDeactivingToken(requestBody types.DenyPermissionRequestBody) (interface{}, error) {

	// validate the body with certain validation set to the struct
	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		return types.AccessTokenProcess{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// generating the url data
	urlDataForToken, err := pkg.GetUrlDataFromGivenStruct(requestBody)
	if err != nil {
		return types.AccessTokenProcess{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// create the request base info's
	requestUrl := pkg.GetUrlForDeactivating(cfg.BaseUrl)
	authorizationHeader := "Basic " + pkg.BasicAuth(cfg.ClientId, cfg.ClientSecret)
	headers := map[string]string{
		"Content-Type":  constanse.ContentTypeForUrlDataPod,
		"Authorization": authorizationHeader,
	}

	// make the request
	var tokenResponse interface{}
	err = pkg.MakeRequestWithUrlData(requestUrl, "POST", urlDataForToken, headers, &tokenResponse)
	if err != nil {
		return types.AccessTokenProcess{}, err
	}
	return tokenResponse, nil
}

// request for token validation check to pod sso
func (cfg *Config) MakeRequestForTokenValidation(requestBody types.AccessTokenProcess) (types.ValidationResponseFromPod, error) {

	// validate the body with certain validation set to the struct
	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		return types.ValidationResponseFromPod{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// generating url data for the request to pod
	urlDataForValidationOfToken := pkg.GetUrlDataForTokenValidationRequest(requestBody)

	// create the request base info's
	requestUrl := pkg.GetUrlForTokenValidation(cfg.BaseUrl)
	authorizationHeader := "Basic " + pkg.BasicAuth(cfg.ClientId, cfg.ClientSecret)
	headers := map[string]string{
		"Content-Type":  constanse.ContentTypeForUrlDataPod,
		"Authorization": authorizationHeader,
	}

	// make the request
	var tokenResponse types.ValidationResponseFromPod
	err := pkg.MakeRequestWithUrlData(requestUrl, "POST", urlDataForValidationOfToken, headers, &tokenResponse)
	if err != nil {
		return types.ValidationResponseFromPod{}, err
	}
	return tokenResponse, nil
}

// request for getting user info
func (cfg *Config) MakeRequestForUserInfo(accessToken string) (types.UserInfoFromPod, error) {

	// geting the url for request
	requestUrl := pkg.GetUrlForUserInfo(cfg.BaseUrl)

	// create the request base info's
	authorizationHeader := "Bearer " + accessToken
	headers := map[string]string{
		"Authorization": authorizationHeader,
	}

	// make the request
	var res types.UserInfoFromPod
	err := pkg.MakeRequestWithNoBody(requestUrl, "GET", headers, &res)
	if err != nil {
		return types.UserInfoFromPod{}, err
	}

	return res, nil

}

// make request for change (add ) user info
func (cfg *Config) MakeRequestForChangeUserInfo(requestBody types.ChangeUserInfoRequestBody, accessToken string) (interface{}, error) {

	//validate the requestbody
	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		return types.ChangeUserInfoRequestBody{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// generating url data for the token grant
	urlData, err := pkg.GetUrlDataFromGivenStruct(requestBody)
	if err != nil {
		return types.ChangeUserInfoRequestBody{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// geting the url for request
	requestUrl := pkg.GetUrlForChangeUserInfo(cfg.BaseUrl)

	// create the request base info's
	authorizationHeader := "Bearer " + accessToken
	headers := map[string]string{
		"Authorization": authorizationHeader,
		"Content-Type":  constanse.ContentTypeForUrlDataPod,
	}

	// make the request
	var response interface{}
	err = pkg.MakeRequestWithUrlData(requestUrl, "POST", urlData, headers, &response)
	if err != nil {
		return types.ChangeUserInfoRequestBody{}, err
	}

	return response, nil

}

// request for geting list of user info
func (cfg *Config) MakeRequestForListOfUsersInfo(requestBody types.UserListRequestBody) (types.ListOfUsersInfo, error) {

	// validate the body with certain validation set to the struct
	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		return types.ListOfUsersInfo{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// generating url data for the request to pod
	if len(requestBody.Identity) == 0 || len(requestBody.IdentityType) == 0 || len(requestBody.Identity) != len(requestBody.IdentityType) {
		return types.ListOfUsersInfo{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// get the request url
	requestUrl := pkg.GetUrlForListOfUserInfo(cfg.BaseUrl)

	// create the request base info's
	authorizationHeader := "Basic " + pkg.BasicAuth(cfg.ClientId, cfg.ClientSecret)
	headers := map[string]string{
		"Content-Type":  constanse.ContentTypeForUrlDataPod,
		"Authorization": authorizationHeader,
	}

	for i := 0; i < len(requestBody.Identity); i++ {
		if i == 0 {
			requestUrl += "?"
		}
		requestUrl += "identityType=" + requestBody.IdentityType[i] + "&" + "identity=" + requestBody.Identity[i]
	}

	// make the request
	var response types.ListOfUsersInfo
	err := pkg.MakeRequestWithNoBody(requestUrl, "GET", headers, &response)
	if err != nil {
		return types.ListOfUsersInfo{}, err
	}
	return response, nil
}

// handshake request for getting user private key
// usage in auto login generate code
func (cfg *Config) SendHandshakeRequestForUserPrivateKey(accessToken string) (types.HandShakeToGetPrivateKeyResponse, error) {

	// geting the url for request
	requestUrl := pkg.GetUrlForPrivateKeyGenerateHandshake(cfg.BaseUrl)

	// create request header
	authorizationHeader := "Bearer " + accessToken
	headers := map[string]string{
		"Content-Type":  constanse.ContentTypeForUrlDataPod,
		"Authorization": authorizationHeader,
	}

	urlData := url.Values{"keyAlgorithm": []string{"RSA"}}

	var response types.HandShakeToGetPrivateKeyResponse
	err := pkg.MakeRequestWithUrlData(requestUrl, "POST", urlData, headers, &response)
	if err != nil {
		return types.HandShakeToGetPrivateKeyResponse{}, err
	}
	return response, nil
}

// request for generating auto login code
func (cfg *Config) MakeRequestForGenerateAutoLoginCode(requestBody types.AutoLoginCodeGenerateRequestBody) (types.AutoLoginCodeGenerateResponse, error) {
	// validate the body with certain validation set to the struct
	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		return types.AutoLoginCodeGenerateResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	timeStamp := time.Now().UnixMilli()

	////// generate signature
	// first we should create the payload
	payLoad := fmt.Sprintf("access_token: %s\nkey_id: %s\ntimestamp: %d", requestBody.AccessToken, requestBody.KeyId, timeStamp)

	// sign private key of user using this payload
	signature, err := pkg.GetSignatureFromString(requestBody.PrivateKey, payLoad)
	if err != nil {
		return types.AutoLoginCodeGenerateResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidSignature,
				"errorDescription": errors.Signature,
			},
		)
	}

	requestBodyToPod := types.AutoLoginCodeGenerateRequestBodyToPod{
		KeyId:       requestBody.KeyId,
		Timestamp:   strconv.FormatInt(timeStamp, 10),
		Signature:   signature,
		AccessToken: requestBody.AccessToken,
	}

	// generating url data for the request to pod
	urlData, err := pkg.GetUrlDataFromGivenStruct(requestBodyToPod)
	if err != nil {
		return types.AutoLoginCodeGenerateResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"error":            errors.ErrInvalidInput,
				"errorDescription": "invalid input",
			},
		)
	}

	// create request header
	authorizationHeader := "Bearer " + cfg.PodAccessToken
	headers := map[string]string{
		"Content-Type":  constanse.ContentTypeForUrlDataPod,
		"Authorization": authorizationHeader,
	}

	// geting the url for request
	requestUrl := pkg.GetUrlRequestForGeneratingAutoLoginCode(cfg.BaseUrl)

	// make request
	var response types.AutoLoginCodeGenerateResponse
	err = pkg.MakeRequestWithUrlData(requestUrl, "POST", urlData, headers, &response)
	if err != nil {
		return types.AutoLoginCodeGenerateResponse{}, err
	}
	return response, nil

}

func (cfg *Config) GetSignature() (string, error) {
	var signature = cfg.Signature
	if cfg.PrivateKeyFile != "" {
		// get signature from given path
		signatureFromFile, err := pkg.GetSignatureFromFile(cfg.PrivateKeyFile, "host: "+pkg.GetHostFromURL(cfg.BaseUrl))
		if err != nil {
			return "", err
		}
		signature = signatureFromFile
	}
	return signature, nil

}
