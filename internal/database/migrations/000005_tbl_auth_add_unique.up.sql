ALTER TABLE auth_tokens
ADD CONSTRAINT unique_user_id UNIQUE (user_id);