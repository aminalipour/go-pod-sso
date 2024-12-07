package pkg

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/aminalipour/go-pod-sso/types"
	"github.com/google/uuid"
)

func GetUrlDataForHandShakeRequest(requestBody types.HandShakeApiAdditionalDataFromClient, deviceUid uuid.UUID) (url.Values, error) {
	v := reflect.ValueOf(requestBody)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %T", v)
	}

	values := url.Values{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("url")

		if tag == "" {
			continue
		}

		fieldValue := v.Field(i)

		if (fieldValue.Kind() == reflect.String && fieldValue.String() == "") ||
			(fieldValue.Kind() == reflect.Int && fieldValue.Int() == 0) ||
			(fieldValue.Kind() == reflect.Float64 && fieldValue.Float() == 0) ||
			(fieldValue.Kind() == reflect.Bool && !fieldValue.Bool()) {
			continue
		}

		var strValue string
		switch fieldValue.Kind() {
		case reflect.String:
			strValue = fieldValue.String()
		case reflect.Int, reflect.Int64, reflect.Int32:
			strValue = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.Float32, reflect.Float64:
			strValue = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
		case reflect.Bool:
			strValue = strconv.FormatBool(fieldValue.Bool())
		default:
			continue
		}

		values.Set(tag, strValue)
	}

	if deviceUid == uuid.Nil {
		values.Set("device_uid", uuid.New().String())
	} else {
		values.Set("device_uid", deviceUid.String())
	}

	return values, nil
}

func GetUrlDataForTokenValidationRequest(requestBody types.AccessTokenProcess) url.Values {
	urlDataForValidationOfToken := url.Values{}

	urlDataForValidationOfToken.Set("token", requestBody.AccessToken)
	urlDataForValidationOfToken.Set("token_type_hint", "access_token")

	return urlDataForValidationOfToken
}

func GetUrlDataFromGivenStruct(requestBody interface{}) (url.Values, error) {
	v := reflect.ValueOf(requestBody)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %T", v)
	}

	values := url.Values{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldValue := v.Field(i)
		if isZeroValue(fieldValue) {
			continue
		}
		fieldName := field.Name
		if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
			fieldName = tag
		} else if tag := field.Tag.Get("url"); tag != "" && tag != "-" {
			fieldName = tag
		}
		if commaIdx := len(fieldName); commaIdx > 0 {
			fieldName = fieldName[:commaIdx]
		}

		var strValue string
		switch fieldValue.Kind() {
		case reflect.String:
			strValue = fieldValue.String()
		case reflect.Int, reflect.Int64, reflect.Int32:
			strValue = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.Float32, reflect.Float64:
			strValue = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
		case reflect.Bool:
			strValue = strconv.FormatBool(fieldValue.Bool())
		case reflect.Slice:
			for j := 0; j < fieldValue.Len(); j++ {
				elem := fieldValue.Index(j)
				values.Add(fieldName, fmt.Sprintf("%v", elem.Interface()))
			}
			continue
		default:
			continue
		}
		values.Set(fieldName, strValue)
	}

	return values, nil
}
