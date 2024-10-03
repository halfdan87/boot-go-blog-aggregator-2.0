# boot-go-blog-aggregator-2.0

This is a simple blog aggregator written in Go. It uses SQLite as a database and SQLC to generate the database client.

## Installation

1. Install SQLC

```bash
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
```

2. Install Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

3. Install SQLite

```bash
brew install sqlite
```

4. Create the database schema

```bash
cd sql/schema
goose -dir ../../sql/schema up
```

5. Create the database client

```bash
sqlc generate
```

6. Run the application

```bash
go run main.go
```

## Commands

### login

Login with a username. This will set the current user in the config file.

### register

Register a new user. This will set the current user in the config file.

### reset

Reset the current user in the config file.

### users

List all users.

### agg

Aggregate the feeds. This will fetch the feeds and create posts for them.

### addfeed

Add a new feed to the database.

### feeds

List all feeds.

### follow

Follow a feed.

### following

List all feeds that you are following.

### unfollow

Unfollow a feed.

### browse

Browse the posts of a user.



