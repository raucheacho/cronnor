# Cronnor ⏰

**Lightweight HTTP Cron Job Scheduler** - A containerized Go service for
scheduling and executing HTTP-based cron jobs with a real-time HTMX web
interface.

## ✨ Features

- 🕒 **Cron-based scheduling** using standard cron expressions
- 🌐 **HTTP job execution** with configurable methods (GET, POST, PUT, etc.)
- 💻 **Modern web interface** built with HTMX and Tailwind CSS v4
- 📊 **Execution history** and detailed logging
- 🔄 **Live status updates** with 3-second polling
- 🐳 **Docker ready** with multi-stage builds
- 💾 **SQLite storage** - no external database required
- ⚡ **Lightweight** - minimal dependencies, pure Go

## 🚀 Quick Start

### Option 1: Using Docker (Recommended)

```bash
# Clone the repository
git clone <repo-url>
cd cronnor

# Start with Docker Compose
make docker-run

# Access the web interface
open http://localhost:8080
```

### Option 2: Local Development

**Prerequisites:**

- Go 1.23 or later
- Node.js & npm (for Tailwind CSS)

```bash
# Download dependencies
make deps

# Run the server
make run

# Access the web interface
open http://localhost:8080
```

## 📖 Usage

### Create a New Job

1. Navigate to `http://localhost:8080`
2. Click **"+ New Job"**
3. Fill in the form:
   - **Name**: Descriptive name for your job
   - **Cron Expression**: Standard cron format (e.g., `*/5 * * * *` for every 5
     minutes)
   - **Target URL**: The HTTP endpoint to call
   - **Method**: HTTP method (GET, POST, PUT, etc.)
   - **Payload**: Optional JSON payload for POST/PUT requests

### Cron Expression Examples

```
*/5 * * * *    # Every 5 minutes
0 * * * *      # Every hour
0 0 * * *      # Daily at midnight
0 9 * * 1      # Every Monday at 9 AM
*/30 9-17 * * * # Every 30 minutes between 9 AM and 5 PM
```

### Managing Jobs

- **Toggle**: Enable/disable jobs without deleting them
- **Run Now**: Execute a job immediately (bypasses the cron schedule)
- **Edit**: Modify job configuration
- **View Details**: See execution history and logs

## 🏗️ Architecture

### Tech Stack

- **Backend**: Go 1.23 with Chi router
- **Scheduler**: robfig/cron v3
- **Database**: SQLite with modernc.org/sqlite (pure Go)
- **Frontend**: HTMX, Tailwind CSS v4
- **Containerization**: Docker multi-stage builds

### Project Structure

```
cronnor/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── http/            # HTTP server and handlers
│   ├── jobs/            # Scheduler and executor
│   ├── models/          # Data models
│   └── storage/         # Database layer
├── migrations/          # SQL schema
├── web/
│   ├── node_modules/    # Frontend dependencies
│   ├── static/          # CSS and assets
│   ├── templates/       # HTML templates
│   └── package.json     # Frontend build config
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .air.toml            # Live reload config
```

## ⚙️ Configuration

Configure via environment variables:

| Variable         | Default                               | Description          |
| ---------------- | ------------------------------------- | -------------------- |
| `PORT`           | `8080`                                | HTTP server port     |
| `DB_PATH`        | `./data/cronnor.db`                   | SQLite database path |
| `MIGRATION_PATH` | `./migrations/001_initial_schema.sql` | Migration file path  |

### Example

```bash
export PORT=3000
export DB_PATH=/var/lib/cronnor/db.sqlite
./cronnor
```

## 🐳 Docker Deployment

### Build Image

```bash
make docker-build
```

### Run Container

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  --name cronnor \
  cronnor:latest
```

### Docker Compose

```bash
# Start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

## 🛠️ Development

### Build from Source

```bash
make build
```

### Run Tests

```bash
make test
```

### Format Code

```bash
make fmt
```

### Development Mode (with auto-reload)

Requires [air](https://github.com/cosmtrek/air) and npm:

```bash
make dev
```

This command runs both:

1. **Air**: Recompiles and restarts the Go server on file changes.
2. **Tailwind CLI**: Watches for changes in HTML/JS and rebuilds CSS.

## 📝 API Reference

### Endpoints

| Method | Path                | Description             |
| ------ | ------------------- | ----------------------- |
| GET    | `/jobs`             | Dashboard page          |
| GET    | `/jobs/list`        | Job list partial (HTMX) |
| POST   | `/jobs`             | Create new job          |
| GET    | `/jobs/{id}`        | Job details             |
| GET    | `/jobs/{id}/edit`   | Edit job form           |
| POST   | `/jobs/{id}`        | Update job              |
| POST   | `/jobs/{id}/toggle` | Toggle active status    |
| POST   | `/jobs/{id}/run`    | Execute job immediately |
| DELETE | `/jobs/{id}`        | Delete job              |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see LICENSE file for details.

## 🙏 Acknowledgments

- [Chi](https://github.com/go-chi/chi) - Lightweight HTTP router
- [robfig/cron](https://github.com/robfig/cron) - Cron scheduler for Go
- [HTMX](https://htmx.org) - High power tools for HTML
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) - Pure Go SQLite driver

---

Made with ❤️ using Go and HTMX
