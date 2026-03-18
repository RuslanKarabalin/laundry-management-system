-- +goose Up
create type appliance_type as enum (
    'washing_machine',
    'tumble_dryer'
);

create table appliances (
    appliance_id uuid primary key default gen_random_uuid(),
    name varchar not null,
    type appliance_type not null
);

create table reservations (
    reservation_id uuid primary key default gen_random_uuid(),
    appliance_id uuid not null,
    user_id varchar not null,
    start_time timestamptz not null default now(),
    end_time timestamptz not null default now()
);

alter table reservations add constraint fk_reservations_appliances
foreign key (appliance_id) references appliances(appliance_id)
on delete cascade;

-- +goose Down
drop table if exists reservations;
drop table if exists appliances;
