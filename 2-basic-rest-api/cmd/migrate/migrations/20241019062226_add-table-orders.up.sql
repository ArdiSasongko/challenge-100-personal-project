create type order_status as enum ('pending', 'completed', 'cancelled');

create table if not exists orders (
    id serial primary key,
    user_id int not null,
    total numeric(10, 2) not null,
    status order_status not null,
    address text not null,
    created_at timestamp not null default current_timestamp,
    constraint fk_user foreign key (user_id) references users(id)
);