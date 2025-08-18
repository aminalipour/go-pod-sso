package pkg

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aminalipour/go-pod-sso/errors"
)

func MakeRequestWithUrlData(url string, method string, urlData url.Values, headers map[string]string, responseBody interface{}) error {
	req, err := http.NewRequest(method, url, strings.NewReader(urlData.Encode()))
	if err != nil {
		return errors.NewCustomError(
			map[string]interface{}{
				"code":             503,
				"error":            errors.ErrServiceUnavailable,
				"errorDescription": "service unavailable",
			},
		)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	tr := &http.Transport{
		TLSNextProto:       make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return errors.NewCustomError(
			map[string]interface{}{
				"code":             503,
				"error":            errors.ErrServiceUnavailable,
				"errorDescription": "service unavailable",
			},
		)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.NewCustomError(
			map[string]interface{}{
				"code":             500,
				"error":            errors.ErrInternalServer,
				"errorDescription": "service error reading response body",
			},
		)
	}

	err = json.Unmarshal(body, responseBody)
	if err != nil && len(body) != 0 {
		return errors.NewCustomError(
			map[string]interface{}{
				"code":             500,
				"error":            errors.ErrInternalServer,
				"errorDescription": "service error unmarshal json",
			},
		)
	}
	if res.StatusCode != http.StatusOK {
		var bodyMap map[string]interface{}
		_ = json.Unmarshal(body, &bodyMap)

		errorMap := map[string]interface{}{
			"code": res.StatusCode,
		}

		for k, v := range bodyMap {
			errorMap[k] = v
		}

		return errors.NewCustomError(errorMap)
	}
	return nil
}

func MakeRequestWithNoBody(url string, method string, headers map[string]string, responseBody interface{}) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return errors.NewCustomError(
			map[string]interface{}{
				"code": 503,
				"message": map[string]string{
					"error":            errors.ErrServiceUnavailable,
					"errorDescription": "service unavailable",
				},
			},
		)
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return errors.NewCustomError(
			map[string]interface{}{
				"code":             503,
				"error":            errors.ErrServiceUnavailable,
				"errorDescription": "service unavailable",
			},
		)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.NewCustomError(
			map[string]interface{}{
				"code":             500,
				"error":            errors.ErrInternalServer,
				"errorDescription": "service error reading response body",
			},
		)
	}

	err = json.Unmarshal(body, responseBody)
	if err != nil {
		return errors.NewCustomError(
			map[string]interface{}{
				"code":             500,
				"error":            errors.ErrInternalServer,
				"errorDescription": "service error unmarshal json",
			},
		)
	}
	if res.StatusCode != http.StatusOK {
		var bodyMap map[string]interface{}
		_ = json.Unmarshal(body, &bodyMap)

		errorMap := map[string]interface{}{
			"code": res.StatusCode,
		}

		for k, v := range bodyMap {
			errorMap[k] = v
		}

		return errors.NewCustomError(errorMap)
	}
	return nil
}
