package server

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"study-planner/internal/auth"
	"study-planner/internal/auth/manager"
	"study-planner/internal/auth/platform"
	"study-planner/internal/auth/token"

	curriculumRepository "study-planner/internal/curriculum/repository"
	disciplineRepository "study-planner/internal/discipline/repository"
	institutionRepository "study-planner/internal/institution/repository"
	taskRepository "study-planner/internal/task/repository"
	userRepository "study-planner/internal/user/repository"

	authDelivery "study-planner/internal/auth/delivery"
	curriculumDelivery "study-planner/internal/curriculum/delivery"
	disciplineDelivery "study-planner/internal/discipline/delivery"
	institutionDelivery "study-planner/internal/institution/delivery"
	taskDelivery "study-planner/internal/task/delivery"
	userDelivery "study-planner/internal/user/delivery"

	"study-planner/internal/user"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

	_ "github.com/go-sql-driver/mysql"
)

var (
	FlagDatabaseHost = &cli.StringFlag{
		Name:    "db-host",
		Usage:   "Address of the database host",
		Value:   "localhost",
		EnvVars: []string{"DB_HOST"},
	}

	FlagDatabaseUser = &cli.StringFlag{
		Name:    "db-user",
		Usage:   "Name of the database user",
		Value:   "planner",
		EnvVars: []string{"DB_USER"},
	}

	FlagDatabasePassword = &cli.StringFlag{
		Name:    "db-password",
		Usage:   "Password of database user provided with --db-user flag",
		EnvVars: []string{"DB_PASSWORD"},
	}

	FlagDatabaseName = &cli.StringFlag{
		Name:    "db-name",
		Usage:   "Name of the database on the server provided with --db-host flag",
		Value:   "planner",
		EnvVars: []string{"DB_NAME"},
	}

	FlagBindAddress = &cli.StringFlag{
		Name:    "bind-address",
		Usage:   "Address the server should be bound to",
		Value:   ":8080",
		EnvVars: []string{"BIND_ADDRESS"},
	}

	FlagAllowedOrigins = &cli.StringSliceFlag{
		Name:    "allowed-origins",
		Usage:   "List of origins that will be allowed to pass CORS checks",
		EnvVars: []string{"ALLOWED_ORIGINS"},
	}

	FlagSigningKey = &cli.StringFlag{
		Name:     "jwt-signing-key",
		Usage:    "base64-encoded key for JWT signatures",
		Required: true,
		EnvVars:  []string{"JWT_SIGNING_KEY"},
	}

	FlagAudience = &cli.StringFlag{
		Name:    "jwt-audience",
		Usage:   "Audience value for JWT tokens",
		Value:   "localhost:5173",
		EnvVars: []string{"JWT_AUDIENCE"},
	}

	FlagTokenLifetime = &cli.DurationFlag{
		Name:    "token-lifetime",
		Usage:   "Lifetime of authorization tokens",
		Value:   24 * time.Hour,
		EnvVars: []string{"TOKEN_LIFETIME"},
	}

	FlagClientId = &cli.StringFlag{
		Name:     "oauth-client-id",
		Usage:    "Client ID for OAuth 2.0 authentication",
		Required: true,
		EnvVars:  []string{"OAUTH_CLIENT_ID"},
	}

	FlagClientSecret = &cli.StringFlag{
		Name:     "oauth-client-secret",
		Usage:    "Client secret for OAuth 2.0 authentication",
		Required: true,
		EnvVars:  []string{"OAUTH_CLIENT_SECRET"},
	}

	FlagRedirectUrl = &cli.StringFlag{
		Name:    "oauth-redirect-url",
		Usage:   "Redirect URL for OAuth 2.0 authentication",
		Value:   "http://localhost:5173/auth/callback",
		EnvVars: []string{"OAUTH_REDIRECT_URL"},
	}
)

func RunApp(ctx *cli.Context) error {
	log.Println("starting up!")

	log.Println("connecting to MySQL @", ctx.String(FlagDatabaseHost.Name))
	db, err := sqlx.Connect("mysql", makeMySqlConfig(ctx).FormatDSN())
	if err != nil {
		return fmt.Errorf("error connecting to MySQL: %s", err)
	}
	log.Println("connected to MySQL")

	log.Println("initializing repositories")
	var (
		curriculumRepo  = curriculumRepository.NewMySqlRepository(db)
		disciplineRepo  = disciplineRepository.NewMySqlRepository(db)
		institutionRepo = institutionRepository.NewMySqlRepository(db)
		taskRepo        = taskRepository.NewMySqlRepository(db)
		userRepo        = userRepository.NewMySqlRepository(db)
	)

	log.Println("initializing auth manager")
	authManager, err := initAuthManager(ctx, userRepo)
	if err != nil {
		return err
	}

	authPlatform, err := initAuthPlatform(ctx)
	if err != nil {
		return err
	}

	log.Println("initializing controllers")
	s := &Server[platform.AuthenticationConfig, platform.CodeRequest]{
		curriculumController:  curriculumDelivery.NewCurriculumController(curriculumRepo),
		disciplineController:  disciplineDelivery.NewDisciplineController(disciplineRepo, userRepo),
		institutionController: institutionDelivery.NewInstitutionController(institutionRepo, curriculumRepo),
		taskController:        taskDelivery.NewTaskController(taskRepo, userRepo),
		userController:        userDelivery.NewUserController(userRepo),
		authController:        authDelivery.NewAuthController(userRepo, authPlatform, authManager),

		authManager:    authManager,
		allowedOrigins: allowedOrigins(ctx),
	}

	log.Println("completed bootstrap process")

	bindAddress := ctx.String(FlagBindAddress.Name)
	log.Println("starting HTTP listener on", bindAddress)
	app := s.MakeApp()
	return app.Listen(bindAddress)
}

func makeMySqlConfig(ctx *cli.Context) *mysql.Config {
	config := mysql.NewConfig()

	config.Net = "tcp"
	config.Addr = ctx.String(FlagDatabaseHost.Name)
	config.User = ctx.String(FlagDatabaseUser.Name)
	config.Passwd = ctx.String(FlagDatabasePassword.Name)
	config.DBName = ctx.String(FlagDatabaseName.Name)
	config.ParseTime = true

	return config
}

func allowedOrigins(ctx *cli.Context) map[string]bool {
	origins := ctx.StringSlice(FlagAllowedOrigins.Name)
	result := make(map[string]bool, len(origins))

	for _, origin := range origins {
		result[origin] = true
	}

	return result
}

func initAuthManager(ctx *cli.Context, userRepo user.Repository) (auth.Manager, error) {
	signingKey, err := base64.StdEncoding.DecodeString(ctx.String(FlagSigningKey.Name))
	if err != nil {
		return nil, fmt.Errorf("invalid signing key: %s", err)
	}

	tokenProvider := token.NewJwtTokenProvider(
		signingKey,
		ctx.String(FlagAudience.Name),
		ctx.Duration(FlagTokenLifetime.Name),
	)

	return manager.NewAuthManager(userRepo, tokenProvider), nil
}

func initAuthPlatform(ctx *cli.Context) (auth.Platform[platform.AuthenticationConfig, platform.CodeRequest], error) {
	// Hardcoded config for Discord OAuth
	cfg := &oauth2.Config{
		ClientID:     ctx.String(FlagClientId.Name),
		ClientSecret: ctx.String(FlagClientSecret.Name),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		RedirectURL: ctx.String(FlagRedirectUrl.Name),
		Scopes:      []string{"identify"},
	}

	return platform.NewOAuthPlatform(cfg, platform.NewDiscordUserSupplier()), nil
}
