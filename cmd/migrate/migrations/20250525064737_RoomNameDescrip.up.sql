-- add room name and room description to room table
ALTER TABLE rooms 
ADD COLUMN room_name VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN room_description VARCHAR(1000) NOT NULL DEFAULT '';