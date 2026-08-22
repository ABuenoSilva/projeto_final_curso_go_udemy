CREATE DATABASE IF NOT EXISTS devbook;
USE devbook;

DROP TABLE IF EXISTS seguidores;
DROP TABLE IF EXISTS usuarios;

CREATE TABLE usuarios (
  id int auto_increment primary key,
  nome varchar(50) not null,
  nick varchar(50) not null unique,
  email varchar(100) not null unique,
  senha varchar(100) not null,
  criadoEm timestamp default current_timestamp
) ENGINE=INNODB;

CREATE TABLE seguidores (
  usuario_id int not null,
  seguidor_id int not null,
  PRIMARY KEY (usuario_id, seguidor_id),
  CONSTRAINT fk_usuario FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
  CONSTRAINT fk_seguidor FOREIGN KEY (seguidor_id) REFERENCES usuarios(id) ON DELETE CASCADE
) ENGINE=INNODB;