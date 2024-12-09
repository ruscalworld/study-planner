CREATE TABLE IF NOT EXISTS `institutions`
(
    `id`   bigint(20)   NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS `curriculums`
(
    `id`             bigint(20)   NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name`           varchar(255) NOT NULL,
    `semester`       int(11)      NOT NULL,
    `institution_id` bigint(20) REFERENCES `institutions` (`id`)
);

CREATE TABLE IF NOT EXISTS `curriculum_codes`
(
    `id`            bigint(20)                         NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `code`          varchar(50)                        NOT NULL,
    `role`          enum ('Owner', 'Editor', 'Viewer') NOT NULL DEFAULT 'Viewer',
    `curriculum_id` bigint(20)                         NOT NULL REFERENCES `curriculums` (`id`),
    `user_id`       bigint(20)                         NOT NULL REFERENCES `users` (`id`),
    `expires_at`    timestamp                                   DEFAULT NULL,
    `created_at`    timestamp                          NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (`code`)
);

CREATE TABLE IF NOT EXISTS `disciplines`
(
    `id`            bigint(20)   NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name`          varchar(255) NOT NULL,
    `curriculum_id` bigint(20)   NOT NULL REFERENCES `curriculums` (`id`)
);

CREATE TABLE IF NOT EXISTS `discipline_links`
(
    `id`            bigint(20)    NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `discipline_id` bigint(20)    NOT NULL REFERENCES `disciplines` (`id`),
    `name`          varchar(255)  NOT NULL,
    `url`           varchar(1000) NOT NULL
);

CREATE TABLE IF NOT EXISTS `task_groups`
(
    `id`            bigint(20)   NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name`          varchar(255) NOT NULL,
    `discipline_id` bigint(20)   NOT NULL REFERENCES `disciplines` (`id`)
);

CREATE TABLE IF NOT EXISTS `tasks`
(
    `id`            bigint(20)                         NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name`          varchar(255)                       NOT NULL,
    `external_name` varchar(255)                                DEFAULT NULL,
    `description`   text                                        DEFAULT NULL,
    `task_group_id` bigint(20)                         NOT NULL REFERENCES `task_groups` (`id`),
    `status`        enum ('NotPublished', 'Available') NOT NULL DEFAULT 'NotPublished',
    `difficulty`    int(10) unsigned                   NOT NULL DEFAULT 1,
    `deadline`      timestamp                          NULL     DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS `task_links`
(
    `id`      bigint(20)    NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `task_id` bigint(20)    NOT NULL REFERENCES `tasks` (`id`),
    `name`    varchar(255)  NOT NULL,
    `url`     varchar(1000) NOT NULL
);

CREATE TABLE IF NOT EXISTS `users`
(
    `id`          bigint(20)   NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name`        varchar(255) NOT NULL,
    `avatar_url`  varchar(1000)         DEFAULT NULL,
    `platform`    varchar(255) NOT NULL,
    `external_id` varchar(255) NOT NULL,
    `created_at`  timestamp    NOT NULL DEFAULT current_timestamp(),

    UNIQUE (`platform`, `external_id`)
);

CREATE TABLE IF NOT EXISTS `user_curriculums`
(
    `id`            bigint(20)                         NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `user_id`       bigint(20)                         NOT NULL REFERENCES `users` (`id`),
    `curriculum_id` bigint(20)                         NOT NULL REFERENCES `curriculums` (`id`),
    `role`          enum ('Owner', 'Editor', 'Viewer') NOT NULL DEFAULT 'Viewer',

    UNIQUE `user_id` (`user_id`, `curriculum_id`)
);

CREATE TABLE IF NOT EXISTS `user_goals`
(
    `id`            bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `user_id`       bigint(20) NOT NULL REFERENCES `users` (`id`),
    `task_group_id` bigint(20) NOT NULL REFERENCES `task_groups` (`id`),
    `min_completed` int(11)    NOT NULL,

    UNIQUE `user_id` (`user_id`, `task_group_id`)
);

CREATE TABLE IF NOT EXISTS `user_task_progress`
(
    `id`           bigint(20)                                                        NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `user_id`      bigint(20)                                                        NOT NULL REFERENCES `users` (`id`),
    `task_id`      bigint(20)                                                        NOT NULL REFERENCES `tasks` (`id`),
    `status`       enum ('NotStarted', 'InProgress', 'NeedsProtection', 'Completed') NOT NULL DEFAULT 'NotStarted',
    `grade`        enum ('Excellent', 'Good', 'Satisfactory', 'Credited')                     DEFAULT NULL,
    `started_at`   timestamp                                                         NULL     DEFAULT NULL,
    `completed_at` timestamp                                                         NULL     DEFAULT NULL,

    UNIQUE `task_id` (`task_id`, `user_id`)
);
