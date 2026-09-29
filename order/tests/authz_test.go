//go:build apitest

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	authv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"
)

// seededManagerLogin — менеджер из миграции migrations/iam/*_seed_test_manager.sql.
const seededManagerLogin = "testmanager"

// doAs выполняет HTTP-запрос к Order от имени указанной сессии.
func doAs(t *testing.T, sessionUUID, method, path string, body any) *http.Response {
	t.Helper()

	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, orderBaseURL()+path, reader)
	require.NoError(t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+sessionUUID)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	t.Cleanup(func() { _ = resp.Body.Close() })

	return resp
}

func loginManager(t *testing.T) string {
	t.Helper()

	resp, err := authSvcClient.Login(context.Background(), &authv1.LoginRequest{
		Login:    seededManagerLogin,
		Password: "password123",
	})
	require.NoError(t, err)

	return resp.GetSessionUuid()
}

// createDefaultOrder создаёт заказ от имени дефолтного клиента.
func createDefaultOrder(t *testing.T) string {
	t.Helper()

	created, resp := createOrder(t, &CreateOrderRequest{HullUUID: HullAluminumUUID, EngineUUID: EngineIonCUUID})
	_ = resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	return created.OrderUUID
}

func TestAuthz_Client_ForeignOrderForbidden(t *testing.T) {
	orderUUID := createDefaultOrder(t)
	strangerSession, _ := registerAndLogin(t, "authz-stranger", "password123")

	orderPath := "/api/v1/orders/" + orderUUID

	require.Equal(t, http.StatusForbidden, doAs(t, strangerSession, http.MethodGet, orderPath, nil).StatusCode)
	require.Equal(t, http.StatusForbidden,
		doAs(t, strangerSession, http.MethodPost, orderPath+"/pay", PayOrderRequest{PaymentMethod: "CARD"}).StatusCode)
	require.Equal(t, http.StatusForbidden, doAs(t, strangerSession, http.MethodPost, orderPath+"/cancel", nil).StatusCode)

	// Отказ не должен менять заказ: владелец видит его в исходном статусе.
	order, resp := getOrder(t, orderUUID)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "PENDING_PAYMENT", order.Status)
}

func TestAuthz_Manager(t *testing.T) {
	orderUUID := createDefaultOrder(t)
	managerSession := loginManager(t)

	orderPath := "/api/v1/orders/" + orderUUID

	// Менеджер читает любой заказ.
	require.Equal(t, http.StatusOK, doAs(t, managerSession, http.MethodGet, orderPath, nil).StatusCode)

	// Менеджер не создаёт и не оплачивает заказы.
	require.Equal(t, http.StatusForbidden, doAs(t, managerSession, http.MethodPost, "/api/v1/orders",
		CreateOrderRequest{HullUUID: HullAluminumUUID, EngineUUID: EngineIonCUUID}).StatusCode)
	require.Equal(t, http.StatusForbidden,
		doAs(t, managerSession, http.MethodPost, orderPath+"/pay", PayOrderRequest{PaymentMethod: "CARD"}).StatusCode)

	// Менеджер отменяет любой заказ.
	require.Equal(t, http.StatusOK, doAs(t, managerSession, http.MethodPost, orderPath+"/cancel", nil).StatusCode)

	order, resp := getOrder(t, orderUUID)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "CANCELLED", order.Status)
}
