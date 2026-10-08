package sender

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/domain"
)

type HTTPSender struct {
	client  *http.Client
	timeout time.Duration
}

func NewHTTPSender(client *http.Client, timeout time.Duration) *HTTPSender {
	return &HTTPSender{
		client:  client,
		timeout: timeout,
	}
}

func (s *HTTPSender) Send(ctx context.Context, job domain.DeliveryJob) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, job.WebhookURL, bytes.NewReader(job.Payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-ID", job.DeliveryID)

	mac := hmac.New(sha256.New, []byte(job.Secret))
	mac.Write(job.Payload)
	signature := hex.EncodeToString(mac.Sum(nil))
	req.Header.Set("X-Webhook-Signature", signature)

	resp, err := s.client.Do(req)

	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrDeliveryFailed, err)
	}
	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body) // Очищаем запрос для переиспользования нашей системой

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%w: status %d", domain.ErrDeliveryFailed, resp.StatusCode)
	}

	return nil
}
