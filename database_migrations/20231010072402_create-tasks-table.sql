-- migrate:up
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'on_going',
    "order" INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (user_id) REFERENCES users(id) 
);

CREATE INDEX tasks_status_idx
ON tasks (status);

-- migrate:down
DROP INDEX IF EXISTS tasks_status_idx;

DROP TABLE IF EXISTS tasks;

