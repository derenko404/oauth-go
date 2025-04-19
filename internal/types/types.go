package types

type AppConfig struct {
	AppPort     string `env:"APP_PORT" env_default:"8080"`
	AppHost     string `env:"APP_HOST" env_default:"localhost"`
	AppLogLevel string `env:"APP_LOG_LEVEL" env_default:"debug"`

	DBPort     string `env:"DB_PORT" env_default:"5432"`
	DBHost     string `env:"DB_HOST" env_default:"localhost"`
	DBUser     string `env:"DB_USER" env_default:"postgres"`
	DBPassword string `env:"DB_PASSWORD" env_default:"postgres"`
	DBName     string `env:"DB_NAME" env_default:"gocommerce"`

	RedisHost string `env:"REDIS_HOST" env_default:"localhost"`
	RedisPort string `env:"REDIS_PORT" env_default:"6379"`

	JwtSecret string `env:"JWT_SECRET"`

	GoogleClientId     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `env:"GOOGLE_REDIRECT_URL"`

	GithubClientId     string `env:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	GithubRedirectURL  string `env:"GITHUB_REDIRECT_URL"`
}
