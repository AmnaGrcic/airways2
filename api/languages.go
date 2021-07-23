package api

import (
	"airways/repository"
)

type Languages struct {
	LanguageId        int    `gorm:"primaryKey"`
	LanguageShortName string `gorm:"unique"`
}

type AirportTranslations struct {
	Airport_id  int `gorm:"unique"`
	Language_id int `gorm:"unique"`
}

func GetLanguages() ([]Languages, error) {
	languages := []Languages{}

	qry := `SELECT * FROM languages`

	err := repository.Db.Debug().Raw(qry).Scan(&languages).Error
	if err != nil {
		return languages, err
	}

	// fmt.Printf("languages: %+v \n",languages)
	return languages, nil
}
