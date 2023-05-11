package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type Configuration struct {
	Db struct {
		Ip          string `json:"DbIp"`
		Port        string `json:"DbPort"`
		Name        string `json:"DbName"`
		User        string `json:"DbUser"`
		Password    string `json:"DbPassword"`
		AutoMigrate bool   `json:"DbMigrate"`
	}

	Listen struct {
		Ip   string `json:"HttpIp"`
		Port string `json:"HttpPort"`
	}

	Template struct {
		BaseName string `json:"BaseTemplate"`
	}
}

func LoadConfig() Configuration {
	c := flag.String("c", "env.json", "File JSON to Configuration")
	flag.Parse()

	file, err := os.Open(*c)

	if err != nil {
		log.Fatal("Error al abrir el archivo env.json")
	}

	defer func(file *os.File) {

		err := file.Close()
		if err != nil {
			log.Fatal("Error para cerrar nuestro archivo", err)
		}

	}(file)

	decoder := json.NewDecoder(file)
	Config := Configuration{}

	err = decoder.Decode(&Config)

	if err != nil {
		log.Fatal("No se pudo condificar", err)
	}

	return Config

}
