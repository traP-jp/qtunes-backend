DROP DATABASE IF EXISTS 21hack02;
CREATE DATABASE 21hack02;
USE 21hack02;

CREATE TABLE IF NOT EXISTS `users` (
  `id` char(36) PRIMARY KEY NOT NULL,
  `name` varchar(32) NOT NULL UNIQUE,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `favorites` (
  `user_id` char(36) NOT NULL,
  `composer_id` char(36) NOT NULL,
  `sound_id` char(36) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `files` (
  `id` char(36) PRIMARY KEY,
  `title` text NOT NULL,
  `composer_id` char(36) NOT NULL,
  `composer_name` varchar(32) NOT NULL,
  `message_id` char(36) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
