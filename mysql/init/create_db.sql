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
  `id` char(36) NOT NULL,
  `title` text NOT NULL,
  `composer_id` char(36) NOT NULL,
  `composer_name` varchar(32) NOT NULL,
  `message_id` char(36) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/* テストデータ */
INSERT INTO files (id, title, composer_id, composer_name, message_id, created_at) VALUE ("7309d123-d4b6-4abc-90b2-be2312f2ee5b", "Hello, World!", "b8d23ec2-c0a1-4567-9604-efc28f178c69", "yukumo", "548ef43a-0de8-4f8a-88d6-0d98292640ab", "2017-10-23 01:07:11");
INSERT INTO files (id, title, composer_id, composer_name, message_id, created_at) VALUE ("fe85b2a9-8604-4f75-8189-975506321a56", "Master", "e0dfd481-910a-48e6-bc77-8fc10dbc6868", "neg", "bf2266ad-0f26-472b-87de-2d4cc84069bb", "2017-10-22 19:03:36");
INSERT INTO files (id, title, composer_id, composer_name, message_id, created_at) VALUE ("9b11c23e-7718-4338-a9f5-36246dfaf60e", "20171021", "8acde99f-7977-4112-a1c3-4cdd2ce9283d", "xylo", "a96432dd-fa3d-4cf5-8df9-e540ea3494a3", "2017-10-21 01:03:00");
INSERT INTO files (id, title, composer_id, composer_name, message_id, created_at) VALUE ("e0794850-5e0d-4445-8e3d-6f548c36e7cb", "traP2dtm35", "01bef5a7-034d-4c6d-82c1-36a5add6d208", "OrangeStar", "0d39fd04-8fbd-4f92-b341-c9719bbe4a94", "2017-10-20 22:54:38");
INSERT INTO files (id, title, composer_id, composer_name, message_id, created_at) VALUE ("36b4ab59-6421-42d1-9e20-7426d41919f3", "f", "fa06c6d2-5fc3-45b3-82f4-1b014a26de28", "uynet", "004a20ab-2085-4810-9e26-2cd37dc9648b", "2017-10-20 22:47:37");
