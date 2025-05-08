package server

import (
	"net/http"
	"time"

	"github.com/adhocore/gronx"
	"github.com/labstack/echo"
	"github.com/noodahl-org/erp/api/models"
)

func (s *Server) FetchMaintenanceTasks(e echo.Context) error {
	ctx := e.Request().Context()

	query := models.MaintenanceTask{}
	if err := e.Bind(&query); err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to bind requests", err))
	}
	results := []models.MaintenanceTask{}
	if err := s.db.FetchMany(ctx, query, &results); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to fetch maintenance tasks", err))
	}
	return e.JSON(http.StatusAccepted, results)
}

func (s *Server) UpsertMaintenanceTask(e echo.Context) error {
	var req models.MaintenanceTask
	ctx := e.Request().Context()
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to bind request", err))
	}

	if err := s.db.Upsert(ctx, &req); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable save maintenance task", err))
	}

	return e.JSON(http.StatusAccepted, req)
}

func (s *Server) UpsertUserMaintenanceTask(e echo.Context) error {
	ctx := e.Request().Context()

	req := models.UserMaintenanceTask{}
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to bind request", err))
	}

	if err := s.db.Fetch(ctx, models.MaintenanceTask{ID: req.MaintenanceTaskID}, &req.MaintenanceTask); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to fetch maintenance task", err))
	}

	startTime, err := gronx.NextTickAfter(req.MaintenanceTask.Cron, time.Now(), true)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to parse cron", err))
	}
	req.StartTime = &startTime

	if err := s.db.Create(ctx, &req); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable save user scheduled task", err))
	}

	return e.JSON(http.StatusAccepted, req)
}

func (s *Server) FetchUserMaintenanceTasks(e echo.Context) error {
	ctx := e.Request().Context()

	query := models.UserMaintenanceTask{}
	if err := e.Bind(&query); err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to bind requests", err))
	}
	results := []models.UserMaintenanceTask{}
	if err := s.db.Fetch(ctx, query, &results); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to fetch user scheduled tasks", err))
	}

	return e.JSON(http.StatusAccepted, results)
}
