package app

import (
	"api/utils/mongodb"
	"api/utils/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"io/ioutil"
)

func Index(c *gin.Context) {
	ginResponse.OKWithString(c, "Howdy! Help yourself doing some simple API stuff.")
}

// TODO, u.Action useless for now
func GET(c *gin.Context, u *UrlDivide) {
	if u.Id != "" {
		var result interface{}
		e := mongodb.FindOne(u.Collection, bson.M{"_id": bson.ObjectIdHex(u.Id)}, &result)
		ginResponse.Retrieve(c, result, e)
		return
	}

	switch u.Action {
	case "":
		break
	case "test":

	default:
		return
	}

	var result []interface{}
	count, e := mongodb.FindAll(u.Collection, bson.M{}, &result, mongodb.ExtParams{})
	ginResponse.List(c, result, e, count)

}

func POST(c *gin.Context, u *UrlDivide) {
	if u.Id != "" {
		ginResponse.MethodNotAllowed(c)
		return
	}

	switch u.Action {
	case "":
		break
	case "test":

	default:
		return
	}

	var jsonData bson.M

	data, _ := ioutil.ReadAll(c.Request.Body)
	if e:=json.Unmarshal(data, &jsonData);e!=nil{
		ginResponse.BadRequest(c, e)
		return
	}
	if len(jsonData) == 0{
		ginResponse.BadRequest(c, errors.New("can not create a empty entry"))
		return
	}

	fmt.Println(jsonData)

	var result []interface{}
	count, e := mongodb.FindAll(u.Collection, bson.M{}, &result, mongodb.ExtParams{})
	ginResponse.List(c, result, e, count)

}

func PUT(c *gin.Context, u *UrlDivide) {

}

func PATCH(c *gin.Context, u *UrlDivide) {

}

func DELETE(c *gin.Context, u *UrlDivide) {

}
