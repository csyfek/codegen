CREATE TABLE `world_definition` (
  `id`                  VARCHAR(256) NOT NULL,
  `name`                VARCHAR(256) NOT NULL,
  `width`               BLOB         NOT NULL,
  `height`              BLOB         NOT NULL,
  `percent_water`       DOUBLE       NOT NULL,
  `percent_grass`       DOUBLE       NOT NULL,
  `percent_mountain`    DOUBLE       NOT NULL,
  `percent_tree`        DOUBLE       NOT NULL,
  `percent_npc_village` DOUBLE       NOT NULL,
  `percent_pc_village`  DOUBLE       NOT NULL,
  `speed_modifier`      DOUBLE       NOT NULL,
  `open_date`           DATETIME     NOT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPRESSED;
