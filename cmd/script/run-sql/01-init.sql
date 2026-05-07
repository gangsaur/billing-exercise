CREATE TABLE users (
  id bigserial PRIMARY KEY,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE loans (
  id bigserial PRIMARY KEY,
  duration smallint NOT NULL,
  principal_amount integer NOT NULL,
  outstanding_amount integer NOT NULL,
  interest real NOT NULL,
  user_id bigint NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX loans_user_id_idx ON loans (user_id);

CREATE TABLE loan_payments (
  id bigserial PRIMARY KEY,
  period smallint NOT NULL,
  amount integer NOT NULL,
  due_date timestamptz NOT NULL,
  paid_at timestamptz,
  status smallint NOT NULL,
  loan_id bigint NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX loan_payments_loan_id_status_idx ON loan_payments (loan_id, status, due_date);
