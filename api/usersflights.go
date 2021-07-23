package api

import (
	"airways/common"
	"airways/repository"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UsersFlights struct {
	UserId   int `json:"user_id"`
	FlightId int `json:"flight_id"`
}

type Rating struct {
	RatingId int    `json:"rating_id" gorm:"primaryKey"`
	UserId   int    `json:"user_id" gorm:"uniqueIndex:idx_users_flights_id"`
	FlightId int    `json:"flight_id" gorm:"uniqueIndex:idx_users_flights_id"`
	Comment  string `json:"comment"`
	Rating   int    `json:"rating"`
}

func GetUserFlights(c *gin.Context) {
	flightsarray := []Flights{}

	user_id, err := GetUserIdFromUserClaims(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	qry := `SELECT f.*
	FROM users as u
		join users_flights uf on uf.user_id = u.user_id 
		join flights f on f.flight_id = uf.flight_id	
	WHERE u.user_id=? `

	err = repository.Db.Debug().Raw(qry, user_id).Scan(&flightsarray).Error
	if err != nil {
		println(err)
		return
	}
	fmt.Printf("flight: %+v \n", flightsarray)
	c.JSON(200, flightsarray)

}

func (r *Rating) Validate() error {

	if r.Rating < 1 || r.Rating > 5 {
		return errors.New("invalid rating")
	}
	return nil
}

func CreateReview(c *gin.Context) {

	flight_id := c.Param("flight_id")
	user_id, err := GetUserIdFromUserClaims(c)

	if err != nil {
		println(err)
		return
	}

	ratings := Rating{}

	err = c.BindJSON(&ratings)
	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}

	validateErr := ratings.Validate()

	if validateErr != nil {
		common.LogError(validateErr)
		c.JSON(400, "Db Error")
		return
	}

	qry := `INSERT INTO ratings (user_id,flight_id,comment,rating)
	      VALUES (?,?,?,?)
		 `

	err = repository.Db.Debug().Exec(qry, user_id, flight_id, ratings.Comment, ratings.Rating).Error

	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}

	c.JSON(200, "Ok")

}
