package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baobei23/go-avdb/internal/store"
)

func TestCreateActor(t *testing.T) {
	// 1. Setup Mock
	mockStore := store.NewMockStore()

	app := &Application{
		Store: mockStore,
	}

	// 2. Prepare Request
	payload := actorPayload{
		Name: "Test Actor",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/actor", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// 3. Execute Handler
	app.createActor(w, req)

	// 4. Verification
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", resp.StatusCode)
	}
}
