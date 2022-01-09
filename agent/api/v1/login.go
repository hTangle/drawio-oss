package v1

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"super-markdown-editor-web/model"
	"time"
)

func (c *Controller) Login(ctx *gin.Context) {
	var loginReq model.Account
	if err := ctx.ShouldBind(&loginReq); err == nil {
		loginReq.UserName = loginReq.Email
		if loginReq.CheckAccount() {
			generateToken(ctx, loginReq)
		}
	} else {
		c.Error(ctx, ERROR, ERROR, "error")
	}
}

func generateToken(ctx *gin.Context, account model.Account) {
	j := model.NewJWT()
	claims := model.CustomClaims{
		Name:  account.UserName,
		Email: account.Email,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),
			ExpiresAt: int64(time.Now().Unix() + 3600000),
			Issuer:    model.LocalEditorConf.SignedKey(),
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": err.Error(),
			"data":    nil,
		})
	}
	logrus.Warnf("generate new token: %s", token)
	ctx.SetCookie("token", token, 10000000, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  1,
		"code":    200,
		"message": "login success",
		"data": map[string]string{
			"user_name": account.UserName,
			"token":     token,
		},
	})
}
