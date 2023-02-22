module shortlink

go 1.20

replace (
	shortlink/internal/adapters => ./internal/adapters
	shortlink/internal/ports => ./internal/ports
	shortlink/internal/services => ./internal/services
	shortlink/internal/models => ./internal/models
)

require (
	github.com/caarlos0/env/v7 v7.0.0
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
