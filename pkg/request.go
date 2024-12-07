package pkg

import (
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
				"code":    503,
				"message": errors.ErrServiceUnavailable,
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
				"code":    503,
				"message": errors.ErrServiceUnavailable,
			},
		)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.NewCustomError(
			map[string]interface{}{
				"code":    500,
				"message": errors.ErrInternalServer,
			},
		)
	}

	err = json.Unmarshal(body, responseBody)
	if err != nil {
		return errors.NewCustomError(
			map[string]interface{}{
				"code":    500,
				"message": errors.ErrInternalServer,
			},
		)
	}
	if res.StatusCode != http.StatusOK {
		return errors.NewCustomError(
			map[string]interface{}{
				"code":    res.StatusCode,
				"message": string(body),
			},
		)
	}
	return nil
}

func MakeRequestWithNoBody(url string, method string, headers map[string]string, responseBody interface{}) (interface{}, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.NewCustomError(
			map[string]interface{}{
				"code":    503,
				"message": errors.ErrServiceUnavailable,
			},
		)
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.NewCustomError(
			map[string]interface{}{
				"code":    503,
				"message": errors.ErrServiceUnavailable,
			},
		)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.NewCustomError(
			map[string]interface{}{
				"code":    503,
				"message": errors.ErrServiceUnavailable,
			},
		)
	}

	err = json.Unmarshal(body, responseBody)
	if err != nil {
		return nil, errors.NewCustomError(
			map[string]interface{}{
				"code":    503,
				"message": errors.ErrServiceUnavailable,
			},
		)
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.NewCustomError(
			map[string]interface{}{
				"code":    res.StatusCode,
				"message": string(body),
			},
		)
	}
	return responseBody, nil
}
