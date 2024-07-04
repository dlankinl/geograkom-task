package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"os"
	"task/internal/app"
	"task/internal/delivery"
)

type config struct {
	srvAddr string
	connStr string
}

func checkVersion(db *pgx.Conn) (ok bool, err error) {
	query := `
		select count(1)
		from information_schema.tables 
		where table_name = 'routes'
		limit 1`

	var val int
	err = db.QueryRow(context.Background(), query).Scan(&val)
	if err != nil {
		return false, fmt.Errorf("check table: %w", err)
	}
	if val == 1 {
		ok = true
	}

	return ok, nil
}

func newConn(ctx context.Context, connStr string) (db *pgx.Conn, err error) {
	db, err = pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("database connecting: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("database ping: %w", err)
	}

	return db, nil
}

func parseVariables() (cfg config, err error) {
	srvAddressFlag := flag.String("a", "", "Server address")
	connStrFlag := flag.String("b", "", "Database connection string")
	flag.Parse()

	srvAddr := os.Getenv("SERVER_ADDRESS")
	connStr := os.Getenv("CONNECTION_STRING")

	if srvAddr == "" {
		srvAddr = *srvAddressFlag
	}
	if connStr == "" {
		connStr = *connStrFlag
	}

	if srvAddr == "" {
		return config{}, fmt.Errorf(`set env variable SERVER_ADDRESS or use "-a" flag`)
	}
	if connStr == "" {
		return config{}, fmt.Errorf(`set env variable CONNECTION_STRING or use "-b" flag`)
	}

	return config{
		srvAddr: srvAddr,
		connStr: connStr,
	}, nil
}

func main() {
	cfg, err := parseVariables()
	if err != nil {
		log.Fatal("reading config: %w", err)
	}

	db, err := newConn(context.Background(), cfg.connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	ok, err := checkVersion(db)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("you need to run migrations before running server")
	}

	a := app.NewApp(db)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(middleware.Logger)

	router.Route("/api/route", func(r chi.Router) {
		r.Post("/register", delivery.RegisterHandler(a))
		r.Get("/{id}", delivery.GetHandler(a))
		r.Delete("/", delivery.DeleteHandler(a))
	})

	fmt.Printf("server is running on %s\n", cfg.srvAddr)
	err = http.ListenAndServe(cfg.srvAddr, router)
	if err != nil {
		log.Fatal(err)
	}
}
