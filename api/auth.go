package api

import (
	"airways/common"
	"airways/repository"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(permission string) gin.HandlerFunc {
	var fn = func(c *gin.Context) {

		userId, err := GetUserIdFromUserClaims(c)
		if err != nil {
			println(err)
			return
		}
		type Test struct {
			Id int `json:"id"`
		}

		test := Test{}

		qry := `SELECT ur.user_id
		FROM permissions p 
		JOIN role_permissions rp ON rp.permission_id=p.permission_id
		JOIN user_roles ur ON ur.role_id=rp.role_id 
		WHERE permission=? AND user_id=?`

		err1 := repository.Db.Debug().Raw(qry, permission, userId).Scan(&test).Error
		if err1 != nil {
			c.AbortWithStatus(403)
			common.LogError(err1)
			return
		}

		c.Next()
	}

	return fn
}


