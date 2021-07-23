package api

import (
	"airways/repository"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	UserId         int    `json:"user_id"`
	UserFirstName  string `json:"user_first_name"`
	UserLastName   string `json:"user_last_name"`
	UserEmail      string `json:"user_email"`
	UserContact    int    `json:"user_contact"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
}

type UserFilters struct {
	CityName string
}

var TokenSecret []byte = ([]byte("fbgnjfttbhrtru"))

func (u *Users) Validate() error {
	if u.Password != u.PasswordRepeat {
		return errors.New("Incorrect")
	}
	return nil
}

func GetUserFilters(uf UserFilters) (string, []interface{}, string, []interface{}) {
	var filtersQry string
	var filterParams []interface{}
	var joinQry string
	var joinParams []interface{}

	if uf.CityName != "" {
		if !strings.Contains(joinQry, " JOIN users_flights") {
			joinQry += ` JOIN users_flights uf ON u.user_id=uf.user_id `
		}
		if !strings.Contains(joinQry, " JOIN flights") {
			joinQry += ` JOIN flights f ON f.flight_id=uf.flight_id `
		}
		if !strings.Contains(joinQry, " JOIN airports") {
			joinQry += ` JOIN airports a ON f.end_air_port_id=a.airport_id `
		}
		if !strings.Contains(joinQry, " JOIN cities") {
			joinQry += ` JOIN cities c ON c.city_id=a.city_id`
		}
		filtersQry += ` AND c.city_name=?`
		filterParams = append(filterParams, uf.CityName)
	}

	return filtersQry, filterParams, joinQry, joinParams
}

func GetAllUsers(c *gin.Context) {
	user_id, err := GetUserIdFromUserClaims(c)
	filters := UserFilters{
		CityName: c.Query("city_name"),
	}

	filtersQry, filterParams, joinQry, joinParams := GetUserFilters(filters)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(user_id)

	users := []Users{}
	qry := `
	SELECT u.*
	FROM users u %s
	WHERE 1=1 %s`

	qry = fmt.Sprintf(qry, joinQry, filtersQry)

	var params []interface{}
	params = append(params, joinParams...)
	params = append(params, filterParams...)
	err = repository.Db.Debug().Raw(qry, params...).Scan(&users).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("users: %+v /n", users)
	c.JSON(200, users)
}

func GetUser(c *gin.Context) {
	user_id := c.Param("user_id")
	users := Users{}
	qry := "SELECT * FROM users WHERE user_id=?"
	err := repository.Db.Debug().Raw(qry, user_id).Scan(&users).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("users: %+v /n", users)
	c.JSON(200, users)
}

func DeleteUser(c *gin.Context) {
	user_id, err := GetUserIdFromUserClaims(c)

	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}

	qry := `UPDATE users
	SET user_first_name='',
	user_last_name='',
	user_email=CONCAT('deleted',user_email)
	WHERE user_id=? `

	err = repository.Db.Debug().Exec(qry, user_id).Error

	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}

	c.JSON(200, "Ok")
}

func HashAndSalt(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err

	}
	return string(hash), nil

}

func RegisterUser(c *gin.Context) {
	users := Users{}
	err := c.BindJSON(&users)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}

	validateErr := users.Validate()
	if validateErr != nil {
		fmt.Println(validateErr)
		c.JSON(400, "Db Error")
		return
	}

	hashPassword, passErr := HashAndSalt([]byte(users.Password))
	if passErr != nil {
		fmt.Println(passErr)
		c.JSON(400, "Db Error")
		return

	}

	qry := `INSERT INTO users (user_first_name,
		user_last_name,
		user_email,
		user_contact,
		password,
		password_repeat)
	VALUES (?,?,?,?,?,?)`
	err = repository.Db.Debug().Exec(qry,
		users.UserFirstName,
		users.UserLastName,
		users.UserEmail,
		users.UserContact,
		hashPassword,
		hashPassword).Error
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}
	c.JSON(200, "Ok")
}

func CreateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
	})

	tokenstring, err := token.SignedString(TokenSecret)
	if err != nil {
		return "", err
	}

	return tokenstring, nil
}

func AuthorizationCheck(c *gin.Context) {
	token, err := jwt.Parse(c.GetHeader("Authorization"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return TokenSecret, nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("user_claims", claims)
		c.Next()
	} else {
		c.JSON(401, "Unauthorized")
	}

}

func GetUserIdFromUserClaims(c *gin.Context) (int, error) {
	userClaims := c.MustGet("user_claims").(jwt.MapClaims)
	userId := userClaims["user_id"].(float64)

	userID := int(userId)

	return userID, nil
}

func Login(c *gin.Context) {
	// user body {password, email}
	// db   select password where email
	// bcrypt.CompareHashAndPassword() hash from db, password
	payload := Users{}

	err := c.BindJSON(&payload)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}

	users := Users{}

	qry := "SELECT user_id,password FROM users WHERE user_email=? "

	err = repository.Db.Debug().Raw(qry, payload.UserEmail).Scan(&users).Error
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Db Error")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(payload.Password))
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Wrong password")
		return
	}

	token, err := CreateToken(users.UserId)
	if err != nil {
		fmt.Println(400, err)
		return
	}
	c.JSON(200, token)

}
