package main

import (
	_ "github.com/junminhong/thrurl/docs"
	"github.com/junminhong/thrurl/router"
)

// @title           Thrurl API
// @version         1.0
// @description     一個簡單易用且強大的縮網址服務
// @termsOfService  http://swagger.io/terms/

// @contact.name   junmin.hong
// @contact.url    https://github.com/junminhong
// @contact.email  junminhong1110@gmail.com

// @license.name  MIT
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      127.0.0.1:9020
// @BasePath  /api/v1
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
func main() {
	router.Setup()
}
