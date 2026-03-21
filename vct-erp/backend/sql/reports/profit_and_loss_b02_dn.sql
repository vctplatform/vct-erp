WITH pnl_source AS (
    SELECT
        a.code,
        CASE
            WHEN a.normal_side = 'credit' THEN ji.amount_signed * -1
            ELSE ji.amount_signed
        END AS natural_delta
    FROM journal_items AS ji
    INNER JOIN journal_entries AS je ON je.id = ji.journal_entry_id
    INNER JOIN accounts AS a ON a.id = ji.account_id
    WHERE je.company_code = $1
      AND je.status = 'posted'
      AND je.posting_date BETWEEN $2 AND $3
      AND (
          a.code LIKE '511%' OR
          a.code LIKE '515%' OR
          a.code LIKE '521%' OR
          a.code LIKE '632%' OR
          a.code LIKE '635%' OR
          a.code LIKE '641%' OR
          a.code LIKE '642%' OR
          a.code LIKE '711%' OR
          a.code LIKE '811%' OR
          a.code LIKE '821%'
      )
),
classified AS (
    SELECT
        CASE
            WHEN code LIKE '511%' THEN '01'
            WHEN code LIKE '521%' THEN '02'
            WHEN code LIKE '632%' THEN '11'
            WHEN code LIKE '515%' THEN '21'
            WHEN code LIKE '635%' THEN '22'
            WHEN code LIKE '641%' THEN '24'
            WHEN code LIKE '642%' THEN '25'
            WHEN code LIKE '711%' THEN '31'
            WHEN code LIKE '811%' THEN '32'
            WHEN code LIKE '821%' THEN '51'
        END AS line_code,
        CASE
            WHEN code LIKE '511%' THEN 'Doanh thu ban hang va cung cap dich vu'
            WHEN code LIKE '521%' THEN 'Cac khoan giam tru doanh thu'
            WHEN code LIKE '632%' THEN 'Gia von hang ban'
            WHEN code LIKE '515%' THEN 'Doanh thu hoat dong tai chinh'
            WHEN code LIKE '635%' THEN 'Chi phi tai chinh'
            WHEN code LIKE '641%' THEN 'Chi phi ban hang'
            WHEN code LIKE '642%' THEN 'Chi phi quan ly doanh nghiep'
            WHEN code LIKE '711%' THEN 'Thu nhap khac'
            WHEN code LIKE '811%' THEN 'Chi phi khac'
            WHEN code LIKE '821%' THEN 'Chi phi thue TNDN hien hanh'
        END AS line_name,
        natural_delta
    FROM pnl_source
),
windowed AS (
    SELECT DISTINCT
        line_code,
        line_name,
        COALESCE(SUM(natural_delta) OVER (PARTITION BY line_code), 0) AS amount
    FROM classified
),
aggregated AS (
    SELECT line_code, line_name, MAX(amount) AS amount
    FROM windowed
    GROUP BY line_code, line_name
),
derived AS (
    SELECT '10' AS line_code, 'Loi nhuan gop ve ban hang va cung cap dich vu' AS line_name,
           COALESCE((SELECT amount FROM aggregated WHERE line_code = '01'), 0)
         - COALESCE((SELECT amount FROM aggregated WHERE line_code = '02'), 0)
         - COALESCE((SELECT amount FROM aggregated WHERE line_code = '11'), 0) AS amount
    UNION ALL
    SELECT '30', 'Loi nhuan thuan tu hoat dong kinh doanh',
           COALESCE((SELECT amount FROM aggregated WHERE line_code = '10'), 0)
         + COALESCE((SELECT amount FROM aggregated WHERE line_code = '21'), 0)
         - COALESCE((SELECT amount FROM aggregated WHERE line_code = '22'), 0)
         - COALESCE((SELECT amount FROM aggregated WHERE line_code = '24'), 0)
         - COALESCE((SELECT amount FROM aggregated WHERE line_code = '25'), 0)
    UNION ALL
    SELECT '40', 'Tong loi nhuan ke toan truoc thue',
           (
               COALESCE((SELECT amount FROM aggregated WHERE line_code = '01'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '02'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '11'), 0)
             + COALESCE((SELECT amount FROM aggregated WHERE line_code = '21'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '22'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '24'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '25'), 0)
             + COALESCE((SELECT amount FROM aggregated WHERE line_code = '31'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '32'), 0)
           )
    UNION ALL
    SELECT '60', 'Loi nhuan sau thue thu nhap doanh nghiep',
           (
               COALESCE((SELECT amount FROM aggregated WHERE line_code = '01'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '02'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '11'), 0)
             + COALESCE((SELECT amount FROM aggregated WHERE line_code = '21'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '22'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '24'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '25'), 0)
             + COALESCE((SELECT amount FROM aggregated WHERE line_code = '31'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '32'), 0)
             - COALESCE((SELECT amount FROM aggregated WHERE line_code = '51'), 0)
           )
)
SELECT line_code, line_name, amount
FROM aggregated
UNION ALL
SELECT line_code, line_name, amount
FROM derived
ORDER BY line_code;
