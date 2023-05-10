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

	JWT string `mapstructure:"JWT_CODE"`

	JWTAdmin string `mapstructure:"JWT_SECRET_ADMIN"`
	JWTUser  string `mapstructure:"JWT_SECRET_USER"`

	AUTHTOKEN  string `mapstructure:"AUTH_TOKEN"`
	ACCOUNTSID string `mapstructure:"ACCOUNT_SID"`
	SERVICESID string `mapstructure:"SERVICE_SID"`

	RazorPayKey    string `mapstructure:"RAZOR_PAY_KEY"`
	RazorPaySecret string `mapstructure:"RAZOR_PAY_SECRET"`

	StripSecretKey      string `mapstructure:"STRIPE_SECRET"`
	StripPublishKey     string `mapstructure:"STRIPE_PUBLISH_KEY"`
	StripeWebhookSecret string `mapstructure:"STRIPE_WEBHOOK"`

	GoathClientID      string `mapstructure:"GOAUTH_CLIENT_ID"`
	GoauthClientSecret string `mapstructure:"GOAUTH_CLIENT_SECRET"`
	GoauthCallbackUrl  string `mapstructure:"GOAUTH_CALL_BACK_URL"`
}

// to hold all names of env variables
var envsNames = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_PORT", // databse
	"JWT_CODE", "JWT_SECRET_ADMIN", "JWT_SECRET_USER", // jwt
	"AUTH_TOKEN", "ACCOUNT_SID", "SERVICE_SID", // twillio
	"RAZOR_PAY_KEY", "RAZOR_PAY_SECRET", // razor pay
	"STRIPE_SECRET", "STRIPE_PUBLISH_KEY", "STRIPE_WEBHOOK", // stripe

	"GOAUTH_CLIENT_ID", "GOAUTH_CLIENT_SECRET", "GOAUTH_CALL_BACK_URL",
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
func GetJWTSecret() JwtSecret {
	return JwtSecret{
		JWTAdmin: config.JWTAdmin,
		JWTUser:  config.JWTUser,
	}
}
func GetJWTCofig() string {
	return config.JWT
}

type JwtSecret struct {
	JWTAdmin string
	JWTUser  string
}

func GetConfig() Config {

	return config
}
