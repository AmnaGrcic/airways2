package api

import (
	"airways/common"
	"airways/repository"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Flights struct {
	FlightId       int       `json:"flight_id"`
	FlightDateTime time.Time `json:"flight_date_time"`
	FlightPrice    int       `json:"flight_price"`
	AeroplaneId    int       `json:"aeroplane_id"`
	StartAirPortId int       `json:"start_air_port_id"`
	EndAirPortId   int       `json:"end_air_port_id"`
	AverageRating  float64   `json:"average_rating"`
}

type FlightsFilters struct {
	StartPrice     string
	EndPrice       string
	AeroplaneId    string
	StartAirportId string
	StartCity      string
	EndCity        string
	StartDate      string
	EndDate        string
}

func GetFlightsFilters(flightsfilter FlightsFilters) (string, []interface{}, string, []interface{}) {
	var filtersQry string
	var filterParams []interface{}
	var joinQry string
	var joinParams []interface{}

	if flightsfilter.StartPrice != "" && flightsfilter.EndPrice != "" {
		filtersQry += ` AND flight_price >= ? AND flight_price <= ?`
		filterParams = append(filterParams, flightsfilter.StartPrice, flightsfilter.EndPrice)
	}

	if flightsfilter.AeroplaneId != "" {
		filtersQry += ` AND aeroplane_id=?`
		filterParams = append(filterParams, flightsfilter.AeroplaneId)
	}

	if flightsfilter.StartAirportId != "" {
		filtersQry += ` AND start_air_port_id=?`
		filterParams = append(filterParams, flightsfilter.StartAirportId)
	}

	if flightsfilter.StartCity != "" {

		if !strings.Contains(joinQry, "JOIN airports") {
			joinQry += ` JOIN airports a ON f.start_air_port_id=a.airport_id`
		}

		if !strings.Contains(joinQry, "JOIN cities") {
			joinQry += ` JOIN cities c ON c.city_id=a.city_id`
		}

		filtersQry += ` AND c.city_name=?`
		filterParams = append(filterParams, flightsfilter.StartCity)
	}

	if flightsfilter.EndCity != "" {

		if !strings.Contains(joinQry, "JOIN airports") {
			joinQry += ` JOIN airports a ON f.end_air_port_id=a.airport_id`
		}

		if !strings.Contains(joinQry, "JOIN cities") {
			joinQry += ` JOIN cities c ON c.city_id=a.city_id`
		}

		filtersQry += ` AND c.city_name=?`
		filterParams = append(filterParams, flightsfilter.EndCity)
	}

	if flightsfilter.StartDate != "" {
		filtersQry += ` AND f.flight_date_time=?`
		filterParams = append(filterParams, flightsfilter.StartDate)

	}

	if flightsfilter.EndDate != "" {
		filtersQry += ` AND f.flight_date_time=?`
		filterParams = append(filterParams, flightsfilter.EndDate)
	}

	return filtersQry, filterParams, joinQry, joinParams
}



func GetAllFlights(c *gin.Context) {
	flights := []Flights{}
	flightsfilter := FlightsFilters{
		StartPrice:     c.Query("start_price"),
		EndPrice:       c.Query("end_price"),
		AeroplaneId:    c.Query("aeroplane_id"),
		StartAirportId: c.Query("start_air_port_id"),
		StartCity:      c.Query("start_city"),
		EndCity:        c.Query("end_city"),
		StartDate:      c.Query("start_date"),
		EndDate:        c.Query("end_date"),
	}

	filtersQry, filterParams, joinQry, joinParams := GetFlightsFilters(flightsfilter)


	qry := `SELECT f.*, (SUM(r.rating)/COUNT(r.rating)) AS average_rating 
	FROM flights f
	LEFT JOIN ratings r ON f.flight_id=r.flight_id
	            %s 
	WHERE 1=1 %s
	GROUP BY flight_id`

	qry = fmt.Sprintf(qry, joinQry, filtersQry)
	var params []interface{}
	params = append(params, joinParams...)
	params = append(params, filterParams...)

	err := repository.Db.Debug().Raw(qry, params...).Scan(&flights).Error
	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}

	fmt.Printf("flights: %+v \n", flights)
	c.JSON(200, flights)
}

func GetFlight(c *gin.Context) {
	flight_id := c.Param("flight_id")
	flight := Flights{}
	qry := `SELECT f.*,
	      (SUM(r.rating)/COUNT(r.rating)) AS average_rating
	      FROM flights f
		  LEFT JOIN ratings r ON f.flight_id=r.flight_id
		  WHERE f.flight_id=?
		  GROUP BY flight_id`
	err := repository.Db.Debug().Raw(qry, flight_id).Scan(&flight).Error
	if err != nil {
		println(err)
		return
	}
	fmt.Printf("flight: %+v \n", flight)
	c.JSON(200, flight)
}

func DeleteFlights(c *gin.Context) {
	flight_id := c.Param("flight_id")

	qry := "DELETE FROM flights WHERE flight_id=?"
	err := repository.Db.Debug().Exec(qry, flight_id).Error
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}
	c.JSON(200, "Ok")
}

func CreateFlight(c *gin.Context) {
	flight := Flights{}

	err := c.BindJSON(&flight)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}
	qry := `INSERT INTO flights (flight_date_time,
		flight_price,
		aeroplane_id,
		start_air_port_id,
		end_air_port_id) 
	VALUES (?,?,?,?,?)`
	err = repository.Db.Debug().Exec(qry,
		flight.FlightDateTime,
		flight.FlightPrice,
		flight.AeroplaneId,
		flight.StartAirPortId,
		flight.EndAirPortId).Error
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}
	c.JSON(200, "Ok")

}

func UpdateFlight(c *gin.Context) {
	flight_id := c.Param("flight_id")

	flight := Flights{}

	err := c.BindJSON(&flight)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
	}
	qry := `UPDATE flights SET
	      flight_date_time=?,
		  flight_price=?,
		  aeroplane_id=?,
		  start_air_port_id=?,
		  end_air_port_id=?
		  WHERE flight_id=?`

	err = repository.Db.Debug().Exec(qry,
		flight.FlightDateTime,
		flight.FlightPrice,
		flight.AeroplaneId,
		flight.StartAirPortId,
		flight.EndAirPortId,
		flight_id).Error
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}
	c.JSON(200, "Ok")
}

func FlightCheckIn(c *gin.Context) {
	flight_id := c.Param("flight_id")
	user_id, err := GetUserIdFromUserClaims(c)

	if err != nil {
		println(err)
		return
	}

	qry := `INSERT INTO users_flights (user_id,flight_id)
	VALUES (?,?)`
	err = repository.Db.Debug().Exec(qry,
		user_id,
		flight_id).Error

	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}
	c.JSON(200, "Ok")

}
