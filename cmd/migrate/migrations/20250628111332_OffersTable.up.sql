CREATE TABLE IF NOT EXISTS `offers` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `offer` TEXT NULL,
    `answer` TEXT NULL,
    `offerer_user_id` INT NOT NULL,
    `answerer_user_id` INT NULL,
    `room_id` INT NOT NULL,
    `offer_ice_candidates` TEXT NULL,
    `answer_ice_candidates` TEXT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`offerer_user_id`) REFERENCES `users` (`id`),
    FOREIGN KEY (`answerer_user_id`) REFERENCES `users` (`id`),
    FOREIGN KEY (`room_id`) REFERENCES `rooms` (`id`)
);