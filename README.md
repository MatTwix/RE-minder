# RE-minder

Opensource Go Fiber + JS React + PostgreSQL project

## made by HSE student Fedorov Matvey

## I feel like I've forgot something

From this moment you will not.

## Database structure

### Users Table

Stores user information.

| Column        | Type      | Constraints                  |
|---------------|-----------|------------------------------|
| id            | SERIAL    | PRIMARY KEY                  |
| username      | TEXT      | UNIQUE, NOT NULL             |
| telegram_id   | BIGINT    | UNIQUE                       |
| github_id     | BIGINT    | UNIQUE, NON NULL             |
| created_at    | TIMESTAMP | DEFAULT NOW()                |
| updated_at    | TIMESTAMP | DEFAULT NOW()                |

---

### Habits Table

Stores user habits with reminders.

| Column       | Type      | Constraints                                                   |
|-------------|------------|---------------------------------------------------------------|
| id          | SERIAL     | PRIMARY KEY                                                   |
| user_id     | INTEGER    | NOT NULL, REFERENCES users(id) ON DELETE CASCADE              |
| name        | TEXT       | NOT NULL                                                      |
| description | TEXT       | DEFAULT ''                                                    |
| frequency   | TEXT       | NOT NULL, CHECK (frequency IN ('daily', 'weekly', 'monthly')) |
| remind_time | TIME       | NOT NULL                                                      |
| timezone    | TEXT       | NOT NULL, DEFAULT 'UTC'                                       |
| created_at  | TIMESTAMP  | DEFAULT NOW()                                                 |
| updated_at  | TIMESTAMP  | DEFAULT NOW()                                                 |
