-- Удаляем индекс
DROP INDEX IF EXISTS idx_tasks_status;

-- Удаляем колонку status
ALTER TABLE odrprocesing.tasks
DROP COLUMN status;
