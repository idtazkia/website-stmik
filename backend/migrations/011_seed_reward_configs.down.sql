-- Remove seeded reward configs
DELETE FROM mgm_reward_configs WHERE academic_year = '2025/2026';
DELETE FROM reward_configs WHERE referrer_type IN ('alumni', 'teacher', 'student', 'partner', 'staff');
