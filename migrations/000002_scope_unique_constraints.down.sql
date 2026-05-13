ALTER TABLE todo.columns DROP CONSTRAINT IF EXISTS column_board_name_unique;
ALTER TABLE todo.boards DROP CONSTRAINT IF EXISTS board_user_name_unique;
ALTER TABLE todo.types DROP CONSTRAINT IF EXISTS type_user_name_unique;

ALTER TABLE todo.types
  ADD CONSTRAINT type_name_unique UNIQUE (name);

ALTER TABLE todo.boards
  ADD CONSTRAINT board_name_unique UNIQUE (name);

ALTER TABLE todo.columns
  ADD CONSTRAINT column_name_unique UNIQUE (name);
