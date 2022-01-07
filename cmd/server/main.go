package main

import (
	"context"
	"fmt"
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
			web.Error(c, http.StatusInternalServerError, "Error in store products")
			return
		}

		_, err = customerService.StoreBulk(ctx)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, "Error in store customers")
			return
		}

		_, err = invoiceService.StoreBulk(ctx)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, "Error in store invoices")
			return
		}

		_, err = saleService.StoreBulk(ctx)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, "Error in store sales")
			return
		}

		log.Println("Insert's work fine")

		_, err = invoiceService.UpdateTotal(context.Background())
		if err != nil {
			fmt.Println(err)
			web.Error(c, http.StatusInternalServerError, "Error updating invoice totals!")
			return
		}

		web.Success(c, http.StatusOK, "Data loaded!")
	}
}
