ALTER TABLE messages
ADD COLUMN group_id INT,
ADD COLUMN is_group_msg BOOLEAN DEFAULT FALSE;