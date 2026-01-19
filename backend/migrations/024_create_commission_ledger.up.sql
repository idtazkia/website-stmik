-- Commission Ledger
-- Tracks individual commission entries for referrers

CREATE TABLE commission_ledger (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_referrer UUID NOT NULL REFERENCES referrers(id),
    id_candidate UUID NOT NULL REFERENCES candidates(id),
    trigger_event VARCHAR(20) NOT NULL CHECK (trigger_event IN ('registration', 'commitment', 'enrollment')),
    amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'paid', 'cancelled')),
    approved_at TIMESTAMPTZ,
    id_approved_by UUID REFERENCES users(id),
    paid_at TIMESTAMPTZ,
    id_paid_by UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(id_referrer, id_candidate, trigger_event)
);

-- Index for referrer lookup
CREATE INDEX idx_commission_ledger_referrer ON commission_ledger(id_referrer);

-- Index for candidate lookup
CREATE INDEX idx_commission_ledger_candidate ON commission_ledger(id_candidate);

-- Index for status filtering
CREATE INDEX idx_commission_ledger_status ON commission_ledger(status);

-- Index for pending commissions
CREATE INDEX idx_commission_ledger_pending ON commission_ledger(status) WHERE status = 'pending';

-- Trigger for updated_at
CREATE TRIGGER update_commission_ledger_updated_at
    BEFORE UPDATE ON commission_ledger
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE commission_ledger IS 'Individual commission entries for referrers';
COMMENT ON COLUMN commission_ledger.trigger_event IS 'Event that triggered the commission: registration, commitment, enrollment';
COMMENT ON COLUMN commission_ledger.amount IS 'Commission amount in IDR';
COMMENT ON COLUMN commission_ledger.status IS 'Commission status: pending (awaiting approval), approved (ready for payout), paid (completed), cancelled';
