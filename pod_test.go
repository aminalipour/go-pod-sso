package pod

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aminalipour/go-pod-sso/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPkg struct {
	mock.Mock
}

var cfg *Config

func init() {
	cfg = &Config{
		ClientId:       "",
		ClientSecret:   "",
		PodAccessToken: "",
		BaseUrl:        "",
		Signature:      "",
	}
}

func TestSendHandshakeRequest(t *testing.T) {
	mockPkg := new(MockPkg)

	requestBody := types.HandShakeApiAdditionalDataFromClient{
		DeviceClientIp: "",
	}

	response, err := cfg.SendHandshakeRequest(requestBody, uuid.Nil)

	val := reflect.ValueOf(response)
	typ := reflect.TypeOf(response)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()
		fmt.Printf("%s: %v\n", field.Name, value)
	}

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if reflect.DeepEqual(response, types.HandShakeResponse{}) {
		t.Fatalf("expected valid response, got %v", response)
	}

	mockPkg.AssertExpectations(t)
}

func TestSendOtpRequest(t *testing.T) {
	mockPkg := new(MockPkg)

	requestBody := types.OTPRequestToPodRequestBody{
		ResponseType:     "code",
		IdentityType:     "phone_number",
		ReferrerType:     "username",
		LinkDeliveryType: "SMS",
		NationalCode:     "",
	}
	response, signature, err := cfg.SendOtpRequest(requestBody, "", "")
	fmt.Println("signiture : ", signature)

	val := reflect.ValueOf(response)
	typ := reflect.TypeOf(response)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()
		fmt.Printf("%s: %v\n", field.Name, value)
	}

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if reflect.DeepEqual(response, types.HandShakeResponse{}) {
		t.Fatalf("expected valid response, got %v", response)
	}

	mockPkg.AssertExpectations(t)
}

func TestMakeRequestForOtpVerify(t *testing.T) {
	mockPkg := new(MockPkg)

	requestBody := types.VerifyOtpRequestBody{
		OTP:         "",
		PhoneNumber: "",
	}

	response, err := cfg.MakeRequestForOtpVerify(requestBody, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	val := reflect.ValueOf(response)
	typ := reflect.TypeOf(response)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()
		fmt.Printf("%s: %v\n", field.Name, value)
	}

	if reflect.DeepEqual(response, types.OtpVerifyResponse{}) {
		t.Fatalf("expected valid response, got %v", response)
	}
	mockPkg.AssertExpectations(t)
}

func TestMakeRequestForGetAccessToken(t *testing.T) {
	mockPkg := new(MockPkg)

	requestBody := types.GetTokenRequestBody{
		Code:         "",
		GrantType:    "authorization_code",
		IdentityType: "phone_number",
		Idenify:      "",
	}

	response, err := cfg.MakeRequestForGetAccessToken(requestBody)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	val := reflect.ValueOf(response)
	typ := reflect.TypeOf(response)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()
		fmt.Printf("%s: %v\n", field.Name, value)
	}

	if reflect.DeepEqual(response, types.GetAccessTokenResponse{}) {
		t.Fatalf("expected valid response, got %v", response)
	}
	mockPkg.AssertExpectations(t)

}

func TestMakeRequestForUserInfo(t *testing.T) {
	mockPkg := new(MockPkg)

	response, err := cfg.MakeRequestForUserInfo("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	val := reflect.ValueOf(response)
	typ := reflect.TypeOf(response)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()
		fmt.Printf("%s: %v\n", field.Name, value)
	}

	if reflect.DeepEqual(response, types.UserInfoFromPod{}) {
		t.Fatalf("expected valid response, got %v", response)
	}
	mockPkg.AssertExpectations(t)

}

func TestMakeRequestForChangeUserInfo(t *testing.T) {
	mockPkg := new(MockPkg)

	reqBody := types.ChangeUserInfoRequestBody{
		NationalCode: "",
		BirthDate:    "",
	}

	resp, err := cfg.MakeRequestForChangeUserInfo(reqBody, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fmt.Println(resp)
	mockPkg.AssertExpectations(t)
}

func TestSendHandshakeRequestForUserPrivateKey(t *testing.T) {
	mockPkg := new(MockPkg)

	resp, err := cfg.SendHandshakeRequestForUserPrivateKey("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fmt.Println(resp)
	mockPkg.AssertExpectations(t)
}

func TestMakeRequestForGenerateAutoLoginCode(t *testing.T) {
	mockPkg := new(MockPkg)

	requestBody := types.AutoLoginCodeGenerateRequestBody{
		KeyId:       "",
		AccessToken: "",
		PrivateKey:  ``,
	}

	resp, err := cfg.MakeRequestForGenerateAutoLoginCode(requestBody)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fmt.Println(resp)
	mockPkg.AssertExpectations(t)
}
