CREATE TABLE keys (
    id serial primary key,
    name varchar
);

CREATE TABLE groups (
    id serial primary key,
    name varchar
);

CREATE TABLE partners (
    id serial primary key,
    name varchar,
    code varchar
);

CREATE TABLE groups_to_keys (
    id serial primary key,
    group_id int,
    key_id int,
    FOREIGN KEY(group_id) REFERENCES groups(id),
    FOREIGN KEY(key_id) REFERENCES keys(id)
);

CREATE TABLE partner_mappings (
    id serial primary key,
    partner_id int,
    key_id int,
    FOREIGN KEY(partner_id) REFERENCES keys(id),
    FOREIGN KEY(key_id) REFERENCES keys(id),
    value varchar
);
