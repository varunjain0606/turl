package types

import (
	"api_service/server/cache"
	"api_service/server/database"
	"api_service/server/handler"
)

// App is the main app for TinyURL API
type App struct {
	DB      database.Database
	Handler handler.Handler
	Cache   *cache.Store
}
