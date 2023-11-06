-- migrate:up
CREATE TABLE IF NOT EXISTS tasks_tags (
    task_id INT NOT NULL,
    tag_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (task_id, tag_id),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- migrate:down
DROP TABLE IF EXISTS tasks_tags;

