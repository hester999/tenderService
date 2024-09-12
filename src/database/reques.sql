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


select  * from tender_history