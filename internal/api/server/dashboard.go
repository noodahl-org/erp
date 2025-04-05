package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/noodahl-org/erp/internal/models"
)

func (s *Server) Dashboard(e echo.Context) error {
	ctx := e.Request().Context()
	id, err := uuid.Parse(e.QueryParam("user_id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, serverMsg("unable to parse user id", err))
	}

	var dashboard models.Dashboard
	if err := s.db.Fetch(ctx, models.User{ID: id.String()}, &dashboard.User); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to fetch user", err))
	}

	if err := s.db.FetchMany(ctx, models.UserEquipment{UserID: dashboard.User.ID}, &dashboard.UserEquipment); err != nil {
		return e.JSON(http.StatusInternalServerError, serverMsg("unable to fetch user equipment", err))
	}

	return e.JSON(http.StatusOK, dashboard)
}
