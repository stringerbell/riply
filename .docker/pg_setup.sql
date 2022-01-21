create schema shortlinks;

create user app;
alter role app with password 'password';

alter database shortlinks set search_path to shortlinks;

create table shortlinks
(
    id            bigserial
        constraint shortlinks_pk primary key,
    link          text,
    custom_suffix text,
    hash          text,
    created_at    timestamp with time zone default now() not null
);
grant delete, insert, references, select, trigger, truncate, update on shortlinks to app;

create unique index shortlinks_hash_uindex on shortlinks (hash);

create unique index shortlinks_link_uindex on shortlinks (link);

create unique index shortlinks_custom_suffix_uindex on shortlinks (custom_suffix);

create table stats
(
  id           bigserial
    constraint stats_pk primary key,
  shortlink_id bigint not null,
  created_at   timestamp with time zone default now()
);
grant delete, insert, references, select, trigger, truncate, update on stats to app;

grant all privileges on schema shortlinks to app;
grant usage, select on all sequences in schema shortlinks to app;
