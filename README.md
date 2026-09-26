# RSS Feed manager

This application is a command line tool allows multiple users to subscribe and view RSS feeds. It is the boot.dev project 'gator'.

It is written in Go and uses Postgres for databases.

## Installation instructions

This guide assumes installation on a Linux machine.

Golang, Postgres and Goose are required.

### Clone repo

Clone the repo via:

```bash
git clone ...
cd gator
```

### Dependencies

#### Golang: 

Installation instructions here:

https://go.dev/doc/install

#### Postgres: 

Instructions can be found here https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql, but in short, run the following:

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
psql --version
```

Ensure a version is displayed after the final command.

Finally, set an easy to remember password for your postgres installation:

```bash
sudo passwd postgres
```

After running the above command you will be prompted to type a new password.

#### Goose

Goose is a database migration tool that also happens to be written in Go. May as well use it since we've already installed Go.

Installation instructions are here: https://github.com/pressly/goose#install, but in short:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Set up

Once all dependencies are installed is installed, we need to start the postgres service, create a json file in the home directory for key application settings, and perform the up migrations with goose.

#### Setting up the database

Start the postgres service:

```bash
sudo service postgresql start # Can stop later with `sudo service postgresql stop`
```

Once started, we need to set up the gator database in postgres. First connect to the database with:

```bash
sudo -u postgres psql
```

Once in the postgres terminal, run the following to create the database:

```sql
CREATE DATABASE gator;
```

If the databse was created succesfully, connect to it via:

```sql
\c gator
```

Finally set the database password:

```sql
ALTER USER postgres PASSWORD 'postgres';
```

The database endpoint should be `DATABASE_URL/gator`

To make other steps easier, set this as a variable in your bash terminal session:

```bash
db_endpoint="<your-database-endpoint-here>"
```

Then perform all up migrations for the database (assuming you are already in the project repo):

```bash
cd sql/schema
goose postgres $db_endpoint up
cd ../..
```

#### Setting up config files

Create a file in your home directory called `.gatorconfig.json`:

```bash
touch "$HOME/.gatorconfig.json"
```

Place the following into that file:

```json
{
  "db_url": "DATABASE_URL/gator?sslmode=disable",
  "current_user_name": ""
}
```

#### Tool installation:

Run `go install gator` within the project directory, the tool will now be usable within the terminal

#### Set up Complete!

That should be it for set up.

# Usage

The basic syntax is:

```bash
gator <command> <arguments>
```

The following commands exist:

- login <name> - 1 argument - Logs in as an existing user
- register <name> - 1 argument - Creates and logs in to a new user
- reset - 0 arguments - Deletes the current user
- users - 0 arguments - Lists all users
- agg <time-interval> - 1 argument, time interval - Fetches all RSS feeds the user is subscribed to at the time interval (in seconds) given
- addfeed <name> <url> - 2 argument, feed name, RSS feed url - Adds the given feed to the database (with the given name) and subscribes the user to that feed
- feeds - 0 arguments - Lists all feeds the user is subscribed to
- follow <url> - 1 argument, feed name - Make the current user follow the given feed
- following - 0 arguments - Lists all feeds the current user is following
- unfollow <url> - 1 argument, feed url - Makes the user unfollow a feeds
- browse [limit] - 1 arguments (optional), number of posts to list (defaults to 2 if not given) - Lists posts in database

