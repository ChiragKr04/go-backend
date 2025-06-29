ALTER TABLE `offers` 
  DROP FOREIGN KEY `offers_ibfk_3`,
  MODIFY COLUMN `room_id` VARCHAR(120) NOT NULL,
  ADD CONSTRAINT `offers_room_id_fk` FOREIGN KEY (`room_id`) REFERENCES `rooms` (`room_id`);
