package repository

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var dbErr error

func ConnectToDatabase() {

	databaseLogin := "root:password@tcp(127.0.0.1:3306)/airline?charset=utf8mb4&parseTime=True&loc=Local"

	Db, dbErr = gorm.Open(mysql.Open(databaseLogin), &gorm.Config{})
	if dbErr != nil {
		fmt.Printf("DB connection failed")
		log.Fatal(dbErr)
	}

}

func Migration() {
	Db.AutoMigrate(&Role{},
		&Permission{},
		&RolePermission{},
		&UserRole{},
		&Aeroplane{},
		&Flights{},
		&Airport{},
		&UsersFlights{},
		&Users{},
		&AeroplaneTranslation{},
		&AirportTranslation{},
		&Cities{},
		&Rating{},
		&Languages{})
}

type Role struct {
	RoleID   int    `gorm:"primaryKey"`
	RoleName string `gorm:"unique"`
}

type Permission struct {
	PermissionId int    `gorm:"primaryKey"`
	Permission   string `gorm:"unique"`
}

type RolePermission struct {
	RoleId       int `gorm:"uniqueIndex:idx_rolepermission_id"`
	PermissionId int `gorm:"uniqueIndex:idx_rolepermission_id"`
}

type UserRole struct {
	UserId int `gorm:"uniqueIndex:idx_userrole_id"`
	RoleId int `gorm:"uniqueIndex:idx_userrole_id"`
}

type Aeroplane struct {
	AeroplaneId            int `gorm:"primaryKey"`
	AeroplaneType          string
	AeroplaneNumberOfSeats int
}

type Flights struct {
	FlightId       int `gorm:"primaryKey"`
	FlightDateTime time.Time
	FlightPrice    int
	AeroplaneId    int
	StartAirPortId int
	EndAirPortId   int
}

type Airport struct {
	AirportId        int `gorm:"primaryKey"`
	CityId           int
	AirportShortName string
}

type UsersFlights struct {
	UserId   int `gorm:"uniqueIndex:idx_usersflights_id"`
	FlightId int `gorm:"uniqueIndex:idx_usersflights_id"`
}

type Users struct {
	UserId         int `gorm:"primaryKey"`
	UserFirstName  string
	UserLastName   string
	UserEmail      string
	UserContact    int
	Password       string
	PasswordRepeat string
}

type AeroplaneTranslation struct {
	AeroplaneId int `gorm:"uniqueIndex:idx_aeroplanelanguage_id"`
	LanguageId  int `gorm:"uniqueIndex:idx_aeroplanelanguage_id"`
	Description string
}

type AirportTranslation struct {
	AirportId   int `gorm:"uniqueIndex:idx_airporttranslation_id"`
	LanguageId  int `gorm:"uniqueIndex:idx_airporttranslation_id"`
	Description string
}

type Cities struct {
	CityId   int `gorm:"primaryKey"`
	CityName string
}

type Rating struct {
	RatingId int `gorm:"primaryKey"`
	UserId   int `gorm:"uniqueIndex:idx_users_flights_id"`
	FlightId int `gorm:"uniqueIndex:idx_users_flights_id"`
	Comment  string
	Rating   int
}

type Languages struct {
	LanguageId        int    `gorm:"primaryKey"`
	LanguageShortName string `gorm:"unique"`
}
