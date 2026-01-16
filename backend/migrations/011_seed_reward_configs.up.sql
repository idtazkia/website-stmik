-- Seed default reward configurations

-- External Referrer Rewards (triggered on enrollment)
INSERT INTO reward_configs (referrer_type, reward_type, amount, is_percentage, trigger_event, description) VALUES
    ('alumni', 'cash', 500000, false, 'enrollment', 'Komisi tunai untuk alumni yang mereferensikan mahasiswa baru'),
    ('teacher', 'cash', 750000, false, 'enrollment', 'Komisi tunai untuk guru yang mereferensikan siswa'),
    ('student', 'cash', 300000, false, 'enrollment', 'Komisi tunai untuk siswa yang mereferensikan teman'),
    ('partner', 'cash', 1000000, false, 'enrollment', 'Komisi tunai untuk mitra/bimbel'),
    ('staff', 'cash', 250000, false, 'enrollment', 'Komisi tunai untuk staf internal');

-- MGM Reward for current academic year
INSERT INTO mgm_reward_configs (academic_year, reward_type, referrer_amount, referee_amount, trigger_event, description) VALUES
    ('2025/2026', 'cash', 200000, NULL, 'enrollment', 'Reward tunai untuk mahasiswa yang mereferensikan teman'),
    ('2025/2026', 'tuition_discount', 10, 5, 'enrollment', 'Diskon SPP 10% untuk referrer, 5% untuk referee');
