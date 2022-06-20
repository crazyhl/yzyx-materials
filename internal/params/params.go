package params

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetInt(c *gin.Context, name string) (int, error) {
	val := c.Param(name)
	intVal, err := strconv.Atoi(val)
	return intVal, err
}

func GetUInt(c *gin.Context, name string) (uint, error) {
	val := c.Param(name)
	uintVal, err := strconv.ParseUint(val, 10, 64)
	return uint(uintVal), err
}
