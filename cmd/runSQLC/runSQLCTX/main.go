package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
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
	password = "88658710"
	dbname   = "courses"
)

type CourseDB struct {
	dbConn *sql.DB
	*db.Queries
}

func NewCourseDB(dbConn *sql.DB) *CourseDB {
	return &CourseDB{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

type CourseParams struct {
	ID          string
	Name        string
	Description sql.NullString
	Price float64
}

type CategoryParams struct {
	ID          string
	Name        string
	Description sql.NullString
}

func (c *CourseDB) callTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := c.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := db.New(tx)
	err = fn(q)

	if err != nil {
		if errRB := tx.Rollback(); errRB != nil {
			return fmt.Errorf("error on rollback: %v, original error: %w", errRB, err)
		}
		return err
	}
	return tx.Commit()
}

func (c *CourseDB) CreateCourseAndCategory(ctx context.Context, argsCategory CategoryParams, argsCourse CourseParams) error {
	err := c.callTx(ctx, func(q *db.Queries) error { // Tudo que for executado aqui, vai ser uma transacion

		var err error
		err = q.CreateCategory(ctx, db.CreateCategoryParams{
			ID:          argsCategory.ID,
			Name:        argsCategory.Name,
			Description: argsCategory.Description,
		})
		if err != nil {
			return err
		}
		err = q.CreateCourse(ctx, db.CreateCourseParams{
			ID:          argsCourse.ID,
			Name:        argsCourse.Name,
			Description: argsCourse.Description,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err

	}

	return nil
}

func main() {
	// trunk-ignore(golangci-lint/typecheck)
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

	//queries := db.New(dbConn)

	courseArgs := CourseParams{
		ID: uuid.New().String(),
		Name: "GO",
		Description: sql.NullString{String: "Go Course", Valid: true},
	}

	categoryArgs := CategoryParams{
		ID: uuid.New().String(),
        Name: "GO",
        Description: sql.NullString{String: "Go Category", Valid: true},
	}

	courseDB := NewCourseDB(dbConn)
	err = courseDB.CreateCourseAndCategory(ctx, categoryArgs, courseArgs)
	if err!= nil {
        log.Panic(err)
    }


}
