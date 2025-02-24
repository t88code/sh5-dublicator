package sh5api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type ClientInterface interface {
	Sh5ExecOnlyProcName(ctx context.Context, procName ProcName) (rep *Sh5ExecRep, err error)
	Sh5ExecWithInput(
		ctx context.Context,
		procName ProcName,
		head Head,
		original []FieldPath,
		values [][]interface{},
		status []Sh5ExecStatus,
	) (rep *Sh5ExecRep, err error)
}

type Client struct {
	ClientInterface

	httpClient *http.Client
	config     *Config
}

type Config struct {
	HttpClient *http.Client
	BaseURL    string
	Username   string
	Password   string

	DebugLog bool

	OnQuery func(req *http.Request, body []byte, reqName string)
	OnReply func(res *http.Response, body []byte, reqName string)

	Logger *slog.Logger
}

func New(config *Config) (ClientInterface, error) {
	httpClient := config.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	if config.Logger == nil {
		config.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	return &Client{
		httpClient: httpClient,
		config:     config,
	}, nil
}

// request Отправка запроса в API SH5
func (c *Client) request(ctx context.Context, path string, req interface{}, reqName string) ([]byte, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := strings.Trim(c.config.BaseURL, "/ ") + path

	resBody, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	if c.config.OnQuery != nil {
		c.config.OnQuery(resBody, reqBody, reqName)
	}

	resp, err := c.httpClient.Do(resBody)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, &TransportError{error: err}
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if c.config.OnReply != nil {
		c.config.OnReply(resp, respBody, reqName)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &TransportError{
			error:    fmt.Errorf("http error code=%d", resp.StatusCode),
			HttpCode: resp.StatusCode,
		}
	}

	var errorReply ErrorReply
	if err = json.Unmarshal(respBody, &errorReply); err != nil {
		return nil, err
	}

	if errorReply.ErrMessage != "OK" {
		c.config.Logger.Error("http error",
			"errorCode", errorReply.ErrorCode,
			"errMessage", errorReply.ErrMessage,
			"Version", errorReply.Version,
			"UserName", errorReply.UserName,
			"actionName", errorReply.ActionName,
			"actionType", errorReply.ActionType,
			"moduleId", errorReply.ErrorInfo.ModuleId,
			"moduleName", errorReply.ErrorInfo.ModuleName,
			"errorId", errorReply.ErrorInfo.ErrorId,
			"statusCode", resp.StatusCode,
		)

		return nil, &ProcessingError{
			error:      errors.New(errorReply.ErrMessage),
			ErrorReply: &errorReply,
		}
	}

	return respBody, nil
}

func (c *Client) Sh5ExecOnlyProcName(ctx context.Context, procName ProcName) (rep *Sh5ExecRep, err error) {
	req := new(Sh5ExecReq)
	req.Password = c.config.Password
	req.UserName = c.config.Username
	req.ProcName = procName

	respBody, err := c.request(ctx, "/api/sh5exec", &req, string(procName))
	if err != nil {
		return nil, err
	}

	var resp Sh5ExecRep
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) Sh5ExecWithInput(
	ctx context.Context,
	procName ProcName,
	head Head,
	original []FieldPath,
	values [][]interface{},
	status []Sh5ExecStatus,
) (rep *Sh5ExecRep, err error) {

	req := new(Sh5ExecReq)
	req.Password = c.config.Password
	req.UserName = c.config.Username
	req.ProcName = procName
	req.Input = append(req.Input, Sh5ExecInput{
		Head:     head,
		Original: original,
		Fields:   nil,
		Values:   values,
		Status:   status,
	})

	respBody, err := c.request(ctx, "/api/sh5exec", &req, string(procName))
	if err != nil {
		return nil, err
	}

	var resp Sh5ExecRep
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
