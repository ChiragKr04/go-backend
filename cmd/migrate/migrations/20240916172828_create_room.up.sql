CREATE TABLE `rooms` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `room_id` VARCHAR(120) NOT NULL,
  `created_by` INT NOT NULL,
  `created_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  UNIQUE INDEX `room_id_UNIQUE` (`room_id` ASC) VISIBLE);