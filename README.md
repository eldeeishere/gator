# Gator - RSS Feed Aggregator CLI

Gator is a command-line RSS feed aggregator written in Go. It allows users to manage RSS feeds, follow feeds, and browse posts from their subscribed feeds.

## Features

- **User Management**: Register new users and manage user authentication
- **Feed Management**: Add, follow, and unfollow RSS feeds
- **Post Aggregation**: Automatically fetch and store posts from RSS feeds
- **Browse Posts**: View posts from followed feeds with customizable limits
- **Database Integration**: PostgreSQL database with automatic schema management
- **Real-time Updates**: Background feed fetching with configurable intervals

## Prerequisites

Before running Gator, ensure you have the following installed:

- **Go** (version 1.21 or later)
- **PostgreSQL** (version 12 or later)

### Installing Prerequisites

#### Go Installation
Visit [https://golang.org/dl/](https://golang.org/dl/) and download the latest version for your operating system.

#### PostgreSQL Installation
- **macOS**: `brew install postgresql`
- **Ubuntu/Debian**: `sudo apt-get install postgresql postgresql-contrib`
- **Windows**: Download from [https://www.postgresql.org/download/windows/](https://www.postgresql.org/download/windows/)

## Installation

Install Gator using Go's built-in package manager:

```bash
go install github.com/eldeeishere/gator@latest
```

This will install the `gator` binary to your `$GOPATH/bin` directory (usually `~/go/bin`).

Make sure `~/go/bin` is in your PATH:
```bash
export PATH=$PATH:~/go/bin
```

## Configuration

### Database Setup

1. **Create a PostgreSQL database**:
   ```sql
   createdb gator
   ```

2. **Create the configuration file**:
   Create a file named `.gatorconfig.json` in your home directory with the following content:
   ```json
   {
     "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
     "current_user_name": ""
   }
   ```

   Replace `username`, `password`, and database connection details with your PostgreSQL credentials.

3. **Run database migrations**:
   The application uses goose for database migrations. Ensure your database schema is up to date by running the migrations in the `sql/schema` directory.

## Usage

### Available Commands

#### User Management
- **Register a new user**:
  ```bash
  gator register <username>
  ```

- **Login as an existing user**:
  ```bash
  gator login <username>
  ```

- **List all users**:
  ```bash
  gator users
  ```

#### Feed Management
- **Add a new RSS feed**:
  ```bash
  gator addfeed <feed_name> <feed_url>
  ```

- **List all feeds**:
  ```bash
  gator feeds
  ```

- **Follow an existing feed**:
  ```bash
  gator follow <feed_url>
  ```

- **Unfollow a feed**:
  ```bash
  gator unfollow <feed_url>
  ```

- **List followed feeds**:
  ```bash
  gator following
  ```

#### Post Management
- **Browse posts from followed feeds**:
  ```bash
  gator browse
  ```

- **Browse posts with limit**:
  ```bash
  gator browse limit <number>
  ```

#### Feed Aggregation
- **Start automatic feed fetching**:
  ```bash
  gator agg <duration>
  ```
  Examples:
  - `gator agg 30s` - Fetch feeds every 30 seconds
  - `gator agg 5m` - Fetch feeds every 5 minutes
  - `gator agg 1h` - Fetch feeds every hour

#### Database Management
- **Reset the database** (removes all users and data):
  ```bash
  gator reset
  ```

### Example Workflow

1. **Register and login**:
   ```bash
   gator register john_doe
   ```

2. **Add and follow a feed**:
   ```bash
   gator addfeed "TechCrunch" "https://techcrunch.com/feed/"
   ```

3. **Start feed aggregation**:
   ```bash
   gator agg 1h
   ```

4. **Browse posts**:
   ```bash
   gator browse limit 10
   ```

## Project Structure

```
gator/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── sqlc.yaml              # SQLC configuration
├── internal/
│   ├── commands/          # CLI command handlers
│   ├── config/            # Configuration management
│   ├── database/          # Database models and queries
│   └── rss/               # RSS feed parsing
└── sql/
    ├── queries/           # SQL queries
    └── schema/            # Database schema migrations
```

## Technologies Used

- **Go**: Programming language
- **PostgreSQL**: Database
- **SQLC**: SQL code generation
- **Goose**: Database migrations (schema management)
- **UUID**: Unique identifiers
- **RSS/XML**: Feed parsing

## Development

### Building from Source

1. **Clone the repository**:
   ```bash
   git clone https://github.com/eldeeishere/gator.git
   cd gator
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Build the application**:
   ```bash
   go build -o gator
   ```

4. **Run the application**:
   ```bash
   ./gator <command> [args...]
   ```

### Database Schema

The application uses the following database tables:
- `users`: User management
- `feeds`: RSS feed information
- `feed_follows`: User-feed relationships
- `posts`: Stored RSS posts

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).

## Troubleshooting

### Common Issues

1. **Database connection errors**:
   - Verify PostgreSQL is running
   - Check database credentials in `.gatorconfig.json`
   - Ensure the database exists

2. **Command not found**:
   - Make sure `~/go/bin` is in your PATH
   - Verify the installation completed successfully

3. **Feed parsing errors**:
   - Check that the RSS feed URL is valid and accessible
   - Some feeds may require specific user agents or headers

For more help, please open an issue on the GitHub repository.