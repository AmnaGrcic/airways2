package api

import (
	"airways/common"
	"airways/repository"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Airport struct {
	AirportId        int    `json:"airport_id"`
	CityId           int    `json:"city_id"`
	AirportShortName string `json:"airport_short_name"`
}

type AirportPayload struct {
	AirportId        int                `json:"airport_id"`
	Description      string             `json:"description"`
	CityId           int                `json:"city_id"`
	CityName         string             `json:"city_name"`
	AirportShortName string             `json:"airport_short_name"`
	StartFlightsRaw  json.RawMessage    `json:"-"`
	EndFlightsRaw    json.RawMessage    `json:"-"`
	StartFlights     []FlightsUnmarshal `json:"startflights,omitempty" gorm:"-"`
	EndFlights       []FlightsUnmarshal `json:"endflights,omitempty" gorm:"-"`
}

type AirportFilters struct {
	AirportCity string
	Limit       string
	Offset      string
}

func GetAirportFilters(af AirportFilters) (string, []interface{}, string, []interface{}) {
	var filtersQry string
	var filterParams []interface{}
	var joinQry string
	var joinParams []interface{}

	if af.AirportCity != "" {
		joinQry += ` JOIN cities c ON a.city_id=c.city_id`
		filtersQry += ` AND c.city_name=? `
		filterParams = append(filterParams, af.AirportCity)

	}

	if af.Limit != "" && af.Offset != "" {
		filtersQry += ` LIMIT ? OFFSET ?`
		filterParams = append(filterParams, af.Limit, af.Offset)
	}

	return filtersQry, filterParams, joinQry, joinParams
}

func GetAllAirports(c *gin.Context) {
	airports := []AirportPayload{}
	filters := AirportFilters{
		AirportCity: c.Query("airport_city"),
		Limit:       c.DefaultQuery("limit", "5"),
		Offset:      c.DefaultQuery("offset", "0"),
	}

	filtersQry, filterParams, joinQry, joinParams := GetAirportFilters(filters)

	var params []interface{}

	params = append(params, joinParams...)
	params = append(params, filterParams...)

	qry := `SELECT 
			a.*, c.city_name,
			(SELECT json_arrayagg(
				json_object(
						'flight_id', f.flight_id,
						'flight_price',f.flight_price
						)
				)
	 			    FROM flights f 
				WHERE a.airport_id=f.start_air_port_id
			) AS start_flights_raw,
			(SELECT json_arrayagg(
				json_object(
					    'flight_id',f.flight_id,
				        'flight_price',f.flight_price)
					    )
				    FROM flights f 
				WHERE a.airport_id=f.end_air_port_id
			) AS end_flights_raw
			FROM airports a 
			JOIN cities c ON a.city_id=c.city_id
			%s
			WHERE 1=1 %s`

	qry = fmt.Sprintf(qry, joinQry, filtersQry)
	err := repository.Db.Debug().Raw(qry, params...).Scan(&airports).Error

	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}

	for i := range airports {

		if len(airports[i].StartFlightsRaw) > 0 {
			err = json.Unmarshal(airports[i].StartFlightsRaw, &airports[i].StartFlights)
			if err != nil {
				common.LogError(err)
				c.JSON(400, "Db Error")
				return
			}
			fmt.Println(airports[i].StartFlights)
		}

		if len(airports[i].EndFlightsRaw) > 0 {
			err = json.Unmarshal(airports[i].EndFlightsRaw, &airports[i].EndFlights)
			if err != nil {
				common.LogError(err)
				c.JSON(400, "Db Error")
				return
			}
			fmt.Println(airports[i].EndFlights)
		}

	}

	fmt.Printf("airports: %+v \n", airports)
	c.JSON(200, airports)
}

func GetAirport(c *gin.Context) {
	airport_id := c.Param("airport_id")

	qry := `SELECT a.*, c.city_name
	 FROM airports a 
	 JOIN cities c ON a.city_id=c.city_id
	 WHERE a.airport_id=?`
	airport := AirportPayload{}
	err := repository.Db.Debug().Raw(qry, airport_id).Scan(&airport).Error
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}
	fmt.Printf("airport: %+v \n", airport)
	c.JSON(200, airport)
}

func DeleteAirport(c *gin.Context) {
	airport_id := c.Param("airport_id")

	qry := "DELETE FROM airports WHERE airport_id=?"
	err := repository.Db.Debug().Exec(qry, airport_id).Error
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}
	c.JSON(200, "Ok")
}

func CreateAirport(c *gin.Context) {
	airport := AirportPayload{}

	err := c.BindJSON(&airport)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}

	languages, langErr := GetLanguages()
	if langErr != nil {
		println(langErr)
		c.JSON(400, "Db Error")
		return
	}

	tx := repository.Db.Begin()

	qry := `INSERT INTO airports (city_id,airport_short_name)
	 VALUES (?,?)`
	err = tx.Debug().Exec(qry,
		airport.CityId,
		airport.AirportShortName).Error
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		tx.Rollback()
		return
	}

	qry = `SELECT LAST_INSERT_ID() AS airport_id`

	err = tx.Debug().Raw(qry).Scan(&airport).Error

	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		tx.Rollback()
		return
	}

	qry = `INSERT INTO airport_translations (airport_id, language_id, description)
	     VALUES %s`

	var valueStrings []string
	var params []interface{}

	for _, l := range languages {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		params = append(params, airport.AirportId, l.LanguageId, airport.Description)
	}

	qry = fmt.Sprintf(qry, strings.Join(valueStrings, ","))

	err = repository.Db.Raw(qry, params...).Scan(&airport).Error
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		tx.Rollback()
		return
	}

	err = tx.Commit().Error

	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}

	c.JSON(200, "Ok")

}

func UpdateAirports(c *gin.Context) {
	airport_id := c.Param("airport_id")

	airport := AirportPayload{}

	err := c.BindJSON(&airport)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
	}

	lang := Languages{}

	header := c.GetHeader("Accept-Language")
	qry := `SELECT language_id 
	    FROM languages
		WHERE language_short_name=?`
	err = repository.Db.Debug().Raw(qry, header).Scan(&lang).Error

	if err != nil {
		println(err)
		c.JSON(400, "Db Error")
		return
	}

	tx := repository.Db.Begin()

	qry = `UPDATE airports SET
	      city_id=?,
		  airport_short_name=?
		  WHERE airport_id=?`

	err = tx.Debug().Exec(qry, airport.CityId, airport.AirportShortName, airport_id).Error

	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		tx.Rollback()
		return
	}

	qry = `UPDATE airport_translations SET
		description=?
		WHERE language_id=? AND airport_id=? `

	err = repository.Db.Debug().Exec(qry, airport.Description, lang.LanguageId, airport_id).Error

	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		tx.Rollback()
		return
	}

	err = tx.Commit().Error

	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}

	c.JSON(200, "Ok")
}
