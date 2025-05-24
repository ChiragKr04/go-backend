ALTER TABLE `go_backend`.`rooms` 
ADD COLUMN `invitations` INT NULL AFTER `short_room_id`,
CHANGE COLUMN `short_room_id` `short_room_id` VARCHAR(255) NOT NULL ;
