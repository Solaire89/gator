# gator

## Description:
Gator is an aggregation tool used to collect RSS feeds.

## Prerequisites:
Before installing gator, make sure you have the following installed:
- [Go](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)

## Installing gator:
1. Run the command `go install github.com/Solaire89/gator` to install the repo.
2. Verify (or set up) that your Go bin directory is on your PATH.
3. Confirm the program works by running gator from any directory.

## Configuring gator:
1. Create a file named .gatorconfig.json in your home directory.
2. The json file needs the following structure: 
    `{
        "db_url":"postgres://<username>:<password>@localhost:5432/gator",
        "current_user_name":"<username>"
    }`

## Goose
1. Install goose using the command `go install github.com/pressly/goose/v3/cmd/goose@latest`
2. After installing goose, navigate to gator/sql/schema.
3. Migrate up the servers using the command `goose postgress "postgres://<username>:<password>@localhost:5432/gator" up` 

## Commands to get you started
1. `gator register <username>`: This registers a user name to the system and automatically logs them in.
2. `gator login <username>`: This logs the user into the system.
- Before inputting the following commands, a user must be logged into the system.
3. `gator addfeed <name> <url>`: This adds the url of the feed to the system, assigning the name to the feed.
4. `gator follow <url>`: This follows the feed with the url in the system.
5. `gator agg <timeduration>`: This is the aggregation loop that runs in the background collecting feeds. The command takes in a time unit, such as 30s, 1m, etc. The command will fetch feeds every <timeduration> that was passed to the command. In order to stop the loop, press CTRL+C.