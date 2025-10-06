CREATE TABLE IF NOT EXISTS "subscriptions" (
    id UUID PRIMARY KEY,
    service_name TEXT NOT NULL,
    price INTEGER NOT NULL,
    user_id UUID NOT NULL,
    start_year INTEGER NOT NULL,
    start_month INTEGER NOT NULL
);
