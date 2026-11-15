package config

type Config struct {
	GeminiAPIKey    string `env:"GEMINI_API_KEY"`
	GeminiModel     string `env:"GEMINI_MODEL"`
	Port            string `env:"PORT"`
	DiscordBotToken string `env:"DISCORD_BOT_TOKEN"`
}
