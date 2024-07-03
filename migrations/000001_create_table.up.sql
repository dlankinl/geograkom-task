create table if not exists routes(
    route_id int primary key,
    route_name varchar(128) not null,
    load float not null,
    cargo_type varchar(64) not null,
    is_actual bool default true
);
