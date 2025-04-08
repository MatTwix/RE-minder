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

| Column      | Type      | Constraints                                                    |
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

## API and Auth endpionts

### API `http://localhost:3000/api`

1. Users `+ /users`
*Get all `/`,  metod: GET
*Get one by id `/:id`,  metod: GET
*Create one `/`,  method: POST,  body (json):  

```structure
   {
    "username": "test username",
    "github_id": 1234567890,
    "telegram_id": 1234567890
    }
```

*Update one by id `/:id`,  method: PUT, body (json):  

```structure
    "username": "updated test username",
    "github_id": 12345678,
    "telegram_id": 12345678
```

*Delete one by id `/:id`,  method: DELETE

1. Habits `+ /habits`
*Get all `/`,  metod: GET
*Get one by id `/:id`,  metod: GET
*Create one `/`,  method: POST,  body (json):  

```structure
{
    "user_id": 1,
    "name": "test habit",
    "description": "test description",
    "frequency": "weekly",
    "remind_time": "13:00",
    "timezone": "MSK"
}```

*Update one by id `/:id`,  method: PUT, body (json):  

```sructure
    "user_id": 2,
    "name": "updated test habit",
    "description": "updated test description",
    "frequency": "daily",
    "remind_time": "12:00",
    "timezone": "UPD"
```

*Delete one by id `/:id`,  method: DELETE
