# Finance Simulation Snapshot

Company: `VCT_SIM`

## Trial Balance Snapshot

| Account | Name | Normal | Opening | Debit | Credit | Closing |
| --- | --- | --- | --- | --- | --- | --- |
| 1111 | Tien mat Viet Nam | debit | 0 | 2746341585.0000 | 43164154.9500 | 2703177430.0500 |
| 1121 | Tien gui ngan hang Viet Nam dong | debit | 0 | 7310612755.0000 | 2262115587.5800 | 5048497167.4200 |
| 131 | Phai thu cua khach hang | debit | 0 | 453695518.0000 | 255416151.0000 | 198279367.0000 |
| 3387 | Doanh thu chua thuc hien | credit | 0 | 6006316003.0000 | 9530237952.0000 | 3523921949.0000 |
| 3388 | Phai tra, phai nop khac | credit | 0 | 43263677.0000 | 43263677.0000 | 0.0000 |
| 5111 | Doanh thu ban hang hoa | credit | 0 | 8028500.0000 | 228036560.0000 | 220008060.0000 |
| 5113 | Doanh thu cung cap dich vu | credit | 0 | 0 | 6460011521.0000 | 6460011521.0000 |
| 5211 | Chiet khau thuong mai | debit | 0 | 11556130.0000 | 0 | 11556130.0000 |
| 632 | Gia von hang ban | debit | 0 | 1991918128.7200 | 0 | 1991918128.7200 |
| 6421 | Chi phi nhan vien quan ly | debit | 0 | 113538187.0000 | 0 | 113538187.0000 |
| 6422 | Chi phi vat lieu quan ly | debit | 0 | 144248309.0000 | 0 | 144248309.0000 |
| 711 | Thu nhap khac | credit | 0 | 0 | 7273189.1900 | 7273189.1900 |

## P&L B02-DN Snapshot

| Line | Description | Amount |
| --- | --- | --- |
| 01 | Doanh thu ban hang va cung cap dich vu | 6680019581.0000 |
| 02 | Cac khoan giam tru doanh thu | 11556130.0000 |
| 10 | Loi nhuan gop ve ban hang va cung cap dich vu | 4676545322.2800 |
| 11 | Gia von hang ban | 1991918128.7200 |
| 25 | Chi phi quan ly doanh nghiep | 257786496.0000 |
| 30 | Loi nhuan thuan tu hoat dong kinh doanh | -257786496.0000 |
| 31 | Thu nhap khac | 7273189.1900 |
| 40 | Tong loi nhuan ke toan truoc thue | 4426032015.4700 |
| 60 | Loi nhuan sau thue thu nhap doanh nghiep | 4426032015.4700 |

## Gross Profit By Cost Center

| Cost Center | Gross Revenue | Deductions | Other Income | COGS | Gross Profit |
| --- | --- | --- | --- | --- | --- |
| dojo | 453695518.0000 | 0 | 0 | 139557562.5600 | 314137955.4400 |
| rental | 0 | 0 | 7273189.1900 | 0 | 7273189.1900 |
| retail | 220008060.0000 | 11556130.0000 | 0 | 136917734.8000 | 71534195.2000 |
| saas | 6006316003.0000 | 0 | 0 | 1715442831.3600 | 4290873171.6400 |

## Revenue Stream JSON Snapshot

| Cost Center | Net Revenue |
| --- | --- |
| dojo | 453695518.0000 |
| rental | 7273189.1900 |
| retail | 208451930.0000 |
| saas | 6006316003.0000 |

## Cash Runway Snapshot

Current cash: `7759703097.4700`

Average monthly burn: `138457054.9867`

| Month | Opening Cash | Contracted Inflow | Projected Burn | Projected Ending |
| --- | --- | --- | --- | --- |
| 2026-03 | 7759703097.4700 | 0.0000 | 138457054.9867 | 7621246042.4833 |
| 2026-04 | 7621246042.4833 | 592409290.0000 | 138457054.9867 | 8075198277.4966 |
| 2026-05 | 8075198277.4966 | 541610514.0000 | 138457054.9867 | 8478351736.5099 |

## Partition Distribution

| Partition | Rows |
| --- | --- |
| journal_items_2026_q1 | 8048 |

## Consistency Log

- Base capture operations seeded: 1015
- Idempotency duplicates attempted: 101
- Idempotency duplicates blocked as replay: 101
- Void cases executed: 5
- Refund cases executed: 30
- SaaS recognitions posted through 2026-03-31: 2294
- Trial balance validation: total debit equals total credit and closing balance formula passed.
