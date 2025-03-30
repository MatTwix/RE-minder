# RE-minder
Opensource Go Fiber + JS React + PostgreSQL project
## made by HSE student Fedorov Matvey

## I feel like I've forgot something..
From this moment you will not.

## Database structure
### Users Table
Stores user information.

| Column        | Type      | Constraints                  |
|--------------|----------|------------------------------|
| id           | SERIAL   | PRIMARY KEY                  |
| username     | TEXT     | UNIQUE, NOT NULL             |
| email        | TEXT     | UNIQUE, NOT NULL             |
| password_hash | TEXT    | NOT NULL                      |
| telegram_id  | BIGINT   | UNIQUE                        |
| created_at   | TIMESTAMP | DEFAULT NOW()               |
| updated_at   | TIMESTAMP | DEFAULT NOW()               |

---

### Habits Table
Stores user habits with reminders.

| Column       | Type      | Constraints                   |
|-------------|----------|-------------------------------|
| id          | SERIAL   | PRIMARY KEY                   |
| user_id     | INTEGER  | REFERENCES users(id) ON DELETE CASCADE |
| name        | TEXT     | NOT NULL                       |
| description | TEXT     | NULLABLE                       |
| frequency   | TEXT     | NOT NULL (e.g., "daily", "weekly") |
| remind_time | TIME     | NOT NULL                       |
| timezone    | TEXT     | NOT NULL                       |
| created_at  | TIMESTAMP | DEFAULT NOW()                |
| updated_at  | TIMESTAMP | DEFAULT NOW()                |