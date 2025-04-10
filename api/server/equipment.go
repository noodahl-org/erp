package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/noodahl-org/erp/api/models"
)

func (s *Server) FetchEquipmentByID(e echo.Context) error {
	ctx := e.Request().Context()
	parsed, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to parse id", err))
	}
	var equipment models.Equipment
	if err := s.db.Fetch(ctx, &models.Equipment{ID: parsed.String()}, &equipment); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to fetch equipment", err))
	}
	return e.JSON(http.StatusOK, equipment)

}

func (s *Server) FetchEquipment(e echo.Context) error {
	var req models.Equipment
	ctx := e.Request().Context()

	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to bind request", err))
	}
	results := []models.Equipment{}
	if err := s.db.FetchMany(ctx, &req, &results); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to fetch equipment", err))
	}

	return e.JSON(http.StatusOK, results)
}

func (s *Server) UpsertEquipment(e echo.Context) error {
	var req models.Equipment
	ctx := e.Request().Context()
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to bind request", err))
	}

	if err := s.db.Upsert(ctx, &req); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable save equipment", err))
	}

	return e.JSON(http.StatusAccepted, req)
}

func (s *Server) DeleteEquipment(e echo.Context) error {
	ctx := e.Request().Context()
	id, err := uuid.Parse(e.QueryParam("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to parse id", err))
	}

	if err := s.db.Delete(ctx, &models.Equipment{
		ID: id.String(),
	}); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to delete equipment", err))
	}

	return e.JSON(http.StatusResetContent, nil)
}

func (s *Server) FetchUserEquipment(e echo.Context) error {
	ctx := e.Request().Context()
	req := models.UserEquipment{}
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to bind request", err))
	}

	results := []models.UserEquipment{}
	if err := s.db.FetchMany(ctx, req, &results, []string{"Equipment"}...); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to fetch user equipment", err))
	}

	return e.JSON(http.StatusOK, results)

}

func (s *Server) UpsertUserEquipment(e echo.Context) error {
	ctx := e.Request().Context()
	req := models.UserEquipment{}
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to bind request", err))
	}

	if err := s.db.Upsert(ctx, &req, []string{"Equipment"}...); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable save user equipment", err))
	}

	if err := s.db.Fetch(ctx, models.Equipment{ID: req.EquipmentID}, &req.Equipment); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable fetch equipment", err))
	}

	return e.JSON(http.StatusAccepted, req)
}

func (s *Server) DeleteUserEquipment(e echo.Context) error {
	ctx := e.Request().Context()
	parsed, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to parse id", err))
	}
	if s.db.Delete(ctx, models.UserEquipment{ID: parsed.String()}); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to delete user equipment", err))
	}
	return e.JSON(http.StatusResetContent, nil)
}
