CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    username      VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE todo_lists
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description VARCHAR(255)
);

CREATE TABLE users_lists
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE      NOT NULL,
    list_id INTEGER REFERENCES todo_lists (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE todo_items
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description VARCHAR(255)
);

CREATE TABLE lists_items
(
    id      SERIAL PRIMARY KEY,
    list_id INTEGER REFERENCES todo_lists (id) ON DELETE CASCADE NOT NULL,
    item_id INTEGER REFERENCES todo_items (id) ON DELETE CASCADE NOT NULL
);
