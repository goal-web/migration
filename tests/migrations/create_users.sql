CREATE TABLE IF NOT EXISTS users
(
    `id`       INT UNSIGNED AUTO_INCREMENT,
    name       varchar(20),
    created_at timestamp,
    PRIMARY KEY (`id`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4;