-- Add password_reset to allowed token types
ALTER TABLE verification_tokens DROP CONSTRAINT IF EXISTS verification_tokens_token_type_check;
ALTER TABLE verification_tokens ADD CONSTRAINT verification_tokens_token_type_check
    CHECK (token_type IN ('email', 'phone', 'password_reset'));
