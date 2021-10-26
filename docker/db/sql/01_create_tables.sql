-- Table for tasks
DROP TABLE IF EXISTS `tasks`;

CREATE TABLE `tasks` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL DEFAULT 0,
    `title` varchar(50) NOT NULL,
    `detail` varchar(200) NOT NULL,
    `priority` varchar(50) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deadline` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `is_done` boolean NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8mb4;
