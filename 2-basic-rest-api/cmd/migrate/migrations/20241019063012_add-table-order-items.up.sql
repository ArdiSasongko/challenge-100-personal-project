create table if not exists order_items(
    id serial primary key,
    order_id int not null,
    product_id int not null,
    quantity int not null,
    price numeric(10, 2) not null,
    created_at timestamp not null default current_timestamp,
    constraint fk_order foreign key (order_id) references orders(id),
    constraint fk_product foreign key (product_id) references products(id)
);