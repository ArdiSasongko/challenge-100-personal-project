create table if not exists users (
    id serial primary key,
    username varchar(255) not null,
    password varchar(255) not null,
    name varchar(255) not null,
    email varchar(255) not null unique,
    created_at timestamp not null default current_timestamp
);