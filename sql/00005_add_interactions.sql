-- auto-generated definition
create table users_directors_interactions
(
    user_id    varchar not null,
    entity_id  integer not null,
    created_at timestamp        default CURRENT_TIMESTAMP,
    alpha      double precision default 1.0,
    beta       double precision default 1.0,
    constraint users_directors_interactions_pk
        unique (user_id, entity_id)
);

alter table users_directors_interactions
    owner to "mabrc-admin";

create index users_directors_interactions_user_id_index
    on users_directors_interactions (user_id);

-- auto-generated definition
create table users_genres_interactions
(
    user_id    varchar not null,
    entity_id  integer not null,
    created_at timestamp        default CURRENT_TIMESTAMP,
    alpha      double precision default 1.0,
    beta       double precision default 1.0,
    constraint users_genres_interactions_pk
        unique (user_id, entity_id)
);

alter table users_genres_interactions
    owner to "mabrc-admin";

create index users_genres_interactions_user_id_index
    on users_genres_interactions (user_id);

-- auto-generated definition
create table users_movies_interactions
(
    user_id    varchar not null,
    entity_id  integer not null,
    created_at timestamp        default CURRENT_TIMESTAMP,
    alpha      double precision default 1.0,
    beta       double precision default 1.0,
    constraint users_movies_interactions_pk
        unique (user_id, entity_id)
);

alter table users_movies_interactions
    owner to "mabrc-admin";

create index users_movies_interactions_user_id_index
    on users_movies_interactions (user_id);

create table recommendations
(
    user_id     varchar not null,
    movie_id    integer,
    entity_type varchar not null,
    created_at  timestamp
);

create index recommendations_user_id_index
    on recommendations (user_id);


