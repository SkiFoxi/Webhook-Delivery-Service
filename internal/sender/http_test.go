package sender_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/domain"
	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/sender"
)

func TestHTTPSender_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	s := sender.NewHTTPSender(http.DefaultClient, 5*time.Second)

	err := s.Send(context.Background(), domain.DeliveryJob{
		DeliveryID: "d-1",
		WebhookURL: srv.URL,
		Payload:    []byte(`{"hello":"world"}`),
		Secret:     "secret",
	})

	if err != nil {
		t.Errorf("err = %v; want nil", err)
	}
}

func TestHTTPSender_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	sender := sender.NewHTTPSender(http.DefaultClient, 5*time.Second)
	err := sender.Send(context.Background(), domain.DeliveryJob{
		WebhookURL: srv.URL,
		Payload:    []byte(`{}`),
	})

	if !errors.Is(err, domain.ErrDeliveryFailed) {
		t.Errorf("err = %v; want ErrDeliveryFailed", err)
	}
}
