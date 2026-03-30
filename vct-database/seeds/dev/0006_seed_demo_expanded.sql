-- ═══════════════════════════════════════════════════════════════
-- VCT PLATFORM — Seed 0006: Expanded Demo Data
-- Sample clubs + Tournament history (for demo/reporting)
-- ═══════════════════════════════════════════════════════════════

BEGIN;

DO $seed$
BEGIN
  IF to_regclass('public.entity_records') IS NULL THEN
    RAISE NOTICE 'entity_records table not found, skip seed 0006';
    RETURN;
  END IF;

  EXECUTE $sql$
    INSERT INTO entity_records(entity, id, payload)
    VALUES
      -- ════════════════════════════════════════════════════════
      -- SAMPLE CLUBS
      -- ════════════════════════════════════════════════════════
      ('sample-clubs', 'CLB-BD-01',  '{"id":"CLB-BD-01","name":"CLB Bình Định Quyền","province":"Quảng Ngãi","school":"Bình Định","founded":1995,"members":85,"grade":"A","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-BD-02',  '{"id":"CLB-BD-02","name":"CLB Tây Sơn Võ Đạo","province":"Quảng Ngãi","school":"Bình Định","founded":2001,"members":60,"grade":"A","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-HCM-01', '{"id":"CLB-HCM-01","name":"CLB Lý Gia Quyền","province":"TP.HCM","school":"Hồng Gia","founded":1998,"members":120,"grade":"A","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-HCM-02', '{"id":"CLB-HCM-02","name":"CLB Tân Khánh Bà Trà","province":"TP.HCM","school":"TKBT","founded":1992,"members":95,"grade":"A","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-HCM-03', '{"id":"CLB-HCM-03","name":"CLB Sa Long Cương HCM","province":"TP.HCM","school":"Sa Long Cương","founded":2005,"members":70,"grade":"B","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-HN-01',  '{"id":"CLB-HN-01","name":"CLB Võ Thăng Long","province":"Hà Nội","school":"Tổng hợp","founded":1997,"members":110,"grade":"A","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-HN-02',  '{"id":"CLB-HN-02","name":"CLB Thiếu Lâm Hà Nội","province":"Hà Nội","school":"Thiếu Lâm","founded":2003,"members":65,"grade":"B","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-HN-03',  '{"id":"CLB-HN-03","name":"CLB Kim Kê Hà Nội","province":"Hà Nội","school":"Kim Kê","founded":2010,"members":40,"grade":"B","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-DN-01',  '{"id":"CLB-DN-01","name":"CLB VCT Đà Nẵng","province":"Đà Nẵng","school":"Tổng hợp","founded":2000,"members":75,"grade":"A","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-DN-02',  '{"id":"CLB-DN-02","name":"CLB Bạch Mi Đà Nẵng","province":"Đà Nẵng","school":"Bạch Mi","founded":2008,"members":45,"grade":"B","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-TH-01',  '{"id":"CLB-TH-01","name":"CLB Thanh Hóa VCT","province":"Thanh Hóa","school":"Tổng hợp","founded":1999,"members":90,"grade":"A","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-NA-01',  '{"id":"CLB-NA-01","name":"CLB Nghệ An VCT","province":"Nghệ An","school":"Tổng hợp","founded":2002,"members":65,"grade":"B","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-HUE-01', '{"id":"CLB-HUE-01","name":"CLB Huế VCT","province":"Huế","school":"Tổng hợp","founded":2004,"members":55,"grade":"B","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-CT-01',  '{"id":"CLB-CT-01","name":"CLB Cần Thơ VCT","province":"Cần Thơ","school":"Tổng hợp","founded":2006,"members":50,"grade":"B","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-DNI-01', '{"id":"CLB-DNI-01","name":"CLB Đồng Nai VCT","province":"Đồng Nai","school":"Tổng hợp","founded":2001,"members":80,"grade":"A","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-KH-01',  '{"id":"CLB-KH-01","name":"CLB Khánh Hòa VCT","province":"Khánh Hòa","school":"Bình Định","founded":2007,"members":55,"grade":"B","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-AG-01',  '{"id":"CLB-AG-01","name":"CLB An Giang VCT","province":"An Giang","school":"Tổng hợp","founded":2009,"members":40,"grade":"B","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-DL-01',  '{"id":"CLB-DL-01","name":"CLB Đắk Lắk VCT","province":"Đắk Lắk","school":"Tổng hợp","founded":2011,"members":35,"grade":"C","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-HP-01',  '{"id":"CLB-HP-01","name":"CLB Hải Phòng VCT","province":"Hải Phòng","school":"Tổng hợp","founded":2003,"members":60,"grade":"B","status":"active"}'::jsonb),
      ('sample-clubs', 'CLB-HCM-04', '{"id":"CLB-HCM-04","name":"CLB Vovinam Q.7","province":"TP.HCM","school":"Vovinam","founded":2015,"members":85,"grade":"A","status":"active"}'::jsonb),

      -- ════════════════════════════════════════════════════════
      -- TOURNAMENT HISTORY
      -- ════════════════════════════════════════════════════════
      ('tournament-history', 'HIST-2024', '{
        "id":"HIST-2024","code":"VCT-2024",
        "name":"Giải VĐ VCT Toàn Quốc 2024",
        "level":"quoc_gia","status":"ket_thuc",
        "start_date":"2024-08-10","end_date":"2024-08-15",
        "location":"Nghệ An","venue":"Nhà thi đấu TP Vinh",
        "stats":{"teams":28,"athletes":450,"events":20,"matches":180},
        "medal_top5":[
          {"rank":1,"province":"TP.HCM","gold":8,"silver":5,"bronze":7},
          {"rank":2,"province":"Quảng Ngãi","gold":6,"silver":7,"bronze":4},
          {"rank":3,"province":"Hà Nội","gold":5,"silver":4,"bronze":6},
          {"rank":4,"province":"Thanh Hóa","gold":3,"silver":3,"bronze":5},
          {"rank":5,"province":"Nghệ An","gold":3,"silver":2,"bronze":4}
        ]
      }'::jsonb),

      ('tournament-history', 'HIST-2025', '{
        "id":"HIST-2025","code":"VCT-2025",
        "name":"Giải VĐ VCT Toàn Quốc 2025",
        "level":"quoc_gia","status":"ket_thuc",
        "start_date":"2025-09-05","end_date":"2025-09-10",
        "location":"TP.HCM","venue":"Nhà thi đấu Phú Thọ",
        "stats":{"teams":32,"athletes":520,"events":22,"matches":210},
        "medal_top5":[
          {"rank":1,"province":"Quảng Ngãi","gold":9,"silver":6,"bronze":5},
          {"rank":2,"province":"TP.HCM","gold":7,"silver":8,"bronze":6},
          {"rank":3,"province":"Hà Nội","gold":5,"silver":5,"bronze":7},
          {"rank":4,"province":"Đồng Nai","gold":4,"silver":3,"bronze":5},
          {"rank":5,"province":"Đà Nẵng","gold":3,"silver":4,"bronze":3}
        ]
      }'::jsonb),

      ('tournament-history', 'HIST-KV7-2025', '{
        "id":"HIST-KV7-2025","code":"KV7-2025",
        "name":"Giải Khu vực VII - Miền Trung Tây Nguyên 2025",
        "level":"khu_vuc","status":"ket_thuc",
        "start_date":"2025-05-20","end_date":"2025-05-23",
        "location":"Đà Nẵng","venue":"Cung thể thao Tiên Sơn",
        "stats":{"teams":11,"athletes":180,"events":16,"matches":85},
        "medal_top5":[
          {"rank":1,"province":"Quảng Ngãi","gold":7,"silver":4,"bronze":3},
          {"rank":2,"province":"Đà Nẵng","gold":4,"silver":5,"bronze":4},
          {"rank":3,"province":"Khánh Hòa","gold":3,"silver":2,"bronze":5},
          {"rank":4,"province":"Đắk Lắk","gold":2,"silver":3,"bronze":2},
          {"rank":5,"province":"Huế","gold":2,"silver":1,"bronze":3}
        ]
      }'::jsonb)

    ON CONFLICT (entity, id) DO UPDATE
      SET payload = EXCLUDED.payload,
          updated_at = NOW();
  $sql$;
END
$seed$;

COMMIT;
