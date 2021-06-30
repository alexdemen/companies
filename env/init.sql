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

CREATE TABLE organization_phone_numbers(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    number varchar(20) NOT NULL,
    organization_id bigint NOT NULL,
    deleted boolean default false  NOT NULL,
    CONSTRAINT fk_organization_phone FOREIGN KEY (organization_id) REFERENCES organizations(id)
);

INSERT INTO categories (name) VALUES ('Спорт');
INSERT INTO categories (name) VALUES ('Фитнес');
INSERT INTO category_hierarchy (category_id, category_up_id, parent) VALUES (
    (SELECT id FROM categories WHERE name = 'Фитнес'),
    (SELECT id FROM categories WHERE name = 'Спорт'),
    true
                                                                            );
INSERT INTO categories (name) VALUES ('Тренажерные залы');
INSERT INTO category_hierarchy (category_id, category_up_id, parent) VALUES (
    (SELECT id FROM categories WHERE name = 'Тренажерные залы'),
    (SELECT id FROM categories WHERE name = 'Спорт'),
    true
                                                                            );

INSERT INTO categories (name) VALUES ('Продукты');
INSERT INTO categories (name) VALUES ('Мясо');
INSERT INTO category_hierarchy (category_id, category_up_id, parent) VALUES (
    (SELECT id FROM categories WHERE name = 'Мясо'),
    (SELECT id FROM categories WHERE name = 'Продукты'),
    true
);
INSERT INTO categories (name) VALUES ('Морепродукты');
INSERT INTO category_hierarchy (category_id, category_up_id, parent) VALUES (
    (SELECT id FROM categories WHERE name = 'Морепродукты'),
    (SELECT id FROM categories WHERE name = 'Продукты'),
    true
                                                                            );
INSERT INTO categories (name) VALUES ('Овощи');
INSERT INTO category_hierarchy (category_id, category_up_id, parent) VALUES (
    (SELECT id FROM categories WHERE name = 'Овощи'),
    (SELECT id FROM categories WHERE name = 'Продукты'),
    true
                                                                            );

INSERT INTO categories (name) VALUES ('Медицина');
INSERT INTO categories (name) VALUES ('Аптеки');
INSERT INTO category_hierarchy (category_id, category_up_id, parent) VALUES (
    (SELECT id FROM categories WHERE name = 'Аптеки'),
    (SELECT id FROM categories WHERE name = 'Медицина'),
    true
                                                                            );
INSERT INTO categories (name) VALUES ('Больницы');
INSERT INTO category_hierarchy (category_id, category_up_id, parent) VALUES (
    (SELECT id FROM categories WHERE name = 'Больницы'),
    (SELECT id FROM categories WHERE name = 'Медицина'),
    true
                                                                            );

CREATE INDEX org_cat_idx ON organization_by_category (category_id);
CREATE INDEX org_cat_org_id_idx ON organization_by_category(organization_id);
CREATE INDEX org_phone_idx ON organization_phone_numbers (organization_id);
