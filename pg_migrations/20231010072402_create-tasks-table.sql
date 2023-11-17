-- migrate:up
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'on_going',
    "order" INT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (user_id) REFERENCES users(id) 
);

CREATE INDEX tasks_status_idx
ON tasks (status);

CREATE INDEX tasks_order_idx
ON tasks ("order");

-- migrate:down
DROP INDEX IF EXISTS tasks_order_idx;

DROP INDEX IF EXISTS tasks_status_idx;

DROP TABLE IF EXISTS tasks;

