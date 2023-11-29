package config

import (
	"log"
	"os"
)

type Config struct {
	DbDsn    string
	CieloUrl string
	RedeUrl  string
	StoneUrl string
	CieloKey string
	RedeKey  string
	StoneKey string
}

var config Config

func init() {
	dbDsn, ok := os.LookupEnv("DB_DSN")
	if !ok || dbDsn == "" {
		log.Fatal("env var DB_DSN is required")
	}

	cieloUrl, ok := os.LookupEnv("CIELO_URL")
	if !ok || cieloUrl == "" {
		log.Fatal("env var CIELO_URL is required")
	}

	redeUrl, ok := os.LookupEnv("REDE_URL")
	if !ok || redeUrl == "" {
		log.Fatal("env var REDE_URL is required")
	}

	stoneUrl, ok := os.LookupEnv("STONE_URL")
	if !ok || stoneUrl == "" {
		log.Fatal("env var STONE_URL is required")
	}

	cieloKey, ok := os.LookupEnv("CIELO_KEY")
	if !ok || cieloKey == "" {
		log.Fatal("env var CIELO_KEY is required")
	}

	redeKey, ok := os.LookupEnv("REDE_KEY")
	if !ok || redeKey == "" {
		log.Fatal("env var REDE_KEY is required")
	}

	stoneKey, ok := os.LookupEnv("STONE_KEY")
	if !ok || stoneKey == "" {
		log.Fatal("env var STONE_KEY is required")
	}

	config = Config{
		DbDsn:    dbDsn,
		CieloUrl: cieloUrl,
		RedeUrl:  redeUrl,
		StoneUrl: stoneUrl,
		CieloKey: cieloKey,
		RedeKey:  redeKey,
		StoneKey: stoneKey,
	}
}

func GetConfig() Config {
	return config
}
