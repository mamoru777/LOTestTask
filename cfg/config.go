package cfg

import (
	"LOTestTask/tools/cache"
	"LOTestTask/tools/logger"
	"os"
	"strconv"
)

type ServerConfig struct {
	Addr string
	Port string
}

type Config struct {
	CacheCfg  cache.Config
	ServerCfg ServerConfig
	LoggerCfg logger.Config
}

func (cfg *Config) LoadConfig() error {
	cacheSpace, err := strconv.Atoi(os.Getenv("CACHE_SPACE"))
	if err != nil {
		return err
	}

	cfg.CacheCfg = cache.Config{
		Space: cacheSpace,
	}

	cfg.ServerCfg = ServerConfig{
		Addr: os.Getenv("SERVER_ADDR"),
		Port: os.Getenv("SERVER_PORT"),
	}

	cfg.LoggerCfg = logger.Config{
		LogLevel: logger.LogLevel(os.Getenv("LOGGER_LOG_LEVEL")),
	}

	/*cfg.CacheCfg = cache.Config{
		Space: 30,
	}

	cfg.ServerCfg = ServerConfig{
		Addr: "localhost",
		Port: "9191",
	}

	cfg.LoggerCfg = logger.Config{
		LogLevel: "info",
	}*/

	return nil
}
