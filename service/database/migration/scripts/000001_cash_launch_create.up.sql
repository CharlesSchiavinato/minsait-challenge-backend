CREATE TABLE "cash_launch" (
    "id" bigserial PRIMARY KEY,
    "reference_date" date NOT NULL,
    "type" varchar(1) NOT NULL CHECK ("type" in ('C', 'D')),
    "description" varchar(100) NOT NULL,
    "value" real NOT NULL,
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
