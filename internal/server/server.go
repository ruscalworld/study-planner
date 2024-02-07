package server

import (
	"study-planner/internal/curriculum"
	"study-planner/internal/discipline"
	"study-planner/internal/institution"
	"study-planner/internal/task"
	"study-planner/internal/user"

	"study-planner/pkg/httputil"

	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	curriculumController  curriculum.Controller
	disciplineController  discipline.Controller
	institutionController institution.Controller
	taskController        task.Controller
	userController        user.Controller
}

func (s *Server) MakeApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: httputil.ErrorHandler,
	})

	app.Use(
		healthcheck.New(),
		logger.New(),
		recover.New(),
	)

	app.Route("/v1", func(r fiber.Router) {
		r.Route("/curriculums", func(r fiber.Router) {
			r.Get("/", httputil.MakeSimpleHandler(s.curriculumController.GetCurriculums))

			r.Route("/:curriculum_id", func(r fiber.Router) {
				r.Get("/", httputil.MakeSimpleHandler(s.curriculumController.GetCurriculum))

				r.Route("/disciplines", func(r fiber.Router) {
					r.Get("/", httputil.MakeSimpleHandler(s.disciplineController.GetDisciplines))
					r.Get("/:discipline_id", httputil.MakeSimpleHandler(s.disciplineController.GetDiscipline))
				})
			})
		})

		r.Route("/disciplines/:discipline_id", func(r fiber.Router) {
			r.Get("/links", httputil.MakeSimpleHandler(s.disciplineController.GetDisciplineLinks))

			r.Route("/groups", func(r fiber.Router) {
				r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTaskGroups))
				r.Get("/:group_id", httputil.MakeSimpleHandler(s.taskController.GetTaskGroup))
			})

			r.Route("/tasks", func(r fiber.Router) {
				r.Get("/", httputil.MakeSimpleHandler(s.taskController.GetTasks))
				r.Get("/:task_id", httputil.MakeSimpleHandler(s.taskController.GetTask))
				r.Get("/:task_id/links", httputil.MakeSimpleHandler(s.taskController.GetTaskLinks))
			})
		})

		r.Route("/institutions", func(r fiber.Router) {
			r.Get("/", httputil.MakeSimpleHandler(s.institutionController.GetInstitutions))

			r.Route("/:institution_id", func(r fiber.Router) {
				r.Get("/", httputil.MakeSimpleHandler(s.institutionController.GetInstitution))
				r.Get("/curriculums", httputil.MakeSimpleHandler(s.institutionController.GetCurriculums))
			})
		})
	})

	return app
}