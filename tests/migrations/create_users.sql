CREATE TABLE IF NOT EXISTS users
(
    `id`       INT UNSIGNED AUTO_INCREMENT,
    name       varchar(20),
    age        int,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;