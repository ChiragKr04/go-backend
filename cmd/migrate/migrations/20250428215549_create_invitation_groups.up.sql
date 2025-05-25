-- Create a new table for invitation groups
CREATE TABLE IF NOT EXISTS `go_backend`.`invitation_groups` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `room_id` VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `room_id_idx` (`room_id` ASC) VISIBLE
);
