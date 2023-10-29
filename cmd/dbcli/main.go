package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file" //no lint
	_ "github.com/jackc/pgx/stdlib"                   //no lint
	corepg "github.com/rubengomes8/HappyCore/pkg/postgres"
	"github.com/urfave/cli"
)

func setup() (string, corepg.Config, error) {

	dbPortEnv := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortEnv)
	if err != nil {
		return "", corepg.Config{}, err
	}
	c := corepg.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PWD"),
	}

	return os.Getenv("DB_MIGRATIONS_PATH"), c, nil
}

func instance(pool *sql.DB, source string) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(pool, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", source),
		"postgres", driver,
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func buildInstance(c corepg.Config, migrationsPath string) (*migrate.Migrate, error) {
	db, err := corepg.InitDB(c)
	if err != nil {
		return nil, err
	}

	return instance(db, migrationsPath)
}

func migrateDB(*cli.Context) error {
	migrationsPath, config, err := setup()
	if err != nil {
		return err
	}

	m, err := buildInstance(config, migrationsPath)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func rollbackDB(*cli.Context) error {

	migrationsPath, config, err := setup()
	if err != nil {
		return err
	}

	m, err := buildInstance(config, migrationsPath)
	if err != nil {
		return err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func main() {
	c := cli.NewApp()
	c.Commands = []cli.Command{
		{
			Name:   "migrate",
			Usage:  "migrate",
			Action: migrateDB,
		},
		{
			Name:   "rollback",
			Usage:  "rollback",
			Action: rollbackDB,
		},
	}

	err := c.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
