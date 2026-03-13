---
description: Check and synchronize i18n translation keys between Vietnamese and English
---

# /i18n-sync - Internationalization Sync Workflow

// turbo-all

## When to Use
- After adding new UI text/labels
- Before a release to ensure all translations are complete
- When translation files drift between languages

## Translation File Structure
```
i18n/
├── vi/
│   ├── common.json       ← Shared strings
│   ├── athlete.json      ← Module-specific
│   ├── club.json
│   ├── tournament.json
│   ├── federation.json
│   └── validation.json   ← Form validation messages
└── en/
    ├── common.json
    ├── athlete.json
    ├── club.json
    ├── tournament.json
    ├── federation.json
    └── validation.json
```

## Steps

### Step 1: Find missing translation keys
```bash
# Script to compare VI and EN keys
node -e "
const fs = require('fs');
const path = require('path');

const viDir = 'i18n/vi';
const enDir = 'i18n/en';

const viFiles = fs.readdirSync(viDir).filter(f => f.endsWith('.json'));
const enFiles = fs.readdirSync(enDir).filter(f => f.endsWith('.json'));

// Check for files only in one language
const viOnly = viFiles.filter(f => !enFiles.includes(f));
const enOnly = enFiles.filter(f => !viFiles.includes(f));

if (viOnly.length) console.log('⚠️  Files only in VI:', viOnly);
if (enOnly.length) console.log('⚠️  Files only in EN:', enOnly);

// Compare keys in each file
const getKeys = (obj, prefix = '') =>
  Object.keys(obj).reduce((keys, key) => {
    const fullKey = prefix ? prefix + '.' + key : key;
    if (typeof obj[key] === 'object' && obj[key] !== null)
      return [...keys, ...getKeys(obj[key], fullKey)];
    return [...keys, fullKey];
  }, []);

const commonFiles = viFiles.filter(f => enFiles.includes(f));
let totalMissing = 0;

for (const file of commonFiles) {
  const viData = JSON.parse(fs.readFileSync(path.join(viDir, file), 'utf8'));
  const enData = JSON.parse(fs.readFileSync(path.join(enDir, file), 'utf8'));
  const viKeys = getKeys(viData);
  const enKeys = getKeys(enData);
  const missingInEn = viKeys.filter(k => !enKeys.includes(k));
  const missingInVi = enKeys.filter(k => !viKeys.includes(k));
  if (missingInEn.length || missingInVi.length) {
    console.log('\n📄', file);
    missingInEn.forEach(k => console.log('  ❌ Missing in EN:', k));
    missingInVi.forEach(k => console.log('  ❌ Missing in VI:', k));
    totalMissing += missingInEn.length + missingInVi.length;
  }
}

if (totalMissing === 0) console.log('✅ All translation keys are in sync!');
else console.log('\n⚠️  Total missing keys:', totalMissing);
"
```

### Step 2: Check for untranslated values (EN values in VI files)
```bash
node -e "
const fs = require('fs');
const path = require('path');

const checkUntranslated = (viObj, enObj, prefix = '') => {
  for (const key of Object.keys(viObj)) {
    const fullKey = prefix ? prefix + '.' + key : key;
    if (typeof viObj[key] === 'object') {
      checkUntranslated(viObj[key], enObj?.[key] || {}, fullKey);
    } else if (viObj[key] === enObj?.[key] && /^[a-zA-Z\s]+$/.test(viObj[key])) {
      console.log('⚠️  Possibly untranslated:', fullKey, '=', viObj[key]);
    }
  }
};

const viDir = 'i18n/vi';
const enDir = 'i18n/en';
const files = fs.readdirSync(viDir).filter(f => f.endsWith('.json'));

for (const file of files) {
  const viData = JSON.parse(fs.readFileSync(path.join(viDir, file), 'utf8'));
  const enData = JSON.parse(fs.readFileSync(path.join(enDir, file), 'utf8'));
  console.log('\n📄', file);
  checkUntranslated(viData, enData);
}
"
```

### Step 3: Check for hardcoded strings in source code
```bash
# Search for Vietnamese text directly in TSX/Go files (should use i18n keys)
grep -rn "[àáảãạăắằẳẵặâấầẩẫậèéẻẽẹêếềểễệ]" apps/web/src/ --include="*.tsx" --include="*.ts" | grep -v "i18n" | grep -v "test" | head -20

# Check Go files for hardcoded Vietnamese
grep -rn "[àáảãạăắằẳẵặâấầẩẫậèéẻẽẹêếềểễệ]" backend/internal/ --include="*.go" | grep -v "_test.go" | head -20
```

### Step 4: Verify i18n config
```bash
# Check react-i18next setup
cat apps/web/src/i18n/index.ts

# Check Expo i18n setup
cat apps/mobile/i18n/index.ts
```

### Step 5: Generate translation report
```bash
echo "=== i18n Sync Report ==="
echo "Date: $(date)"
echo ""
echo "VI files: $(ls i18n/vi/*.json | wc -l)"
echo "EN files: $(ls i18n/en/*.json | wc -l)"
echo ""
echo "VI total keys: $(node -e "
const fs = require('fs');
let count = 0;
fs.readdirSync('i18n/vi').filter(f=>f.endsWith('.json')).forEach(f=>{
  const data = JSON.parse(fs.readFileSync('i18n/vi/'+f,'utf8'));
  count += JSON.stringify(data).split('\":\"').length - 1;
});
console.log(count);
")"
echo "EN total keys: $(node -e "
const fs = require('fs');
let count = 0;
fs.readdirSync('i18n/en').filter(f=>f.endsWith('.json')).forEach(f=>{
  const data = JSON.parse(fs.readFileSync('i18n/en/'+f,'utf8'));
  count += JSON.stringify(data).split('\":\"').length - 1;
});
console.log(count);
")"
```

## Translation Style Guide

| Rule | Vietnamese | English |
|------|-----------|---------|
| Button labels | Verb first: "Thêm VĐV" | Verb first: "Add Athlete" |
| Form labels | Noun: "Họ và tên" | Noun: "Full Name" |
| Error messages | Polite: "Vui lòng nhập..." | Direct: "Please enter..." |
| Success messages | "Thêm thành công" | "Added successfully" |
| Placeholders | "Nhập email..." | "Enter email..." |
| Date format | DD/MM/YYYY | MM/DD/YYYY |
| Number format | 1.234.567 | 1,234,567 |
