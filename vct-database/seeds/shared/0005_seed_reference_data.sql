-- ═══════════════════════════════════════════════════════════════
-- VCT PLATFORM — Seed 0005: Standard Forms & Templates
-- Bài quyền chuẩn, notification templates
-- ═══════════════════════════════════════════════════════════════

BEGIN;

DO $seed$
BEGIN
  IF to_regclass('public.entity_records') IS NULL THEN
    RAISE NOTICE 'entity_records table not found, skip seed 0005';
    RETURN;
  END IF;

  EXECUTE $sql$
    INSERT INTO entity_records(entity, id, payload)
    VALUES
      -- ── Bài Quyền tay không ──
      ('standard-forms', 'quyen-ngoc-tran',   '{"id":"quyen-ngoc-tran","code":"ngoc_tran","name":"Ngọc Trản quyền","category":"quyen_tay_khong","origin":"Bình Định","difficulty":5,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'quyen-lao-mai',     '{"id":"quyen-lao-mai","code":"lao_mai","name":"Lão Mai quyền","category":"quyen_tay_khong","origin":"Bình Định","difficulty":6,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'quyen-hung-ke',     '{"id":"quyen-hung-ke","code":"hung_ke","name":"Hùng Kê quyền","category":"quyen_tay_khong","origin":"Bình Định","difficulty":6,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'quyen-tu-hai',      '{"id":"quyen-tu-hai","code":"tu_hai","name":"Tứ Hải quyền","category":"quyen_tay_khong","origin":"Bình Định","difficulty":5,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'quyen-lao-ho',      '{"id":"quyen-lao-ho","code":"lao_ho","name":"Lão Hổ Thượng Sơn quyền","category":"quyen_tay_khong","origin":"Bình Định","difficulty":7,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'quyen-than-dong',   '{"id":"quyen-than-dong","code":"than_dong","name":"Thần Đồng quyền","category":"quyen_tay_khong","origin":"Trung Bộ","difficulty":4,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'quyen-bach-hac',    '{"id":"quyen-bach-hac","code":"bach_hac","name":"Bạch Hạc Lượng Xí quyền","category":"quyen_tay_khong","origin":"Nam Bộ","difficulty":6,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'quyen-ngu-mon',     '{"id":"quyen-ngu-mon","code":"ngu_mon","name":"Ngũ Môn quyền","category":"quyen_tay_khong","origin":"Tổng hợp","difficulty":5,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'quyen-tu-linh',     '{"id":"quyen-tu-linh","code":"tu_linh","name":"Tứ Linh quyền","category":"quyen_tay_khong","origin":"Tổng hợp","difficulty":7,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'quyen-nhat-lo-mai', '{"id":"quyen-nhat-lo-mai","code":"nhat_lo_mai","name":"Nhất Lộ Mai Hoa quyền","category":"quyen_tay_khong","origin":"Bắc Bộ","difficulty":5,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'quyen-tu-tru',      '{"id":"quyen-tu-tru","code":"tu_tru","name":"Tứ Trụ quyền","category":"quyen_tay_khong","origin":"Bình Định","difficulty":6,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'quyen-hanh',        '{"id":"quyen-hanh","code":"hanh_quyen","name":"Hành Quyền","category":"quyen_tay_khong","origin":"Tổng hợp","difficulty":3,"gender":"both","type":"ca_nhan"}'::jsonb),

      -- ── Binh khí ──
      ('standard-forms', 'bk-roi-thuan',      '{"id":"bk-roi-thuan","code":"roi_thuan","name":"Roi Thuận Truyền","category":"binh_khi","weapon":"roi","origin":"Bình Định","difficulty":6,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'bk-roi-mua',        '{"id":"bk-roi-mua","code":"roi_mua","name":"Roi Múa","category":"binh_khi","weapon":"roi","origin":"Tổng hợp","difficulty":4,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'bk-bat-quai-con',   '{"id":"bk-bat-quai-con","code":"bat_quai_con","name":"Bát Quái Côn","category":"binh_khi","weapon":"con","origin":"Tổng hợp","difficulty":7,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'bk-don-kiem',       '{"id":"bk-don-kiem","code":"don_kiem","name":"Đơn Kiếm","category":"binh_khi","weapon":"kiem","origin":"Tổng hợp","difficulty":5,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'bk-song-kiem',      '{"id":"bk-song-kiem","code":"song_kiem","name":"Song Kiếm","category":"binh_khi","weapon":"kiem","origin":"Bình Định","difficulty":7,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'bk-tu-linh-dao',    '{"id":"bk-tu-linh-dao","code":"tu_linh_dao","name":"Tứ Linh Đao","category":"binh_khi","weapon":"dao","origin":"Bình Định","difficulty":6,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'bk-song-dao',       '{"id":"bk-song-dao","code":"song_dao","name":"Song Đao Phá Trận","category":"binh_khi","weapon":"dao","origin":"Nam Bộ","difficulty":7,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'bk-thuong',         '{"id":"bk-thuong","code":"thuong","name":"Thương pháp","category":"binh_khi","weapon":"thuong","origin":"Bình Định","difficulty":6,"gender":"both","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'bk-dai-dao',        '{"id":"bk-dai-dao","code":"dai_dao","name":"Đại Đao","category":"binh_khi","weapon":"dao","origin":"Bình Định","difficulty":8,"gender":"nam","type":"ca_nhan"}'::jsonb),
      ('standard-forms', 'bk-con-nhi-khuc',   '{"id":"bk-con-nhi-khuc","code":"con_nhi_khuc","name":"Côn Nhị Khúc","category":"binh_khi","weapon":"nhi_khuc","origin":"Tổng hợp","difficulty":6,"gender":"both","type":"ca_nhan"}'::jsonb),

      -- ── Song luyện ──
      ('standard-forms', 'sl-tay-khong',      '{"id":"sl-tay-khong","code":"sl_tay_khong","name":"Song luyện tay không","category":"song_luyen","difficulty":5,"gender":"both","type":"doi"}'::jsonb),
      ('standard-forms', 'sl-vu-khi',         '{"id":"sl-vu-khi","code":"sl_vu_khi","name":"Song luyện vũ khí","category":"song_luyen","difficulty":7,"gender":"both","type":"doi"}'::jsonb),
      ('standard-forms', 'sl-con-dao',        '{"id":"sl-con-dao","code":"sl_con_vs_dao","name":"Đối luyện Côn vs Đao","category":"song_luyen","difficulty":7,"gender":"both","type":"doi"}'::jsonb),
      ('standard-forms', 'sl-kiem-thuong',    '{"id":"sl-kiem-thuong","code":"sl_kiem_vs_thuong","name":"Đối luyện Kiếm vs Thương","category":"song_luyen","difficulty":8,"gender":"both","type":"doi"}'::jsonb),
      ('standard-forms', 'sl-tay-dao',        '{"id":"sl-tay-dao","code":"sl_tay_vs_dao","name":"Đối luyện Tay không vs Đao","category":"song_luyen","difficulty":6,"gender":"both","type":"doi"}'::jsonb),

      -- ── Đồng đội ──
      ('standard-forms', 'dd-quyen-nam',      '{"id":"dd-quyen-nam","code":"dd_quyen_nam","name":"Đồng đội quyền nam","category":"dong_doi","team_size":5,"gender":"nam","type":"dong_doi"}'::jsonb),
      ('standard-forms', 'dd-quyen-nu',       '{"id":"dd-quyen-nu","code":"dd_quyen_nu","name":"Đồng đội quyền nữ","category":"dong_doi","team_size":5,"gender":"nu","type":"dong_doi"}'::jsonb),
      ('standard-forms', 'dd-quyen-hh',       '{"id":"dd-quyen-hh","code":"dd_quyen_hh","name":"Đồng đội quyền hỗn hợp","category":"dong_doi","team_size":5,"gender":"hon_hop","type":"dong_doi"}'::jsonb),
      ('standard-forms', 'dd-binh-khi',       '{"id":"dd-binh-khi","code":"dd_binh_khi","name":"Đồng đội binh khí","category":"dong_doi","team_size":3,"gender":"both","type":"dong_doi"}'::jsonb),

      -- ── Notification Templates ──
      ('notification-templates', 'tpl-registration-submitted', '{"id":"tpl-registration-submitted","code":"registration_submitted","type":"info","category":"registration","title":"Đăng ký thi đấu đã được gửi","body":"Đoàn {{team_name}} đã gửi đăng ký cho VĐV {{athlete_name}} tại nội dung {{event_name}}.","recipients":["delegate","btc"],"channels":["in_app","email"]}'::jsonb),
      ('notification-templates', 'tpl-registration-approved',  '{"id":"tpl-registration-approved","code":"registration_approved","type":"success","category":"registration","title":"Đăng ký đã được duyệt","body":"Đăng ký của VĐV {{athlete_name}} tại {{event_name}} đã được phê duyệt.","recipients":["delegate","athlete"],"channels":["in_app","email","push"]}'::jsonb),
      ('notification-templates', 'tpl-registration-rejected',  '{"id":"tpl-registration-rejected","code":"registration_rejected","type":"warning","category":"registration","title":"Đăng ký bị từ chối","body":"Đăng ký của VĐV {{athlete_name}} tại {{event_name}} bị từ chối. Lý do: {{reason}}.","recipients":["delegate"],"channels":["in_app","email"]}'::jsonb),
      ('notification-templates', 'tpl-weigh-in-reminder',      '{"id":"tpl-weigh-in-reminder","code":"weigh_in_reminder","type":"info","category":"weigh_in","title":"Nhắc nhở cân kỹ thuật","body":"VĐV {{athlete_name}} cần cân kỹ thuật lúc {{time}} ngày {{date}} tại {{location}}.","recipients":["delegate","athlete"],"channels":["in_app","push"]}'::jsonb),
      ('notification-templates', 'tpl-schedule-published',     '{"id":"tpl-schedule-published","code":"schedule_published","type":"info","category":"schedule","title":"Lịch thi đấu đã công bố","body":"Lịch thi đấu ngày {{date}} đã được công bố.","recipients":["delegate","referee","athlete"],"channels":["in_app","push"]}'::jsonb),
      ('notification-templates', 'tpl-match-result',           '{"id":"tpl-match-result","code":"match_result","type":"info","category":"match","title":"Kết quả trận đấu","body":"{{match_name}}: {{winner_name}} ({{winner_team}}) thắng {{result_detail}}.","recipients":["delegate","athlete"],"channels":["in_app"]}'::jsonb),
      ('notification-templates', 'tpl-appeal-submitted',       '{"id":"tpl-appeal-submitted","code":"appeal_submitted","type":"warning","category":"appeal","title":"Khiếu nại mới","body":"Đoàn {{team_name}} đã gửi {{appeal_type}}: {{reason}}.","recipients":["btc","federation_admin"],"channels":["in_app","email"]}'::jsonb),
      ('notification-templates', 'tpl-discipline-action',      '{"id":"tpl-discipline-action","code":"discipline_action","type":"error","category":"discipline","title":"Quyết định kỷ luật","body":"{{subject_name}} bị xử lý: {{sanction_type}}. Lý do: {{reason}}.","recipients":["delegate","athlete","federation_admin"],"channels":["in_app","email"]}'::jsonb),
      ('notification-templates', 'tpl-cert-issued',            '{"id":"tpl-cert-issued","code":"certification_issued","type":"success","category":"certification","title":"Chứng chỉ đã cấp","body":"{{holder_name}} được cấp {{cert_type}}: {{cert_number}}.","recipients":["athlete","coach","province_admin"],"channels":["in_app","email"]}'::jsonb),
      ('notification-templates', 'tpl-cert-expiring',          '{"id":"tpl-cert-expiring","code":"certification_expiring","type":"warning","category":"certification","title":"Chứng chỉ sắp hết hạn","body":"Chứng chỉ {{cert_type}} ({{cert_number}}) sẽ hết hạn vào {{valid_until}}.","recipients":["athlete","coach","province_admin"],"channels":["in_app","email"]}'::jsonb),
      ('notification-templates', 'tpl-tournament-announce',    '{"id":"tpl-tournament-announce","code":"tournament_announcement","type":"info","category":"tournament","title":"Thông báo giải đấu","body":"{{tournament_name}} diễn ra từ {{start_date}} đến {{end_date}} tại {{venue}}.","recipients":["province_admin","delegate","coach"],"channels":["in_app","email","push"]}'::jsonb)

    ON CONFLICT (entity, id) DO UPDATE
      SET payload = EXCLUDED.payload,
          updated_at = NOW();
  $sql$;
END
$seed$;

COMMIT;
