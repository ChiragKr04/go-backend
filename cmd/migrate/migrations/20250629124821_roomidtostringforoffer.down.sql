ALTER TABLE `offers` 
  DROP FOREIGN KEY `offers_room_id_fk`,
  MODIFY COLUMN `room_id` INT NOT NULL,
  ADD CONSTRAINT `offers_ibfk_3` FOREIGN KEY (`room_id`) REFERENCES `rooms` (`id`);
