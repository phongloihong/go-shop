# Libraries should use

- connect-go: http + Grpc + web_rpc 3 in 1
- nats: lightweight message broker
- sqlc + pgx + goose

## logging system

- Slog + Zap backend
- Slog: standard library, standard interface, context-aware
- Separate logging interface and implementation
- Future-proof: slog will become standard in go ecosystem
