-- ══════════════════════════════════════════════════════════════════
-- VCT Platform — Migration 0100: Federation Hierarchy & 34 Provinces (v3.0)
-- Bổ sung hệ thống cấp bậc LĐ (parent_id) và Seed tự động 34 LĐ mới nhất
-- ══════════════════════════════════════════════════════════════════

BEGIN;

-- ═══════════════════════════════════════════════
-- 1. ADD PARENT_ID TO FEDERATIONS
-- ═══════════════════════════════════════════════

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'core'
          AND table_name = 'federations'
          AND column_name = 'parent_id'
    ) THEN
        ALTER TABLE core.federations
        ADD COLUMN parent_id UUID REFERENCES core.federations(id) ON DELETE SET NULL;
        
        COMMENT ON COLUMN core.federations.parent_id IS 'Mã của Liên đoàn cấp trên (Quốc gia quản lý Tỉnh/Thành)';
    END IF;
END $$;

-- ═══════════════════════════════════════════════
-- 2. SEED NATIONAL FEDERATION (VIỆT NAM)
-- ═══════════════════════════════════════════════

DO $$
DECLARE
    v_national_id UUID;
    v_prov RECORD;
    v_name TEXT;
	v_type TEXT;
    
    -- Danh sách 6 Thành phố trực thuộc trung ương
    v_cities TEXT[] := ARRAY['hanoi', 'haiphong', 'danang', 'hochiminh', 'cantho', 'hue'];
    
    -- Danh sách 28 Tỉnh (sau sáp nhập)
    v_provinces TEXT[] := ARRAY[
        'quangninh', 'bacninh', 'hungyen', 'ninhbinh', 'phutho', 'thainguyen', 'tuyenquang', 'laocai', 'caobang',
        'langson', 'sonla', 'dienbien', 'laichau', 'thanhhoa', 'nghean', 'hatinh', 'quangtri', 'quangngai',
        'khanhhoa', 'gialai', 'daklak', 'lamdong', 'dongnai', 'angiang', 'camau', 'dongthap', 'tayninh', 'vinhlong'
    ];
    
    v_city_names JSONB := '{"hanoi": "Hà Nội", "haiphong": "Hải Phòng", "danang": "Đà Nẵng", "hochiminh": "Hồ Chí Minh", "cantho": "Cần Thơ", "hue": "Huế"}';
    v_province_names JSONB := '{
        "quangninh": "Quảng Ninh", "bacninh": "Bắc Ninh", "hungyen": "Hưng Yên", "ninhbinh": "Ninh Bình",
        "phutho": "Phú Thọ", "thainguyen": "Thái Nguyên", "tuyenquang": "Tuyên Quang", "laocai": "Lào Cai",
        "caobang": "Cao Bằng", "langson": "Lạng Sơn", "sonla": "Sơn La", "dienbien": "Điện Biên",
        "laichau": "Lai Châu", "thanhhoa": "Thanh Hóa", "nghean": "Nghệ An", "hatinh": "Hà Tĩnh",
        "quangtri": "Quảng Trị", "quangngai": "Quảng Ngãi", "khanhhoa": "Khánh Hòa", "gialai": "Gia Lai",
        "daklak": "Đắk Lắk", "lamdong": "Lâm Đồng", "dongnai": "Đồng Nai", "angiang": "An Giang",
        "camau": "Cà Mau", "dongthap": "Đồng Tháp", "tayninh": "Tây Ninh", "vinhlong": "Vĩnh Long"
    }';

BEGIN
    -- Upsert Liên Đoàn Quốc Gia
    INSERT INTO core.federations (name, region, province_code, is_active)
    VALUES ('Liên Đoàn Võ Thuật Cổ Truyền Việt Nam', 'quocgia', 'vietnam', true)
    ON CONFLICT (province_code) DO UPDATE 
    SET name = EXCLUDED.name, region = EXCLUDED.region
    RETURNING id INTO v_national_id;

    RAISE NOTICE 'Upserted National Federation (ID: %)', v_national_id;

    -- ═══════════════════════════════════════════════
    -- 3. SEED 34 LOCAL FEDERATIONS
    -- ═══════════════════════════════════════════════

    -- 3.1 6 Thành phố trực thuộc Trung ương
    FOR i IN 1 .. array_length(v_cities, 1) LOOP
        v_name := 'Liên Đoàn Võ Thuật Cổ Truyền Thành phố ' || (v_city_names ->> v_cities[i]);
        
        -- Tạo Schema riêng cho liên đoàn này (chức năng từ Migration 0095)
        PERFORM core.create_federation_schema(v_cities[i]);
        
        -- Ghi dữ liệu vào core.federations với parent_id = Quốc gia
        INSERT INTO core.federations (name, region, province_code, parent_id, is_active)
        VALUES (v_name, 'tinh', v_cities[i], v_national_id, true)
        ON CONFLICT (province_code) DO UPDATE 
        SET name = EXCLUDED.name, parent_id = EXCLUDED.parent_id;
    END LOOP;

    -- 3.2 28 Tỉnh
    FOR i IN 1 .. array_length(v_provinces, 1) LOOP
        v_name := 'Liên Đoàn Võ Thuật Cổ Truyền Tỉnh ' || (v_province_names ->> v_provinces[i]);
        
        -- Tạo Schema riêng cho liên đoàn này
        PERFORM core.create_federation_schema(v_provinces[i]);
        
        -- Ghi dữ liệu vào core.federations với parent_id = Quốc gia
        INSERT INTO core.federations (name, region, province_code, parent_id, is_active)
        VALUES (v_name, 'tinh', v_provinces[i], v_national_id, true)
        ON CONFLICT (province_code) DO UPDATE 
        SET name = EXCLUDED.name, parent_id = EXCLUDED.parent_id;
    END LOOP;

END $$;

COMMIT;
