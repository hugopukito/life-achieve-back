package main

import (
	"lifeAchieve/repository"
	"lifeAchieve/router"
)

// @title Life Achieve
// @version 1.0
func main() {
	repository.InitDB()
	router.InitRouter()
}
