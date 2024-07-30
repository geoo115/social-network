# social-network


backend/
├── cmd/
│   ├── main.go
├── pkg/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── auth.go
│   │   │   ├── profile.go
│   │   │   ├── posts.go
│   │   │   ├── groups.go
│   │   │   ├── chat.go
│   │   │   ├── notifications.go
│   │   ├── middlewares/
│   │   │   ├── auth.go
│   │   ├── router.go
│   ├── db/
│   │   ├── migrations/
│   │   │   ├── 000001_create_users_table.up.sql
│   │   │   ├── 000001_create_users_table.down.sql
│   │   │   ├── ...
│   │   ├── sqlite/
│   │   │   ├── sqlite.go
│   ├── models/
│   │   ├── user.go
│   │   ├── post.go
│   │   ├── group.go
│   │   ├── chat.go
│   │   ├── notification.go
│   ├── services/
│   │   ├── auth.go
│   │   ├── profile.go
│   │   ├── posts.go
│   │   ├── groups.go
│   │   ├── chat.go
│   │   ├── notifications.go
├── Dockerfile
└── .env
