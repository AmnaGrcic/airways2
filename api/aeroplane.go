package api

import (
	"airways/common"
	"airways/repository"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Aeroplane struct {
	AeroplaneId            int    `json:"aeroplane_id"`
	AeroplaneType          string `json:"aeroplane_type"`
	AeroplaneNumberOfSeats int    `json:"aeroplane_number_of_seats"`
}

type AeroplaneTranslations struct {
	AeroplaneId            int    `json:"aeroplane_id"`
	AeroplaneType          string `json:"aeroplane_type"`
	AeroplaneNumberOfSeats int    `json:"aeroplane_number_of_seats"`
	Description            string `json:"description"`
}

type AeroplanePayload struct {
	AeroplaneId            int                `json:"aeroplane_id" `
	AeroplaneType          string             `json:"aeroplane_type"`
	AeroplaneNumberOfSeats int                `json:"aeroplane_number_of_seats"`
	FlightsRaw             json.RawMessage    `json:"-"`
	Flights                []FlightsUnmarshal `json:"flights" gorm:"-"`
}

type FlightsUnmarshal struct {
	FlightId    int `json:"flight_id"`
	FlightPrice int `json:"flight_price"`
}

func GetAllAeroplanes(c *gin.Context) {
	aeroplanes := []AeroplanePayload{}

	limit := c.DefaultQuery("limit", "5")
	offset := c.DefaultQuery("offset", "0")

	qry := `SELECT 
	a.aeroplane_id,
	a.aeroplane_type,
	a.aeroplane_number_of_seats,
	json_arrayagg(
		json_object(
			'flight_id', f.flight_id,
			'flight_price',f.flight_price
			)
	) AS flights_raw
		FROM aeroplanes a
			LEFT JOIN flights f ON a.aeroplane_id=f.aeroplane_id
			GROUP BY a.aeroplane_id, a.aeroplane_type, a.aeroplane_number_of_seats
			LIMIT ? OFFSET ?`

	err := repository.Db.Debug().Raw(qry, limit, offset).Scan(&aeroplanes).Error

	if err != nil {
		common.LogError(err)
		return
	}

	for i := range aeroplanes {
		err = json.Unmarshal(aeroplanes[i].FlightsRaw, &aeroplanes[i].Flights)
		if err != nil {
			common.LogError(err)
			c.JSON(400, "Db Error")
			return
		}

		fmt.Println(aeroplanes[i].Flights)
	}

	fmt.Printf("aeroplanes: %+v \n", aeroplanes)
	c.JSON(200, aeroplanes)
}

func GetAeroplane(c *gin.Context) {
	aeroplane_id := c.Param("aeroplane_id")
	qry := "SELECT * FROM aeroplanes where aeroplane_id=?"
	aeroplane := Aeroplane{}
	err := repository.Db.Debug().Raw(qry, aeroplane_id).Scan(&aeroplane).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("aeroplanes: %+v \n", aeroplane)
	c.JSON(200, aeroplane)

}

func DeleteAeroplane(c *gin.Context) {
	aeroplane_id := c.Param("aeroplane_id")

	qry := "DELETE FROM aeroplanes WHERE aeroplane_id=?"

	err := repository.Db.Debug().Exec(qry, aeroplane_id).Error
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}

	c.JSON(200, "Ok")

}

func CreateAeroplane(c *gin.Context) {
	aeroplane := AeroplaneTranslations{}

	err := c.BindJSON(&aeroplane)

	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}

	languages, LangErr := GetLanguages()
	if LangErr != nil {
		fmt.Println(LangErr)
		return
	}

	tx := repository.Db.Begin()

	qry := `INSERT INTO aeroplanes (aeroplane_type, aeroplane_number_of_seats) 
	VALUES (?,?)`

	err = tx.Debug().Exec(qry,
		aeroplane.AeroplaneType,
		aeroplane.AeroplaneNumberOfSeats).Error
	if err != nil {
		fmt.Println(err)
		c.JSON(500, "Db Error")
		tx.Rollback()
		return
	}

	qry = `SELECT LAST_INSERT_ID() AS aeroplane_id`

	err = tx.Debug().Raw(qry).Scan(&aeroplane).Error

	if err != nil {
		fmt.Println(err)
		c.JSON(500, "Db Error")
		tx.Rollback()
		return
	}

	qry = `INSERT INTO aeroplane_translations (aeroplane_id,language_id,description)
	VALUES %s`

	var valueStrings []string
	var params []interface{}

	for _, l := range languages {
		valueStrings = append(valueStrings, "(?,?,?)")
		params = append(params, aeroplane.AeroplaneId, l.LanguageId, aeroplane.Description)
	}

	qry = fmt.Sprintf(qry, strings.Join(valueStrings, ","))

	err = tx.Exec(qry, params...).Error

	if err != nil {
		common.LogError(err)
		c.JSON(500, "Db Error")
		tx.Rollback()
		return
	}

	err = tx.Commit().Error
	if err != nil {
		common.LogError(err)
		c.JSON(500, "Db Error")
		return
	}

	c.JSON(200, "Ok")
}

func UpdateAeroplane(c *gin.Context) {
	aeroplane_id := c.Param("aeroplane_id")

	aeroplane := AeroplaneTranslations{}

	err := c.BindJSON(&aeroplane)

	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}
	lang_id, err := common.GetLanguageIdFromHeader(c)
	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}

	tx := repository.Db.Begin()

	qry := `UPDATE aeroplanes SET 
				aeroplane_type=?, 
				aeroplane_number_of_seats=? 
			WHERE aeroplane_id=?`

	err = tx.Debug().Exec(qry,
		aeroplane.AeroplaneType,
		aeroplane.AeroplaneNumberOfSeats,
		aeroplane_id).Error

	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		tx.Rollback()
		return
	}

	qry = `UPDATE aeroplane_translations
	SET description=?
	WHERE language_id=? AND aeroplane_id=?`

	err = tx.Debug().Exec(qry, aeroplane.Description, lang_id, aeroplane_id).Error

	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		tx.Rollback()
		return
	}
	err = tx.Commit().Error

	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, "Ok")

}
