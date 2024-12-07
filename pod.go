package pod

import (
	"fmt"

	"git.mohaami.ir/efa/go-pod-sso/constanse"
	"git.mohaami.ir/efa/go-pod-sso/errors"
	"git.mohaami.ir/efa/go-pod-sso/pkg"
	"git.mohaami.ir/efa/go-pod-sso/types"
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
				"code":    400,
				"message": errors.ErrInvalidInput,
			},
		)
	}

	// create the url data
	urlData, err := pkg.GetUrlDataForHandShakeRequest(requestBody, uuid.Nil)
	if err != nil {
		return types.HandShakeResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"code":    400,
				"message": errors.ErrInvalidInput,
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
				"code":    400,
				"message": errors.ErrInvalidInput,
			},
		)
	}

	// create the url data
	urlData, err := pkg.GetUrlDataFromGivenStruct(requestBody)
	if err != nil {
		return types.OtpResponse{}, "", errors.NewCustomError(
			map[string]interface{}{
				"code":    400,
				"message": errors.ErrInvalidInput,
			},
		)
	}
	if cfg.Signature == "" && cfg.PrivateKeyFile == "" {
		return types.OtpResponse{}, "", errors.NewCustomError(
			map[string]interface{}{
				"code":    400,
				"message": errors.ErrSignatureKeyOrFileIsMissing,
			},
		)
	}

	signature, sigErr := cfg.GetSignature()
	if sigErr != nil {
		return types.OtpResponse{}, "", errors.NewCustomError(
			map[string]interface{}{
				"code":    400,
				"message": errors.Signature,
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
				"code":    400,
				"message": errors.ErrInvalidInput,
			},
		)
	}

	// generating url data for the request to pod
	urlDataForOtp, err := pkg.GetUrlDataFromGivenStruct(requestBody)
	if err != nil {
		return types.OtpVerifyResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"code":    400,
				"message": errors.ErrInvalidInput,
			},
		)
	}

	signature, sigErr := cfg.GetSignature()
	if sigErr != nil {
		return types.OtpVerifyResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"code":    400,
				"message": errors.Signature,
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
				"code":    400,
				"message": errors.ErrInvalidInput,
			},
		)
	}

	// generating url data for the request to pod
	urlData, err := pkg.GetUrlDataFromGivenStruct(requestBody)
	if err != nil {
		return types.GetAccessTokenResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"code":    400,
				"message": errors.ErrInvalidInput,
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
				"code":    400,
				"message": errors.ErrInvalidInput,
			},
		)
	}

	// generating url data for the token grant
	urlDataForToken, err := pkg.GetUrlDataFromGivenStruct(requestBody)
	if err != nil {
		return types.GetAccessTokenResponse{}, errors.NewCustomError(
			map[string]interface{}{
				"code":    400,
				"message": errors.ErrInvalidInput,
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
				"code":    400,
				"message": errors.ErrInvalidInput,
			},
		)
	}

	// generating the url data
	urlDataForToken, err := pkg.GetUrlDataFromGivenStruct(requestBody)
	if err != nil {
		return types.AccessTokenProcess{}, errors.NewCustomError(
			map[string]interface{}{
				"code":    400,
				"message": errors.ErrInvalidInput,
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
				"code":    400,
				"message": errors.ErrInvalidInput,
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
	_, err := pkg.MakeRequestWithNoBody(requestUrl, "GET", headers, &res)
	if err != nil {
		return types.UserInfoFromPod{}, err
	}
	return res, nil

}

// request for geting list of user
// not completed , may contain bug
func (cfg *Config) MakeRequestForListOfUsersInfo(requestBody types.UserListRequestBody) (interface{}, error) {

	// validate the body with certain validation set to the struct
	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		return types.UserListRequestBody{}, errors.NewCustomError(
			map[string]interface{}{
				"code":    400,
				"message": errors.ErrInvalidInput,
			},
		)
	}

	// generating url data for the request to pod
	urlDataForValidationOfToken, err := pkg.GetUrlDataFromGivenStruct(requestBody)
	if err != nil {
		return types.UserListRequestBody{}, errors.NewCustomError(
			map[string]interface{}{
				"code":    400,
				"message": errors.ErrInvalidInput,
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

	// make the request
	var temp interface{}
	err = pkg.MakeRequestWithUrlData(requestUrl, "GET", urlDataForValidationOfToken, headers, &temp)
	if err != nil {
		return types.ValidationResponseFromPod{}, err
	}
	return temp, nil
}

func (cfg *Config) GetSignature() (string, error) {
	var signature = cfg.Signature
	if cfg.PrivateKeyFile != "" {
		// get signature from given path
		signatureFromFile, err := pkg.Getsignature(cfg.PrivateKeyFile, "host: "+pkg.GetHostFromURL(cfg.BaseUrl))
		if err != nil {
			return "", err
		}
		signature = signatureFromFile
	}
	return signature, nil

}