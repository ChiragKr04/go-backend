ALTER TABLE `rooms` MODIFY COLUMN `short_room_id` VARCHAR(255) NULL;
ALTER TABLE `rooms` DROP COLUMN `invitations`;
