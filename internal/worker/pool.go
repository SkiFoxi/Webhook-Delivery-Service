package worker

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/domain"
	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/sender"
)

type Pool struct {
	jobs          chan domain.DeliveryJob   //Очередь для задач
	workers      int                       //Сколько воркеров
	sender       sender.DeliverySender     //Кто отправляет
	deliveryRepo domain.DeliveryRepository //Куда писать результат
	wg           sync.WaitGroup
}

func NewPool(
	workers int,
	queueSize int,
	sender sender.DeliverySender,
	deliveryRepo domain.DeliveryRepository,
) *Pool {
	jobs := make(chan domain.DeliveryJob, queueSize)
	return &Pool{
		jobs:          jobs,
		workers:      workers,
		sender:       sender,
		deliveryRepo: deliveryRepo,
	}
}

func (p *Pool) process(ctx context.Context, job domain.DeliveryJob) {
	if err := p.sender.Send(ctx, job); err != nil {
		updateCtx := context.WithoutCancel(ctx)
		if errA := p.deliveryRepo.UpdateStatus(updateCtx, job.DeliveryID, domain.DelivFailed, err.Error()); errA != nil {
			slog.Error("cannot update status", "err", err, "delivery_id", job.DeliveryID)
		}

		slog.Error("delivery failed", "err", err, "delivery_id", job.DeliveryID)
		return
	}
	updateCtx := context.WithoutCancel(ctx)
	if err := p.deliveryRepo.UpdateStatus(updateCtx, job.DeliveryID, domain.DelivSuccess, ""); err != nil {
		slog.Error("cannot update status", "err", err, "delivery_id", job.DeliveryID)
	}
}

func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			slog.Info("worker stopped", "id", id)
			return
		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			p.process(ctx, job)
		}
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}
}

func (p *Pool) Submit(job domain.DeliveryJob) error {
	select {
	case p.jobs <- job:
		return nil
	default:
		return errors.New("queue full")
	}
}

func (p *Pool) Stop() {
	close(p.jobs)
	p.wg.Wait()
}
