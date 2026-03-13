---
name: ui-ux-designer
description: UI/UX Designer role - Design system, component library, accessibility, responsive design, color schemes, typography, and user experience patterns for VCT Platform.
---

# UI/UX Designer - VCT Platform

## Role Overview
Designs the user experience and visual interface of the VCT Platform for both web (React 20 + Tailwind CSS v4) and mobile (Expo React Native). Responsible for design system, accessibility, responsive design, cross-platform consistency, and ensuring a premium, modern look and feel.

## Design System

### 1. Tailwind CSS v4 Theme Configuration

```css
/* apps/web/src/styles/app.css */
@import "tailwindcss";

@theme {
  /* === Brand Colors (Vietnamese Flag inspired + Sports Energy) === */

  /* Primary - Energetic Red (Vietnamese flag) */
  --color-primary-50: #FFF5F5;
  --color-primary-100: #FED7D7;
  --color-primary-200: #FEB2B2;
  --color-primary-300: #FC8181;
  --color-primary-400: #F56565;
  --color-primary-500: #E53E3E;
  --color-primary-600: #C53030;
  --color-primary-700: #9B2C2C;
  --color-primary-800: #822727;
  --color-primary-900: #63171B;

  /* Secondary - Gold (Achievement/Medal) */
  --color-secondary-50: #FFFFF0;
  --color-secondary-500: #D69E2E;
  --color-secondary-700: #975A16;

  /* Accent - Ocean Blue (Triathlon/Water) */
  --color-accent-50: #EBF8FF;
  --color-accent-500: #3182CE;
  --color-accent-700: #2B6CB0;

  /* Semantic */
  --color-success: #38A169;
  --color-warning: #DD6B20;
  --color-error: #E53E3E;
  --color-info: #3182CE;

  /* Dark Mode */
  --color-dark-bg: #0F1117;
  --color-dark-surface: #1A1D2E;
  --color-dark-border: #2D3748;

  /* === Typography === */
  --font-heading: 'Inter', 'Segoe UI', system-ui, sans-serif;
  --font-body: 'Inter', 'Segoe UI', system-ui, sans-serif;
  --font-mono: 'JetBrains Mono', 'Fira Code', monospace;

  /* === Spacing & Layout === */
  --radius-sm: 0.25rem;
  --radius-md: 0.5rem;
  --radius-lg: 0.75rem;
  --radius-xl: 1rem;
  --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

/* Dark mode variant */
@variant dark (&:where(.dark, .dark *));
```

### 2. Cross-Platform Typography

#### Web (Tailwind v4 classes)
```
text-xs   → 12px    text-sm  → 14px    text-base → 16px
text-lg   → 18px    text-xl  → 20px    text-2xl  → 24px
text-3xl  → 30px    text-4xl → 36px

font-normal → 400    font-medium → 500
font-semibold → 600  font-bold → 700
```

#### Mobile (React Native / NativeWind)
```tsx
// Same Tailwind classes via NativeWind v4
<Text className="text-lg font-semibold text-gray-900 dark:text-white">
  {athlete.name}
</Text>
```

### 3. Component Design Patterns

#### Admin Dashboard Layout (Web)
```
┌──────────────────────────────────────────────────┐
│  🏠 VCT Platform           🔔  👤 Admin  🌐 VI │  ← Top Bar
├──────┬───────────────────────────────────────────┤
│      │                                           │
│  📊  │  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐    │  ← Stat Cards
│  Nav │  │ VĐV  │ │ CLB  │ │ Giải │ │ Kết  │    │
│      │  │ 1,234│ │  56  │ │  12  │ │ quả  │    │
│  🏃  │  └──────┘ └──────┘ └──────┘ └──────┘    │
│  🏢  │                                           │
│  🏆  │  ┌─────────────────┐ ┌─────────────────┐ │  ← Charts
│  ⚙️  │  │ Registration    │ │ Performance     │ │
│      │  │ Trend Chart     │ │ Analytics       │ │
│      │  └─────────────────┘ └─────────────────┘ │
│      │                                           │
│      │  ┌─────────────────────────────────────┐ │  ← Data Table
│      │  │ Recent Registrations               │ │
│      │  │ Name | Club | Date | Status         │ │
│      │  └─────────────────────────────────────┘ │
└──────┴───────────────────────────────────────────┘
```

#### Mobile Layout (Expo)
```
┌─────────────────────┐
│  VCT Platform    🔔 │  ← Header
├─────────────────────┤
│  ┌────┐ ┌────┐      │
│  │VĐV │ │CLB │      │  ← Stat Cards (2-col grid)
│  │1234│ │ 56 │      │
│  └────┘ └────┘      │
│  ┌────┐ ┌────┐      │
│  │Giải│ │Kết │      │
│  │ 12 │ │quả │      │
│  └────┘ └────┘      │
│                      │
│  Recent Athletes     │  ← FlatList
│  ┌──────────────┐   │
│  │ 🏃 Nguyễn A  │   │
│  │ CLB Hà Nội   │   │
│  └──────────────┘   │
│  ┌──────────────┐   │
│  │ 🏃 Trần B    │   │
│  │ CLB TP.HCM   │   │
│  └──────────────┘   │
├─────────────────────┤
│  🏠  🏃  🏆  👤   │  ← Tab Bar (Expo Router)
└─────────────────────┘
```

### 4. Micro-Animations (Tailwind v4 utilities)

```html
<!-- Button hover effect -->
<button class="transition-all duration-200 ease-in-out hover:-translate-y-0.5 hover:shadow-md active:translate-y-0">
  Submit
</button>

<!-- Card hover effect -->
<div class="transition-all duration-300 ease-[cubic-bezier(0.4,0,0.2,1)] hover:-translate-y-1 hover:shadow-xl">
  Card content
</div>

<!-- Skeleton loading -->
<div class="animate-pulse bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 bg-[length:200%_100%]">
  Loading...
</div>
```

### 5. Accessibility (WCAG 2.1 AA)

| Requirement | Standard | Implementation |
|-------------|----------|---------------|
| Color contrast | ≥ 4.5:1 (text), ≥ 3:1 (large) | Test with contrast checker |
| Focus indicators | Visible on all interactive elements | `focus-visible:ring-2 focus-visible:ring-primary-500` |
| Screen reader | ARIA labels, landmarks, live regions | `aria-label`, `role` attributes |
| Keyboard nav | Tab through all interactive elements | `tabindex`, focus management |
| Alt text | All meaningful images | `alt` attribute |
| Form labels | All inputs have labels | `<label>` with `htmlFor` |
| Error messages | Associated with form fields | `aria-describedby` |
| Vietnamese text | Proper diacritics rendering | UTF-8, proper font support |
| Mobile a11y | VoiceOver (iOS), TalkBack (Android) | `accessibilityLabel` in RN |

### 6. Design Checklist
- [ ] Tailwind v4 `@theme` colors follow brand palette
- [ ] Typography uses Tailwind v4 scale
- [ ] Spacing follows 4px/8px grid (Tailwind defaults)
- [ ] All states designed: default, hover, active, focus, disabled, loading, error, empty
- [ ] Dark mode supported (Tailwind `dark:` variant)
- [ ] Web tested at 320px, 768px, 1024px, 1440px
- [ ] Mobile tested on iOS (iPhone SE+) and Android (Pixel 4+)
- [ ] Accessibility: contrast ratio ≥ 4.5:1
- [ ] Accessibility: keyboard navigation works (web)
- [ ] Accessibility: VoiceOver/TalkBack works (mobile)
- [ ] Vietnamese text renders correctly (ă, ê, ơ, ư, etc.)
- [ ] Consistent icon style (Lucide for web, Expo Vector Icons for mobile)
- [ ] Loading skeletons for async content
- [ ] Toast notifications for actions
- [ ] Confirmation dialogs for destructive actions
- [ ] Cross-platform visual consistency (web ↔ mobile)
