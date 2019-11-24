package ginResponse

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
)

func get(c *gin.Context, result interface{}, list bool, e error, count ...int) {
	if e != nil {
		if e.Error() == "not found" {
			NotFound(c)
		} else {
			BadRequest(c, e)
		}
		return
	}

	if list {
		c.Header("X-Total", fmt.Sprint(count[0]))
		if len(result.([]interface{})) == 0 {
			result = make([]string, 0)
		}
	}

	if !reflect.ValueOf(result).IsValid(){
		result = make(map[string]string, 0)
	}
	c.JSON(http.StatusOK, result)
}

func List(c *gin.Context, result []interface{}, e error, count int) {
	get(c, result, true, e, count)
}

func Retrieve(c *gin.Context, result interface{}, e error) {
	get(c, result, false, e)
}

func NotFound(c *gin.Context) {
	Msg(c, http.StatusNotFound, "not found")
}

func Created(c *gin.Context, e error) {
	if e != nil {
		BadRequest(c, e)
	} else {
		c.Status(http.StatusCreated)
	}
}

func BadRequest(c *gin.Context, e error) {
	status := http.StatusBadRequest
	log.Printf("[%v]: %s", fmt.Sprint(status), e.Error())
	Msg(c, status, e.Error())
}

func ServerError(c *gin.Context, e error) {
	status := http.StatusInternalServerError
	log.Printf("[%v]: %s", fmt.Sprint(status), e.Error())
	Msg(c, status, e.Error())
}

func OK(c *gin.Context, e error) {
	if e != nil {
		BadRequest(c, e)
		return
	}
	c.Status(http.StatusOK)
	return
}

func OKWithString(c *gin.Context, msg string) {
	Msg(c, http.StatusOK, msg)
}

func RouterNotExists(c *gin.Context, r *gin.Engine) {
	Msg(c, http.StatusNotFound, "Resource not exists")
}

func MethodNotAllowed(c *gin.Context) {
	Msg(c, http.StatusMethodNotAllowed, "Method Not Allowed")
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Msg(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"message": msg})
}
