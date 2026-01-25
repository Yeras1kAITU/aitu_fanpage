package app

import (
	"aitu_fanpage/internal/config"
	"aitu_fanpage/internal/handlers"
	"aitu_fanpage/internal/service"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	config   *config.Config
	server   *http.Server
	db       *mongo.Database
	services *service.Container
	handlers *handlers.Container
}

func New(cfg *config.Config) (*App, error) {

}
