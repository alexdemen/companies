CREATE TABLE categories(
                           id BIGSERIAL PRIMARY KEY NOT NULL,
                           name varchar(256) NOT NULL,
                           deleted boolean default false NOT NULL
);

CREATE TABLE buildings(
                          id BIGSERIAL PRIMARY KEY NOT NULL,
                          address varchar(256) NOT NULL,
                          coordinate point NOT NULL,
                          deleted boolean default false NOT NULL
);

CREATE TABLE organizations(
                              id BIGSERIAL PRIMARY KEY NOT NULL,
                              name varchar(256) NOT NULL,
                              building_id bigint NOT NULL,
                              deleted boolean default false NOT NULL,
                              CONSTRAINT fk_building FOREIGN KEY (building_id) REFERENCES buildings(id)
);

CREATE TABLE organization_by_category(
                                         organization_id bigint NOT NULL,
                                         category_id bigint NOT NULL,
                                         CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organizations(id),
                                         CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE category_hierarchy(
                                   category_id bigint NOT NULL,
                                   category_up_id bigint NOT NULL,
                                   parent boolean default false NOT NULL
);