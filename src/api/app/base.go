package app

import (
	"api/utils/model"
	"api/utils/mongodb"
	"api/utils/response"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"io/ioutil"
)

func Index(c *gin.Context) {
	ginResponse.OKWithString(c, "Howdy! Help yourself doing some simple API stuff.")
}

// TODO, u.Action useless for now
func GET(c *gin.Context, u *model.UrlDivide) {
	if u.Id != "" {
		var result interface{}
		e := mongodb.FindOne(u.Collection, bson.M{"_id": u.Id}, &result)
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
	count, e := mongodb.FindAll(u.Collection, bson.M{}, &result, model.ExtParams{})
	ginResponse.List(c, result, e, count)

}

func POST(c *gin.Context, u *model.UrlDivide) {
	switch u.Action {
	case "":
		break
	case "test":

	default:
		return
	}

	var jsonData bson.M

	data, _ := ioutil.ReadAll(c.Request.Body)
	if e := json.Unmarshal(data, &jsonData); e != nil {
		ginResponse.BadRequest(c, e)
		return
	}
	if len(jsonData) == 0 {
		ginResponse.BadRequest(c, errors.New("empty entry"))
		return
	}
	delete(jsonData, "_id")
	e := mongodb.Insert(u.Collection, jsonData)
	ginResponse.OK(c, e)

}

func PUT(c *gin.Context, u *model.UrlDivide) {
	switch u.Action {
	case "":
		break
	default:
		return
	}
	var jsonData bson.M

	data, _ := ioutil.ReadAll(c.Request.Body)
	if e := json.Unmarshal(data, &jsonData); e != nil {
		ginResponse.BadRequest(c, e)
		return
	}

	e := mongodb.Update(u.Collection, bson.M{"_id": u.Id}, jsonData)
	ginResponse.OK(c, e)
}

func PATCH(c *gin.Context, u *model.UrlDivide) {
	switch u.Action {
	case "":
		break
	default:
		return
	}
	var patchData []model.PatchParams
	if e := c.Bind(&patchData); e != nil {
		ginResponse.BadRequest(c, e)
		return
	}
	e := mongodb.Patch(u.Collection, u.Id, &patchData)
	ginResponse.OK(c, e)

}

func DELETE(c *gin.Context, u *model.UrlDivide) {
	switch u.Action {
	case "":
		break
	default:
		return
	}
	e := mongodb.Delete(u.Collection, u.Id)
	ginResponse.OK(c, e)
}
