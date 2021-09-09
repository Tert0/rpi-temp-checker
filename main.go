package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stianeikeland/go-rpio/v4"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type tempStatus string

const (
	tempStatusNormal tempStatus = "NORMGeAL"
	tempStatusHigh tempStatus = "HIGH"
	tempStatusCritical tempStatus = "CRITICAL"
)

var gpioPin int64

func getTemp() float64 {
	content, err := ioutil.ReadFile("/sys/class/hwmon/hwmon0/device/temp")
	if err != nil {
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

func checkJob(timeout time.Duration) {
	pin := rpio.Pin(gpioPin)

	switch getTempStatus() {
	case tempStatusNormal:
		pin.Low()
	case tempStatusHigh:
		pin.High()
	case tempStatusCritical:
		pin.High()
	}
	time.Sleep(time.Second * timeout)
	go checkJob(timeout)
}

func main() {
	err := rpio.Open()

	if err != nil {
		fmt.Println("Error: ", err.Error())
		panic("RPi GPIO Error")
	}

	gpioPin, err = strconv.ParseInt(os.Getenv("GPIO_PIN"), 10, 64)
	if err != nil {
		gpioPin = 3
	}

	timeout, err := strconv.ParseInt(os.Getenv("TIMEOUT"), 10, 64)
	if err != nil {
		timeout = 10
	}
	go checkJob(time.Duration(timeout))

	r := gin.Default()

	r.GET("/temp", func(c *gin.Context) {
		c.String(200, strconv.FormatFloat(getTemp(), 'f', 2, 64))
	})

	r.GET("/status", func(c *gin.Context) {
		c.String(200, string(getTempStatus()))
	})


	err = r.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error: ", err.Error())
		panic("Can't start Webserver!")
	}
}
