package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cin/daily-spark/internal/combination"
	"github.com/cin/daily-spark/internal/database"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/cin/daily-spark/docs"
)

// @title Daily-Spark API
// @version 1.0
// @description Daily-Spark 后端接口文档
// @BasePath /api/v1
// @schemes http
func main() {
    addr := ":8080"
    if p := os.Getenv("PORT"); p != "" {
        addr = ":" + p
    }

    db, err := database.InitDB()
    if err != nil {
        log.Fatal(err)
    }

    userModule := combination.NewUserModule(db)

    r := gin.Default()
    combination.RegisterUserRoutes(r, userModule.Handler)

    // 注册 Swagger 路由
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    r.GET("/healthz", func(c *gin.Context) {
        c.String(http.StatusOK, "ok")
    })

    log.Printf("Daily-Spark backend listening on %s", addr)
    if err := r.Run(addr); err != nil {
        log.Fatal(err)
    }
}