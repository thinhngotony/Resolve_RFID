package utils

import (
	"github.com/magiconair/properties"
)

type Config struct {
	ApiKey string
}

type Database struct {
	ApiKey   string
	Username string
	Password string
	Hostname string
	Dbname   string
}

func LoadConfig(path string) string {
	P := properties.MustLoadFile(path, properties.UTF8)
	v := P.GetString("API.Key", "")
	return v

}

func LoadDatabase(path string) Database {

	P := properties.MustLoadFile(path, properties.UTF8)

	x := P.GetString("API.Key", "")
	y := P.GetString("USERNAME", "")
	z := P.GetString("PASSWORD", "")
	t := P.GetString("HOSTNAME", "")
	w := P.GetString("DBNAME", "")

	data := Database{
		ApiKey:   x,
		Username: y,
		Password: z,
		Hostname: t,
		Dbname:   w,
	}

	return data
}
