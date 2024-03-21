package utils

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	// "github.com/goombaio/namegenerator"
)

func GenerateUsername() string {
	//	seed := time.Now().UTC().UnixNano()
	//	nameGenerator := namegenerator.NewNameGenerator(seed)
	//	username := nameGenerator.Generate()
	//	fmt.Println("Username: " + username)
	username := randomdata.SillyName()
	fmt.Println("Username: " + username)

	return username
}
