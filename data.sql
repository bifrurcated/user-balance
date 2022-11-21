DROP TABLE IF EXISTS balance CASCADE;
DROP TABLE IF EXISTS reserve CASCADE;

CREATE TABLE balance
(
    id      BIGSERIAL PRIMARY KEY,
    user_id BIGINT                   NOT NULL,
    amount  DECIMAL(13, 4) DEFAULT 0 NOT NULL
);
CREATE TABLE reserve
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT         NOT NULL,
    service_id BIGINT         NOT NULL,
    order_id   BIGINT         NOT NULL,
    amount     DECIMAL(13, 4) NOT NULL
);

-- CREATE TABLE usr
-- (
--     id     BIGSERIAL PRIMARY KEY,
--     name   VARCHAR(64)              NOT NULL,
--     amount DECIMAL(13, 4) DEFAULT 0 NOT NULL
-- );
-- CREATE TABLE history_operations
-- (
--     id               BIGSERIAL PRIMARY KEY,
--     sender_user_id   BIGINT                                       NOT NULL,
--     receiver_user_id BIGINT                                       NOT NULL,
--     amount           DECIMAL(13, 4)                               NOT NULL,
--     type             VARCHAR(32)                                  NOT NULL,
--     datetime         TIMESTAMP DEFAULT (now() AT TIME ZONE 'utc') NOT NULL,
--
--
--     CONSTRAINT sender_user_fk FOREIGN KEY (sender_user_id) REFERENCES balance (id),
--     CONSTRAINT receiver_user_fk FOREIGN KEY (receiver_user_id) REFERENCES balance (id),
-- );
-- CREATE TABLE service
-- (
--     id    BIGSERIAL PRIMARY KEY,
--     name  VARCHAR(100)  NOT NULL,
--     price DECIMAL(10, 2) NOT NULL,
-- );
-- CREATE TABLE order
-- (
--     id BIGSERIAL PRIMARY KEY,
-- );
