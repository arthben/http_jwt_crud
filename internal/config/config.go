package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

func LoadConfig() (cfg EnvParams, err error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("development")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	// read OS env
	for _, k := range viper.AllKeys() {
		v := viper.GetString(k)
		viper.Set(k, os.ExpandEnv(v))
	}

	if err = viper.Unmarshal(&cfg); err != nil {
		return
	}

	// validate value
	if len(cfg.Port) == 0 {
		err = errors.New("Parameter Port is empty!")
		return
	}

	if _, err = strconv.Atoi(cfg.Port); err != nil {
		err = errors.New("Parameter Port invalid value")
		return
	}

	if len(cfg.ServerTimeout) == 0 {
		err = errors.New("Parameter serverTimeout is empty")
		return
	}

	if _, err = strconv.Atoi(cfg.ServerTimeout); err != nil {
		err = errors.New("Parameter serverTimeout invalid value")
		return
	}

	if len(cfg.Token.Expire) == 0 {
		err = errors.New("Parameter Token Expire is empty")
		return
	}

	if _, err = strconv.Atoi(cfg.Token.Expire); err != nil {
		err = errors.New("Parameter Token Expire invalid value")
		return
	}

	if len(cfg.Token.Issuer) == 0 {
		err = errors.New("Parameter Token Issuer is empty")
		return
	}

	if len(cfg.Server.PublicKey) == 0 {
		err = errors.New("Parameter Server Public Key is empty")
		return
	}

	rawSrvPub, err := os.ReadFile(cfg.Server.PublicKey)
	if err != nil {
		err = errors.New("Parameter Server Public Key not valid file")
		return
	}
	cfg.Server.PublicKey = string(rawSrvPub)

	if len(cfg.Server.PrivateKey) == 0 {
		err = errors.New("Parameter Server Private Key is empty")
		return
	}

	rawSrvPri, err := os.ReadFile(cfg.Server.PrivateKey)
	if err != nil {
		err = errors.New("Parameter Server Private Key not valid file")
		return
	}
	cfg.Server.PrivateKey = string(rawSrvPri)

	if len(cfg.DB.Host) == 0 {
		err = errors.New("Parameter DB Host is empty")
		return
	}

	if len(cfg.DB.Username) == 0 {
		err = errors.New("Parameter DB Username is empty")
		return
	}

	if len(cfg.DB.Name) == 0 {
		err = errors.New("Parameter DB Name is empty")
		return
	}

	if len(cfg.DB.MinPool) == 0 {
		err = errors.New("Parameter DB MinPool is empty")
		return
	}

	if _, err = strconv.Atoi(cfg.DB.MinPool); err != nil {
		err = errors.New("Parameter DB MinPool invalid value")
		return
	}

	if len(cfg.DB.MaxPool) == 0 {
		err = errors.New("Parameter DB MaxPool is empty")
		return
	}

	if _, err = strconv.Atoi(cfg.DB.MaxPool); err != nil {
		err = errors.New("Parameter DB MaxPool invalid value")
		return
	}

	if len(cfg.Client.Key) == 0 {
		err = errors.New("Parameter Client Key is empty")
		return
	}

	if len(cfg.Client.Secret) == 0 {
		err = errors.New("Parameter Client Secret is empty")
		return
	}

	if len(cfg.Client.PublicKey) == 0 {
		err = errors.New("Parameter Client Public Key is empty")
		return
	}

	rawSrvPubClient, err := os.ReadFile(cfg.Client.PublicKey)
	if err != nil {
		err = errors.New("Parameter Client Public Key not valid file")
		return
	}
	cfg.Client.PublicKey = string(rawSrvPubClient)

	return
}

type EnvParams struct {
	Port          string `yaml:"port"`
	ServerTimeout string `yaml:"serverTimeout"`
	AppMode       string `yaml:"appMode"`
	Token         struct {
		Expire string `yaml:"expire"`
		Issuer string `yaml:"issuer"`
	} `yaml:"token"`
	Server struct {
		PublicKey  string `yaml:"publicKey"`
		PrivateKey string `yaml:"privateKey"`
	} `yaml:"server"`
	DB struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		MinPool  string `yaml:"minPool"`
		MaxPool  string `yaml:"maxPool"`
	} `yaml:"db"`
	Client struct {
		Key       string `yaml:"key"`
		Secret    string `yaml:"secret"`
		PublicKey string `yaml:"publicKey"`
	} `yaml:"client"`
}
