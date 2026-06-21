package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"time"

	// 1. Updated to v5 import path
	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := newConn()
	if err != nil {
		fmt.Println("newDB", err)
		return
	}

	defer func() {
		_ = conn.Close(context.Background())
		fmt.Println("closed")
	}()

	now := time.Now()

	if err := insertsUsers(conn); err != nil {
		fmt.Println("failed", err)
		return
	}

	fmt.Println("total", time.Since(now))
}

func insertsUsers(conn *pgx.Conn) error {
	f, err := os.Open(path.Join("..", "data.csv"))
	if err != nil {
		return fmt.Errorf("os.Open %w", err)
	}
	defer f.Close() // Good practice to close file handles

	cr := csv.NewReader(f)

	// 2. Migrated from conn.BeginTxFunc to package-level pgx.BeginTxFunc(ctx, conn, options, ...)
	if err := pgx.BeginTxFunc(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		for {
			record, err := cr.Read()
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				return fmt.Errorf("cr.Read %w", err)
			}

			// 3. FIXED BUG: Changed 'conn.Exec' to 'tx.Exec' to actually use the transaction context!
			_, err = tx.Exec(context.Background(), "INSERT INTO users(first_name, last_name) VALUES($1, $2)", record[0], record[1])
			if err != nil {
				return fmt.Errorf("tx.Exec %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("pgx.BeginTxFunc %w", err)
	}

	return nil
}

func newConn() (*pgx.Conn, error) {
	dsn := url.URL{
		Scheme: "postgres",
		Host:   "localhost:5432",
		User:   url.UserPassword("postgres", "postgres"),
		Path:   "dbname",
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")

	dsn.RawQuery = q.Encode()

	conn, err := pgx.Connect(context.Background(), dsn.String())
	if err != nil {
		// 4. Added missing space in formatting string
		return nil, fmt.Errorf("pgx.Connect: %w", err)
	}

	return conn, nil
}
