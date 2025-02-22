package accrual

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sonikq/gophermart/internal/models"
	"log"
	"net/http"
	"time"
)

// Client - структура Accrual клиента
type Client struct {
	Client     *resty.Client
	WorkerPool chan WorkerPool
}

// WorkerPool - структура запроса
type WorkerPool struct {
	URL      string
	Request  *resty.Request
	Err      chan error
	Response chan *resty.Response
}

// NewClient - конструктор для Accrual сервиса
func NewClient(serverAddress string, pool chan WorkerPool) *Client {
	return &Client{
		Client:     resty.New().SetBaseURL(serverAddress),
		WorkerPool: pool,
	}
}

// MakeRequest - функция для создания запроса по определенному URL
func (w *WorkerPool) MakeRequest() (*resty.Response, error) {
	return w.Request.Get(w.URL)
}

// Run - запуск сервиса Accrual
func (c *Client) Run() {
	for w := range c.WorkerPool {
		resp, err := w.MakeRequest()
		w.Err <- err
		w.Response <- resp
	}
}

// GetAccrualInfo - контролирует количество запросов к сервису Accrual и формирует сам запрос, потом пишет все это в канал
func (c *Client) GetAccrualInfo(orderNum string) (models.AccrualInfo, error) {
	const (
		source           = "accrual.GetAccrualInfo"
		maxAttemptPerSec = 5
	)
	// Создаем тикер, чтобы можно было контролировать кол-во запросов в определенный промежуток времени
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for i := 0; i < maxAttemptPerSec; i++ {
		respChan := make(chan *resty.Response)
		errChan := make(chan error)

		c.WorkerPool <- WorkerPool{
			URL:      "/api/orders/" + orderNum,
			Request:  c.Client.R(),
			Response: respChan,
			Err:      errChan,
		}
		err := <-errChan
		close(errChan)
		if err != nil {
			return models.AccrualInfo{}, fmt.Errorf("%s: %w", source, err)
		}

		resp := <-respChan
		close(respChan)

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
			log.Println(fmt.Errorf("%s: %s", source, "the number of requests to the service has been exceeded"))
			i--
			<-ticker.C
			continue
		default:
			return models.AccrualInfo{}, fmt.Errorf("unexpected status code: %s: code: %d", source, resp.StatusCode())
		}
	}
	return models.AccrualInfo{}, fmt.Errorf("%s: %s", source, "could not reach the service")
}
