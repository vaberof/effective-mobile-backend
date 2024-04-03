CREATE TABLE IF NOT EXISTS cars
(
    id      SERIAL PRIMARY KEY,
    reg_num VARCHAR(9) UNIQUE NOT NULL,
    mark    VARCHAR(30),
    model   VARCHAR(50),
    year    INT,
    owner   JSONB
);

CREATE INDEX IF NOT EXISTS reg_num_idx ON cars (reg_num);

--Тестовые данные--
INSERT INTO cars (reg_num, mark, model, year, owner) VALUES ('X121XX121', 'Mark1', 'Model1', 2001, '{"name": {"Valid": true, "String": "Name1"}, "surname": {"Valid": true, "String": "Surname1"}, "patronymic": {"Valid": false, "String": ""}}');
INSERT INTO cars (reg_num, mark, model, year, owner) VALUES ('X122XX122', 'Mark2', 'Model2', 2002, '{"name": {"Valid": true, "String": "Name2"}, "surname": {"Valid": true, "String": "Surname2"}, "patronymic": {"Valid": true, "String": "Patronymic2"}}');
INSERT INTO cars (reg_num, mark, model, year, owner) VALUES ('X123XX123', 'Mark3', 'Model3', 2003, '{"name": {"Valid": true, "String": "Name3"}, "surname": {"Valid": true, "String": "Surname3"}, "patronymic": {"Valid": true, "String": "Patronymic3"}}');
INSERT INTO cars (reg_num, mark, model, year, owner) VALUES ('X124XX124', 'Mark4', 'Model4', 2004, '{"name": {"Valid": true, "String": "Name4"}, "surname": {"Valid": true, "String": "Surname4"}, "patronymic": {"Valid": true, "String": "Patronymic4"}}');
INSERT INTO cars (reg_num, mark, model, year, owner) VALUES ('X125XX125', 'Mark5', 'Model5', 2005, '{"name": {"Valid": true, "String": "Name5"}, "surname": {"Valid": true, "String": "Surname5"}, "patronymic": {"Valid": true, "String": "Patronymic5"}}');