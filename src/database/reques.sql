CREATE TABLE tender_history (
                                history_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                tender_id UUID NOT NULL REFERENCES tender(id) ON DELETE CASCADE,
                                name VARCHAR(255),
                                description TEXT,
                                service_type VARCHAR(100),
                                status VARCHAR(50),
                                version INTEGER NOT NULL,
                                organization_id uuid,
                                creator_user_id uuid,
                                creator_username VARCHAR(50),
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                UNIQUE (tender_id, version)
);





delete from tender where name = 'тест';
--
--
select  * from  tender_history ;


select  * from employee;
select  * from organization;


select  * from  tender;

INSERT INTO organization_responsible (id, organization_id, user_id)
VALUES (
           uuid_generate_v4(),
           (SELECT id FROM organization WHERE name = 'My Organization'),
           (SELECT id FROM employee WHERE username = 'test_user1')
       );


select  * from tender_history;












create table organization
(
    id          uuid      default uuid_generate_v4() not null
        primary key,
    name        varchar(100)                         not null,
    description text,
    type        organization_type,
    created_at  timestamp default CURRENT_TIMESTAMP,
    updated_at  timestamp default CURRENT_TIMESTAMP
);


create table employee
(
    id         uuid      default uuid_generate_v4() not null
        primary key,
    username   varchar(50)                          not null
        unique,
    first_name varchar(50),
    last_name  varchar(50),
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_at timestamp default CURRENT_TIMESTAMP
);



create table organization_responsible
(
    id              uuid default uuid_generate_v4() not null
        primary key,
    organization_id uuid
        references organization
            on delete cascade,
    user_id         uuid
        references employee
            on delete cascade
);




create table tender
(
    id               uuid      default uuid_generate_v4() not null
        primary key,
    name             varchar(255)                         not null,
    description      text,
    service_type     varchar(100),
    status           varchar(50),
    organization_id  uuid
        references organization
            on delete cascade,
    creator_user_id  uuid
        references employee
            on delete cascade,
    creator_username varchar(50)
        references employee (username),
    version          integer   default 1,
    created_at       timestamp default CURRENT_TIMESTAMP,
    updated_at       timestamp default CURRENT_TIMESTAMP
);


create table tender_history
(
    history_id       uuid      default uuid_generate_v4() not null
        primary key,
    tender_id        uuid                                 not null
        references tender
            on delete cascade,
    name             varchar(255),
    description      text,
    service_type     varchar(100),
    status           varchar(50),
    version          integer                              not null,
    organization_id  uuid,
    creator_user_id  uuid,
    creator_username varchar(50),
    created_at       timestamp default CURRENT_TIMESTAMP,
    updated_at       timestamp default CURRENT_TIMESTAMP,
    unique (tender_id, version)
);



-- Create table for bid (formerly proposal)
create table bid
(
    id              uuid            default uuid_generate_v4() not null
        primary key,
    tender_id       uuid                                       not null
        references tender
            on delete cascade,
    organization_id uuid                                       not null
        references organization
            on delete cascade,
    creator_user_id uuid                                       not null
        references employee
            on delete cascade,
    name            varchar(255)                               not null,
    description     text,
    status          proposal_status default 'CREATED'::proposal_status,
    version         integer         default 1,
    created_at      timestamp       default CURRENT_TIMESTAMP,
    updated_at      timestamp       default CURRENT_TIMESTAMP,
    unique (tender_id, organization_id, version)
);

-- Create table for bid history (formerly proposal_history)
create table bid_history
(
    history_id      uuid      default uuid_generate_v4() not null
        primary key,
    bid_id          uuid                                 not null
        references bid
            on delete cascade,
    tender_id       uuid                                 not null
        references tender
            on delete set null,
    organization_id uuid                                 not null
        references organization
            on delete set null,
    creator_user_id uuid                                 not null
        references employee
            on delete set null,
    name            varchar(255),
    description     text,
    service_type    varchar(100),
    status          varchar(50),
    version         integer                              not null,
    created_at      timestamp default CURRENT_TIMESTAMP,
    updated_at      timestamp default CURRENT_TIMESTAMP,
    unique (bid_id, version)
);







