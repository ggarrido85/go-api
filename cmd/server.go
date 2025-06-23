package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ggarrido85/api-backend/api"
	//"github.com/ggarrido85/api-backend/infra"
	"github.com/ggarrido85/api-backend/utils"
//	"github.com/ggarrido85/api-backend/repositories"

	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
)

func RunServer(config CompiledConfig) error {
	// This is where we read the environment variables and set up the configuration for the application.
	apiConfig := api.Configuration{
		Env:                 utils.GetEnv("ENV", "development"),
		AppName:             "marble-backend",
		AppUrl:        utils.GetEnv("APP_URL", "127.0.0.1"),
		BackofficeUrl: utils.GetEnv("BACKOFFICE_URL", ""),
		Port:                "8080", //"utils.GetRequiredEnv[string]("PORT"),
		RequestLoggingLevel: utils.GetEnv("REQUEST_LOGGING_LEVEL", "all"),
		TokenLifetimeMinute: utils.GetEnv("TOKEN_LIFETIME_MINUTE", 60*2),
		SegmentWriteKey:     utils.GetEnv("SEGMENT_WRITE_KEY", config.SegmentWriteKey),
		DefaultTimeout:      time.Duration(utils.GetEnv("DEFAULT_TIMEOUT_SECOND", 5)) * time.Second,
	}
/*
	pgConfig := infra.PgConfig{
		ConnectionString:   utils.GetEnv("PG_CONNECTION_STRING", ""),
		Database:           utils.GetEnv("PG_DATABASE", "marble"),
		Hostname:           utils.GetEnv("PG_HOSTNAME", "127.0.0.1"),
		Password:           utils.GetEnv("PG_PASSWORD", "marble"),
		Port:               utils.GetEnv("PG_PORT", "5433"),
		User:               utils.GetEnv("PG_USER", "marble"),
		MaxPoolConnections: utils.GetEnv("PG_MAX_POOL_SIZE", infra.DEFAULT_MAX_CONNECTIONS),
		ClientDbConfigFile: utils.GetEnv("CLIENT_DB_CONFIG_FILE", ""),
		SslMode:            utils.GetEnv("PG_SSL_MODE", "prefer"),
	}
*/
	serverConfig := struct {
		batchIngestionMaxSize            int
		caseManagerBucket                string
		ingestionBucketUrl               string
		offloadingBucketUrl              string
		jwtSigningKey                    string
		jwtSigningKeyFile                string
		loggingFormat                    string
		sentryDsn                        string
		transferCheckEnrichmentBucketUrl string
		firebaseEmulatorHost             string
	}{
		batchIngestionMaxSize:            utils.GetEnv("BATCH_INGESTION_MAX_SIZE", 0),
		caseManagerBucket:                utils.GetEnv("CASE_MANAGER_BUCKET_URL", ""),
		ingestionBucketUrl:               utils.GetEnv("INGESTION_BUCKET_URL", ""),
		offloadingBucketUrl:              utils.GetEnv("OFFLOADING_BUCKET_URL", ""),
		jwtSigningKey:                    utils.GetEnv("AUTHENTICATION_JWT_SIGNING_KEY", ""),
		jwtSigningKeyFile:                utils.GetEnv("AUTHENTICATION_JWT_SIGNING_KEY_FILE", ""),
		loggingFormat:                    utils.GetEnv("LOGGING_FORMAT", "text"),
		sentryDsn:                        utils.GetEnv("SENTRY_DSN", ""),
		transferCheckEnrichmentBucketUrl: utils.GetEnv("TRANSFER_CHECK_ENRICHMENT_BUCKET_URL", ""), // required for transfercheck
		firebaseEmulatorHost:             utils.GetEnv("FIREBASE_AUTH_EMULATOR_HOST", ""),
	}

	logger := utils.NewLogger(serverConfig.loggingFormat)

	ctx := utils.StoreLoggerInContext(context.Background(), logger)
	// marbleJwtSigningKey := infra.ReadParseOrGenerateSigningKey(ctx, serverConfig.jwtSigningKey, serverConfig.jwtSigningKeyFile)
	

	defer sentry.Flush(3 * time.Second)

	/*tracingConfig := infra.TelemetryConfiguration{
		ApplicationName: apiConfig.AppName,
		Enabled:         gcpConfig.EnableTracing,
		ProjectID:       gcpConfig.ProjectId,
	}*/
	//telemetryRessources, err := infra.InitTelemetry(tracingConfig, config.Version)
	/*if err != nil {
		utils.LogAndReportSentryError(ctx, err)
	}*/

	/*pool, err := infra.NewPostgresConnectionPool(ctx, pgConfig.GetConnectionString(),
		telemetryRessources.TracerProvider, pgConfig.MaxPoolConnections)
	if err != nil {
		utils.LogAndReportSentryError(ctx, err)
	}*/

	/*clientDbConfig, err := infra.ParseClientDbConfig(pgConfig.ClientDbConfigFile)
	if err != nil {
		utils.LogAndReportSentryError(ctx, err)
		// return err
	}*/

	/*repositories := repositories.NewRepositories(
		pool,
		"", 
	)*/


	deps := api.InitDependencies(ctx, apiConfig/*, pool /*marbleJwtSigningKey, nil*/)

	router := api.InitRouterMiddlewares(ctx, apiConfig, /*apiConfig.DisableSegment  false,*/
		)
	server := api.NewServer(router, apiConfig, deps.TokenHandler, logger)

	notify, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.InfoContext(ctx, "starting server", slog.String("version", config.Version), slog.String("port", apiConfig.Port))
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			//utils.LogAndReportSentryError(ctx, errors.Wrap(err, "Error while serving the app"))
		}
		logger.InfoContext(ctx, "server returned")
	}()

	<-notify.Done()
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	deps.SegmentClient.Close()

	err := server.Shutdown(shutdownCtx);
	if  err != nil {
		/*utils.LogAndReportSentryError(
			ctx,
			errors.Wrap(err, "Error while shutting down the server"),
		)*/
		return err
	}

	return err
}
