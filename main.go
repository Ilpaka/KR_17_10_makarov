// @title PRoducts API KR_17_03
// @version 1.0
// @description This is simple example CRUD API for Peoducts
// @host localhost:8090
// @securityDefinitions.bearerauth BearerAuth
// @description JWT token
// @bearerformat JWT
package main

func main() {
	r := setupRouter()
	r.Run(":8090")
}
