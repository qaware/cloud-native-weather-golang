package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func main() {
	engine := gin.Default()

	// get global Monitor object
	m := ginmetrics.GetMonitor()
	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	// set middleware for gin
	m.Use(engine)

	// configuration for static files and templates
	engine.LoadHTMLFiles("templates/index.html")
	engine.StaticFile("/favicon.ico", "favicon.ico")

	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Weather Service",
		})
	})

	engine.GET("/api/weather", queryWeather) // get weather for city

	engine.GET("/healthz", healthz)
	engine.GET("/readyz", readyz)

	ConnectDatabase()
	engine.Run(port())
}

func queryWeather(c *gin.Context) {
	weather := GetOrRetrieveWeather(c.Query("city"))
	c.JSON(http.StatusOK, weather)
}

func healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}

func readyz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}
