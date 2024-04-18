package gophermart

import (
	"context"
	"fmt"
	"github.com/sonikq/gophermart/internal/models"
	"github.com/sonikq/gophermart/pkg/validator"
	"log"
	"regexp"
	"sync"
)

func (s *Service) UploadOrder(ctx context.Context, orderNum string, username string) error {
	if err := validateOrderNum(orderNum); err != nil {
		return fmt.Errorf("%w: %w", models.ErrInvalidOrderNum, err)
	}

	order, err := s.storage.GetOrder(ctx, orderNum)
	if err != nil {
		return err
	}

	if order == nil {
		return s.storage.UploadOrder(ctx, orderNum, username)
	}

	if order.Username != username {
		return models.ErrOrderAlreadyUploadedByAnotherUser
	}

	return models.ErrOrderAlreadyUploadedByThisUser
}

func validateOrderNum(orderNumStr string) error {
	var re = regexp.MustCompile(`^[0-9]+$`)

	if !re.MatchString(orderNumStr) {
		return fmt.Errorf("order number should consist only of digits")
	}

	if !validator.CheckLuhn(orderNumStr) {
		return fmt.Errorf("order number did not pass the validity check using the luhn algorithm")
	}

	return nil
}

func (s *Service) UpdateUserOrders(ctx context.Context, username string) error {
	orders, err := s.storage.ListUserOrders(ctx, username)
	if err != nil {
		return err
	}

	if len(orders) == 0 {
		return models.ErrEmptyOrderList
	}

	actualOrders := make([]models.Order, 0, len(orders))
	for _, order := range orders {
		if order.Status != models.InvalidOrder &&
			order.Status != models.ProcessedOrder {
			actualOrders = append(actualOrders, order)
		}
	}

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	// создаем слайс уже с известной емкостью равной len(actualOrders)
	accrualInfos := make([]models.AccrualInfo, 0, len(actualOrders))
	for _, order := range actualOrders {
		// добавляем горутину для обработки каждого заказа
		wg.Add(1)

		go func(order models.Order) {
			defer wg.Done()
			accrualInfo, err := s.accrualClient.GetAccrualInfo(order.Number)
			if err != nil {
				log.Println(err)
				return
			}

			// захватываем мьютекс, чтобы не было одновременного чтения и записи в слайс
			mu.Lock()
			accrualInfos = append(accrualInfos, accrualInfo)
			mu.Unlock()
		}(order)
	}
	// ждем пока все горутины закончат свою работу
	wg.Wait()

	if err = s.storage.UpdateOrders(ctx, username, accrualInfos); err != nil {
		return err
	}

	return nil
}

func (s *Service) ListUserOrders(ctx context.Context, username string) ([]models.Order, error) {
	if err := s.UpdateUserOrders(ctx, username); err != nil {
		return nil, err
	}

	orders, err := s.storage.ListUserOrders(ctx, username)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, models.ErrEmptyOrderList
	}

	return orders, nil
}
