BEGIN;

create table t_bark_device
(
    device_key   uuid         not null
        primary key,
    device_token varchar(256) not null
);

create table t_device
(
    id         serial
        primary key,
    user_id    integer     not null,
    device_id  varchar(64) not null
        unique,
    type       varchar(32),
    is_clip    smallint,
    name       varchar(32),
    created_at timestamp   not null,
    updated_at timestamp
);

create table t_message
(
    id       serial
        primary key,
    user_id  integer     not null,
    text     text,
    type     varchar(32),
    note     varchar(64),
    push_key varchar(64) not null,
    url      varchar(255),
    send_at  timestamp   not null
);

create table t_push_key
(
    id         serial
        primary key,
    user_id    integer     not null,
    name       varchar(32),
    key        varchar(64) not null
        unique,
    created_at timestamp   not null,
    updated_at timestamp
);

create table t_user
(
    id         serial
        primary key,
    apple_id   varchar(255) not null
        unique,
    email      varchar(127),
    name       varchar(127),
    uuid       uuid         not null
        unique,
    created_at timestamp    not null,
    updated_at timestamp
);

COMMIT;
