CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    priority VARCHAR(50),
    status VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT now(),
    due_date TIMESTAMPTZ DEFAULT now()
);

