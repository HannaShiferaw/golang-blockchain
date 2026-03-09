package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"coffee-consortium/backend/internal/contract"
	"coffee-consortium/backend/internal/repo/couchdb"
	"coffee-consortium/backend/internal/repo/memstore"
	"coffee-consortium/backend/internal/repo/postgres"
	ledgerSvc "coffee-consortium/backend/internal/service/ledger"
	"coffee-consortium/backend/internal/service/identity"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	srv *http.Server
}

func NewApp() (*App, error) {
	var (
		state contract.StateStore
		bs    ledgerSvc.BlockStore
		idRepo identity.Repo
	)
	couchURL := os.Getenv("COUCHDB_URL")
	if couchURL != "" {
		user := os.Getenv("COUCHDB_USER")
		pass := os.Getenv("COUCHDB_PASSWORD")
		db := os.Getenv("COUCHDB_DB")
		if db == "" {
			db = "coffee_ledger"
		}
		client, err := couchdb.New(couchURL, user, pass, db)
		if err != nil {
			return nil, err
		}
		if err := retry(context.Background(), 25, 2*time.Second, func() error { return client.EnsureDB(context.Background()) }); err != nil {
			return nil, err
		}
		state = couchdb.NewStateStore(client)
		bs = couchdb.NewBlockStore(client)
		idRepo = couchdb.NewPKIRepo(client)
	} else {
		state = memstore.New()
		bs = memstore.NewBlockStore()
	}

	ids, err := identity.NewService(idRepo)
	if err != nil {
		return nil, err
	}
	if os.Getenv("SEED_ON_START") == "true" {
		if err := ids.SeedDefaults(); err != nil {
			return nil, err
		}
	}

	var ix ledgerSvc.TxIndexer
	if os.Getenv("POSTGRES_HOST") != "" {
		var db *postgres.DB
		if err := retry(context.Background(), 25, 2*time.Second, func() error {
			var err error
			db, err = postgres.ConnectFromEnv(context.Background())
			return err
		}); err != nil {
			return nil, err
		}
		if err := retry(context.Background(), 25, 2*time.Second, func() error { return db.EnsureSchema(context.Background()) }); err != nil {
			return nil, err
		}
		ix = postgres.NewIndexer(db, state)
	}

	led := ledgerSvc.New(ids, state, bs, ix)

	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	v1 := r.Group("/api/v1")
	registerV1Routes(v1, ids, led)

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{srv: srv}, nil
}

func (a *App) Run() error {
	return a.srv.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}

