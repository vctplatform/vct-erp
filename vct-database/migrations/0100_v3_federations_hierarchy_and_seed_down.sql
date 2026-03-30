-- ══════════════════════════════════════════════════════════════════
-- VCT Platform — Migration 0100_down (Undo Federation Hierarchy)
-- Xóa cột parent_id, nhưng giữ nguyên data cũ (do là schema dev).
-- ══════════════════════════════════════════════════════════════════

BEGIN;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'core'
          AND table_name = 'federations'
          AND column_name = 'parent_id'
    ) THEN
        ALTER TABLE core.federations DROP COLUMN parent_id;
    END IF;
END $$;

COMMIT;
