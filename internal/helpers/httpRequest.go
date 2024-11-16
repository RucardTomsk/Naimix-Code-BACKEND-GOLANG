package helpers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/common"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"io"
	"net/http"
	"net/url"
)

type HttpClient struct {
	url             string
	queryParameters []QueryParameter
	Client          *http.Client
}

func NewHttpClient(config *common.HttpClientConfig) (*HttpClient, error) {
	if config.URL == "" {
		return nil, errors.New("request path empty")
	}

	return &HttpClient{
		url:    config.URL,
		Client: http.DefaultClient}, nil
}

type QueryParameter struct {
	Key       string
	Parameter string
}

func (h *HttpClient) CreateNewClientWithQueryParams(queryParameters []QueryParameter) HttpClient {
	t := *h
	t.queryParameters = queryParameters
	return t
}

func (h *HttpClient) HttpRequest(method string, request string, body io.Reader, ctx context.Context, parameter ...QueryParameter) ([]byte, *base.ServiceError) {
	fullUrl, err := url.JoinPath(h.url, request)
	if err != nil {
		return nil, base.NewPathError(err)
	}
	req, err := http.NewRequest(method, fullUrl, body)
	if err != nil {
		return nil, base.NewHttpServerConnectError(err)
	}
	return h.baseHttpRequest(method, req, ctx, parameter)
}
func (h *HttpClient) HttpRequestWithBasicAuth(method, request, login, password string, body io.Reader, ctx context.Context, parameter ...QueryParameter) ([]byte, *base.ServiceError) {
	fullUrl, err := url.JoinPath(h.url, request)
	if err != nil {
		return nil, base.NewPathError(err)
	}
	req, err := http.NewRequest(method, fullUrl, body)
	if err != nil {
		return nil, base.NewHttpServerConnectError(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", basicAuth(login, password)))

	return h.baseHttpRequest(method, req, ctx, parameter)
}

func (h *HttpClient) baseHttpRequest(method string, req *http.Request, ctx context.Context, parameter []QueryParameter) ([]byte, *base.ServiceError) {
	if len(parameter) != 0 {
		q := req.URL.Query()
		for _, query := range parameter {
			q.Add(query.Key, query.Parameter)
		}
		req.URL.RawQuery = q.Encode()
	}

	if h.queryParameters != nil {
		q := req.URL.Query()
		for _, query := range h.queryParameters {
			q.Add(query.Key, query.Parameter)
		}
		req.URL.RawQuery = q.Encode()
	}

	if method == http.MethodPost {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := h.Client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		if res != nil {
			return nil, base.NewHttpServerRequestError(err, res.StatusCode)
		}
		return nil, base.NewHttpServerRequestError(err, 500)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, base.NewReadByteError(err)
	}

	h.queryParameters = nil
	return resBody, nil
}

func basicAuth(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
}
