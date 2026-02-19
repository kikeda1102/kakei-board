CREATE TABLE IF NOT EXISTS events (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    aggregate_id   VARCHAR(36)     NOT NULL,
    aggregate_type VARCHAR(64)     NOT NULL,
    version        INT UNSIGNED    NOT NULL,
    event_type     VARCHAR(128)    NOT NULL,
    payload        JSON            NOT NULL,
    recorded_by    VARCHAR(36)     NOT NULL DEFAULT 'anonymous',
    occurred_at    DATETIME(6)     NOT NULL DEFAULT (UTC_TIMESTAMP(6)),
    PRIMARY KEY (id),
    INDEX idx_aggregate (aggregate_type, aggregate_id, version),
    UNIQUE KEY uk_aggregate_version (aggregate_id, aggregate_type, version)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
