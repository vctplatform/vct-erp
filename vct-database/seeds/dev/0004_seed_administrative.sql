-- ═══════════════════════════════════════════════════════════════
-- VCT PLATFORM — Seed 0004: Administrative & Reference Data
-- Dữ liệu hành chính và tham chiếu nền tảng
-- Load từ JSON files trong backend/data/
-- ═══════════════════════════════════════════════════════════════

BEGIN;

-- ════════════════════════════════════════════════════════
-- SYSTEM CONFIGURATION KEY-VALUE
-- ════════════════════════════════════════════════════════

DO $seed$
BEGIN
  IF to_regclass('public.entity_records') IS NULL THEN
    RAISE NOTICE 'entity_records table not found, skip seed 0004';
    RETURN;
  END IF;

  EXECUTE $sql$
    INSERT INTO entity_records(entity, id, payload)
    VALUES
      -- ── Hệ thống Đai / Đẳng ──
      ('reference-data', 'belt-ranks', '{
        "id": "belt-ranks",
        "type": "belt_rank_system",
        "data": [
          {"code":"white","name":"Đai trắng","color":"#FFFFFF","order":1},
          {"code":"yellow","name":"Đai vàng","color":"#FFD700","order":2},
          {"code":"green","name":"Đai xanh lá","color":"#228B22","order":3},
          {"code":"blue","name":"Đai xanh dương","color":"#0066CC","order":4},
          {"code":"brown","name":"Đai nâu","color":"#8B4513","order":5},
          {"code":"black_1","name":"Đai đen sơ đẳng","color":"#000000","order":6},
          {"code":"black_2","name":"Đai đen nhị đẳng","color":"#000000","order":7},
          {"code":"black_3","name":"Đai đen tam đẳng","color":"#000000","order":8},
          {"code":"black_4","name":"Đai đen tứ đẳng","color":"#000000","order":9},
          {"code":"black_5","name":"Đai đen ngũ đẳng","color":"#000000","order":10},
          {"code":"black_6","name":"Đai đen lục đẳng","color":"#000000","order":11},
          {"code":"black_7","name":"Đai đen thất đẳng","color":"#000000","order":12},
          {"code":"black_8","name":"Đai đen bát đẳng","color":"#000000","order":13},
          {"code":"black_9","name":"Đai đen cửu đẳng","color":"#000000","order":14}
        ]
      }'::jsonb),

      -- ── Hạng cân chuẩn ──
      ('reference-data', 'weight-classes-standard', '{
        "id": "weight-classes-standard",
        "type": "standard_weight_classes",
        "nam": [
          {"code":"M48","name":"Nam 48kg","min":0,"max":48},
          {"code":"M52","name":"Nam 52kg","min":48.1,"max":52},
          {"code":"M55","name":"Nam 55kg","min":52.1,"max":55},
          {"code":"M60","name":"Nam 60kg","min":55.1,"max":60},
          {"code":"M65","name":"Nam 65kg","min":60.1,"max":65},
          {"code":"M70","name":"Nam 70kg","min":65.1,"max":70},
          {"code":"M75","name":"Nam 75kg","min":70.1,"max":75},
          {"code":"M80","name":"Nam 80kg","min":75.1,"max":80},
          {"code":"M85","name":"Nam 85kg","min":80.1,"max":85},
          {"code":"M85P","name":"Nam +85kg","min":85.1,"max":null}
        ],
        "nu": [
          {"code":"F42","name":"Nữ 42kg","min":0,"max":42},
          {"code":"F45","name":"Nữ 45kg","min":42.1,"max":45},
          {"code":"F48","name":"Nữ 48kg","min":45.1,"max":48},
          {"code":"F52","name":"Nữ 52kg","min":48.1,"max":52},
          {"code":"F56","name":"Nữ 56kg","min":52.1,"max":56},
          {"code":"F60","name":"Nữ 60kg","min":56.1,"max":60},
          {"code":"F65","name":"Nữ 65kg","min":60.1,"max":65},
          {"code":"F65P","name":"Nữ +65kg","min":65.1,"max":null}
        ]
      }'::jsonb),

      -- ── Nhóm lứa tuổi chuẩn ──
      ('reference-data', 'age-groups-standard', '{
        "id": "age-groups-standard",
        "type": "standard_age_groups",
        "data": [
          {"code":"thieu_nhi_a","name":"Thiếu nhi A","min":8,"max":10},
          {"code":"thieu_nhi_b","name":"Thiếu nhi B","min":11,"max":13},
          {"code":"thieu_nien","name":"Thiếu niên","min":14,"max":16},
          {"code":"thanh_nien","name":"Thanh niên","min":17,"max":35},
          {"code":"trung_nien","name":"Trung niên","min":36,"max":50},
          {"code":"cao_nien","name":"Cao niên","min":51,"max":null}
        ]
      }'::jsonb),

      -- ── Cấu hình hệ thống ──
      ('reference-data', 'system-config', '{
        "id": "system-config",
        "type": "system_configuration",
        "platform_name": "VCT Platform",
        "default_language": "vi-VN",
        "timezone": "Asia/Ho_Chi_Minh",
        "currency": "VND",
        "date_format": "DD/MM/YYYY",
        "phone_country_code": "+84",
        "tournament_defaults": {
          "max_athletes_per_team_per_event": 3,
          "weigh_in_tolerance_kg": 0.5,
          "max_rounds_combat": 3,
          "round_duration_seconds": 120,
          "min_judges_form": 5,
          "drop_highest_score": true,
          "drop_lowest_score": true
        }
      }'::jsonb),

      -- ── Danh mục vi phạm ──
      ('reference-data', 'violation-types', '{
        "id": "violation-types",
        "type": "violation_type_catalog",
        "data": [
          {"code":"weight_fraud","name":"Khai gian cân nặng","severity":"high"},
          {"code":"unsportsmanlike","name":"Hành vi phi thể thao","severity":"medium"},
          {"code":"doping","name":"Sử dụng doping","severity":"critical"},
          {"code":"competition_rule_violation","name":"Vi phạm quy chế thi đấu","severity":"medium"},
          {"code":"admin_breach","name":"Vi phạm hành chính","severity":"low"},
          {"code":"referee_abuse","name":"Xúc phạm trọng tài","severity":"high"},
          {"code":"age_fraud","name":"Khai gian tuổi","severity":"high"},
          {"code":"document_forgery","name":"Giả mạo giấy tờ","severity":"critical"},
          {"code":"match_fixing","name":"Dàn xếp kết quả","severity":"critical"},
          {"code":"safety_violation","name":"Vi phạm an toàn","severity":"high"}
        ]
      }'::jsonb),

      -- ── Danh mục hình thức xử lý ──
      ('reference-data', 'sanction-types', '{
        "id": "sanction-types",
        "type": "sanction_type_catalog",
        "data": [
          {"code":"warning","name":"Cảnh cáo"},
          {"code":"fine","name":"Phạt tiền"},
          {"code":"suspension","name":"Đình chỉ thi đấu"},
          {"code":"ban","name":"Cấm vĩnh viễn"},
          {"code":"result_annulment","name":"Hủy kết quả"},
          {"code":"probation","name":"Quản chế"},
          {"code":"legal_referral","name":"Chuyển cơ quan pháp luật"}
        ]
      }'::jsonb)

    ON CONFLICT (entity, id) DO UPDATE
      SET payload = EXCLUDED.payload,
          updated_at = NOW();
  $sql$;
END
$seed$;

COMMIT;
