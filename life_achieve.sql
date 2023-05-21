CREATE TABLE IF NOT EXISTS `user` (
  `id`        binary(16) PRIMARY KEY NOT NULL,
  `first_name` varchar(255) NOT NULL,
  `last_name`  varchar(255) NOT NULL,
  `email`     varchar(255) NOT NULL,
  `password`  varchar(255) NOT NULL
);