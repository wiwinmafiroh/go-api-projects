package product_postgres

import (
	"06-go-api-with-unittest/entity"
	"06-go-api-with-unittest/pkg/errs"
	"06-go-api-with-unittest/repository/product_repository"
	"database/sql"
	"time"
)

const (
	createProductQuery = `
		INSERT INTO products (name, description, price, image_url, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, price, image_url, user_id, created_at;
	`

	getAllProducts = `
		SELECT * FROM products
		ORDER BY id ASC;
	`

	getUserProducts = `
		SELECT * FROM products
		WHERE user_id = $1
		ORDER BY id ASC;
	`

	getProductByIdQuery = `
		SELECT * FROM products
		WHERE id = $1;
	`

	updateProductByIdQuery = `
		UPDATE products
		SET name = $2, description = $3, price = $4, image_url = $5, updated_at = $6
		WHERE id = $1
		RETURNING id, name, description, price, image_url, user_id, updated_at;
	`

	deleteProductByIdQuery = `
		DELETE FROM products
		WHERE id = $1;
	`
)

type productPostgres struct {
	db *sql.DB
}

func NewProductPostgres(db *sql.DB) product_repository.ProductRepository {
	return &productPostgres{
		db: db,
	}
}

func (p *productPostgres) CreateProduct(productEntity entity.Product) (*entity.Product, errs.ErrorMessage) {
	row := p.db.QueryRow(createProductQuery, productEntity.Name, productEntity.Description, productEntity.Price, productEntity.ImageUrl, productEntity.UserId)

	var createdProduct entity.Product

	err := row.Scan(&createdProduct.Id, &createdProduct.Name, &createdProduct.Description, &createdProduct.Price, &createdProduct.ImageUrl, &createdProduct.UserId, &createdProduct.CreatedAt)
	if err != nil {
		return nil, errs.NewInternalServerError("Failed to create new product")
	}

	return &createdProduct, nil
}

func (p *productPostgres) GetAllProducts() ([]*entity.Product, errs.ErrorMessage) {
	rows, err := p.db.Query(getAllProducts)
	if err != nil {
		return nil, errs.NewInternalServerError("Failed to retrieve products data")
	}

	var retrievedProducts []*entity.Product

	defer rows.Close()

	for rows.Next() {
		var retrievedProduct entity.Product

		err = rows.Scan(&retrievedProduct.Id, &retrievedProduct.Name, &retrievedProduct.Description, &retrievedProduct.Price, &retrievedProduct.ImageUrl, &retrievedProduct.UserId, &retrievedProduct.CreatedAt, &retrievedProduct.UpdatedAt)
		if err != nil {
			return nil, errs.NewInternalServerError("Failed to retrieve product details")
		}

		retrievedProducts = append(retrievedProducts, &retrievedProduct)
	}

	return retrievedProducts, nil
}

func (p *productPostgres) GetUserProducts(userId uint) ([]*entity.Product, errs.ErrorMessage) {
	rows, err := p.db.Query(getUserProducts, userId)
	if err != nil {
		return nil, errs.NewInternalServerError("Failed to retrieve products data")
	}

	var retrievedProducts []*entity.Product

	defer rows.Close()

	for rows.Next() {
		var retrievedProduct entity.Product

		err = rows.Scan(&retrievedProduct.Id, &retrievedProduct.Name, &retrievedProduct.Description, &retrievedProduct.Price, &retrievedProduct.ImageUrl, &retrievedProduct.UserId, &retrievedProduct.CreatedAt, &retrievedProduct.UpdatedAt)
		if err != nil {
			return nil, errs.NewInternalServerError("Failed to retrieve product details")
		}

		retrievedProducts = append(retrievedProducts, &retrievedProduct)
	}

	return retrievedProducts, nil
}

func (p *productPostgres) GetProductById(productId uint) (*entity.Product, errs.ErrorMessage) {
	row := p.db.QueryRow(getProductByIdQuery, productId)

	var retrievedProduct entity.Product

	err := row.Scan(&retrievedProduct.Id, &retrievedProduct.Name, &retrievedProduct.Description, &retrievedProduct.Price, &retrievedProduct.ImageUrl, &retrievedProduct.UserId, &retrievedProduct.CreatedAt, &retrievedProduct.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Product not found")
		}

		return nil, errs.NewInternalServerError("Failed to retrieve product data")
	}

	return &retrievedProduct, nil
}

func (p *productPostgres) UpdateProductById(productEntity entity.Product) (*entity.Product, errs.ErrorMessage) {
	row := p.db.QueryRow(updateProductByIdQuery, productEntity.Id, productEntity.Name, productEntity.Description, productEntity.Price, productEntity.ImageUrl, time.Now())

	var updatedProduct entity.Product

	err := row.Scan(&updatedProduct.Id, &updatedProduct.Name, &updatedProduct.Description, &updatedProduct.Price, &updatedProduct.ImageUrl, &updatedProduct.UserId, &updatedProduct.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Product not found")
		}

		return nil, errs.NewInternalServerError("Failed to update product")
	}

	return &updatedProduct, nil
}

func (p *productPostgres) DeleteProductById(productId uint) errs.ErrorMessage {
	_, err := p.db.Exec(deleteProductByIdQuery, productId)
	if err != nil {
		return errs.NewInternalServerError("Failed to delete product")
	}

	return nil
}
