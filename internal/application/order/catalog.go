package order

import "context"

// Product representa la información del producto que necesitamos.
type CatalogProduct struct {
    ID   int64
    Name string
}

// ProductCatalog define las operaciones que necesitamos del catálogo.
type ProductCatalog interface {
    GetProducts(ctx context.Context, productIDs []int64) ([]CatalogProduct, error)
}