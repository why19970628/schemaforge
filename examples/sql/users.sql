CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'primary id',
  `email` varchar(255) NOT NULL COMMENT 'login email',
  `enabled` tinyint(1) NOT NULL DEFAULT 1 COMMENT 'enabled flag',
  `created_at` datetime NOT NULL COMMENT 'created time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='users';
