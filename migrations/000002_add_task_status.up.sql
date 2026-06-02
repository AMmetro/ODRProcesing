ALTER TABLE odrprocesing.tasks
ADD COLUMN status VARCHAR(50) NOT NULL DEFAULT 'unconfirmed';

-- Создаём индекс для быстрого поиска по статусу
CREATE INDEX idx_tasks_status ON odrprocesing.tasks(status);
