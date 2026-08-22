-- Senha: "senha"
insert into usuarios (nome, nick, email, senha)
values
("bueno","bueno","bueno@test.com","$2a$10$0vXL36AS2zJHL9aGrF1A1u7ELh1n8VVYrxQuZv/NIDH/MhttdjcsC"),
("zanza","zanza","zanza@test.com","$2a$10$0vXL36AS2zJHL9aGrF1A1u7ELh1n8VVYrxQuZv/NIDH/MhttdjcsC"),
("felipe","felipe","felipe@test.com","$2a$10$0vXL36AS2zJHL9aGrF1A1u7ELh1n8VVairYrxQuZv/NIDH/MhttdjcsC");
insert into seguidores (usuario_id, seguidor_id)
values
(1, 2),
(3, 1),
(1, 3);