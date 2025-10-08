package router

import (
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"subs/subservice/internal/config"
	"subs/subservice/internal/handlers"
	"subs/subservice/internal/middleware"
	"time"
)

func Router(cfg *config.Config, db *sqlx.DB) *http.Server {
	mux := http.NewServeMux()
	h := handlers.NewHandler(db)

	mux.HandleFunc("/create_sub", h.HandleCreateSub)
	mux.HandleFunc("/sum_subs", h.HandleSumSubs)
	mux.HandleFunc("/get_all_subs", h.GetSubs)
	mux.HandleFunc("/update_sub", h.UpdateSub)
	mux.HandleFunc("/delete_sub", h.DeleteSub)

	srv := &http.Server{
		Addr:         cfg.SubServiceConfig.Port,
		Handler:      middleware.LoggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	slog.Info("Http server configured on", "port", cfg.SubServiceConfig)
	return srv
}
