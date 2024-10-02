cd sql/schema
goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" down
cd ../..


