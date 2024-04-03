package apicar

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
)

type GetCarResponseBody struct {
	Description string  `json:"description"`
	Content     Content `json:"content"`
}

type Content struct {
	GetCarResponse
}

type GetCarResponse struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int16  `json:"year,omitempty"`
	Owner  *Owner `json:"owner"`
}

type Owner struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type ErrorResponse struct {
	Description any `json:"description"`
}

func (httpClient *HttpClient) GetCarInfo(regNum string) (*GetCarResponse, error) {
	return httpClient.getCarInfo(regNum)
}

func (httpClient *HttpClient) getCarInfo(regNum string) (*GetCarResponse, error) {
	httpClient.logger.Debug("getting car info from api", "regNum", regNum)

	responseBodyReader, err := httpClient.callCarApi(regNum)
	if err != nil {
		httpClient.logger.Debug("failed to get car info from api", "error", err.Error())

		return nil, err
	}

	car, err := httpClient.parseCarInfoApiResponse(responseBodyReader)
	if err != nil {
		httpClient.logger.Debug("failed to parse car info api response", "error", err.Error())

		return nil, err
	}

	httpClient.logger.Debug("received car info response from api")

	return car, nil
}

func (httpClient *HttpClient) callCarApi(regNum string) (io.Reader, error) {
	httpClient.logger.Debug("calling car api", "regNum", regNum)

	resp, err := httpClient.makeRequest(regNum)
	if err != nil {
		httpClient.logger.Debug("failed to make a request to car api", "error", err.Error())

		return nil, err
	}
	responseBodyReader := bytes.NewReader(resp.Body())

	httpClient.logger.Debug("successfully called car api")

	return responseBodyReader, nil
}

func (httpClient *HttpClient) makeRequest(regNum string) (*resty.Response, error) {
	httpClient.logger.Debug("making a request", "regNum", regNum)

	resp, err := httpClient.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("regNum", regNum).
		Get(httpClient.host)

	if err != nil {
		httpClient.logger.Debug("failed to make a request to api", "error", err.Error())

		return nil, fmt.Errorf("failed to make request to api: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		errorResponse, err := httpClient.getErrorResponse(bytes.NewReader(resp.Body()))
		if err != nil {
			httpClient.logger.Debug("failed to get error response from car api", "error", err.Error())

			return nil, err
		}

		httpClient.logger.Debug("failed to get response from car api", "error-description", errorResponse.Description)

		return nil, errors.New(fmt.Sprintf("%s", errorResponse.Description))
	}

	httpClient.logger.Debug("successfully made a request to car api")

	return resp, nil
}

func (httpClient *HttpClient) getErrorResponse(responseBodyReader io.Reader) (*ErrorResponse, error) {
	var errorResponse ErrorResponse
	err := json.NewDecoder(responseBodyReader).Decode(&errorResponse)
	if err != nil {
		httpClient.logger.Debug("failed to get an error response from car api", "error", err.Error())

		return nil, err
	}

	httpClient.logger.Debug("successfully got an error response from car api")

	return &errorResponse, nil
}

func (httpClient *HttpClient) parseCarInfoApiResponse(r io.Reader) (*GetCarResponse, error) {
	var getCarRespBody GetCarResponseBody
	err := json.NewDecoder(r).Decode(&getCarRespBody)
	if err != nil {
		httpClient.logger.Debug("failed to decode a response from car api", "error", err.Error())

		return nil, err
	}

	httpClient.logger.Debug("successfully parsed a response from car api")

	return &getCarRespBody.Content.GetCarResponse, nil
}
