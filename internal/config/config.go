package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type MainConfig struct {
	Server     ServerConfig   `yaml:"server"`
	PostgresDB PostgresConfig `yaml:"postgres"`
	Email      EmailConfig    `yaml:"email"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}
type PostgresConfig struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Host   string `yaml:"host"`
	Name   string `yaml:"name"`
	Reload bool   `yaml:"reload"`
	Port   string `yaml:"port"`
}

type EmailConfig struct {
	Sender   string `yaml:"sender"`
	Password string `yaml:"password"`
	Reciever string `yaml:"reciever"`
}

// загрузка конфигов из yaml файла в структуры
func (cfg *MainConfig) ConfigMustLoad(name string) {

	path := "./config/" + name + ".yaml"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		pathBackup := "../config/" + name + ".yaml"
		if _, err := os.Stat(pathBackup); os.IsNotExist(err) {
			panic("nothing from this path")
		}
		path = pathBackup
	}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("failed to read")
	}

}
