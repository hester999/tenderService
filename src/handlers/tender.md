```
id  - id тендера в бд
name - название тендера
description - описание тендера
status -  статус тендера
serviceType - тип сервиса 
version - версия тендера
createdAt - дата создания 

```



```
CREATE TABLE tender (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(255) NOT NULL,
                        description TEXT,
                        service_type VARCHAR(100),
                        status VARCHAR(50),
                        organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
                        creator_user_id INT REFERENCES employee(id) ON DELETE CASCADE,
                        creator_username VARCHAR(50) REFERENCES employee(username),
                        version INTEGER DEFAULT 1,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);




INSERT INTO organization (name, description, type) VALUES ('My Organization', 'Description', 'LLC');
INSERT INTO employee (username, first_name, last_name) VALUES ('test_user', 'Test', 'User');

delete from tender where id=4;
select * from  tender;

```
