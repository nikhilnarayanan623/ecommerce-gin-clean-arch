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

	AwsAccessKeyID string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretKey   string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AwsRegion      string `mapstructure:"AWS_REGION"`
	AwsBucketName  string `mapstructure:"AWS_BUCKET_NAME"`
}

// name of envs and used to read from system envs
var envsNames = []string{
	"ADMIN_EMAIL", "ADMIN_USER_NAME", "ADMIN_PASSWORD",
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_PORT", // database
	"ADMIN_AUTH_KEY", "USER_AUTH_KEY", // token auth
	"AUTH_TOKEN", "ACCOUNT_SID", "SERVICE_SID", // twilio
	"RAZOR_PAY_KEY", "RAZOR_PAY_SECRET", // razor pay
	"STRIPE_SECRET", "STRIPE_PUBLISH_KEY", "STRIPE_WEBHOOK", // stripe
	"GOAUTH_CLIENT_ID", "GOAUTH_CLIENT_SECRET", "GOAUTH_CALL_BACK_URL", //goath
	"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_REGION", "AWS_BUCKET_NAME", // aws s3
}

func LoadConfig() (config Config, err error) {

	// read from .env file
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	// if there is an error to read from config means user using system envs instead of .env file
	if err != nil {
		// bind from system envs
		for _, env := range envsNames {
			if err := viper.BindEnv(env); err != nil {
				return config, err
			}
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
