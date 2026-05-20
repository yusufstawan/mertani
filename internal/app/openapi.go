package app

import (
	"github.com/labstack/echo/v5"

	"mertani/internal/shared/response"
)

const swaggerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Mertani API Docs</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: "/openapi.yaml",
        dom_id: "#swagger-ui"
      });
    };
  </script>
</body>
</html>`

func (s *Server) registerOpenAPI() {
	s.echo.GET("/openapi.yaml", openAPI)
	s.echo.GET("/swagger", swagger)
	s.echo.GET("/swagger/", swagger)
}

func openAPI(c *echo.Context) error {
	return c.File("docs/openapi.yaml")
}

func swagger(c *echo.Context) error {
	return c.HTML(response.StatusOK, swaggerHTML)
}
