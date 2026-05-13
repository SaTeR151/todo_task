CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

--
-- Todo
--

CREATE SCHEMA todo;

--
-- Todo.Users
--

CREATE TABLE todo.users (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  login text NOT NULL,
  password text NOT NULL,
  refresh_token text,

  CONSTRAINT login_unique UNIQUE (login),
  CONSTRAINT users_pkey PRIMARY KEY (id)
);

-- 
-- Todo.Types 
--

CREATE TABLE todo.types (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  user_id uuid NOT NULL,
  name text NOT NULL,
  color text NOT NULL,

  FOREIGN KEY (user_id) REFERENCES todo.users (id) ON DELETE CASCADE,
  CONSTRAINT types_pkey PRIMARY KEY (id),
  CONSTRAINT type_name_unique UNIQUE (name)
);

--
-- Todo.Boards
--

CREATE TABLE todo.boards (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  user_id uuid NOT NULL,
  name text NOT NULL,

  FOREIGN KEY (user_id) REFERENCES todo.users (id) ON DELETE CASCADE,
  CONSTRAINT boards_pkey PRIMARY KEY (id),
  CONSTRAINT board_name_unique UNIQUE (name)
);

--
-- Todo.Columns
--

CREATE TABLE todo.columns (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  board_id uuid NOT NULL,
  name text NOT NULL,
  order_number integer NOT NULL,

  FOREIGN KEY (board_id) REFERENCES todo.boards (id) ON DELETE CASCADE,
  CONSTRAINT columns_pkey PRIMARY KEY (id),
  CONSTRAINT column_name_unique UNIQUE (name)
);

--
-- Todo.Tasks
--

CREATE TABLE todo.tasks (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  type_id uuid NOT NULL,
  column_id uuid NOT NULL,
  label text NOT NULL,
  description text NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),

  FOREIGN KEY (column_id) REFERENCES todo.columns (id) ON DELETE CASCADE,
  FOREIGN KEY (type_id) REFERENCES todo.types (id),
  CONSTRAINT tasks_pkey PRIMARY KEY (id)
);

-- 
-- Todo.MoveEvents
--

CREATE TABLE todo.move_events (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  task_id uuid NOT NULL,
  from_column_id uuid NOT NULL,
  to_column_id uuid NOT NULL,
  timestamp timestamp with time zone NOT NULL DEFAULT now(),

  FOREIGN KEY (task_id) REFERENCES todo.tasks (id) ON DELETE CASCADE,
  CONSTRAINT move_events_pkey PRIMARY KEY (id)
);


