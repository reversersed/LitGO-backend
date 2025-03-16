package swagger

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const swaggerUrl = "https://raw.githubusercontent.com/reversersed/LitGO-proto/main/gen/docs/swagger/swagger.json"

func fetchSwaggerJSON() ([]byte, error) {
	resp, err := http.Get(swaggerUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func InitiateSwagger(router *gin.Engine) {
	swaggerJSON, err := fetchSwaggerJSON()
	if err != nil {
		log.Fatalf("Failed to fetch Swagger JSON: %v", err)
	}

	err = os.WriteFile("doc.json", swaggerJSON, 0600)
	if err != nil {
		log.Fatalf("Failed to write Swagger JSON to file: %v", err)
	}

	router.GET("/doc.json", func(c *gin.Context) {
		c.File("doc.json")
	})
	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/doc.json"),
		ginSwagger.InstanceName("LitGO Swagger")))
}
