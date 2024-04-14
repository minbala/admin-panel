package main

import (
	"admin-panel/cmd"
	"admin-panel/dependency_manager"
	"log"
)

//	@title			AdminPanel  API
//	@version		1.0
//	@description	AdminPanel API in Go using Gin framework
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.minbala.com/support
//	@contact.email	support@devxmm.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host	localhost:8080

//	@securityDefinitions.apiKey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	container := dependency_manager.InitializeContainer()
	if cmdErr := cmd.Run(container); cmdErr != nil {
		log.Fatalf(cmdErr.Error())
	}
}
