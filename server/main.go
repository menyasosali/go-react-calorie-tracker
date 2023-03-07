package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/menyasosali/go-react-calorie-tracker/routes"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	router.POST("/entry/create", routes.AddEntry)
	router.GET("/entries", routes.GetEntries)
	router.GET("/entry/:id/", routes.GetEntryById)
	router.GET("/ingredient/:ingredient", routes.GetEntriesByIngredient)

	router.PUT("/entry/update/:id", routes.UpdateEntry)
	router.PUT("/ingredient/update/:id", routes.UpdateIngredient)
	router.DELETE("/entry/delete/:id", routes.DeleteEntry)

	router.Logger.Fatal(router.Start(":" + port))
}
