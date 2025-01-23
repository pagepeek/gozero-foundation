package oauth

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/pagepeek/gozero-foundation/pkg/utils"
)

var ErrDataFormatterNotFound = errors.New("unable to process request body, please set \"DataFormatter\"")

type BaseClient struct {
	Endpoint      string
	Schema        string
	HttpClient    *http.Client
	Options       []RequestOption
	DataFormatter func(req *http.Request, data any) error
}

func JsonClient(endpoint string, cli *http.Client, opts ...RequestOption) *BaseClient {
	if cli == nil {
		cli = http.DefaultClient
	}

	return &BaseClient{
		Endpoint:   endpoint,
		Schema:     "https",
		HttpClient: cli,
		Options:    opts,
		DataFormatter: func(req *http.Request, data any) error {
			body, err := utils.JsonEncode(data)
			if err != nil {
				return err
			}

			req.Body = io.NopCloser(body)
			req.ContentLength = int64(body.Len())
			req.Header.Add("Content-Type", "application/json")
			return nil
		},
	}
}

func QueryClient(endpoint string, cli *http.Client, opts ...RequestOption) *BaseClient {
	if cli == nil {
		cli = http.DefaultClient
	}

	return &BaseClient{
		Endpoint:   endpoint,
		Schema:     "https",
		HttpClient: cli,
		Options:    opts,
		DataFormatter: func(req *http.Request, data any) error {
			query := req.URL.Query()

			switch params := data.(type) {
			case map[string]string:
				for k, v := range params {
					query.Add(k, v)
				}
			case url.Values:
				for k := range params {
					query.Add(k, params.Get(k))
				}
			}

			req.URL.RawQuery = query.Encode()

			return nil
		},
	}
}

func (cli *BaseClient) Get(path string, data map[string]string, options ...RequestOption) (*http.Response, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	for k := range data {
		query.Add(k, data[k])
	}

	u.RawQuery = query.Encode()

	return cli.Send(http.MethodGet, u.String(), nil, options...)
}

func (cli *BaseClient) Post(path string, data any, options ...RequestOption) (*http.Response, error) {
	return cli.Send(http.MethodPost, path, data, options...)
}

func (cli *BaseClient) Put(path string, data any, options ...RequestOption) (*http.Response, error) {
	return cli.Send(http.MethodPut, path, data, options...)
}

func (cli *BaseClient) Delete(path string, data any, options ...RequestOption) (*http.Response, error) {
	return cli.Send(http.MethodDelete, path, data, options...)
}

func (cli *BaseClient) Send(method, path string, data any, options ...RequestOption) (*http.Response, error) {
	target, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	if target.Scheme == "" {
		target.Scheme = cli.Schema
	}
	if target.Host == "" {
		target.Host = cli.Endpoint
	}

	var body io.Reader
	if _, ok := data.(io.Reader); ok {
		body = data.(io.Reader)
	}

	req, err := http.NewRequest(method, target.String(), body)
	if err != nil {
		return nil, err
	}

	if req.Body == nil && data != nil {
		if cli.DataFormatter == nil {
			return nil, ErrDataFormatterNotFound
		}

		if err := cli.DataFormatter(req, data); err != nil {
			return nil, err
		}
	}

	for i := range options {
		options[i](req)
	}

	for i := range cli.Options {
		cli.Options[i](req)
	}

	return cli.HttpClient.Do(req)
}
