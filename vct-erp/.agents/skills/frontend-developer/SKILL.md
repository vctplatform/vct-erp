---
name: frontend-developer
description: Frontend Developer role - React 20/TypeScript, Expo React Native, Tailwind CSS v4, Supabase integration, component patterns, state management, i18n, and cross-platform development for VCT Platform.
---

# Frontend Developer - VCT Platform

## Role Overview
Implements the user interface using React 20 and TypeScript for web, and Expo React Native for mobile. Responsible for component architecture, state management, API integration, i18n, cross-platform development, and frontend testing.

## Technology Stack
- **Runtime**: Node.js 25.6.1 (ESM by default, native fetch, built-in test runner)
- **Web Framework**: React 20 with TypeScript (React Compiler, Server Components, Actions)
- **Mobile Framework**: Expo React Native (SDK 53+, Expo Router, file-based routing)
- **Build Tool**: Vite 7+
- **Styling (Web)**: Tailwind CSS v4 (CSS-first config, `@theme` directive, zero-JS)
- **Styling (Mobile)**: NativeWind v4 (Tailwind for React Native)
- **State Management**: React Context + useReducer / Zustand / Supabase Realtime
- **Backend-as-a-Service**: Supabase (Auth, Database, Realtime, Storage)
- **HTTP Client**: TanStack Query v5 + Supabase JS Client
- **Routing (Web)**: React Router v7+ / TanStack Router
- **Routing (Mobile)**: Expo Router v4 (file-based)
- **i18n**: react-i18next / expo-localization
- **Forms**: React Hook Form + Zod validation
- **Charts**: Recharts (web) / Victory Native (mobile)
- **Testing**: Vitest + React Testing Library (web) / Jest + @testing-library/react-native (mobile)
- **Linting**: ESLint 9+ (flat config) + Prettier

## Core Patterns

### 1. Monorepo Project Structure
```
vct-platform/
├── apps/
│   ├── web/                    # React 20 web app
│   │   ├── src/
│   │   │   ├── app/
│   │   │   │   ├── App.tsx
│   │   │   │   ├── Router.tsx
│   │   │   │   └── providers/
│   │   │   │       ├── AuthProvider.tsx      # Supabase Auth
│   │   │   │       ├── ThemeProvider.tsx
│   │   │   │       └── I18nProvider.tsx
│   │   │   ├── modules/
│   │   │   │   └── athlete/
│   │   │   │       ├── pages/
│   │   │   │       ├── components/
│   │   │   │       └── hooks/
│   │   │   └── styles/
│   │   │       └── app.css              # Tailwind v4 entry
│   │   ├── tailwind.config.ts           # Minimal (v4 CSS-first)
│   │   ├── vite.config.ts
│   │   └── package.json
│   └── mobile/                 # Expo React Native app
│       ├── app/                # Expo Router (file-based routing)
│       │   ├── (tabs)/
│       │   │   ├── index.tsx
│       │   │   ├── athletes.tsx
│       │   │   └── tournaments.tsx
│       │   ├── athlete/[id].tsx
│       │   └── _layout.tsx
│       ├── components/
│       ├── app.json
│       └── package.json
├── packages/
│   ├── shared/                 # Shared between web & mobile
│   │   ├── types/
│   │   │   ├── athlete.ts
│   │   │   ├── api.ts
│   │   │   └── common.ts
│   │   ├── hooks/
│   │   │   ├── useAthletes.ts
│   │   │   ├── useAuth.ts      # Supabase auth hook
│   │   │   └── useSupabase.ts
│   │   ├── services/
│   │   │   ├── supabase.ts     # Supabase client init
│   │   │   └── athleteApi.ts
│   │   └── utils/
│   │       ├── formatters.ts
│   │       └── validators.ts
│   └── ui/                     # Shared UI primitives
│       ├── Button.tsx
│       ├── Input.tsx
│       └── Modal.tsx
├── i18n/
│   ├── vi/
│   │   ├── common.json
│   │   └── athlete.json
│   └── en/
├── supabase/                   # Supabase config & migrations
│   ├── config.toml
│   ├── migrations/
│   └── seed.sql
└── package.json                # Workspace root
```

### 2. Tailwind CSS v4 Configuration
```css
/* apps/web/src/styles/app.css */
@import "tailwindcss";

@theme {
  /* Brand Colors - Vietnamese Flag inspired + Sports Energy */
  --color-primary-50: #FFF5F5;
  --color-primary-100: #FED7D7;
  --color-primary-500: #E53E3E;
  --color-primary-600: #C53030;
  --color-primary-700: #9B2C2C;

  /* Secondary - Gold (Achievement/Medal) */
  --color-secondary-500: #D69E2E;
  --color-secondary-700: #975A16;

  /* Accent - Ocean Blue */
  --color-accent-500: #3182CE;
  --color-accent-700: #2B6CB0;

  /* Fonts */
  --font-heading: 'Inter', system-ui, sans-serif;
  --font-body: 'Inter', system-ui, sans-serif;
  --font-mono: 'JetBrains Mono', monospace;

  /* Radius */
  --radius-lg: 0.75rem;
  --radius-xl: 1rem;
}

/* Dark mode using Tailwind v4 automatic dark variant */
@variant dark (&:where(.dark, .dark *));
```

### 3. React 20 Component Pattern
```tsx
// AthleteCard.tsx - Using React 20 features
import { use, useOptimistic } from 'react';
import { useTranslation } from 'react-i18next';

interface AthleteCardProps {
  athlete: Athlete;
  onEdit?: (id: string) => void;
  onDelete?: (id: string) => void;
  className?: string;
}

export function AthleteCard({ athlete, onEdit, onDelete, className }: AthleteCardProps) {
  const { t } = useTranslation('athlete');

  return (
    <div className={`rounded-xl border bg-white p-6 shadow-sm transition-all hover:-translate-y-1 hover:shadow-lg dark:bg-gray-800 ${className}`}>
      <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
        {athlete.firstName} {athlete.lastName}
      </h3>
      <p className="text-sm text-gray-500">
        {t('status')}: {t(`status.${athlete.status}`)}
      </p>
      <div className="mt-4 flex gap-2">
        {onEdit && (
          <button onClick={() => onEdit(athlete.id)}
            className="rounded-lg bg-primary-500 px-4 py-2 text-white hover:bg-primary-600">
            {t('edit')}
          </button>
        )}
        {onDelete && (
          <button onClick={() => onDelete(athlete.id)}
            className="rounded-lg bg-red-500 px-4 py-2 text-white hover:bg-red-600">
            {t('delete')}
          </button>
        )}
      </div>
    </div>
  );
}
```

### 4. Supabase Integration Pattern
```tsx
// packages/shared/services/supabase.ts
import { createClient } from '@supabase/supabase-js';
import type { Database } from '../types/supabase';

const supabaseUrl = import.meta.env.VITE_SUPABASE_URL;
const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY;

export const supabase = createClient<Database>(supabaseUrl, supabaseAnonKey);

// packages/shared/hooks/useAuth.ts
import { useEffect, useState } from 'react';
import { supabase } from '../services/supabase';
import type { User, Session } from '@supabase/supabase-js';

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [session, setSession] = useState<Session | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session);
      setUser(session?.user ?? null);
      setLoading(false);
    });

    const { data: { subscription } } = supabase.auth.onAuthStateChange(
      (_event, session) => {
        setSession(session);
        setUser(session?.user ?? null);
      }
    );

    return () => subscription.unsubscribe();
  }, []);

  return {
    user,
    session,
    loading,
    signIn: (email: string, password: string) =>
      supabase.auth.signInWithPassword({ email, password }),
    signUp: (email: string, password: string) =>
      supabase.auth.signUp({ email, password }),
    signOut: () => supabase.auth.signOut(),
    signInWithGoogle: () =>
      supabase.auth.signInWithOAuth({ provider: 'google' }),
  };
}
```

### 5. API Hook Pattern (TanStack Query + Supabase)
```tsx
// packages/shared/hooks/useAthletes.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { supabase } from '../services/supabase';

export const useAthletes = (filter: AthleteFilter) => {
  return useQuery({
    queryKey: ['athletes', filter],
    queryFn: async () => {
      let query = supabase
        .from('athletes')
        .select('*, clubs(name)', { count: 'exact' })
        .is('deleted_at', null);

      if (filter.clubId) query = query.eq('club_id', filter.clubId);
      if (filter.status) query = query.eq('status', filter.status);
      if (filter.search) query = query.or(
        `first_name.ilike.%${filter.search}%,last_name.ilike.%${filter.search}%`
      );

      const { data, error, count } = await query
        .range((filter.page - 1) * filter.perPage, filter.page * filter.perPage - 1)
        .order(filter.sortBy ?? 'created_at', { ascending: filter.sortDir === 'asc' });

      if (error) throw error;
      return { data, total: count ?? 0 };
    },
    staleTime: 5 * 60 * 1000,
  });
};

export const useCreateAthlete = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: CreateAthleteDTO) => {
      const { data, error } = await supabase
        .from('athletes')
        .insert(input)
        .select()
        .single();
      if (error) throw error;
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['athletes'] });
    },
  });
};
```

### 6. Supabase Realtime Subscription
```tsx
// Real-time scoring updates
import { useEffect } from 'react';
import { supabase } from '../services/supabase';

export function useRealtimeScores(tournamentId: string) {
  const [scores, setScores] = useState<Score[]>([]);

  useEffect(() => {
    const channel = supabase
      .channel(`scores:${tournamentId}`)
      .on('postgres_changes', {
        event: '*',
        schema: 'public',
        table: 'scores',
        filter: `tournament_id=eq.${tournamentId}`,
      }, (payload) => {
        if (payload.eventType === 'INSERT') {
          setScores(prev => [...prev, payload.new as Score]);
        } else if (payload.eventType === 'UPDATE') {
          setScores(prev => prev.map(s =>
            s.id === payload.new.id ? payload.new as Score : s
          ));
        }
      })
      .subscribe();

    return () => { supabase.removeChannel(channel); };
  }, [tournamentId]);

  return scores;
}
```

### 7. Expo React Native Pattern
```tsx
// apps/mobile/app/(tabs)/athletes.tsx
import { FlatList, View, Text, Pressable } from 'react-native';
import { useRouter } from 'expo-router';
import { useAthletes } from '@vct/shared/hooks/useAthletes';

export default function AthletesScreen() {
  const router = useRouter();
  const { data, isLoading } = useAthletes({ page: 1, perPage: 20 });

  if (isLoading) return <LoadingSkeleton />;

  return (
    <FlatList
      data={data?.data}
      keyExtractor={(item) => item.id}
      renderItem={({ item }) => (
        <Pressable
          onPress={() => router.push(`/athlete/${item.id}`)}
          className="mx-4 mb-3 rounded-xl bg-white p-4 shadow-sm dark:bg-gray-800"
        >
          <Text className="text-lg font-semibold text-gray-900 dark:text-white">
            {item.first_name} {item.last_name}
          </Text>
          <Text className="text-sm text-gray-500">{item.clubs?.name}</Text>
        </Pressable>
      )}
    />
  );
}
```

### 8. i18n Pattern
```json
// i18n/vi/athlete.json
{
  "title": "Quản lý Vận động viên",
  "list": "Danh sách VĐV",
  "create": "Thêm VĐV mới",
  "edit": "Chỉnh sửa VĐV",
  "fields": {
    "firstName": "Họ",
    "lastName": "Tên",
    "dateOfBirth": "Ngày sinh",
    "gender": "Giới tính",
    "email": "Email",
    "phone": "Số điện thoại",
    "club": "Câu lạc bộ",
    "status": "Trạng thái"
  },
  "status": {
    "active": "Hoạt động",
    "inactive": "Không hoạt động",
    "suspended": "Tạm ngưng"
  },
  "messages": {
    "created": "Thêm VĐV thành công",
    "updated": "Cập nhật VĐV thành công",
    "deleted": "Xóa VĐV thành công",
    "confirmDelete": "Bạn có chắc muốn xóa VĐV này?"
  }
}
```

### 9. Form Validation Pattern
```tsx
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';

const athleteSchema = z.object({
  firstName: z.string().min(2, 'Tối thiểu 2 ký tự').max(100),
  lastName: z.string().min(2, 'Tối thiểu 2 ký tự').max(100),
  dateOfBirth: z.date().refine(
    (d) => differenceInYears(new Date(), d) >= 6,
    'VĐV phải từ 6 tuổi trở lên'
  ),
  gender: z.enum(['male', 'female']),
  email: z.string().email('Email không hợp lệ'),
  phone: z.string().regex(/^\+?[0-9]{10,15}$/, 'Số điện thoại không hợp lệ').optional(),
});

type AthleteFormData = z.infer<typeof athleteSchema>;
```

### 10. Responsive Design Rules
- Mobile first: Design for 320px and scale up
- Breakpoints: `sm: 640px`, `md: 768px`, `lg: 1024px`, `xl: 1280px`
- Sidebar: Collapsible on mobile, fixed on desktop
- Tables: Horizontal scroll on mobile, cards view alternative
- Forms: Single column on mobile, multi-column on desktop
- Native mobile uses Expo's SafeAreaView + platform-specific layouts

### 11. Development Checklist
- [ ] TypeScript strict mode, no `any` types
- [ ] All text uses i18n (`useTranslation`)
- [ ] Responsive design tested at 320px, 768px, 1024px, 1440px
- [ ] Mobile tested on iOS (iPhone SE+) and Android (Pixel 4+)
- [ ] Loading states for async operations
- [ ] Error boundaries for component error handling
- [ ] Empty states for lists
- [ ] Accessibility: proper ARIA labels, keyboard navigation
- [ ] Tailwind v4 classes only (no inline styles, no CSS modules)
- [ ] ESLint 9 flat config passes with zero warnings
- [ ] Component tests with React Testing Library
- [ ] No console.log in production code
- [ ] Supabase RLS policies verified for data access
