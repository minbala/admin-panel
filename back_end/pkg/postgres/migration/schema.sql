
create type user_role_type as enum ('admin','normal user');

Create Table if not exists users (
                                     id bigserial primary key,
                                     created_at timestamp default current_timestamp,
                                     updated_at timestamp default current_timestamp,
                                     email varchar(255) unique not null ,
                                     name varchar(255) not null ,
                                     password text not null ,
                                     user_role user_role_type not null
);


Create Table if not exists user_logs (
                                         id bigserial primary key,
                                         created_at timestamp default current_timestamp,
                                         updated_at timestamp default current_timestamp,
                                         user_id bigint,
                                         event varchar(20) not null ,
                                         request_url varchar(255) not null ,
                                         data jsonb ,
                                         status int not null,
                                         error_message text,
                                         FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);




Create Table if not exists sessions (
                                        id bigserial primary key,
                                        created_at timestamp default current_timestamp,
                                        updated_at timestamp default current_timestamp,
                                        access_token text not null ,
                                        user_id bigint not null,
                                        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE

);



