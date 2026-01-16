-- Remove seed data (in reverse order of dependencies)
DELETE FROM prodis WHERE code IN ('SI', 'TI');
DELETE FROM obstacles;
DELETE FROM interaction_categories;
DELETE FROM fee_types WHERE code IN ('registration', 'tuition', 'dormitory');
