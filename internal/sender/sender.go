package sender

import (
	"context"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/domain"
)

type DeliverySender interface {
	Send(ctx context.Context, job domain.DeliveryJob) error
}

