<div align="center">

# ✅ Todo List — Full-Stack

A secure, production-minded **Todo application** built with **Go (Gin)** and **React + TypeScript**.
Register, sign in, and manage your personal tasks with powerful search and filtering.

[![Go](https://img.shields.io/badge/Go-Gin-00ADD8?logo=go&logoColor=white)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-DC382D?logo=redis&logoColor=white)](https://redis.io/)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Vite](https://img.shields.io/badge/Vite-646CFF?logo=vite&logoColor=white)](https://vite.dev/)
[![Tailwind](https://img.shields.io/badge/Tailwind-v4-06B6D4?logo=tailwindcss&logoColor=white)](https://tailwindcss.com/)

</div>

---

## 📖 Overview

This project is a full-stack Todo application. To create todos you must **sign up and log in** — every todo belongs to its owner, and each user only sees and manages their own tasks.

Todos are created through a form (title, priority, category, due date, tags) and can be **searched and filtered** by title, category and priority, so you can always find what you need quickly.

The goal of the project is not just a todo list, but a **clean reference implementation** of a modern, secure auth flow and a layered Go backend paired with a modern React frontend.

## ✨ Features

### Application
- 🔐 User registration and login
- 📝 Create, edit, complete and delete todos via a validated form
- 🏷️ Per-todo **priority** (low / medium / high), **category**, **due date** and **tags**
- 🔎 Live search (debounced) + filter by **category** and **priority**
- 🌗 Light / dark theme
- 🔔 Toast notifications and friendly error handling

### Backend security & architecture
- **JWT** authentication with short-lived **access token** (15 min) and long-lived **refresh token** (7 days)
- Tokens are stored in **HttpOnly, SameSite=Strict cookies** (never exposed to JavaScript)
- Automatic **token refresh** and **logout / logout-all** (sessions tracked in **Redis**)
- **CSRF protection** — signed CSRF token bound to the session (`X-CSRF-Token` header) plus origin / `Sec-Fetch-Site` checks
- **CORS** with an explicit allow-list and credentials support
- **Rate limiting** — global limiter plus a much stricter one on auth endpoints
- **Argon2** password hashing with a configurable password policy
- **Role-based authorization** (`admin`, `user`, `guest`) with per-action rules
- Secure headers, request ID, panic recovery, request timeout and body-size limit
- Generic, reusable **CRUD handlers** with a **filter / sort / pagination** query engine
- **Layered (clean) architecture**, soft deletes, and SQL **migrations**

## 🧱 Tech Stack

| Layer | Technologies |
|-------|--------------|
| **Backend** | Go, [Gin](https://gin-gonic.com/), GORM, PostgreSQL, Redis, JWT, Viper, golang-migrate |
| **Frontend** | React 19, TypeScript, Vite, React Router, TanStack Query, shadcn/ui, Tailwind CSS v4, React Hook Form, Zod, Axios, Zustand, Sonner |
| **Tooling** | pnpm, Docker Compose, Makefile, ESLint, Prettier |

## 🗂️ Project Structure

```text
todo-list-FullStack/
├── Backend/
│   ├── cmd/api/              # Application entry point
│   ├── config/               # YAML configuration
│   ├── migrations/           # SQL migrations (schema + seed roles)
│   ├── internal/
│   │   ├── api/              # HTTP layer: router, handlers, middleware, validation, responses
│   │   ├── usecase/          # Business logic (auth, user, todo, authz) + DTOs
│   │   ├── domain/           # Models, authz rules, filter spec, repository interfaces
│   │   ├── infra/            # PostgreSQL, Redis, repository implementations
│   │   ├── config/           # Config loading (YAML + env)
│   │   └── bootstrap/        # Dependency wiring (DI container)
│   ├── pkg/                  # Reusable packages: security (JWT/CSRF/Argon2), ratelimit, apperror
│   ├── docker-compose.yml    # PostgreSQL + Redis + migrations
│   └── Makefile
└── Frontend/
    └── src/
        ├── page/             # Route pages (Home, Sign in, Sign up, Todo list, 404)
        ├── features/         # auth/ and todo/ (api, form, components, hooks, types)
        ├── components/       # UI components (shadcn/ui + custom)
        ├── services/         # Axios instance (CSRF header + auto refresh)
        ├── routes/           # React Router configuration
        └── hooks/ utils/ lib/
```

### Backend layers

```text
HTTP request
   │
   ▼
 api (router → middleware → handler)      ← validation, authn/authz, CSRF, rate limit
   │
   ▼
 usecase (business logic, DTOs)
   │
   ▼
 domain (models, rules, interfaces)       ← no dependency on frameworks
   ▲
   │ implements
 infra (PostgreSQL / Redis repositories)
```

## 🔑 Authentication Flow

1. **Register / Login** → server sets three cookies: `access_token` (HttpOnly), `refresh_token` (HttpOnly) and `csrf_token` (readable by the frontend).
2. For every state-changing request, the frontend sends the CSRF token in the `X-CSRF-Token` header.
3. When the access token expires, the Axios interceptor calls `/auth/refresh` **once** (concurrent requests share the same refresh), then retries the original request.
4. If the refresh fails, the user is redirected to the sign-in page.
5. **Logout** clears cookies and revokes the session in Redis; **logout-all** revokes every session of the user.

## 🌐 API Reference

Base path: `/api/v1`

| Method | Endpoint | Auth | Description |
|--------|----------|:----:|-------------|
| `GET` | `/healthz` | ❌ | Health check (outside the base path) |
| `POST` | `/auth/register` | ❌ | Create a new user |
| `POST` | `/auth/login` | ❌ | Log in and receive auth cookies |
| `POST` | `/auth/refresh` | ❌ (refresh cookie) | Rotate tokens |
| `POST` | `/auth/logout` | ❌ | Log out current session |
| `POST` | `/auth/logout-all` | ✅ | Log out from all sessions |
| `GET` | `/me` | ✅ | Current user and permissions |
| `GET` | `/todos` | ✅ | List todos (search / filter / sort / paginate) |
| `GET` | `/todos/:id` | ✅ | Get a todo |
| `POST` | `/todos` | ✅ | Create a todo |
| `PATCH` | `/todos/:id` | ✅ | Update a todo |
| `DELETE` | `/todos/:id` | ✅ | Delete a todo |

**List query parameters:** `page`, `pageSize`, `sort`, `search`, plus field filters such as `category`, `priority`, `done`, `tags`, `dueDate`.

```http
GET /api/v1/todos?search=report&priority=high&category=Work&page=1&pageSize=20
```

**Todo fields**

| Field | Type | Notes |
|-------|------|-------|
| `title` | string | required, max 200 chars |
| `priority` | `low` \| `medium` \| `high` | default `medium` |
| `category` | string \| null | max 50 chars |
| `dueDate` | ISO datetime \| null | |
| `tags` | string[] | |
| `done` | boolean | default `false` |

> Passwords must be at least **12 characters** and include uppercase, lowercase and a digit (configurable in `config.yml`).

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) (version from `Backend/go.mod`)
- [Node.js](https://nodejs.org/) and [pnpm](https://pnpm.io/)
- [Docker](https://www.docker.com/) & Docker Compose (for PostgreSQL and Redis)

### 1. Clone

```bash
git clone https://github.com/MohammadrezaNadirkhanloo/todo-list-FullStack.git
cd todo-list-FullStack
```

### 2. Backend

```bash
cd Backend

# Create your environment file and set passwords / secrets
cp .env.example .env
```

Edit `.env` — at minimum set these:

```env
APP_POSTGRES_PASSWORD=your_db_password
APP_POSTGRES_DBNAME=tododb
APP_REDIS_PASSWORD=your_redis_password
APP_JWT_ACCESSSECRET=<random secret>      # openssl rand -base64 48
APP_JWT_REFRESHSECRET=<different secret>  # openssl rand -base64 48
APP_CORS_ALLOWORIGINS=http://localhost:5173
```

Start PostgreSQL, Redis and run migrations, then start the API:

```bash
docker compose up -d      # postgres + redis + migrate
make run                  # or: go run ./cmd/api
```

The API will be available at `http://localhost:<APP_SERVER_PORT>` (e.g. `5005`, or `8080` if not overridden).

### 3. Frontend

```bash
cd Frontend
cp .env.example .env
```

Set the API URL in `.env`:

```env
VITE_API_BASE_URL=http://localhost:5005/api/v1
```

```bash
pnpm install
pnpm dev
```

Open **http://localhost:5173**, create an account at `/sign-up`, sign in, and go to `/todo-list`.

## ⚙️ Configuration

Backend settings live in `Backend/config/config.yml` and can be overridden with environment variables prefixed with `APP_` (for example `APP_SERVER_PORT`, `APP_POSTGRES_HOST`). Key sections: `server`, `postgres`, `redis`, `jwt`, `authz`, `password`, `pagination`, `cors`, `csrf`, `ratelimit`.

> ⚠️ **Production checklist:** set `APP_SERVER_INSECURECOOKIES=false` (cookies become `Secure` / `__Host-` prefixed, HTTPS required), use strong unique JWT secrets, enable PostgreSQL SSL, and keep `APP_AUTHZ_ENFORCE=true`.

## 🛠️ Useful Commands

**Backend (`make`)**

| Command | Description |
|---------|-------------|
| `make run` | Run the dev server |
| `make migrate-up` | Apply migrations |
| `make migrate-down` | Roll back the last migration |
| `make migrate-create name=add_x` | Create a new migration |
| `make help` | List all commands |

**Frontend (`pnpm`)**

| Command | Description |
|---------|-------------|
| `pnpm dev` | Start the dev server |
| `pnpm build` | Type-check and build for production |
| `pnpm lint` | Run ESLint |
| `pnpm format` | Format code with Prettier |
| `pnpm typecheck` | Run the TypeScript compiler |

## 🗺️ Roadmap

- [ ] Route guards for protected pages on the frontend
- [ ] Pagination and sorting controls in the UI
- [ ] Drag & drop task ordering
- [ ] Dockerfiles for the API and the frontend
- [ ] Unit / integration tests and CI pipeline

## 🤝 Contributing

Issues and pull requests are welcome. Please open an issue first to discuss major changes.

## 👤 Author

**Mohammadreza Nadirkhanloo** — [@MohammadrezaNadirkhanloo](https://github.com/MohammadrezaNadirkhanloo)

---

<div align="center">If you find this project useful, please consider giving it a ⭐</div>
