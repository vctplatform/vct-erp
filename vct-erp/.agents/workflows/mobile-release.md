---
description: Mobile app release workflow - Expo EAS Build to TestFlight/Play Console and App Store/Play Store submission
---

# /mobile-release - Mobile App Release Workflow

## Prerequisites
- Expo EAS CLI installed: `npm install -g eas-cli`
- EAS account authenticated: `eas login`
- Apple Developer account (iOS)
- Google Play Console account (Android)
- `app.json` / `app.config.ts` configured

## Steps

### Step 1: Update version in app config
```bash
cd apps/mobile

# Update version in app.json or app.config.ts
# "version": "1.2.0"
# "ios.buildNumber": "42"
# "android.versionCode": 42
```

### Step 2: Run pre-release checks
```bash
cd apps/mobile

# Check Expo doctor
npx expo-doctor

# Run tests
npm run test -- --run

# Type check
npx tsc --noEmit

# Lint
npm run lint
```

### Step 3: Configure EAS build profiles
```json
// eas.json
{
  "cli": { "version": ">= 14.0.0" },
  "build": {
    "development": {
      "developmentClient": true,
      "distribution": "internal",
      "ios": { "simulator": true },
      "env": {
        "EXPO_PUBLIC_SUPABASE_URL": "https://dev-project.supabase.co",
        "EXPO_PUBLIC_API_URL": "https://staging-api.vct-platform.com"
      }
    },
    "preview": {
      "distribution": "internal",
      "env": {
        "EXPO_PUBLIC_SUPABASE_URL": "https://staging-project.supabase.co",
        "EXPO_PUBLIC_API_URL": "https://staging-api.vct-platform.com"
      }
    },
    "production": {
      "env": {
        "EXPO_PUBLIC_SUPABASE_URL": "https://prod-project.supabase.co",
        "EXPO_PUBLIC_API_URL": "https://api.vct-platform.com"
      }
    }
  },
  "submit": {
    "production": {
      "ios": {
        "appleId": "your@email.com",
        "ascAppId": "1234567890"
      },
      "android": {
        "serviceAccountKeyPath": "./play-store-key.json"
      }
    }
  }
}
```

### Step 4: Build for internal testing (Preview)
```bash
cd apps/mobile

# Build for both platforms (internal distribution)
eas build --profile preview --platform all

# Or specific platform
eas build --profile preview --platform ios
eas build --profile preview --platform android
```

### Step 5: Test internal build
```bash
# Install on device via QR code from EAS dashboard
echo "EAS Dashboard: https://expo.dev/accounts/YOUR_ACCOUNT/projects/vct-platform/builds"

# Verify checklist:
# - [ ] App launches without crash
# - [ ] Supabase Auth works (login/register)
# - [ ] Core navigation works
# - [ ] Data loads correctly
# - [ ] Push notifications work
# - [ ] Offline handling works
# - [ ] Vietnamese text renders correctly
# - [ ] Deep links work
```

### Step 6: Build for production
```bash
cd apps/mobile

# Production build
eas build --profile production --platform all

# Wait for build to complete
eas build:list --platform all --status in-progress
```

### Step 7: Submit to App Store (iOS)
```bash
cd apps/mobile

# Auto-submit to App Store Connect
eas submit --platform ios --latest

# Or submit specific build
eas submit --platform ios --id BUILD_ID

# Manual steps in App Store Connect:
# 1. Add release notes (Vietnamese + English)
# 2. Update screenshots if needed
# 3. Set pricing and availability
# 4. Submit for review
```

### Step 8: Submit to Play Store (Android)
```bash
cd apps/mobile

# Auto-submit to Google Play Console
eas submit --platform android --latest

# Manual steps in Play Console:
# 1. Upload to Internal Testing first
# 2. Promote to Closed Testing (beta)
# 3. Add release notes (Vietnamese + English)
# 4. Promote to Production
```

### Step 9: OTA Update (for JS-only changes)
```bash
cd apps/mobile

# Publish OTA update (no app store review needed)
eas update --branch production --message "Fix: scoring display bug"

# Check update status
eas update:list --branch production
```

### Step 10: Post-release verification
```bash
echo "=== Mobile Release Checklist ==="
echo ""
echo "iOS:"
echo "  - [ ] App Store review submitted"
echo "  - [ ] TestFlight build available"
echo "  - [ ] Crashlytics monitoring active"
echo ""
echo "Android:"
echo "  - [ ] Play Console review submitted"
echo "  - [ ] Internal testing verified"
echo "  - [ ] ANR/Crash monitoring active"
echo ""
echo "Both:"
echo "  - [ ] Version number matches backend expectations"
echo "  - [ ] Supabase project URL is production"
echo "  - [ ] API URL points to production"
echo "  - [ ] Push notification certificates valid"
echo "  - [ ] Deep link domains verified"
```

## Release Timeline

| Day | Activity |
|-----|---------|
| D-5 | Feature freeze, start QA |
| D-3 | Preview build for internal testing |
| D-2 | Fix critical bugs found in QA |
| D-1 | Production build, final testing |
| D-0 | Submit to App Store + Play Console |
| D+1 | iOS review (typically 24-48h) |
| D+3 | Both platforms live, monitor crashes |

## Troubleshooting

| Issue | Solution |
|-------|---------|
| EAS build fails | Check `eas build:inspect` for errors |
| iOS provisioning error | `eas credentials` to manage certificates |
| Android signing error | Verify keystore in `eas.json` credentials |
| App Store rejection | Check rejection reason, fix, resubmit |
| OTA update not applied | Check `expo-updates` config, verify branch |
| Large app size | Use `npx expo-optimize`, check assets |
