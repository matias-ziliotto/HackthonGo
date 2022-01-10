package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matias-ziliotto/HackthonGo/internal/customer"
	"github.com/matias-ziliotto/HackthonGo/internal/invoice"
	"github.com/matias-ziliotto/HackthonGo/internal/product"
	"github.com/matias-ziliotto/HackthonGo/internal/sale"
	"github.com/matias-ziliotto/HackthonGo/pkg/database/sql"
	"github.com/matias-ziliotto/HackthonGo/pkg/web"
)

func main() {
	router := gin.Default()

	router.GET("/load-files", LoadData())
	router.GET("/customers/total-by-condition", GetCustomersTotalByCondition())
	router.GET("/products/top/most-selled", GetProductsMostSelled())
	router.GET("/customers/top/cheaper-products", GetCustomersCheaperProducts())

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}

func LoadData() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		dbProduct := sql.MySqlDB
		dbCustomer := sql.MySqlDB
		dbInvoice := sql.MySqlDB
		dbSale := sql.MySqlDB

		// Products
		productRepository := product.NewProductRepository(dbProduct)
		productService := product.NewProductService(productRepository)

		// Customers
		customerRepository := customer.NewCustomerRepository(dbCustomer)
		customerService := customer.NewCustomerService(customerRepository)

		// Invoices
		invoiceRepository := invoice.NewInvoiceRepository(dbInvoice)
		invoiceService := invoice.NewInvoiceService(invoiceRepository)

		// Sales
		saleRepository := sale.NewSaleRepository(dbSale)
		saleService := sale.NewSaleService(saleRepository)

		_, err := productService.StoreBulk(ctx)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = customerService.StoreBulk(ctx)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = invoiceService.StoreBulk(ctx)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = saleService.StoreBulk(ctx)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = invoiceService.UpdateTotal(context.Background())
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		web.Success(c, http.StatusOK, "Data loaded!")
	}
}

func GetCustomersTotalByCondition() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		dbCustomer := sql.MySqlDB

		// Customers
		customerRepository := customer.NewCustomerRepository(dbCustomer)
		customerService := customer.NewCustomerService(customerRepository)

		customerTotalByConditionDTO, err := customerService.GetTotalByCondition(ctx)

		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		web.Success(c, http.StatusOK, customerTotalByConditionDTO)
	}
}

func GetProductsMostSelled() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		dbProducts := sql.MySqlDB

		// Products
		productRepository := product.NewProductRepository(dbProducts)
		productService := product.NewProductService(productRepository)

		productsMostSelled, err := productService.GetProductsMostSelled(ctx)

		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		web.Success(c, http.StatusOK, productsMostSelled)
	}
}

func GetCustomersCheaperProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		dbCustomer := sql.MySqlDB

		// Customers
		customerRepository := customer.NewCustomerRepository(dbCustomer)
		customerService := customer.NewCustomerService(customerRepository)

		customerCheaperProducts, err := customerService.GetCustomerCheaperProducts(ctx)

		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		web.Success(c, http.StatusOK, customerCheaperProducts)
	}
}
