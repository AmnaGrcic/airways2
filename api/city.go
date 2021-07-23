package api

import (
	"airways/common"
	"airways/redis"
	"airways/repository"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Cities struct {
	CityId   int    `json:"city_id"`
	CityName string `json:"city_name"`
}

func GetAllCities(c *gin.Context) {
	cities := []Cities{}

	value, err := redis.GetRedisKey("cities")
	if err != nil {
		fmt.Println("Redis Error")
	} else if value != "" {
		fmt.Println(value)

		err = json.Unmarshal([]byte(value), &cities)
		if err != nil {
			fmt.Println(err)
		} else {
			c.JSON(200, cities)
			return
		}
	}

	qry := ` SELECT * 
	       FROM cities`

	err = repository.Db.Debug().Raw(qry).Scan(&cities).Error

	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}

	byteCities, _ := json.Marshal(cities)

	err = redis.SetRedisKey("cities", string(byteCities), 0)
	if err != nil {
		common.LogError(err)
	}

	fmt.Printf("cities %+v \n", cities)

	c.JSON(200, cities)
}

func GetCity(c *gin.Context) {
	city := Cities{}
	cityId := c.Param("city_id")

	var redisKey = "city_" + cityId
	fmt.Println(redisKey)

	value, err := redis.GetRedisKey(redisKey)
	if err != nil {
		fmt.Println("Log Error")
	} else if value != "" {
		fmt.Println(value)
		err = json.Unmarshal([]byte(value), &city)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(city)
			c.JSON(200, &city)
			return
		}

	}
	qry := `SELECT *
	      FROM cities
		  WHERE city_id=?`
	err = repository.Db.Debug().Raw(qry, cityId).Scan(&city).Error

	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}
	byteCity, _ := json.Marshal(city)

	err = redis.SetRedisKey(redisKey, string(byteCity), 0)

	if err != nil {
		common.LogError(err)
	}

	fmt.Printf("city %+v \n", city)

	c.JSON(200, city)

}

func DeleteCity(c *gin.Context) {
	city := c.Param("city_id")

	qry := `DELETE
	      FROM cities
		  WHERE city_id=?`

	err := repository.Db.Debug().Exec(qry, city).Error

	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}

	c.JSON(200, "Ok")
}

func CreateCity(c *gin.Context) {
	cities := Cities{}

	err := c.BindJSON(&cities)

	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}

	qry := `INSERT INTO cities (
		    city_id,
		    city_name
		) 
		VALUES (
			?,
			?
		)`

	err = repository.Db.Debug().Raw(qry, cities.CityId, cities.CityName).Scan(&cities).Error

	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}

	err = redis.DeleteRedisKey("cities")
	if err != nil {
		common.LogError(err)
	}

	c.JSON(200, "Ok")
}

func UpdateCity(c *gin.Context) {
	cityId := c.Param("city_id")

	cities := Cities{}

	err := c.BindJSON(&cities)

	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}

	qry := `UPDATE cities
	       SET 
		      city_name=?
		   WHERE  
		      city_id=?`

	err = repository.Db.Debug().Exec(qry, cities.CityName, cityId).Error

	if err != nil {
		common.LogError(err)
		c.JSON(400, "Db Error")
		return
	}

	var redisKey = "city_" + cityId

	err = redis.DeleteRedisKey("cities", redisKey)
	if err != nil {
		common.LogError(err)
	}

	c.JSON(200, "Ok")
}
