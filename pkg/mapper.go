package pkg

import "github.com/aminalipour/go-pod-sso/types"

func UserInfoFromPodToUserInfoConverted(userInfo types.UserInfoFromPod) types.UserInfoConverted {
	return types.UserInfoConverted{
		EmailVerified:              userInfo.EmailVerified,
		FamilyName:                 userInfo.FamilyName,
		GivenName:                  userInfo.GivenName,
		HasNotinouProfile:          userInfo.HasNotinouProfile,
		HasPassword:                userInfo.HasPassword,
		Id:                         userInfo.Id,
		LegalNationalCodeVerified:  userInfo.LegalNationalCodeVerified,
		NationalcodeSerialVerified: userInfo.NationalcodeSerialVerified,
		NationalcodeVerified:       userInfo.NationalcodeVerified,
		PhoneNumberVerified:        userInfo.PhoneNumberVerified,
		PhysicalVerified:           userInfo.PhysicalVerified,
		PreferredUsername:          userInfo.PreferredUsername,
		RegisterTime:               userInfo.RegisterTime,
		RegisterTimeShamsi:         userInfo.RegisterTimeShamsi,
		Sope:                       userInfo.Sope,
		Sub:                        userInfo.Sub,
		UpdatedAt:                  userInfo.UpdatedAt,
		UpdatedAtShamsi:            userInfo.UpdatedAtShamsi,
	}
}
