package payment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/PabloGolobaro/cosmic_factory/payment/internal/errors"
)

func TestPaySuccess(t *testing.T) {
	svc := NewPaymentService()
	txUUID, err := svc.Pay(context.Background(), uuid.NewString(), "CARD")
	require.NoError(t, err)
	assert.NotEmpty(t, txUUID)
	_, parseErr := uuid.Parse(txUUID)
	assert.NoError(t, parseErr)
}

func TestPayEmptyPaymentMethod(t *testing.T) {
	svc := NewPaymentService()
	_, err := svc.Pay(context.Background(), uuid.NewString(), "")
	require.ErrorIs(t, err, errs.ErrInvalidPaymentMethod)
}

func TestPayInvalidUUID(t *testing.T) {
	svc := NewPaymentService()
	_, err := svc.Pay(context.Background(), "not-a-uuid", "CARD")
	require.ErrorIs(t, err, errs.ErrInvalidUUID)
}
