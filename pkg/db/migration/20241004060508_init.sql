-- +goose Up
-- +goose StatementBegin

create table backlite_tasks
(
    id               text              not null primary key,
    created_at       integer           not null,
    queue            text              not null,
    task             blob              not null,
    wait_until       integer,
    claimed_at       integer,
    last_executed_at integer,
    attempts         integer default 0 not null
)
    strict;

create index backlite_tasks_wait_until
    on backlite_tasks (wait_until) where wait_until IS NOT NULL;

create table backlite_tasks_completed
(
    id                  text    not null primary key,
    created_at          integer not null,
    queue               text    not null,
    last_executed_at    integer,
    attempts            integer not null,
    last_duration_micro integer,
    succeeded           integer,
    task                blob,
    expires_at          integer,
    error               text
)
    strict;

create table users
(
    id         integer            not null primary key autoincrement,
    name       text               not null,
    email      text               not null,
    password   text               not null,
    verified   bool default false not null,
    created_at datetime           not null
);

create table password_tokens
(
    id                  integer  not null primary key autoincrement,
    hash                text     not null,
    created_at          datetime not null,
    password_token_user integer  not null
        constraint password_tokens_users_user
            references users
);

create unique index users_email_key
    on users (email);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table backlite_tasks;
drop table backlite_tasks_completed;
drop table password_tokens;
drop table users;
-- +goose StatementEnd
