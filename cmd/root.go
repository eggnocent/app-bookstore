package cmd

import (
	"app-bookstore/config"
	"app-bookstore/lib"
	"net/http"
	"os"

	"app-bookstore/database/seeder"

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
	cobra.OnInitialize(initConfig, initLogger, initDatabase, initJWT)
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

	seeder.SeedUsers(dbPool)

	log.Info().Msg("Database connect successfullly...")
}

func initJWT() {
	jwtService = lib.NewJWT(cfg)
}

func startServer(cmd *cobra.Command, args []string) {
	r := mux.NewRouter()

	log.Info().Msg("Starting server...")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the Home Page!"))
	}).Methods("GET")

	log.Info().Msgf("Server running on port %s", cfg.App.AppPort)
	if err := http.ListenAndServe(":"+cfg.App.AppPort, r); err != nil {
		log.Fatal().Msgf("Failed to start server: %v", err)
	}

}
