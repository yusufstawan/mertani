package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"mertani/internal/config"
)

const migrationsPath = "file://migrations"

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg := config.Load()

	m, err := migrate.New(migrationsPath, databaseURL(cfg))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		sourceErr, databaseErr := m.Close()
		if sourceErr != nil {
			log.Printf("close source: %v", sourceErr)
		}
		if databaseErr != nil {
			log.Printf("close database: %v", databaseErr)
		}
	}()

	if err := run(m, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(m *migrate.Migrate, args []string) error {
	switch args[0] {
	case "up":
		return ignoreNoChange(m.Up())
	case "down":
		if len(args) == 1 {
			return ignoreNoChange(m.Down())
		}

		steps, err := parseSteps(args[1])
		if err != nil {
			return err
		}

		return ignoreNoChange(m.Steps(-steps))
	case "steps":
		if len(args) < 2 {
			return errors.New("steps requires a number")
		}

		steps, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid steps value: %w", err)
		}

		return ignoreNoChange(m.Steps(steps))
	case "force":
		if len(args) < 2 {
			return errors.New("force requires a version number")
		}

		version, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid version value: %w", err)
		}

		return m.Force(version)
	case "version":
		version, dirty, err := m.Version()
		if errors.Is(err, migrate.ErrNilVersion) {
			fmt.Println("version: nil")
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Printf("version: %d dirty: %t\n", version, dirty)
		return nil
	default:
		printUsage()
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func parseSteps(value string) (int, error) {
	steps, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid steps value: %w", err)
	}
	if steps <= 0 {
		return 0, errors.New("steps must be greater than zero")
	}

	return steps, nil
}

func ignoreNoChange(err error) error {
	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("no migration changes")
		return nil
	}

	return err
}

func databaseURL(cfg config.Config) string {
	databaseURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:   net.JoinHostPort(cfg.DBHost, cfg.DBPort),
		Path:   cfg.DBName,
	}

	query := databaseURL.Query()
	query.Set("sslmode", cfg.DBSSLMode)
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String()
}

func printUsage() {
	fmt.Println("usage: go run ./cmd/migrate <up|down [steps]|steps <n>|force <version>|version>")
}
