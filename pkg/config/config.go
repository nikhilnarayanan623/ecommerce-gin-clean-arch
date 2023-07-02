package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// to store env variables
type Config struct {
	AdminEmail    string `mapstructure:"ADMIN_EMAIL"`
	AdminUserName string `mapstructure:"ADMIN_USER_NAME"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBName        string `mapstructure:"DB_NAME"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBPort        string `mapstructure:"DB_PORT"`

	AdminAuthKey string `mapstructure:"ADMIN_AUTH_KEY"`
	UserAuthKey  string `mapstructure:"USER_AUTH_KEY"`

	TwilioAuthToken  string `mapstructure:"AUTH_TOKEN"`
	TwilioAccountSID string `mapstructure:"ACCOUNT_SID"`
	TwilioServiceID  string `mapstructure:"SERVICE_SID"`

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
	"ADMIN_EMAIL", "ADMIN_USER_NAME", "ADMIN_PASSWORD",
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_PORT", // database
	"ADMIN_AUTH_KEY", "USER_AUTH_KEY", // token auth
	"AUTH_TOKEN", "ACCOUNT_SID", "SERVICE_SID", // twilio
	"RAZOR_PAY_KEY", "RAZOR_PAY_SECRET", // razor pay
	"STRIPE_SECRET", "STRIPE_PUBLISH_KEY", "STRIPE_WEBHOOK", // stripe
	"GOAUTH_CLIENT_ID", "GOAUTH_CLIENT_SECRET", "GOAUTH_CALL_BACK_URL", //goath
}

var config Config

func LoadConfig() (Config, error) {

	// set-up viper
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	for _, env := range envsNames {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(config); err != nil {
		return config, err
	}
	return config, nil
}

func GetConfig() Config {
	return config
}
