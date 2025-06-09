package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
)

var Supabase *supabase.Client

func Init() {
	_ = godotenv.Load()

	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_ANON_KEY")
	Supabase = supabase.CreateClient(url, key)
}
