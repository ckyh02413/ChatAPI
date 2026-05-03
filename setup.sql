CREATE TABLE users (
    username TEXT PRIMARY KEY,
    mail TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE chatrooms (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    creator TEXT NOT NULL REFERENCES users(username),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chatroom_id INTEGER NOT NULL REFERENCES chatrooms(id) ON DELETE CASCADE,
    creator TEXT NOT NULL REFERENCES users(username),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
