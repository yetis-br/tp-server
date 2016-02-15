package util

import (
	"fmt"
	"log"

	"github.com/go-ini/ini"
)

var cfg *ini.File

func init() {
	cfg = ini.Empty()
}

//AppendConfigFile insert a new file with configurations
func AppendConfigFile(path string) {
	var err error
	err = cfg.Append("config.ini")
	if err != nil {
		log.Println(err)
	}
}

//GetKeyValue return the value of a key from the configuration file
func GetKeyValue(section string, key string) string {
	return cfg.Section(section).Key(key).String()
}

//FailOnError log fatal and panic if error is not null
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

//LogOnError log if error is not null
func LogOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
