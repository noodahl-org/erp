package server

import (
	"fmt"
	"net/http"

	"github.com/adhocore/gronx"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/noodahl-org/erp/api/clients/postgres"
	api "github.com/noodahl-org/erp/api/conf"
)

type Msg struct {
	Msg   string `json:"msg"`
	Error string `json:"err"`
}

func serverMsg(msg string, err error) Msg {
	return Msg{
		Msg:   msg,
		Error: err.Error(),
	}
}

type Server struct {
	conf  *api.ApiConf
	e     *echo.Echo
	db    postgres.DBClient
	valid *validator.Validate
	cron  *gronx.Gronx
}

func NewServer(opts ...func(*Server)) *Server {
	echoServer := echo.New()
	s := &Server{
		e:    echoServer,
		cron: gronx.New(),
	}
	for _, opt := range opts {
		opt(s)
	}

	//routes
	//equipment
	s.e.GET("/liveness", s.Liveliness)
	s.e.GET("/equipment", s.FetchEquipment)
	s.e.GET("/equipment/:id", s.FetchEquipmentByID)
	s.e.POST("/equipment", s.UpsertEquipment)
	s.e.DELETE("/equipment", s.DeleteEquipment)

	//maintenance
	s.e.POST("/maintenance", s.UpsertMaintenanceTask)
	s.e.GET("/maintenance", s.FetchMaintenanceTasks)

	//user scheduled tasks
	s.e.POST("/scheduledtasks", s.UpsertUserMaintenanceTask)
	s.e.GET("/scheduledtasks", s.FetchUserMaintenanceTasks)

	//users
	s.e.GET("/users", s.FetchUser)
	s.e.POST("/users", s.UpsertUser)
	s.e.DELETE("/users", s.DeleteUser)

	//user equipment
	s.e.GET("/equipment/user/:id", s.FetchUserEquipment)
	s.e.POST("/equipment/user", s.UpsertUserEquipment)
	s.e.DELETE("/equipment/user/:id", s.DeleteUserEquipment)

	//user maintenance task
	s.e.GET("/maintenance/user/:id", s.FetchUserMaintenanceTasks)
	s.e.POST("/maintenance/user", s.UpsertUserMaintenanceTask)

	return s
}

func WithMiddleware(f echo.MiddlewareFunc) func(*Server) {
	return func(s *Server) {
		s.e.Use(f)
	}
}

func WithConfig(conf *api.ApiConf) func(*Server) {
	return func(s *Server) {
		s.conf = conf
	}
}

func WithDB(db postgres.DBClient) func(*Server) {
	return func(s *Server) {
		s.db = db
	}
}

func WithValidator(v *validator.Validate) func(*Server) {
	return func(s *Server) {
		s.valid = v
	}
}

func (s *Server) Start() {
	s.e.Logger.Fatal(s.e.Start(fmt.Sprintf(":%v", s.conf.Port)))
}

func (s *Server) MigrateDomainModel(model *interface{}) error {
	return s.db.AutoMigrate(model)
}

func (s *Server) Health(e echo.Context) error {
	return e.JSON(http.StatusOK, nil)
}

func (s *Server) Liveliness(e echo.Context) error {
	return e.JSON(http.StatusOK, nil)
}
