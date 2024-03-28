package processor

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/vindosVP/loyalty-system/internal/models"
	"github.com/vindosVP/loyalty-system/pkg/logger"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

var ErrTooManyRequests = errors.New("too many requests")

type Storage interface {
	GetUnprocessedOrders(ctx context.Context) ([]int, error)
	UpdateOrder(ctx context.Context, id int, status string, sum float64) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, id int, status string) (*models.Order, error)
}

type Processor struct {
	RequestInterval time.Duration
	ServerAddress   string
	Done            <-chan struct{}
	Storage         Storage
	Client          *resty.Client
}

type job struct {
	id            int
	serverAddress string
	client        *resty.Client
	order         int
	storage       Storage
}

type result struct {
	id    int
	order int
	err   error
}

type accrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

func New(RequestInterval time.Duration, ServerAddress string, Storage Storage) *Processor {
	return &Processor{
		RequestInterval: RequestInterval,
		ServerAddress:   ServerAddress,
		Storage:         Storage,
		Client:          resty.New(),
	}
}

func (p *Processor) Run() {
	tick := time.NewTicker(p.RequestInterval * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-p.Done:
			return
		case <-tick.C:
			p.requestAccruals()
		}
	}
}

func (p *Processor) requestAccruals() {
	ctx := context.Background()
	orders, err := p.Storage.GetUnprocessedOrders(ctx)
	if err != nil {
		logger.Log.Error("Failed to get unprocessed orders", zap.Error(err))
		return
	}
	jobs := p.generateJobs(orders)
	results := make(chan result)
	go listenResults(results)
	p.startWorkers(jobs, results, 10)
}

func listenResults(results <-chan result) {
	for res := range results {
		if res.err != nil {
			logger.Log.Error("worker failed", zap.Error(res.err), zap.Int("id", res.id))
		} else {
			logger.Log.Info("worker finished", zap.Int("id", res.id))
		}
	}
}

func (p *Processor) startWorkers(jobs <-chan job, results chan<- result, workers int) {
	wg := sync.WaitGroup{}
	logger.Log.Info(fmt.Sprintf("Starting %d workers", workers))
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}
	wg.Wait()
	close(results)
}

func worker(jobs <-chan job, results chan<- result, wg *sync.WaitGroup) {
	for j := range jobs {
		err := processOrder(j.client, j.serverAddress, j.order, j.storage)
		results <- result{j.id, j.order, err}
	}
	wg.Done()
}

func (p *Processor) generateJobs(orders []int) chan job {
	jobs := make(chan job)
	go func() {
		id := 0
		for _, order := range orders {
			jobs <- job{
				id:            id,
				serverAddress: p.ServerAddress,
				client:        p.Client,
				order:         order,
				storage:       p.Storage,
			}
			id++
		}
		defer close(jobs)
	}()
	return jobs
}

func processOrder(client *resty.Client, serverAddress string, order int, storage Storage) error {
	var response accrualResponse
	url := fmt.Sprintf("%s/api/orders/%s", serverAddress, strconv.Itoa(order))
	resp, err := client.R().SetResult(&response).Get(url)
	if err != nil {
		logger.Log.Error("Failed to get accruals by order", zap.Int("orderId", order), zap.Error(err))
	}
	if resp.StatusCode() != 200 {
		if resp.StatusCode() == 429 {
			logger.Log.Error("Failed to get accruals by order", zap.Int("orderId", order), zap.Error(ErrTooManyRequests))
		}
		logger.Log.Error("Failed to get accruals by order", zap.Int("orderId", order), zap.Int("statusCode", resp.StatusCode()))
		return err
	}
	id, _ := strconv.Atoi(response.Order)
	if response.Status == models.OrderStatusProcessed {
		_, err = storage.UpdateOrder(context.Background(), id, response.Status, response.Accrual)
	} else {
		_, err = storage.UpdateOrderStatus(context.Background(), id, response.Status)
	}

	if err != nil {
		return fmt.Errorf("storage.UpdateOrder: %w", err)
	}
	return nil
}
