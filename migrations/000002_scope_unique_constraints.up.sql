ALTER TABLE todo.types DROP CONSTRAINT IF EXISTS type_name_unique;
ALTER TABLE todo.boards DROP CONSTRAINT IF EXISTS board_name_unique;
ALTER TABLE todo.columns DROP CONSTRAINT IF EXISTS column_name_unique;

INSERT INTO todo.types (user_id, name, color)
SELECT users.id, 'null', '#FFFFFF'
FROM todo.users AS users
WHERE NOT EXISTS (
  SELECT 1
  FROM todo.types AS types
  WHERE types.user_id = users.id
    AND types.name = 'null'
);

ALTER TABLE todo.types
  ADD CONSTRAINT type_user_name_unique UNIQUE (user_id, name);

ALTER TABLE todo.boards
  ADD CONSTRAINT board_user_name_unique UNIQUE (user_id, name);

ALTER TABLE todo.columns
  ADD CONSTRAINT column_board_name_unique UNIQUE (board_id, name);
