-- Table for tasks
DROP TABLE IF EXISTS `tasks`;

CREATE TABLE `tasks` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL,
    `title` varchar(50) NOT NULL,
    `detail` varchar(200) NOT NULL,
    `priority` varchar(50) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deadline` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `is_done` boolean NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8mb4;

CREATE TABLE `users` (
    `user_id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_name` varchar(50) NOT NULL UNIQUE,
    `password` varchar(200) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `is_deleted` boolean NOT NULL DEFAULT b'0',
    PRIMARY KEY (`user_id`)
) DEFAULT CHARSET=utf8mb4;
