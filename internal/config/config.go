package config

type Config struct {
	Port  string
	DBDSN string
}

func LoadConfig() Config {
	return Config{
		Port:  "8080",
		DBDSN: "postgres://postgres:secret@postgres:5432/soccer_manager?sslmode=disable",
	}
}
