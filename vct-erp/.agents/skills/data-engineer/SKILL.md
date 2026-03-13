---
name: data-engineer
description: Data Engineer role - Reporting, analytics, data export (PDF/Excel/CSV), ETL patterns, dashboards, and data visualization for VCT Platform.
---

# Data Engineer - VCT Platform

## Role Overview
Designs and implements data pipelines, reporting, analytics dashboards, and data export features. Responsible for transforming raw database records into actionable business intelligence for federation administrators, club managers, and organizing committees.

## Technology Stack
- **Database**: Neon PostgreSQL 18+ (analytical queries, JSON_TABLE, window functions)
- **Backend**: Go 1.26 (report generation, export workers)
- **PDF Generation**: wkhtmltopdf / go-pdf / headless Chrome
- **Excel Export**: excelize (Go) / exceljs (Node)
- **CSV Export**: encoding/csv (Go) / native JS
- **Charts (Web)**: Recharts / Chart.js
- **Charts (Mobile)**: react-native-chart-kit / Victory Native
- **Caching**: Redis 7+ (report cache, pre-computed aggregates)
- **Scheduling**: cron / pg_cron (Supabase)
- **BI Integration**: Metabase (optional, connects to Neon)

## Core Patterns

### 1. Analytical SQL Patterns (PostgreSQL 18+)

#### Window Functions for Rankings
```sql
-- Athlete rankings by category and gender
SELECT
    a.id,
    a.first_name || ' ' || a.last_name AS name,
    c.name AS club_name,
    cat.name AS category,
    a.gender,
    SUM(rr.points) AS total_points,
    RANK() OVER (
        PARTITION BY cat.id, a.gender
        ORDER BY SUM(rr.points) DESC
    ) AS rank,
    COUNT(rr.id) AS races_completed
FROM athletes a
JOIN race_results rr ON rr.athlete_id = a.id
JOIN races r ON r.id = rr.race_id
JOIN categories cat ON cat.id = r.category_id
JOIN clubs c ON c.id = a.club_id
WHERE rr.status = 'finished'
AND a.deleted_at IS NULL
GROUP BY a.id, a.first_name, a.last_name, c.name, cat.id, cat.name, a.gender
ORDER BY cat.name, a.gender, rank;
```

#### Tournament Statistics
```sql
-- Tournament summary statistics
SELECT
    t.id,
    t.name,
    t.start_date,
    COUNT(DISTINCT tr.athlete_id) AS total_participants,
    COUNT(DISTINCT r.id) AS total_races,
    COUNT(DISTINCT CASE WHEN rr.status = 'finished' THEN rr.athlete_id END) AS finishers,
    COUNT(DISTINCT CASE WHEN rr.status = 'DNF' THEN rr.athlete_id END) AS dnf_count,
    COUNT(DISTINCT CASE WHEN rr.status = 'DNS' THEN rr.athlete_id END) AS dns_count,
    ROUND(AVG(EXTRACT(EPOCH FROM rr.finish_time - rr.start_time))::numeric, 2) AS avg_time_seconds,
    MIN(rr.finish_time - rr.start_time) AS fastest_time,
    COUNT(DISTINCT c.id) AS clubs_represented,
    COUNT(DISTINCT c.province_id) AS provinces_represented
FROM tournaments t
LEFT JOIN races r ON r.tournament_id = t.id
LEFT JOIN tournament_registrations tr ON tr.tournament_id = t.id
LEFT JOIN race_results rr ON rr.race_id = r.id
LEFT JOIN athletes a ON a.id = rr.athlete_id
LEFT JOIN clubs c ON c.id = a.club_id
WHERE t.id = $1
GROUP BY t.id;
```

#### JSON_TABLE for Complex Reports
```sql
-- Parse tournament scoring rules from JSON config
SELECT t.name, jt.*
FROM tournaments t,
  JSON_TABLE(t.scoring_config, '$.categories[*]'
    COLUMNS (
      category_name TEXT PATH '$.name',
      points_1st INT PATH '$.points[0]',
      points_2nd INT PATH '$.points[1]',
      points_3rd INT PATH '$.points[2]',
      bonus_points INT PATH '$.bonus' DEFAULT '0' ON EMPTY
    )
  ) AS jt
WHERE t.id = $1;
```

### 2. Report Generation (Go Backend)

```go
// internal/modules/report/usecase/generate_report.go

type ReportUseCase struct {
    athleteRepo AthleteRepository
    resultRepo  ResultRepository
    cache       *redis.Client
    pdfGen      PDFGenerator
    excelGen    ExcelGenerator
}

type ReportRequest struct {
    Type       string    `json:"type"`       // "athlete_ranking", "tournament_summary", "club_stats"
    Format     string    `json:"format"`     // "json", "pdf", "xlsx", "csv"
    Filters    Filters   `json:"filters"`
    DateRange  DateRange `json:"date_range"`
    OrgID      string    `json:"org_id"`
}

func (uc *ReportUseCase) Generate(ctx context.Context, req ReportRequest) ([]byte, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("report:%s:%s:%v", req.Type, req.OrgID, req.Filters)
    if cached, err := uc.cache.Get(ctx, cacheKey).Bytes(); err == nil {
        return cached, nil
    }

    // Generate report data
    var data interface{}
    switch req.Type {
    case "athlete_ranking":
        data, _ = uc.generateAthleteRanking(ctx, req)
    case "tournament_summary":
        data, _ = uc.generateTournamentSummary(ctx, req)
    case "club_stats":
        data, _ = uc.generateClubStats(ctx, req)
    }

    // Format output
    var output []byte
    switch req.Format {
    case "json":
        output, _ = json.Marshal(data)
    case "pdf":
        output, _ = uc.pdfGen.Generate(req.Type, data)
    case "xlsx":
        output, _ = uc.excelGen.Generate(req.Type, data)
    case "csv":
        output, _ = generateCSV(data)
    }

    // Cache for 15 minutes
    uc.cache.Set(ctx, cacheKey, output, 15*time.Minute)

    return output, nil
}
```

### 3. Excel Export (Go - excelize)

```go
import "github.com/xuri/excelize/v2"

func (g *ExcelGenerator) GenerateAthleteRanking(rankings []AthleteRanking) ([]byte, error) {
    f := excelize.NewFile()
    sheet := "Bảng Xếp Hạng VĐV"
    f.SetSheetName("Sheet1", sheet)

    // Header styling
    headerStyle, _ := f.NewStyle(&excelize.Style{
        Font:      &excelize.Font{Bold: true, Size: 12, Color: "#FFFFFF"},
        Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#C53030"}},
        Alignment: &excelize.Alignment{Horizontal: "center"},
    })

    // Headers
    headers := []string{"Hạng", "Họ tên", "CLB", "Hạng mục", "Giới tính", "Điểm", "Số giải"}
    for i, h := range headers {
        cell, _ := excelize.CoordinatesToCellName(i+1, 1)
        f.SetCellValue(sheet, cell, h)
        f.SetCellStyle(sheet, cell, cell, headerStyle)
    }

    // Data rows
    for i, r := range rankings {
        row := i + 2
        f.SetCellValue(sheet, fmt.Sprintf("A%d", row), r.Rank)
        f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.Name)
        f.SetCellValue(sheet, fmt.Sprintf("C%d", row), r.ClubName)
        f.SetCellValue(sheet, fmt.Sprintf("D%d", row), r.Category)
        f.SetCellValue(sheet, fmt.Sprintf("E%d", row), r.Gender)
        f.SetCellValue(sheet, fmt.Sprintf("F%d", row), r.TotalPoints)
        f.SetCellValue(sheet, fmt.Sprintf("G%d", row), r.RacesCompleted)
    }

    // Auto-fit columns
    for i := range headers {
        col, _ := excelize.ColumnNumberToName(i + 1)
        f.SetColWidth(sheet, col, col, 18)
    }

    buf, _ := f.WriteToBuffer()
    return buf.Bytes(), nil
}
```

### 4. CSV Export Pattern

```go
func generateCSV(athletes []AthleteRanking) ([]byte, error) {
    var buf bytes.Buffer
    // BOM for UTF-8 Excel compatibility (Vietnamese characters)
    buf.Write([]byte{0xEF, 0xBB, 0xBF})

    writer := csv.NewWriter(&buf)
    writer.Write([]string{"Hạng", "Họ tên", "CLB", "Điểm", "Số giải"})

    for _, a := range athletes {
        writer.Write([]string{
            strconv.Itoa(a.Rank),
            a.Name,
            a.ClubName,
            strconv.Itoa(a.TotalPoints),
            strconv.Itoa(a.RacesCompleted),
        })
    }

    writer.Flush()
    return buf.Bytes(), nil
}
```

### 5. Dashboard Data API

```go
// GET /api/v1/dashboard/stats
type DashboardStats struct {
    TotalAthletes    int            `json:"total_athletes"`
    TotalClubs       int            `json:"total_clubs"`
    ActiveTournaments int           `json:"active_tournaments"`
    RecentResults    int            `json:"recent_results"`
    TrendData        []TrendPoint   `json:"trend_data"`
    TopAthletes      []TopAthlete   `json:"top_athletes"`
    ClubDistribution []ClubDistr    `json:"club_distribution"`
}

type TrendPoint struct {
    Month string `json:"month"`
    Count int    `json:"count"`
}
```

### 6. Scheduled Reports (pg_cron via Supabase)

```sql
-- Auto-generate monthly statistics (runs on 1st of each month)
SELECT cron.schedule(
    'monthly-stats',
    '0 0 1 * *',
    $$
    INSERT INTO monthly_reports (month, year, data)
    SELECT
        EXTRACT(MONTH FROM NOW() - INTERVAL '1 month'),
        EXTRACT(YEAR FROM NOW() - INTERVAL '1 month'),
        json_build_object(
            'new_athletes', (SELECT count(*) FROM athletes WHERE created_at >= date_trunc('month', NOW() - INTERVAL '1 month') AND created_at < date_trunc('month', NOW())),
            'new_clubs', (SELECT count(*) FROM clubs WHERE created_at >= date_trunc('month', NOW() - INTERVAL '1 month') AND created_at < date_trunc('month', NOW())),
            'tournaments_held', (SELECT count(*) FROM tournaments WHERE start_date >= date_trunc('month', NOW() - INTERVAL '1 month') AND start_date < date_trunc('month', NOW())),
            'total_participants', (SELECT count(DISTINCT athlete_id) FROM tournament_registrations tr JOIN tournaments t ON t.id = tr.tournament_id WHERE t.start_date >= date_trunc('month', NOW() - INTERVAL '1 month'))
        );
    $$
);
```

### 7. Frontend Chart Patterns (Recharts)

```tsx
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';

function AthleteRegistrationChart({ data }: { data: TrendPoint[] }) {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={data}>
        <XAxis dataKey="month" />
        <YAxis />
        <Tooltip
          formatter={(value: number) => [`${value} VĐV`, 'Đăng ký mới']}
        />
        <Bar
          dataKey="count"
          fill="var(--color-primary-500)"
          radius={[4, 4, 0, 0]}
        />
      </BarChart>
    </ResponsiveContainer>
  );
}
```

### 8. Data Engineer Checklist
- [ ] Analytical queries use proper indexes (covering, partial)
- [ ] Window functions for rankings and statistics
- [ ] Report caching in Redis (TTL: 15 min for real-time, 1h for historical)
- [ ] CSV export includes UTF-8 BOM for Vietnamese in Excel
- [ ] PDF reports include VCT branding and Vietnamese formatting
- [ ] Dashboard APIs return pre-aggregated data
- [ ] pg_cron scheduled for monthly/weekly reports (Supabase)
- [ ] Export endpoints handle large datasets (streaming)
- [ ] Date/number formatting follows Vietnamese conventions
- [ ] Charts accessible (alt text, color-blind safe palette)
