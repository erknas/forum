package poolcfg

import (
	"fmt"

	"github.com/erknas/forum/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg *config.Config) (*pgxpool.Config, error) {
	dns := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	return pgxpool.ParseConfig(dns)
}
