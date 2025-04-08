# RE-minder

**Open-source project built with Go Fiber + React JS + PostgreSQL**  
Made by HSE student **Fedorov Matvey**

> _"I feel like I've forgot something."_

From this moment, you will not.

---

## 📦 Database Structure

### 🧑 Users Table

Stores user information.

| Column      | Type      | Constraints                  |
|-------------|-----------|------------------------------|
| id          | SERIAL    | PRIMARY KEY                  |
| username    | TEXT      | UNIQUE, NOT NULL             |
| telegram_id | BIGINT    | UNIQUE                       |
| github_id   | BIGINT    | UNIQUE, NOT NULL             |
| created_at  | TIMESTAMP | DEFAULT NOW()                |
| updated_at  | TIMESTAMP | DEFAULT NOW()                |

---

### ✅ Habits Table

Stores user habits with reminders.

| Column      | Type      | Constraints                                                    |
|-------------|-----------|----------------------------------------------------------------|
| id          | SERIAL    | PRIMARY KEY                                                    |
| user_id     | INTEGER   | NOT NULL, REFERENCES users(id) ON DELETE CASCADE              |
| name        | TEXT      | NOT NULL                                                       |
| description | TEXT      | DEFAULT ''                                                     |
| frequency   | TEXT      | NOT NULL, CHECK (frequency IN ('daily', 'weekly', 'monthly')) |
| remind_time | TIME      | NOT NULL                                                       |
| timezone    | TEXT      | NOT NULL, DEFAULT 'UTC'                                        |
| created_at  | TIMESTAMP | DEFAULT NOW()                                                  |
| updated_at  | TIMESTAMP | DEFAULT NOW()                                                  |

---

## 🌐 API and Auth Endpoints

Base URL: `http://localhost:3000/api`

---

### 1. 👤 Users `/users`

#### Endpoints:

- **Get all**  
  `GET /`

- **Get one by ID**  
  `GET /:id`

- **Create new user**  
  `POST /`  
  **Body (JSON):**

  ```json
  {
    "username": "test username",
    "github_id": 1234567890,
    "telegram_id": 1234567890
  }
  ```

- **Update user by ID**  
  `PUT /:id`  
  **Body (JSON):**

  ```json
  {
    "username": "updated test username",
    "github_id": 12345678,
    "telegram_id": 12345678
  }
  ```

- **Delete user by ID**  
  `DELETE /:id`

---

### 2. 🔁 Habits `/habits`

#### Endpoints

- **Get all habits**  
  `GET /`

- **Get one by ID**  
  `GET /:id`

- **Create new habit**  
  `POST /`  
  **Body (JSON):**

  ```json
  {
    "user_id": 1,
    "name": "test habit",
    "description": "test description",
    "frequency": "weekly",
    "remind_time": "13:00",
    "timezone": "MSK"
  }
  ```

- **Update habit by ID**  
  `PUT /:id`  
  **Body (JSON):**

  ```json
  {
    "user_id": 2,
    "name": "updated test habit",
    "description": "updated test description",
    "frequency": "daily",
    "remind_time": "12:00",
    "timezone": "UTC"
  }
  ```

- **Delete habit by ID**  
  `DELETE /:id`

---
