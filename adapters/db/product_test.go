package db_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/andreis3/go-hexagonal/adapters/db"
	"github.com/andreis3/go-hexagonal/application"
)

var Db *sql.DB

func setUp() {
	Db, _ = sql.Open("sqlite3", ":memory:")
	createTable(Db)
	createProduct(Db)
}

func createTable(db *sql.DB) {
	table := `CREATE TABLE products (
		"id" string,
		"name" string,
		"price" float,
		"status" string
	);`

	stmt, err := db.Prepare(table)
	if err != nil {
		log.Fatal(err.Error())
	}

	stmt.Exec()
}

func createProduct(db *sql.DB) {
	insert := `INSERT into products values("abc","Product Test",0,"disabled")`
	stmt, err := db.Prepare(insert)
	if err != nil {
		log.Fatal(err.Error())
	}
	stmt.Exec()
}

func TestProductCreate(t *testing.T) {
	setUp()
	defer Db.Close()

	productDb := db.NewProductDb(Db)
	product, err := productDb.Get("abc")

	require.Nil(t, err)
	require.Equal(t, "Product Test", product.GetName())
	require.Equal(t, 0.0, product.GetPrice())
	require.Equal(t, "disabled", product.GetStatus())
}

func TestProductDb_Save(t *testing.T) {
	setUp()
	defer Db.Close()
	productDb := db.NewProductDb(Db)

	product := application.NewProduct()
	product.Name = "Product Test"
	product.Price = 256

	productResult, err := productDb.Save(product)

	require.Nil(t, err)
	require.Equal(t, productResult.GetName(), product.Name)
	require.Equal(t, productResult.GetPrice(), product.Price)
	require.Equal(t, productResult.GetStatus(), product.Status)

	product.Status = "enabled"
	require.Equal(t, productResult.GetName(), product.Name)
	require.Equal(t, productResult.GetPrice(), product.Price)
	require.Equal(t, productResult.GetStatus(), product.Status)
}
