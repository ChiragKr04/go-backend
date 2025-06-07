CREATE TABLE IF NOT EXISTS `room_users` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `room_id` VARCHAR(255) NOT NULL,
  `user_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  INDEX `user_id_idx` (`user_id` ASC) VISIBLE,
  INDEX `room_id_idx` (`room_id` ASC) VISIBLE,
  CONSTRAINT `user_id`
    FOREIGN KEY (`user_id`)
    REFERENCES `users` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `room_id`
    FOREIGN KEY (`room_id`)
    REFERENCES `rooms` (`room_id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);
