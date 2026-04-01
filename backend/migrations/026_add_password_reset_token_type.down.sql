-- Revert token type constraint
DELETE FROM verification_tokens WHERE token_type = 'password_reset';
ALTER TABLE verification_tokens DROP CONSTRAINT IF EXISTS verification_tokens_token_type_check;
ALTER TABLE verification_tokens ADD CONSTRAINT verification_tokens_token_type_check
    CHECK (token_type IN ('email', 'phone'));
