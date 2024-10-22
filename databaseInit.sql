create table public.account(
    id serial PRIMARY KEY,
    username VARCHAR (255) UNIQUE NOT NULL,
    password VARCHAR (255) NOT NULL,
    created_on TIMESTAMP NOT NULL default CURRENT_TIMESTAMP
);
CREATE TYPE purchase_status AS ENUM (
    'NEW',
    'PROCESSING',
    'INVALID',
    'PROCESSED'
);

create table public.purchase(
	id serial primary key,
	account_id integer not null,
	purchase_num VARCHAR(255) not null unique,
    purchase_status purchase_status not null,
	uploaded_at timestamp not null default CURRENT_TIMESTAMP,
	FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE
);

create table public.balance(
    id serial primary key,
    account_id integer not null,
    amount float NOT NULL DEFAULT 0,
    withdrawn integer NOT NULL DEFAULT 0,
    FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE
);

CREATE TYPE transaction_type AS ENUM (
    'IN',
    'OUT'
);

create table public.balance_history(
    id serial primary key,
    account_id integer not null,
    purchase_id integer not null,
    amount integer not null,
    transaction_type transaction_type not null,
    processed_at TIMESTAMP NOT NULL default CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE,
    FOREIGN KEY (purchase_id) REFERENCES public.purchase (id) ON DELETE CASCADE
);