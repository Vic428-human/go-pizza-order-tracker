package main

import "os"

type Config struct {
	Port   string
	DBPath string
}

func loadConfig() Config {
	return Config{
		Port:   getEnv("PORT", "8080"), // 定義key 跟 value 
		DBPath: getEnv("DATABASE_URL", "./data/orders.db"),
	}
}

// Example: Setting an environment variable in code
// err := os.Setenv("NEW_VAR", "GoLang Rocks!")
// Retrieve the newly set variable
// fmt.Println("NEW_VAR =", os.Getenv("NEW_VAR")) => NEW_VAR = GoLang Rocks!
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key); 
	if(value != "") {
		return value
	}
	return defaultValue
}