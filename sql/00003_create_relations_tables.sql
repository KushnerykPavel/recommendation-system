create table movies_actors_relation
(
    source_id      integer not null,
    destination_id integer not null,
    constraint movies_actors_relation_pk
        unique (source_id, destination_id)
);

alter table movies_actors_relation
    owner to "mabrc-admin";

create table movies_genres_relation
(
    source_id      integer not null,
    destination_id integer not null,
    constraint movies_genres_relation_pk
        unique (source_id, destination_id)
);

alter table movies_genres_relation
    owner to "mabrc-admin";

create table movies_directors_relation
(
    source_id      integer not null,
    destination_id integer not null,
    constraint movies_directors_relation_pk
        unique (source_id, destination_id)
);

alter table movies_directors_relation
    owner to "mabrc-admin";

