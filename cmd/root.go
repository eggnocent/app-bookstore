package cmd

import (
	"app-bookstore/config"
	"app-bookstore/database/seeder"
	v1 "app-bookstore/delivery/v1"
	"app-bookstore/helper"
	"app-bookstore/lib"
	"app-bookstore/router"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	cfg        *config.Config
	dbPool     *sqlx.DB
	jwtService lib.Jwt
)

var rootCmd = &cobra.Command{
	Use:   "app-bookstore",
	Short: "API for bookstore",
	Run:   startServer,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig, initLogger, initJWT, initDatabase)
}

func initConfig() {
	cfg = config.NewConfig()
	if cfg == nil {
		log.Fatal().Msg("Failed to initialize configuration...")
	}
}

func initLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	if cfg.App.AppEnv == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func initDatabase() {
	var err error
	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatal().Msgf("Failed to initialize database: %v", err)
	}

	dbPool = db.DB

	seeder.SeedSuperAdmin(dbPool)

	router.Init(dbPool, jwtService)

	log.Info().Msg("Database connect successfully...")
}

func initJWT() {
	jwtService = lib.NewJWT(cfg)
}

func startServer(cmd *cobra.Command, args []string) {
	r := mux.NewRouter()

	log.Info().Msg("Starting server...")
	mw := helper.NewMiddleware(jwtService, dbPool)

	public := r.PathPrefix("/api/v1").Subrouter()
	v1.NewAPIUser(public)
	v1.NewAPIResetPass(public)

	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(mw.CheckAccess)

	v1.NewAPIRole(protected)
	v1.NewAPIUserRequest(protected)
	v1.NewAPIUserRoles(protected)
	v1.NewAPIResource(protected)
	v1.NewAPIRoleResource(protected)
	v1.NewAPIAuthors(protected)
	v1.NewAPIPublisher(protected)
	v1.NewAPICategories(protected)
	v1.NewAPIBooks(protected)
	v1.NewAPILoans(protected)
	v1.NewAPIRating(protected)

	log.Info().Msgf("Server running on port %s", cfg.App.AppPort)
	if err := http.ListenAndServe(":"+cfg.App.AppPort, r); err != nil {
		log.Fatal().Msgf("Failed to start server: %v", err)
	}

}
