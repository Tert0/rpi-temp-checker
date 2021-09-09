package main

import (
	"fmt"
	_ "github.com/stianeikeland/go-rpio/v4"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)
import "github.com/gin-gonic/gin"

type tempStatus string

const (
	tempStatusNormal tempStatus = "NORMAL"
	tempStatusHigh tempStatus = "HIGH"
	tempStatusCritical tempStatus = "CRITICAL"
)

func getTemp() float64 {
	content, err := ioutil.ReadFile("/sys/class/hwmon/hwmon0/device/temp")
	if err != err {
		fmt.Println("Error: ", err.Error())
		panic("Cannot get CPU Temp!")
	}

	temp, err := strconv.ParseInt(strings.Replace(string(content), "\n", "", -1), 10, 64)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		panic("Cannot get CPU Temp!")
	}
	return float64(temp / 1000)
}

func getTempStatus() tempStatus {
	temp := getTemp()

	highTemp, err := strconv.ParseFloat(os.Getenv("HIGH_TEMP"), 64)
	if err != nil {
		highTemp = 55.0
	}
	criticalTemp, err := strconv.ParseFloat(os.Getenv("CRITICAL_TEMP"), 64)
	if err != nil {
		criticalTemp = 70.0
	}

	switch {
	case temp >= criticalTemp:
		return tempStatusCritical
	case temp >= highTemp:
		return tempStatusHigh
	default:
		return tempStatusNormal
	}
}

func main() {

	r := gin.Default()

	r.GET("/temp", func(c *gin.Context) {
		c.String(200, strconv.FormatFloat(getTemp(), 'f', 2, 64))
	})

	r.GET("/status", func(c *gin.Context) {
		c.String(200, string(getTempStatus()))
	})

	r.Run("0.0.0.0:8080")
}
