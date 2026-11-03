package catalog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProducts(t *testing.T) {
	// Servidor HTTP simulado.
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodGet {
				t.Errorf("se esperaba GET, se recibió %s", r.Method)
			}

			if r.URL.Query().Get("ids") != "1,2" {
				t.Errorf("IDs incorrectos: %s", r.URL.Query().Get("ids"))
			}

			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusOK)

			w.Write([]byte(`[
				{"id": 1, "name": "Teclado"},
				{"id": 2, "name": "Mouse"}
			]`))
		},
	))

	defer server.Close()

	client := NewHTTPClient(server.URL)

	products, err := client.GetProducts(
		context.Background(),
		[]int64{1, 2},
	)

	if err != nil {
		t.Fatalf("no se esperaba un error: %v", err)
	}

	if len(products) != 2 {
		t.Fatalf("se esperaban 2 productos, llegaron %d", len(products))
	}

	if products[0].Name != "Teclado" {
		t.Errorf("nombre incorrecto: %s", products[0].Name)
	}
}
