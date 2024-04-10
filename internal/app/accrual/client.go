package accrual

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sonikq/gophermart/internal/models"
	"net/http"
)

type Client struct {
	Client *resty.Client
}

func NewClient(serverAddress string) *Client {
	return &Client{
		Client: resty.New().SetBaseURL(serverAddress),
	}
}

func (c *Client) GetAccrualInfo(orderNum string) (models.AccrualInfo, error) {
	const source = "accrual.GetAccrualInfo"
	resp, err := c.Client.R().Get("/api/orders/" + orderNum)
	if err != nil {
		return models.AccrualInfo{}, fmt.Errorf("%s: %w", source, err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		var info models.AccrualInfo
		err = json.Unmarshal(resp.Body(), &info)
		if err != nil {
			return models.AccrualInfo{}, fmt.Errorf("%s: %w", source, err)
		}
		return info, nil
	case http.StatusNoContent:
		return models.AccrualInfo{}, fmt.Errorf("%s: %s", source, "the order is not registered in the payment system")
	case http.StatusTooManyRequests:
		return models.AccrualInfo{}, fmt.Errorf("%s: %s", source, "the number of requests to the service has been exceeded")
	default:
		return models.AccrualInfo{}, fmt.Errorf("unexpected status code: %s: code: %d", source, resp.StatusCode())
	}
}
