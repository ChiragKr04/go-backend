-- Create a new table for invitation-user associations
CREATE TABLE `go_backend`.`invitation_users` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `invitation_group_id` INT NOT NULL,
  `user_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `invitation_group_id_idx` (`invitation_group_id` ASC) VISIBLE,
  INDEX `user_id_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `fk_invitation_group_id`
    FOREIGN KEY (`invitation_group_id`)
    REFERENCES `go_backend`.`invitation_groups` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `fk_user_id`
    FOREIGN KEY (`user_id`)
    REFERENCES `go_backend`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
