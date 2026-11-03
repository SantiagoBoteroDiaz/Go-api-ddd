package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	appOrder "api-ddd/internal/application/order"
)

type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Estructura que representa la respuesta JSON del catálogo externo.
type productResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (c *HTTPClient) GetProducts(
	ctx context.Context,
	productIDs []int64,
) ([]appOrder.CatalogProduct, error) {

	if len(productIDs) == 0 {
		return nil, fmt.Errorf("no se proporcionaron productos")
	}

	// Construimos los parámetros de la URL.
	params := url.Values{}
	var ids []string

	for _, id := range productIDs {
		ids = append(ids, strconv.FormatInt(id, 10))
	}

	params.Set("ids", strings.Join(ids, ","))

	endpoint := c.baseURL + "/products?" + params.Encode()

	// Creamos la solicitud asociada al contexto recibido.
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("creando solicitud al catálogo: %w", err)
	}

	// Ejecutamos la solicitud HTTP.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("consultando catálogo: %w", err)
	}
	defer resp.Body.Close()

	// Validamos el código HTTP.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"el catálogo respondió con estado HTTP %d",
			resp.StatusCode,
		)
	}

	// Decodificamos el JSON.
	var productsResponse []productResponse

	if err := json.NewDecoder(resp.Body).Decode(&productsResponse); err != nil {
		return nil, fmt.Errorf("decodificando respuesta del catálogo: %w", err)
	}

	// Convertimos la respuesta externa al modelo que utiliza nuestra aplicación.
	products := make([]appOrder.CatalogProduct, 0, len(productsResponse))

	for _, product := range productsResponse {
		products = append(products, appOrder.CatalogProduct{
			ID:   product.ID,
			Name: product.Name,
		})
	}

	return products, nil
}

// Comprobamos en compilación que implementamos la interfaz.
var _ appOrder.ProductCatalog = (*HTTPClient)(nil)
