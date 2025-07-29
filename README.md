# RE-minder

An open-source project built with Go Fiber + PostgreSQL, created by HSE student **Fedorov Matvey**.

> "I feel like I've forgotten something."

From this moment on, you will not.

---

## 🚀 Features

- **Habit Management:** Create, track, and manage your daily, weekly, or monthly habits.
- **Flexible Notifications:** Receive reminders via Telegram, Discord, or Google Calendar.
- **GitHub Authentication:** Securely log in using your GitHub account.
- **Account Linking:** Link multiple third-party services (Discord, Google) to a single profile.
- **REST API:** A well-documented API for managing users, habits, and settings.
- **Internal API:** A secure API for interacting with external services, such as a notification bot.

---

## ⚙️ Setup and Installation

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.20+)
- [Docker](https://www.docker.com/get-started) and Docker Compose

### Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/RE-minder.git
    cd RE-minder
    ```

2. **Configure environment variables:**
    Copy `.env.example` to `.env` and fill in the required values (API keys, database  connection parameters, etc.).

    ```bash
    cp .env.example .env
    ```

3. **Start the database and RabbitMQ:**
    Use Docker Compose to run PostgreSQL and RabbitMQ in the background.

    ```bash
    docker-compose up -d
    ```

4. **Install dependencies:**

    ```bash
    go mod tidy
    ```

5. **Apply migrations:**
    The application automatically applies migrations on startup. Ensure that the database connection is configured correctly.

6. **Run the application:**

    ```bash
    go run main.go
    ```

    The server will be available at `http://localhost:8080`.

---

## 🔔 Notification Architecture

The notification system is designed for flexibility and scalability.

- **Scheduler (`/scheduler`):** Every minute, the scheduler checks the database for habits that need reminders, taking into account users' timezones.
- **Message Queue (`/queue`):** When it's time for a reminder, the scheduler doesn't send it directly. Instead, it places a message in a **RabbitMQ** queue. This decouples the scheduling logic from the sending logic, increasing system reliability.
- **External Worker (Bot):** A separate service (e.g., a Telegram bot) is expected to listen to this queue, retrieve messages, and send notifications to the end-user through the appropriate channel (Telegram, Discord, etc.).

This approach makes it easy to add new notification methods without changing the core application code.

---

## 🔗 OAuth and Account Linking

The project uses OAuth 2.0 for authentication and authorization.

- **Primary Authentication:** Login is handled via **GitHub**. A new user is created upon the first login.
- **Service Linking:** After logging in, a user can link other services like **Discord** or **Google** to their account. This is done using a dedicated endpoint that generates a unique linking URL.
- **Provider Factory (`/oauth/factory.go`):** A Factory pattern is used to manage different OAuth providers, making it easy to add new ones in the future.

---

## 🏗️ Project Structure

```structure
RE-minder/
│
├── config/                 # Application configuration settings
├── database/               # Database connection and initialization logic
├── handlers/               # API endpoint handlers for users and habits
├── middleware/             # Custom middleware for authentication and user management
├── migrations/             # Database schema and migration files
├── models/                 # Data models defining the application's entities
├── oauth/                  # Logic for OAuth2 providers (Github, Discord, Google)
├── queue/                  # RabbitMQ message queue logic
├── routes/                 # API route definitions and setup
├── scheduler/              # Job scheduler for sending notifications
├── services/               # Business logic and data access layer
├── main.go                 # Main application entry point
└── go.mod                  # Go module dependencies
```

---

## 📦 Database Structure

### 🧑 Users Table

| Column      | Type      | Constraints                  |
|-------------|-----------|------------------------------|
| id          | SERIAL    | PRIMARY KEY                  |
| username    | TEXT      | UNIQUE, NOT NULL             |
| telegram_id | BIGINT    | UNIQUE                       |
| github_id   | BIGINT    | UNIQUE, NOT NULL             |
| is_admin    | BOOLEAN   | DEFAULT FALSE                |
| created_at  | TIMESTAMP | DEFAULT NOW()                |
| updated_at  | TIMESTAMP | DEFAULT NOW()                |

### ✅ Habits Table

| Column      | Type      | Constraints                                                    |
|-------------|-----------|----------------------------------------------------------------|
| id          | SERIAL    | PRIMARY KEY                                                    |
| user_id     | INTEGER   | NOT NULL, REFERENCES users(id) ON DELETE CASCADE              |
| name        | TEXT      | NOT NULL                                                       |
| description | TEXT      | DEFAULT ''                                                     |
| frequency   | TEXT      | NOT NULL, CHECK (frequency IN ('daily', 'weekly', 'monthly')) |
| remind_time | TIME      | NOT NULL                                                       |
| timezone    | TEXT      | NOT NULL, DEFAULT 'UTC'                                        |
| start_date  | DATE      | NOT NULL, DEFAULT CURRENT_DATE                                 |
| created_at  | TIMESTAMP | DEFAULT NOW()                                                  |
| updated_at  | TIMESTAMP | DEFAULT NOW()                                                  |

### ⚙️ Notification Settings Table

| Column                | Type      | Constraints                                      |
|-----------------------|-----------|--------------------------------------------------|
| id                    | SERIAL    | PRIMARY KEY                                      |
| user_id               | INTEGER   | NOT NULL, UNIQUE, REFERENCES users(id) ON DELETE CASCADE |
| telegram_notification | BOOLEAN   | NOT NULL, DEFAULT FALSE                          |
| discord_notification  | BOOLEAN   | NOT NULL, DEFAULT FALSE                          |
| google_notification       | BOOLEAN   | NOT NULL, DEFAULT FALSE                          |
| created_at            | TIMESTAMP | NOT NULL, DEFAULT NOW()                          |
| updated_at            | TIMESTAMP | NOT NULL, DEFAULT NOW()                          |

---

## 🌐 API Endpoints

Base URL: `/api`

### User API Authentication

All API endpoints require a `Bearer` token in the `Authorization` header.

---

### 🔐 Auth

#### `GET /auth/github`

Redirects to the GitHub OAuth page for authentication.

#### `GET /auth/github/callback`

Handles the callback from GitHub after authentication. On success, it creates or updates the user and returns a JWT token.

- **Success Response (200)**

  ```json
  {
    "token": "your_jwt_token_here",
    "user": "github_username"
  }
  ```

---

### 🤖 Bot Linking

These endpoints are used to link third-party accounts (e.g., Discord, Google) to a user's RE-minder profile.

- **Access:** Authenticated users.

#### `GET /auth/bot/:platform`

Redirects the user to the OAuth authorization page of the selected platform (`discord`, `google`, etc.).

- **URL Parameters:**
  - `platform` (string, required): The platform to authorize with. Supported values: `discord`, `google`.

#### `GET /auth/bot/:platform/callback`

Handles the callback from the OAuth provider after successful authorization. Links the external account to the current user.

- **Success Response (200)**
  
  ```json
  {
    "message": "Account linked successfully"
  }
  ```

- **Error Response (400/500)**
  
  ```json
  {
    "error": "Failed to link account"
  }
  ```

---

### 👤 Me

#### `GET /api/me`

Retrieves the profile of the currently authenticated user.

- **Access:** Authenticated users.
- **Success Response (200)**

  ```json
    {
      "id": 1,
      "username": "testuser",
      "telegram_id": 123456789,
      "github_id": 987654321,
      "is_admin": false,
      "created_at": "2025-07-07T12:00:00Z",
      "updated_at": "2025-07-07T12:00:00Z"
    }
  ```

---

### 🧑 Users

#### `GET /api/users`

Retrieves a list of all users.

- **Access:** Admin only.
- **Success Response (200)**

  ```json
    [
      {
        "id": 1,
        "username": "admin",
        "github_id": 111,
        "is_admin": true,
        ...
      },
      {
        "id": 2,
        "username": "user",
        "github_id": 222,
        "is_admin": false,
        ...
      }
    ]
  ```

#### `GET /api/users/:id`

Retrieves a single user by their ID.

- **Access:** Admin or the user themselves.
- **Success Response (200)**

  ```json
    {
      "id": 2,
      "username": "user",
      ...
    }
  ```

#### `POST /api/users`

Creates a new user.

- **Access:** Admin only.
- **Request Body**

  ```json
    {
      "username": "newuser",
      "github_id": 333,
      "telegram_id": 12345
    }
  ```

- **Success Response (200)**
  
  ```json
    {
      "id": 3,
      "username": "newuser",
      "github_id": 333,
      "telegram_id": 12345,
      "is_admin": false,
      ...
    }
  ```

#### `PUT /api/users/:id`

Updates a user's information.

- **Access:** Admin only.
- **Request Body**

  ```json
    {
      "username": "updateduser",
      "github_id": 333,
      "telegram_id": 54321
    }
  ```

- **Success Response (200)**
  
  ```json
    {
      "id": 3,
      "username": "updateduser",
      ...
    }
  ```

#### `PATCH /api/users/:id/telegram_id`

Updates a user's Telegram ID.

- **Access:** Admin or the user themselves.
- **Request Body**
  
  ```json
    {
      "telegram_id": 98765
    }
  ```

- **Success Response (200)**
  
  ```json
    {
      "id": 2,
      "telegram_id": 98765,
      ...
    }
  ```

#### `PATCH /api/users/:id/is_admin`

Toggles a user's admin status.

- **Access:** Admin only.
- **Success Response (200)**
  
  ```json
    {
      "id": 2,
      "is_admin": true,
      ...
    }
  ```

#### `DELETE /api/users/:id`

Deletes a user.

- **Access:** Admin or the user themselves.
- **Success Response (200)**
  
  ```json
    {
      "message": "User deleted successfully"
    }
  ```

---

### ✅ Habits

#### `GET /api/habits`

Retrieves a list of all habits from all users.

- **Access:** Admin only.
- **Success Response (200)**
  
  ```json
    [
      {
        "id": 1,
        "user_id": 1,
        "name": "Read a book",
        ...
      },
      {
        "id": 2,
        "user_id": 2,
        "name": "Go for a run",
        ...
      }
    ]
  ```

#### `GET /api/habits/user/:id`

Retrieves all habits for a specific user.

- **Access:** Admin or the user themselves.
- **Success Response (200)**
  
  ```json
    [
      {
        "id": 2,
        "user_id": 2,
        "name": "Go for a run",
        ...
      }
    ]
  ```

#### `GET /api/habits/:id`

Retrieves a single habit by its ID.

- **Access:** Any authenticated user can view any habit.
- **Success Response (200)**
  
  ```json
    {
      "id": 1,
      "user_id": 1,
      "name": "Read a book",
      "description": "Read 20 pages of a non-fiction book.",
      "frequency": "daily",
      "remind_time": "21:00:00",
      "timezone": "UTC",
      "start_date": "2025-07-01T00:00:00Z",
      ...
    }
  ```

#### `POST /api/habits`

Creates a new habit for the authenticated user.

- **Access:** Authenticated users.
- **Request Body**
  
  ```json
    {
      "name": "Morning Jogging",
      "description": "Jog for 30 minutes.",
      "frequency": "daily",
      "remind_time": "07:00",
      "timezone": "Europe/Moscow",
      "start_date": "2025-08-01T00:00:00Z"
    }
  ```

- **Success Response (200)**
  
  ```json
    {
      "id": 3,
      "user_id": 2,
      "name": "Morning Jogging",
      ...
    }
  ```

#### `PUT /api/habits/:id`

Updates a habit.

- **Access:** The user who created the habit or an admin.
- **Request Body**
  
  ```json
    {
      "name": "Evening Jogging",
      "description": "Jog for 45 minutes.",
      "frequency": "weekly",
      "remind_time": "19:00",
      "timezone": "UTC"
    }
  ```

- **Success Response (200)**
  
  ```json
    {
      "id": 3,
      "name": "Evening Jogging",
      ...
    }
  ```

#### `DELETE /api/habits/:id`

Deletes a habit.

- **Access:** The user who created the habit or an admin.
- **Success Response (200)**
  
  ```json
    {
      "message": "Habit succesfully deleted"
    }
  ```

---

### ⚙️ Notification Settings

#### `GET /api/notifications_settings/user/:id`

Retrieves the notification settings for a specific user.

- **Access:** Admin or the user themselves.
- **Success Response (200)**

  ```json
  {
    "id": 1,
    "user_id": 1,
    "telegram_notification": true,
    "discord_notification": false,
    "google_notification": false,
    "created_at": "2025-07-22T10:00:00Z",
    "updated_at": "2025-07-22T10:00:00Z"
  }
  ```

#### `PUT /api/notifications_settings/user/:id`

Updates the notification settings for a specific user.

- **Access:** Admin or the user themselves.
- **Request Body**

  ```json
  {
    "telegram_notification": true,
    "discord_notification": true,
    "google_notification": false
  }
  ```

- **Success Response (200)**

  ```json
  {
    "id": 1,
    "user_id": 1,
    "telegram_notification": true,
    "discord_notification": true,
    "google_notification": false,
    ...
  }
  ```

---

## 🤖 Internal API

Endpoints intended for service-to-service communication (e.g., with a notification bot).

### Internal API Authentication

All Internal API endpoints require an `ApiKey` token in the `Authorization` header. The `INTERNAL_API_KEY` must be set in the environment variables.

**Example Header:**
`Authorization: ApiKey your_super_secret_and_long_random_string`

---

### ⚙️ Notification Settings (Internal)

#### `GET /internal/notifications_settings/user/:id`

Retrieves the notification settings for a specific user. This endpoint is intended for use by internal services like a notification bot to check a user's delivery preferences.

- **Access:** Any authenticated internal service.
- **Success Response (200)**

  ```json
  {
    "id": 1,
    "user_id": 1,
    "telegram_notification": true,
    "discord_notification": false,
    "google_notification": false,
    ...
  }
  ```

## 🛠️ Technologies Used

- **Backend:**
  - [Go](https://golang.org/)
  - [Fiber](https://gofiber.io/) - Web framework
  - [PostgreSQL](https://www.postgresql.org/) - Database
  - [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
  - [JWT](https://jwt.io/) - For authentication
  - [RabbitMQ](https://www.rabbitmq.com/) - Message broker
- **Other:**
  - [Docker](https://www.docker.com/) - For containerization
