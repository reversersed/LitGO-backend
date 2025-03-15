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

func fetchSwaggerJSON(url string) ([]byte, error) {
	resp, err := http.Get(url)
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
	url := "https://raw.githubusercontent.com/reversersed/LitGO-proto/main/gen/docs/swagger/swagger.json"
	swaggerJSON, err := fetchSwaggerJSON(url)
	if err != nil {
		log.Fatalf("Failed to fetch Swagger JSON: %v", err)
	}

	// Сохраняем Swagger JSON в файл (опционально)
	err = os.WriteFile("doc.json", swaggerJSON, 0644)
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
