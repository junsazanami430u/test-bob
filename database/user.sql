-- MySQL syntax
create table users (
    id binary(16) primary key COMMENT 'ユーザーID ULID',
    name varchar(255) not null COMMENT 'ユーザー名',
    email varchar(255) not null unique COMMENT 'メールアドレス',
    password varchar(255) not null COMMENT 'パスワード',
    created_at datetime not null default current_timestamp COMMENT '作成日時',
    updated_at datetime not null default current_timestamp on update current_timestamp COMMENT '更新日時'
);
