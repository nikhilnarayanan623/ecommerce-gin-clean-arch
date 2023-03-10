package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// to store env variables
type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBPort     string `mapstructure:"DB_PORT"`
	JWT        string `mapstructure:"JWT_CODE"`
	AUTHTOKEN  string `mapstructure:"AUTH_TOKEN"`
	ACCOUNTSID string `mapstructure:"ACCOUNT_SID"`
	SERVICESID string `mapstructure:"SERVICE_SID"`
}

// to hold all names of env variables
var envsNames = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_PORT", "JWT_CODE", "AUTH_TOKEN", "ACCOUNT_SID", "SERVICE_SID",
}

var config Config // create an instance of Config

// func to get env variable and store it on struct Config and retuen it with error as nil or error
func LoadConfig() (Config, error) {

	// set-up viper
	viper.AddConfigPath("./")   // add the config path
	viper.SetConfigFile(".env") // set up the file name to viper
	viper.ReadInConfig()        // read the env file

	// range through through the envNames and take each envName and bind that env variable to viper
	for _, env := range envsNames {
		if err := viper.BindEnv(env); err != nil {
			return config, err // error when binding the env to viper
		}
	}

	// then unmarshel the viper into config variable
	if err := viper.Unmarshal(&config); err != nil {
		return config, err // error when unmarsheling the viper to env
	}

	// atlast validate the config file using validator pakage
	// create instance and direct validte
	if err := validator.New().Struct(config); err != nil {
		return config, err // error when validating struct
	}

	//successfully loaded the env values into struct config
	return config, nil
}

// to get the secred code for jwt
func GetJWTCofig() string {

	return config.JWT
}

func GetCofig() Config {
	return config
}
