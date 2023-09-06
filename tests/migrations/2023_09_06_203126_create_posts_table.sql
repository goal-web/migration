CREATE TABLE IF NOT EXISTS posts
(
    `id`
    INT
    UNSIGNED
    AUTO_INCREMENT,
    created_at
    timestamp,
    updated_at
    timestamp,
    PRIMARY
    KEY
(
    `id`
) ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;