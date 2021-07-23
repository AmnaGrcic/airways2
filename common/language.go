package common

import (
	"airways/repository"

	"github.com/gin-gonic/gin"
)

type Language struct {
	LanguageId        int    `gorm:"primaryKey"`
	LanguageShortName string `gorm:"unique"`
}

func GetLanguageIdFromHeader(c *gin.Context) (int, error) {
	var langId int

	header := c.GetHeader("Accept-Language")

	qry := `SELECT language_id
	FROM languages
	WHERE language_short_name=?`

	err := repository.Db.Debug().Raw(qry, header).Row().Scan(&langId)

	if err != nil {
		return langId, err
	}

	return langId, nil
}
