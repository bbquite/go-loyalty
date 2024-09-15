create table public.account(
    id serial PRIMARY KEY,
    username VARCHAR (50) UNIQUE NOT NULL,
    password VARCHAR (50) NOT NULL,
    balance integer NOT NULL DEFAULT 0,
    created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP
);
CREATE TYPE order_status AS ENUM (
    'NEW',
    'PROCESSING',
    'INVALID',
    'PROCESSED'
);

create table public.purchase(
	id serial primary key,
	account_id integer not null,
	order_num integer not null,
    order_status order_status not null,
	uploaded_at timestamp not null,
	FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE
);

CREATE TYPE transaction_type AS ENUM (
    'IN',
    'OUT'
);

create table public.balance_history(
    id serial primary key,
    account_id integer not null,
    order_id integer not null,
    amount integer not null,
    transaction_type transaction_type not null,
    processed_at TIMESTAMP NOT NULL,
    FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES public.purchase (id) ON DELETE CASCADE
)