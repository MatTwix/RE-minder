# RE-minder

Open-source project built with Go Fiber + PostgreSQL, made by HSE student **Fedorov Matvey**.

> "I feel like I've forgot something."

From this moment, you will not.

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
├── routes/                 # API route definitions and setup
├── services/               # Business logic and data access layer for users and habits
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

---

## 🌐 API Endpoints

Base URL: `/api`

### Authentication

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

## 🛠️ Technologies Used

- **Backend:**
  - [Go](https://golang.org/)
  - [Fiber](https://gofiber.io/) - Web framework
  - [PostgreSQL](https://www.postgresql.org/) - Database
  - [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
  - [JWT](https://jwt.io/) - For authentication
- **Other:**
  - [Docker](https://www.docker.com/) - For containerization
