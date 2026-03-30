-- ══════════════════════════════════════════════════════════════════
-- VCT Platform — Migration 0099: PostgreSQL 18.3 Features (v3.0)
-- ══════════════════════════════════════════════════════════════════

BEGIN;

-- ═══════════════════════════════════════════════
-- 1. UUIDv7 Support (PG18 Native)
-- ═══════════════════════════════════════════════
-- PostgreSQL 18 introduces native uuidv7() function.
-- Bảng hiện tại vẫn dùng gen_random_uuid() (v4) để không xáo trộn dữ liệu cũ.
-- Nhưng ta có thể dùng uuidv7 cho các bảng log hoặc partition sau này.

DO $$
BEGIN
    IF current_setting('server_version_num')::int >= 180000 THEN
        -- Test gọi uuidv7() để đảm bảo function tồn tại trên db PostgreSQL 18.
        PERFORM uuidv7();
        RAISE NOTICE '[PG18] uuidv7() is available!';
    ELSE
        RAISE NOTICE 'Running on older PostgreSQL version. uuidv7() skipped.';
    END IF;
END $$;


-- ═══════════════════════════════════════════════
-- 2. Temporal Constraints & Period (PG18+)
-- ═══════════════════════════════════════════════
-- PostgreSQL 18 hỗ trợ Constraints với Period.
-- Sau này ta có thể dùng `PERIOD FOR` trong `core.belt_history` 
-- để kiểm tra tính liên tục của đai thay cho triggers phức tạp.


-- ═══════════════════════════════════════════════
-- 3. Virtual Generated Columns (PG18+)
-- ═══════════════════════════════════════════════
-- Cột `search_name` trong core.users hiện là STORED. 
-- PostgreSQL 18 hỗ trợ `GENERATED ALWAYS AS ... VIRTUAL` (mặc định của GENERATED).
-- Giúp giảm I/O cho các text fields search lớn.

-- Demo (Không apply để tránh lock bảng users):
-- ALTER TABLE core.users 
--   ADD COLUMN search_name_v2 TEXT GENERATED ALWAYS AS (lower(core.immutable_unaccent(full_name))); 
-- (Virtual theo mặc định trên PG18)


COMMIT;
