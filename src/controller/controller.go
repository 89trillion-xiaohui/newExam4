package controller

import (
	"fmt"

	"test3/src/model"
	"test3/src/service"

	"github.com/gin-gonic/gin"
	"net/http"
)

func Create(c *gin.Context) {
	var GiftCodeInfo = model.GiftCodeInfo{}
	err := c.ShouldBind(&GiftCodeInfo)
	if err != nil {
		fmt.Println("Error : ", err)
	}

	Code := service.CreateCode(GiftCodeInfo)
	GiftCodeInfo.GiftCode = Code
	c.JSON(http.StatusOK, Code)
}

func Inquire(c *gin.Context) {
	GiftCode := c.Query("GiftCode")
	Code := service.Inquire(GiftCode)
	c.JSON(http.StatusOK, Code)
}

func Client(c *gin.Context) {
	GiftCode := c.Query("GiftCode")
	ClientName := c.Query("ClientName")
	Gift := service.Verify(ClientName, GiftCode)
	c.JSON(http.StatusOK, Gift)
}

func Log(c *gin.Context) {
	UserID := c.Query("UserID")
	User := service.Log(UserID)
	c.JSON(http.StatusOK, User)
}
