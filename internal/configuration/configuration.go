package configuration

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	ConnectionString string
	Port             string
}

var env = "local"

func GetConfig(environment string) Configuration {
	configuration := Configuration{}
	if environment != "" {
		env = environment
	}
	fileName := fmt.Sprintf("./settings.%s.json", env)
	fmt.Println("fileName -->" + fileName)

	gonfig.GetConf(fileName, &configuration)

	fmt.Println(configuration)
	return configuration
}
