//	@title			Social API
//	@description	This is a sample server for Social API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/v1

//@securityDefinitions.apikey ApiKeyAuth
//	@in	header
//	@name	Authorization
// @description

package main

import (
	"time"
	"github.com/hannanaarif/Social/internal/db"
	"github.com/hannanaarif/Social/internal/env"
	"github.com/hannanaarif/Social/internal/store"
	"go.uber.org/zap"
)

const version = "1.0.0"

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		apiUrl: env.GetString("EXTERNAL_URL", "http://localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp: time.Hour*24*3,
		},
	}

	logger:=zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		// log.Fatal(err)
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	defer db.Close()
	logger.Info("database connection established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()
	// log.Fatal(app.Run(mux))
	logger.Fatal("failed to run server", zap.Error(app.Run(mux)))
}
