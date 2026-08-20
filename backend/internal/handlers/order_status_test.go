package handlers

import (
	"testing"

	"cupcake-delivery/internal/models"
)

func TestValidAdminTransition(t *testing.T) {
	if !validAdminTransition(models.StatusPending, models.StatusPreparing) ||
		!validAdminTransition(models.StatusPreparing, models.StatusReady) ||
		validAdminTransition(models.StatusPending, models.StatusDelivered) {
		t.Fatal("transições de pedido inválidas")
	}
}
