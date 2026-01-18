-- Drop commission_ledger table
DROP TRIGGER IF EXISTS update_commission_ledger_updated_at ON commission_ledger;
DROP TABLE IF EXISTS commission_ledger;
