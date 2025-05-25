CREATE TABLE IF NOT EXISTS `go_backend`.`invitations` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `room_id` VARCHAR(255) NOT NULL,
  `user_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `user_id_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `id`
    FOREIGN KEY (`user_id`)
    REFERENCES `go_backend`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);
