package types

import "github.com/google/uuid"

type DeviceType int
type Algorithm string

const (
	MobilePhone DeviceType = iota + 1
	Desktop
	Tablet
	Console
	TVDevice
	MobileDevice
	Unknown
)

func (e DeviceType) String() string {
	switch e {
	case MobilePhone:
		return "Mobile Phone"
	case Desktop:
		return "Desktop"
	case Tablet:
		return "Tablet"
	case Console:
		return "Console"
	case TVDevice:
		return "TV Device"
	case MobileDevice:
		return "Mobile Device"
	default:
		return "Unknown"
	}
}

type HandShakeRequestBody struct {
	DeviceUid        string     `url:"device_uid" validate:"required"`
	DeviceLat        string     `url:"device_latوomitempty"`
	DeviceLon        string     `url:"device_lonوomitempty"`
	DeviceOs         string     `url:"device_os,omitempty"`
	DeviceOsVersion  string     `url:"device_os_version,omitempty"`
	DeviceType       DeviceType `url:"device_type"`
	DeviceName       string     `url:"device_name,omitempty"`
	DeviceAppName    string     `url:"device_app_name,omitempty"`
	DeviceAppVersion string     `url:"device_app_version,omitempty"`
	Algorithm        Algorithm  `url:"algorithm" validate:"omitempty,oneof=rsa-sha256 rsa-sha1"`
}

type HandShakeApiAdditionalDataFromClient struct {
	PhoneNumber      string     `url:"phoneNumber,omitempty"`
	DeviceLat        string     `url:"device_lat,omitempty"`
	DeviceLon        string     `url:"device_lon,omitempty"`
	DeviceOs         string     `url:"device_os,omitempty"`
	DeviceOsVersion  string     `url:"device_os_version,omitempty"`
	DeviceType       DeviceType `url:"device_type"`
	DeviceName       string     `url:"device_name,omitempty"`
	DeviceAppName    string     `url:"device_app_name,omitempty"`
	DeviceAppVersion string     `url:"device_app_version,omitempty"`
	DeviceClientIp   string     `url:"device_client_ip" validate:"required"`
	Algorithm        Algorithm  `url:"algorithm" validate:"omitempty,oneof=rsa-sha256 rsa-sha1"`
}

type VerifyOtpRequestBody struct {
	OTP         string `json:"otp" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type GetTokenRequestBody struct {
	GrantType    string `json:"grant_type" validate:"required"`
	Code         string `json:"code"`
	RedirectUrl  string `json:"redirect_uri"`
	RefreshToken string `json:"refresh_token"`
	UserName     string `json:"username"`
	Idenify      string `json:"identity"`
	IdentityType string `json:"identityType"`
	Password     string `json:"password"`
	CodeVerifier string `json::"code_verifier"`
}

type DenyPermissionRequestBody struct {
	Token         string `json:"token"`
	TokenTypeHint string `json:"token_type_hint"`
	CodeVerifier  string `json:"code_verifier"`
}

type OTPRequestToPodRequestBody struct {
	ResponseType                string `url:"response_type"`
	IdentityType                string `url:"identityType"`
	NationalCode                string `url:"nationalcode"`
	NationalCodeSerial          string `url:"nationalcodeSerial"`
	BirthDate                   string `url:"birthdate"`
	LoginAsUserId               string `url:"loginAsUserId"`
	LoginAsRelativeNationalCode string `url:"loginAsRelativeNationalcode"`
	LoginAsUsernameChild        string `url:"loginAsUsernameChild"`
	OtpType                     string `url:"otpType"`
	WebOtpDomain                string `url:"webOtpDomain"`
	CodeLength                  string `url:"codeLength"`
	State                       string `url:"state"`
	ClientId                    string `url:"client_id"`
	RedirectUri                 string `url:"redirect_uri"`
	CallbackUri                 string `url:"callback_uri"`
	Scope                       string `url:"scope"`
	CodeChallange               string `url:"code_challenge"`
	CodeChallengeMethod         string `url:"code_challenge_method"`
	Referrer                    string `url:"referrer"`
	ReferrerType                string `url:"referrerType"`
	LinkDeliveryType            string `url:"linkDeliveryType"`
}

type ClientInfo struct {
	AccessTokenExpireTime   int      `json:"accessTokenExpiryTime"`
	Active                  bool     `json:"active"`
	AllowedGrantTypes       []string `json:"allowedGrantTypes"`
	AllowedRedirectUris     []string `json:"allowedRedirectUris"`
	AllowedScopes           []string `json:"allowedRedirectUris"`
	AutoLoginAs             bool     `json:"autoLoginAs"`
	CaptchaEnable           bool     `json:"captchaEnabled"`
	ClientId                string   `json:"client_id"`
	CssEnabled              bool     `json:"cssEnabled"`
	Id                      int      `json:"id"`
	LimitedLoginAs          bool     `json:"limitedLoginAs"`
	LoginAsDepositEnabled   bool     `json:"loginAsDepositEnabled"`
	LoginUrl                string   `json:"loginUrl"`
	Name                    string   `json:"name"`
	OtpExpieryTime          int      `json:"otpCodeExpiryTime"`
	PkceEnabled             bool     `json:"pkceEnabled"`
	RefreshTokenExpieryTime int      `json:"refreshTokenExpiryTime"`
	Roles                   []string `json:"roles"`
	SignUpEnabled           bool     `json:"signupEnabled"`
	TwoFaEnabled            bool     `json:"twoFAEnabled"`
	Url                     string   `json:"url"`
	UserId                  int      `json:"userId"`
}

type DeviceInfo struct {
	Current        bool                   `json:"current"`
	Id             int                    `json:"id"`
	Ip             string                 `json:"ip"`
	Language       string                 `json:"language"`
	LastAccessTime int                    `json:"lastAccessTime"`
	Location       map[string]interface{} `json:"location"`
	Uid            string                 `json:"uid"`
}

type HandShakeResponse struct {
	Algorithm  string     `json:"algorithm"`
	ClientInfo ClientInfo `json:"client"`
	DeviceInfo DeviceInfo `json:"device"`
	ExpiresIn  int        `json:"expires_in"`
	KeyFromat  string     `json:"keyFormat"`
	KeyId      string     `json:"keyId"`
	PublicKey  string     `json:"publicKey"`
}

type OtpResponse struct {
	CodeLength int    `json:"codeLength"`
	ExpiresIn  int    `json:"expires_in"`
	Identity   string `json:"identity"`
	SentBefore bool   `json:"sent_before"`
	Type       string `json:"type"`
}

type OtpVerifyResponse struct {
	Code      string `json:"code"`
	DeviceUid string `json:"device_uid"`
}

type GetAccessTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	DeviceUid    uuid.UUID `json:"device_uid"`
	ExpiresIn    int       `json:"expires_in"`
	IdToken      string    `json:"id_token"`
	RefreshToken string    `json:"refresh_token"`
	Scope        string    `json:"scope"`
	TokenType    string    `json:"token_type"`
}

type AccessTokenResponseClient struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type AccessTokenProcess struct {
	Token         string `url:"token" validate:"required"`
	TokenTypeHint string `url:"token_type_hint" validate:"required"`
}

type CallerClient struct {
	EmailVerified              bool   `json:"email_verified"`
	FamilyName                 string `json:"family_name"`
	GivenName                  string `json:"given_name"`
	HasPassword                bool   `json:"hasPassword"`
	Id                         int    `json:"id"`
	NationalcodeSerialVerified bool   `json:"nationalcode_serial_verified"`
	NationalcodeVerified       bool   `json:"nationalcode_verified"`
	PhoneNumberVerified        bool   `json:"phone_number_verified"`
	PhysicalVerified           bool   `json:"physical_verified"`
	PreferredUsername          string `json:"preferred_username"`
}

type ValidationResponseFromPod struct {
	Active                 bool         `json:"active"`
	ClientId               string       `json:"client_id"`
	DeviceUid              string       `json:"device_uid"`
	Expiration             int          `json:"exp"`
	IsApiToken             bool         `json:"is_api_token"`
	IssuerClient           int          `json:"issuer_client"`
	LoginAsDepositEnabaled bool         `json:"loginAsDepositEnabled"`
	LoginType              string       `json:"login_type"`
	Scope                  string       `json:"scope"`
	ShamsiExpDate          string       `json:"shamsi_exp_date"`
	Sub                    string       `json:"sub"`
	ActionType             string       `json:"actionType"`
	IssueType              string       `json:"issueType"`
	ClientName             string       `json:"clientName"`
	CallerClient           CallerClient `json:"callerClient"`
}

type UserInfoFromPod struct {
	EmailVerified              bool   `json:"email_verified"`
	FamilyName                 string `json:"family_name"`
	GivenName                  string `json:"given_name"`
	HasNotinouProfile          bool   `json:"hasNotinouProfile"`
	HasPassword                bool   `json:"hasPassword"`
	Id                         int    `json:"id"`
	LegalNationalCodeVerified  bool   `json:"legalNationalCode_verified"`
	NationalcodeSerialVerified bool   `json:"nationalcode_serial_verified"`
	NationalcodeVerified       bool   `json:"nationalcode_verified"`
	PhoneNumberVerified        bool   `json:"phone_number_verified"`
	PhoneNumber                string `json:"phone_number"`
	PhysicalVerified           bool   `json:"physical_verified"`
	PreferredUsername          string `json:"preferred_username"`
	RegisterTime               int    `json:"registerTime"`
	RegisterTimeShamsi         string `json:"registerTimeShamsi"`
	Sope                       string `json:"scope"`
	Sub                        string `json:"sub"`
	UpdatedAt                  int    `json:"updated_at"`
	UpdatedAtShamsi            string `json:"updated_at_shamsi"`
	BpiCustomerNumber          string `json:"bpiCustomerNumber"`
}

type ListOfUsersInfo struct {
	Users []ListOfUsersInfoItem `json:"users"`
}

type ListOfUsersInfoItem struct {
	EmailVerified               bool           `json:"email_verified"`
	FamilyName                  string         `json:"family_name"`
	ForeignCodeVerified         bool           `json:"foreigncode_verified"`
	GivenName                   string         `json:"given_name"`
	HasNotinouProfile           bool           `json:"hasNotinouProfile"`
	HasPassword                 bool           `json:"hasPassword"`
	ID                          int64          `json:"id"`
	LegalInquireStatus          []LegalInquire `json:"legalInquireStatus"`
	LegalNationalCodeVerified   bool           `json:"legalNationalCode_verified"`
	NationalCode                string         `json:"nationalcode"`
	NationalCodeSerial          string         `json:"nationalcode_serial"`
	NationalCodeSerialVerified  bool           `json:"nationalcode_serial_verified"`
	NationalCodeSerialVerifiers []int64        `json:"nationalcode_serial_verifiers"`
	NationalCodeVerified        bool           `json:"nationalcode_verified"`
	NationalCodeVerifiers       []int64        `json:"nationalcode_verifiers"`
	PhoneNumber                 string         `json:"phone_number"`
	PhoneNumberVerified         bool           `json:"phone_number_verified"`
	PhoneNumberVerifiers        []int64        `json:"phone_number_verifiers"`
	PhysicalVerified            bool           `json:"physical_verified"`
	PreferredUsername           string         `json:"preferred_username"`
	RegisterTime                int64          `json:"registerTime"`
	RegisterTimeShamsi          string         `json:"registerTimeShamsi"`
	Scope                       string         `json:"scope"`
	Sub                         string         `json:"sub"`
	UpdatedAt                   int64          `json:"updated_at"`
	UpdatedAtShamsi             string         `json:"updated_at_shamsi"`
	BpiCustomerNumber           string         `json:"bpiCustomerNumber"`
}

type LegalInquire struct {
	AnswerFromCoreTime   string       `json:"answerFromCoreTime"`
	CallerClient         CallerClient `json:"callerClient"`
	InquiryTime          int64        `json:"inquiryTime"`
	InquiryTimeShamsi    string       `json:"inquiryTimeShamsi"`
	InquiryType          string       `json:"inquiryType"`
	InsertTime           int64        `json:"insertTime"`
	InsertTimeShamsi     string       `json:"insertTimeShamsi"`
	LastUpdateTime       int64        `json:"lastUpdateTime"`
	LastUpdateTimeShamsi string       `json:"lastUpdateTimeShamsi"`
	LegalAuthority       string       `json:"legalAuthority"`
	NationalCode         string       `json:"nationalcode"`
	NotifyToCoreTime     string       `json:"notifyToCoreTime"`
	PhoneNumber          string       `json:"phoneNumber"`
	Status               string       `json:"status"`
	TryCount             int          `json:"tryCount"`
}

type UserInfoConverted struct {
	EmailVerified              bool   `json:"emailVerified"`
	FamilyName                 string `json:"familyName"`
	GivenName                  string `json:"givenName"`
	HasNotinouProfile          bool   `json:"hasNotinouProfile"`
	HasPassword                bool   `json:"hasPassword"`
	Id                         int    `json:"id"`
	LegalNationalCodeVerified  bool   `json:"legalNationalCodeVerified"`
	NationalcodeSerialVerified bool   `json:"nationalcodeSerialVerified"`
	NationalcodeVerified       bool   `json:"nationalcodeVerified"`
	PhoneNumberVerified        bool   `json:"phoneNumberVerified"`
	PhysicalVerified           bool   `json:"physicalVerified"`
	PreferredUsername          string `json:"preferredUsername"`
	RegisterTime               int    `json:"registerTime"`
	RegisterTimeShamsi         string `json:"registerTimeShamsi"`
	Sope                       string `json:"scope"`
	Sub                        string `json:"sub"`
	UpdatedAt                  int    `json:"updatedAt"`
	UpdatedAtShamsi            string `json:"updatedAtShamsi"`
}

type UserListRequestBody struct {
	Identity     []string `json:"identity"`
	IdentityType []string `json:"identityType"`
}

type ChangeUserInfoRequestBody struct {
	NationalCode string `url:"nationalcode"`
	BirthDate    string `url:"birthdate"`
}

type HandShakeToGetPrivateKeyResponse struct {
	KeyId      string `json:"keyId"`
	PrivateKey string `json:"privateKey"`
}

type AutoLoginCodeGenerateRequestBody struct {
	KeyId       string `json:"keyId"`
	AccessToken string `json:"accessToken"`
	PrivateKey  string `json:"privateKey"`
}

type AutoLoginCodeGenerateRequestBodyToPod struct {
	KeyId       string `json:"key_id"`
	Timestamp   string `json:"timestamp"`
	Signature   string `json:"signature"`
	AccessToken string `json:"access_token"`
}

type AutoLoginCodeGenerateResponse struct {
	AccessToken   string `json:"access_token"`
	AutoLoginCode string `json:"auto_login_code"`
	KeyId         string `json:"key_id"`
	Signature     string `json:"signature"`
	TimeStamp     string `json:"timestamp"`
}
