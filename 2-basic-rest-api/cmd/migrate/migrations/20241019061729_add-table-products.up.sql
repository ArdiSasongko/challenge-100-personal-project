create table if not exists products (
    id serial primary key,
    name varchar(255) not null,
    description text not null,
    image_url varchar(255) not null,
    price integer not null,
    created_at timestamp not null default current_timestamp
);