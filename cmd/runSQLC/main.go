package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/VictorOliveiraPy/internal/db"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

const (
	host     = "172.27.0.2"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "courses"
)

func main() {
	ctx := context.Background()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://sql/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	queries := db.New(dbConn)
	// err = queries.CreateCategory(ctx, db.CreateCategoryParams{
	// 	ID: uuid.New().String(),
	// 	Name: "Backend",
	// 	Description: sql.NullString{String: "test backend", Valid: true},
	// })
	// if err!= nil {
	//     log.Fatal(err)
	// }
	// categories, err := queries.ListCategories(ctx)
	// if err!= nil {
	//     log.Fatal(err)
	// }

	// for _, category := range categories{
	// 	fmt.Println(category.Name, category.Description, category.ID)
	// }

	// categoriesUpdated, err := queries.UpdateCategory(ctx, db.UpdateCategoryParams{
	// 	ID: "3874abdd-d5d1-4a73-ad58-72218915f6f5",
	//     Name: "Backend",
	//     Description: sql.NullString{String: "test backend updated 1", Valid: true},
	// })

	// if err!= nil {
	//     log.Fatal(err)
	// }
	// fmt.Println(categoriesUpdated)

	//  _, err = queries.DeleteCategory(ctx, "3874abdd-d5d1-4a73-ad58-72218915f6f5")
	// if err!= nil {
	// 		log.Fatal(err)
	// 	}

	getAllCategories, err := queries.GetCategories(ctx, "f36a7c9a-95a2-4f3d-aa1f-01f22bdb2a29")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(getAllCategories)

}
