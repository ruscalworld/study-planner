package server

import (
	"fmt"
	"log"

	curriculumRepository "study-planner/internal/curriculum/repository"
	disciplineRepository "study-planner/internal/discipline/repository"
	institutionRepository "study-planner/internal/institution/repository"
	taskRepository "study-planner/internal/task/repository"
	userRepository "study-planner/internal/user/repository"

	curriculumDelivery "study-planner/internal/curriculum/delivery"
	disciplineDelivery "study-planner/internal/discipline/delivery"
	institutionDelivery "study-planner/internal/institution/delivery"
	taskDelivery "study-planner/internal/task/delivery"
	userDelivery "study-planner/internal/user/delivery"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"

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

	log.Println("initializing controllers")
	s := &Server{
		curriculumController:  curriculumDelivery.NewCurriculumController(curriculumRepo),
		disciplineController:  disciplineDelivery.NewDisciplineController(disciplineRepo),
		institutionController: institutionDelivery.NewInstitutionController(institutionRepo, curriculumRepo),
		taskController:        taskDelivery.NewTaskController(taskRepo),
		userController:        userDelivery.NewUserController(userRepo),

		allowedOrigins: allowedOrigins(ctx),
	}

	log.Println("completed bootstrap process")

	bindAddress := ctx.String(FlagBindAddress.Name)
	log.Println("starting HTTP listener on", bindAddress)
	app := s.MakeApp()
	return app.Listen(bindAddress)
}

func makeMySqlConfig(ctx *cli.Context) *mysql.Config {
	return &mysql.Config{
		Addr:   ctx.String(FlagDatabaseHost.Name),
		User:   ctx.String(FlagDatabaseUser.Name),
		Passwd: ctx.String(FlagDatabasePassword.Name),
		DBName: ctx.String(FlagDatabaseName.Name),

		ParseTime: true,
	}
}

func allowedOrigins(ctx *cli.Context) map[string]bool {
	origins := ctx.StringSlice(FlagAllowedOrigins.Name)
	result := make(map[string]bool, len(origins))

	for _, origin := range origins {
		result[origin] = true
	}

	return result
}
