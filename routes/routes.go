package routes

import (
	"airways/api"

	"github.com/gin-gonic/gin"
)

func Startroutes() {
	r := gin.Default()

	aeroplanes := r.Group("/aeroplanes", api.AuthorizationCheck)
	{
		aeroplanes.GET("", api.GetAllAeroplanes)
		aeroplanes.GET("/:aeroplane_id", api.GetAeroplane)
		aeroplanes.DELETE("/:aeroplane_id", api.DeleteAeroplane)
		aeroplanes.POST("", api.PermissionMiddleware("create_aeroplane"), api.CreateAeroplane)
		aeroplanes.PUT("/:aeroplane_id", api.PermissionMiddleware("update_aeroplane"), api.UpdateAeroplane)
	}

	airports := r.Group("/airports", api.AuthorizationCheck)
	{
		airports.GET("", api.GetAllAirports)
		airports.GET("/:airport_id", api.GetAirport)
		airports.DELETE("/:airport_id", api.DeleteAirport)
		airports.POST("", api.PermissionMiddleware("create_airport"), api.CreateAirport)
		airports.PUT("/:airport_id", api.PermissionMiddleware("update_airport"), api.UpdateAirports)
	}

	flights := r.Group("/flights")
	{
		flights.GET("", api.GetAllFlights)
		flights.GET("/:flight_id", api.GetFlight)
		flights.DELETE("/:flight_id", api.DeleteFlights)
		flights.POST("", api.CreateFlight)
		flights.PUT("/:flight_id", api.AuthorizationCheck, api.UpdateFlight)
		flights.POST("/check-in/:flight_id", api.AuthorizationCheck, api.FlightCheckIn)
	}
	users := r.Group("/users")
	{
		users.GET("", api.AuthorizationCheck, api.GetAllUsers)
		users.GET("/:user_id", api.GetUser)
		users.DELETE("", api.AuthorizationCheck, api.DeleteUser)
		users.POST("/register", api.RegisterUser)
		users.POST("/login", api.Login)
		users.GET("/test", api.AuthorizationCheck)
	}
	usersflights := r.Group("/usersflights")
	{
		usersflights.GET("", api.AuthorizationCheck, api.GetUserFlights)
		usersflights.POST("/:flight_id", api.AuthorizationCheck, api.CreateReview)
	}

	cities := r.Group("/cities")
	{
		cities.GET("", api.GetAllCities)
		cities.GET("/:city_id", api.GetCity)
		cities.DELETE("/:city_id", api.DeleteCity)
		cities.POST("", api.CreateCity)
		cities.PUT("/:city_id", api.UpdateCity)
	}
	r.Run()

}
