package echo

import (
	"net/http"

	"github.com/labstack/echo"
)

func RenderBadRequest(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{"info": err.Error()})
}

func RenderOK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"info": "ok", "data": data})
}
