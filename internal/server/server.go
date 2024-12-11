package server

import (
	"github.com/ruscalworld/study-planner/internal/auth"
	"github.com/ruscalworld/study-planner/internal/auth/delivery"

	"github.com/ruscalworld/study-planner/internal/curriculum"
	"github.com/ruscalworld/study-planner/internal/discipline"
	"github.com/ruscalworld/study-planner/internal/institution"
	"github.com/ruscalworld/study-planner/internal/task"
	"github.com/ruscalworld/study-planner/internal/user"

	"github.com/ruscalworld/study-planner/pkg/httputil"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/fiber/v2"
)

type Server[AC, AT comparable] struct {
	curriculumController  curriculum.Controller
	disciplineController  discipline.Controller
	institutionController institution.Controller
	taskController        task.Controller
	userController        user.Controller
	authController        auth.Controller[AC, AT]

	authManager    auth.Manager
	allowedOrigins map[string]bool
}

func (s *Server[AC, AT]) MakeApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: httputil.ErrorHandler,
	})

	authMiddleware := delivery.NewMiddleware(s.authManager)

	app.Use(
		healthcheck.New(),
		logger.New(),
		recover.New(),

		cors.New(cors.Config{
			AllowOriginsFunc: s.isAllowedOrigin,
		}),
	)

	app.Route("/v1", func(r fiber.Router) {
		r.Route("/auth", func(r fiber.Router) {
			r.Get("/config", httputil.MakeSimpleHandler(s.authController.GetConfig))
			r.Post("/sign-in", httputil.MakeHandler(s.authController.Authenticate))

			r.Route("/refresh", func(r fiber.Router) {
				r.Use(authMiddleware)
				r.Post("/", httputil.MakeSimpleHandler(s.authController.Refresh))
			})
		})

		r.Route("/profile", func(r fiber.Router) {
			r.Use(authMiddleware)
			r.Get("/", httputil.MakeSimpleHandler(s.authController.GetCurrentUser))

			r.Route("/curriculums", func(r fiber.Router) {
				r.Get("/", httputil.MakeSimpleHandler(s.userController.GetUserCurriculums))
				r.Post("/", httputil.MakeHandler(s.userController.CreateUserCurriculum))
				r.Delete("/:curriculum_id", httputil.MakeSimpleHandler(s.userController.DeleteUserCurriculum))
			})
		})

		r.Route("/curriculums", func(r fiber.Router) {
			r.Use(authMiddleware)
			r.Post("/", httputil.MakeHandler(s.curriculumController.CreateCurriculum))

			r.Route("/:curriculum_id", func(r fiber.Router) {
				r.Get("/", httputil.MakeSimpleHandler(s.curriculumController.GetCurriculum))
				r.Put("/", httputil.MakeHandler(s.curriculumController.UpdateCurriculum))
				r.Delete("/", httputil.MakeSimpleHandler(s.curriculumController.DeleteCurriculum))

				r.Route("/codes", func(r fiber.Router) {
					r.Get("/", httputil.MakeSimpleHandler(s.curriculumController.GetCurriculumCodes))
					r.Post("/", httputil.MakeHandler(s.curriculumController.CreateCurriculumCode))
					r.Delete("/:code_id", httputil.MakeSimpleHandler(s.curriculumController.DeleteCurriculumCode))
				})

				r.Route("/disciplines", func(r fiber.Router) {
					r.Get("/", httputil.MakeSimpleHandler(s.disciplineController.GetDisciplines))
					r.Post("/", httputil.MakeHandler(s.disciplineController.CreateDiscipline))

					r.Get("/:discipline_id", httputil.MakeSimpleHandler(s.disciplineController.GetCurriculumDiscipline))
					r.Put("/:discipline_id", httputil.MakeHandler(s.disciplineController.UpdateDiscipline))
					r.Delete("/:discipline_id", httputil.MakeSimpleHandler(s.disciplineController.DeleteDiscipline))
				})

				r.Route("/users", func(r fiber.Router) {
					r.Get("/", httputil.MakeSimpleHandler(s.curriculumController.GetCurriculumUsers))
					r.Delete("/:user_id", httputil.MakeSimpleHandler(s.curriculumController.DeleteCurriculumUser))
				})
			})
		})

		r.Route("/disciplines/:discipline_id", func(r fiber.Router) {
			r.Use(authMiddleware)

			r.Route("/links", func(r fiber.Router) {
				r.Get("/", httputil.MakeSimpleHandler(s.disciplineController.GetDisciplineLinks))
				r.Post("/", httputil.MakeHandler(s.disciplineController.CreateDisciplineLink))
				r.Put("/:discipline_link_id", httputil.MakeHandler(s.disciplineController.UpdateDisciplineLink))
				r.Delete("/:discipline_link_id", httputil.MakeSimpleHandler(s.disciplineController.DeleteDisciplineLink))
			})

			r.Route("/progress", func(r fiber.Router) {
				r.Get("/", httputil.MakeSimpleHandler(s.disciplineController.GetDisciplineProgress))
			})

			r.Route("/stats", func(r fiber.Router) {
				r.Get("/", httputil.MakeSimpleHandler(s.disciplineController.GetDisciplineStats))
			})

			r.Route("/groups", func(r fiber.Router) {
				r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTaskGroups))
				r.Post("/", httputil.MakeHandler(s.taskController.CreateTaskGroup))

				r.Route("/:group_id", func(r fiber.Router) {
					r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTaskGroup))
					r.Put("/", httputil.MakeHandler(s.taskController.UpdateTaskGroup))
					r.Delete("/", httputil.MakeSimpleHandler(s.taskController.DeleteTaskGroup))

					r.Route("/goal", func(r fiber.Router) {
						r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTaskGroupGoal))
						r.Put("/", httputil.MakeHandler(s.taskController.UpdateTaskGroupGoal))
					})

					r.Route("/tasks", func(r fiber.Router) {
						r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTasks))
						r.Post("/", httputil.MakeHandler(s.taskController.CreateTask))

						r.Route("/:task_id", func(r fiber.Router) {
							r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTask))
							r.Put("/", httputil.MakeHandler(s.taskController.UpdateTask))
							r.Delete("/", httputil.MakeSimpleHandler(s.taskController.DeleteTask))

							r.Route("/progress", func(r fiber.Router) {
								r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTaskProgress))
								r.Put("/", httputil.MakeHandler(s.taskController.UpdateTaskProgress))
							})

							r.Route("/links", func(r fiber.Router) {
								r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTaskLinks))
								r.Post("/", httputil.MakeHandler(s.taskController.CreateTaskLink))

								r.Route("/:task_link_id", func(r fiber.Router) {
									r.Put("/", httputil.MakeHandler(s.taskController.UpdateTaskLink))
									r.Delete("/", httputil.MakeSimpleHandler(s.taskController.DeleteTaskLink))
								})
							})
						})
					})
				})
			})

			r.Route("/tasks", func(r fiber.Router) {
				r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTasks))

				r.Route("/:task_id", func(r fiber.Router) {
					r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTask))

					r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTask))
					r.Put("/", httputil.MakeHandler(s.taskController.UpdateTask))
					r.Delete("/", httputil.MakeSimpleHandler(s.taskController.DeleteTask))

					r.Route("/progress", func(r fiber.Router) {
						r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTaskProgress))
						r.Put("/", httputil.MakeHandler(s.taskController.UpdateTaskProgress))
					})

					r.Route("/links", func(r fiber.Router) {
						r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTaskLinks))
						r.Post("/", httputil.MakeHandler(s.taskController.CreateTaskLink))

						r.Route("/:task_link_id", func(r fiber.Router) {
							r.Put("/", httputil.MakeHandler(s.taskController.UpdateTaskLink))
							r.Delete("/", httputil.MakeSimpleHandler(s.taskController.DeleteTaskLink))
						})
					})
				})
			})
		})

		r.Route("/institutions", func(r fiber.Router) {
			r.Get("/", httputil.MakeSimpleHandler(s.institutionController.GetInstitutions))

			r.Route("/:institution_id", func(r fiber.Router) {
				r.Get("/", httputil.MakeSimpleHandler(s.institutionController.GetInstitution))

				r.Route("/curriculums", func(r fiber.Router) {
					r.Use(authMiddleware)
					r.Post("/", httputil.MakeHandler(s.institutionController.CreateCurriculum))
				})
			})
		})
	})

	return app
}

func (s *Server[AC, AT]) isAllowedOrigin(origin string) bool {
	_, ok := s.allowedOrigins[origin]
	return ok
}
