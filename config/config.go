package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Conf struct {
	Name                       string `split_words:"true" default:"Loan Billing Engine"`
	Timeout                    int    `split_words:"true" default:"60"`
	ListenPort                 string `split_words:"true" default:"8080"`
	RootURL                    string `split_words:"true"`
	CallbackTimeout            int64  `split_words:"true" default:"60"`
	DatabaseMaxOpenConn        int    `split_words:"true" default:"20"`
	DatabaseMaxIdleConn        int    `split_words:"true" default:"10"`
	DatabaseSetConnMaxIdleTime int    `split_words:"true" default:"60"`
	DatabaseUrl                string `split_words:"true" default:"host=localhost port=5432 dbname=amartha user=indra password=pass1234 sslmode=disable"`
	PrivateKeyFile             string `split_words:"true" default:"./asset/keys/private_key.pem"`
	PublicKeyFile              string `split_words:"true" default:"./asset/keys/public_key.pem"`
	JwtExpired                 int64  `split_words:"true" default:"1800"`
}

func Load() *Conf {

	var Config Conf

	envconfig.MustProcess("APP", &Config)

	log.WithField("Config", Config).Info("Loaded configs")

	return &Config
}
