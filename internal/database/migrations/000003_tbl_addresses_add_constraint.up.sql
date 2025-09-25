ALTER TABLE addresses 
ADD CONSTRAINT one_default_address_per_user
CHECK (
    (is_default = true AND user_id IS NOT NULL)
    OR is_default = false
);
