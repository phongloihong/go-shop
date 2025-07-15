# Libraries should use

- connect-go: http + Grpc + web_rpc 3 in 1
- nats: lightweight message broker
- sqlc + pgx + goose

# logging system

- Slog + Zap backend
- Slog: standard library, standard interface, context-aware
- separate loggin interface and implementation
- future-proff: slog will become standard in go ecosystem
