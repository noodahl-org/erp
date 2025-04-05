package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/noodahl-org/erp/internal/models"
)

func (s *Server) FetchUser(e echo.Context) error {
	ctx := e.Request().Context()
	var user models.User

	if err := e.Bind(&user); err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to parse user request", err))
	}
	if err := s.db.Fetch(ctx, user, &user); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to fetch user", err))
	}
	return e.JSON(http.StatusOK, user)
}

func (s *Server) UpsertUser(e echo.Context) error {
	ctx := e.Request().Context()
	var req models.User
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to parse user request", err))
	}
	if err := s.db.Upsert(ctx, &req); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to upsert user", err))
	}

	return e.JSON(http.StatusAccepted, req)
}

func (s *Server) DeleteUser(e echo.Context) error {
	ctx := e.Request().Context()
	id, err := uuid.Parse(e.QueryParam("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to parse id", err))
	}

	if err := s.db.Delete(ctx, &models.User{ID: id.String()}); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to delete user", err))
	}
	return e.JSON(http.StatusResetContent, nil)
}
