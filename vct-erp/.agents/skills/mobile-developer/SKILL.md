---
name: mobile-developer
description: Mobile Developer role - Deep Expo React Native patterns, native modules, OTA updates, push notifications, offline-first, app store optimization for VCT Platform.
---

# Mobile Developer - VCT Platform

## Role Overview
Specializes in deep Expo React Native development beyond basic frontend patterns. Handles native module integration, OTA updates, push notifications, offline-first architecture, biometric auth, camera/media, and app store optimization.

## Technology Stack
- **Framework**: Expo SDK 53+ / React Native
- **Navigation**: Expo Router (file-based)
- **Styling**: NativeWind v4 (Tailwind CSS for RN)
- **State**: Zustand + TanStack Query v5
- **Auth**: Supabase Auth (social, biometric)
- **Realtime**: Supabase Realtime
- **Storage**: Supabase Storage + expo-file-system
- **Build**: EAS Build + EAS Submit
- **OTA**: EAS Update
- **Push**: Expo Push Notifications
- **Testing**: Jest + @testing-library/react-native + Detox

## Core Patterns

### 1. Expo Router (File-Based Navigation)

```
apps/mobile/app/
├── _layout.tsx           ← Root layout (auth check)
├── index.tsx             ← Entry / splash
├── (auth)/
│   ├── _layout.tsx       ← Auth stack layout
│   ├── login.tsx
│   └── register.tsx
├── (tabs)/
│   ├── _layout.tsx       ← Tab navigator
│   ├── index.tsx         ← Dashboard tab
│   ├── athletes.tsx      ← Athletes tab
│   ├── tournaments.tsx   ← Tournaments tab
│   └── profile.tsx       ← Profile tab
├── athlete/
│   ├── [id].tsx          ← Dynamic route: athlete detail
│   └── [id]/edit.tsx     ← Nested: edit athlete
├── tournament/
│   ├── [id].tsx          ← Tournament detail
│   └── [id]/live.tsx     ← Live scoring view
└── +not-found.tsx        ← 404
```

```typescript
// app/_layout.tsx
import { Stack } from 'expo-router';
import { useAuth } from '@/hooks/useAuth';

export default function RootLayout() {
  const { session, isLoading } = useAuth();

  if (isLoading) return <SplashScreen />;

  return (
    <Stack screenOptions={{ headerShown: false }}>
      {session ? (
        <Stack.Screen name="(tabs)" />
      ) : (
        <Stack.Screen name="(auth)" />
      )}
    </Stack>
  );
}
```

### 2. Offline-First Architecture

```typescript
// lib/offline.ts
import NetInfo from '@react-native-community/netinfo';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { onlineManager } from '@tanstack/react-query';

// Sync TanStack Query with network state
onlineManager.setEventListener((setOnline) => {
  return NetInfo.addEventListener((state) => {
    setOnline(!!state.isConnected);
  });
});

// Persist query cache for offline access
import { createAsyncStoragePersister } from '@tanstack/query-async-storage-persister';

export const asyncStoragePersister = createAsyncStoragePersister({
  storage: AsyncStorage,
  throttleTime: 1000,
});

// Offline mutation queue
export async function queueMutation(mutation: {
  endpoint: string;
  method: string;
  body: any;
}) {
  const queue = JSON.parse(
    (await AsyncStorage.getItem('mutation-queue')) || '[]'
  );
  queue.push({ ...mutation, timestamp: Date.now() });
  await AsyncStorage.setItem('mutation-queue', JSON.stringify(queue));
}

export async function processMutationQueue() {
  const queue = JSON.parse(
    (await AsyncStorage.getItem('mutation-queue')) || '[]'
  );
  for (const mutation of queue) {
    try {
      await fetch(mutation.endpoint, {
        method: mutation.method,
        body: JSON.stringify(mutation.body),
        headers: { 'Content-Type': 'application/json' },
      });
    } catch {
      break; // Stop processing, will retry later
    }
  }
  await AsyncStorage.removeItem('mutation-queue');
}
```

### 3. Push Notifications (Expo)

```typescript
// lib/notifications.ts
import * as Notifications from 'expo-notifications';
import * as Device from 'expo-device';
import { supabase } from '@/lib/supabase';

// Configure notification behavior
Notifications.setNotificationHandler({
  handleNotification: async () => ({
    shouldShowAlert: true,
    shouldPlaySound: true,
    shouldSetBadge: true,
  }),
});

export async function registerForPushNotifications() {
  if (!Device.isDevice) {
    console.warn('Push notifications require a physical device');
    return null;
  }

  const { status: existingStatus } = await Notifications.getPermissionsAsync();
  let finalStatus = existingStatus;

  if (existingStatus !== 'granted') {
    const { status } = await Notifications.requestPermissionsAsync();
    finalStatus = status;
  }

  if (finalStatus !== 'granted') return null;

  const token = (await Notifications.getExpoPushTokenAsync({
    projectId: 'YOUR_EXPO_PROJECT_ID',
  })).data;

  // Save token to Supabase
  await supabase.from('push_tokens').upsert({
    user_id: (await supabase.auth.getUser()).data.user?.id,
    token,
    platform: Device.osName,
  });

  return token;
}

// Handle notification tap (deep linking)
export function useNotificationObserver() {
  const router = useRouter();

  useEffect(() => {
    const subscription = Notifications.addNotificationResponseReceivedListener(
      (response) => {
        const { data } = response.notification.request.content;
        if (data?.route) {
          router.push(data.route); // e.g., '/tournament/123/live'
        }
      }
    );

    return () => subscription.remove();
  }, [router]);
}
```

### 4. Biometric Authentication

```typescript
// lib/biometrics.ts
import * as LocalAuthentication from 'expo-local-authentication';
import * as SecureStore from 'expo-secure-store';

export async function isBiometricAvailable(): Promise<boolean> {
  const compatible = await LocalAuthentication.hasHardwareAsync();
  const enrolled = await LocalAuthentication.isEnrolledAsync();
  return compatible && enrolled;
}

export async function authenticateWithBiometrics(): Promise<boolean> {
  const result = await LocalAuthentication.authenticateAsync({
    promptMessage: 'Xác thực để tiếp tục',
    cancelLabel: 'Hủy',
    fallbackLabel: 'Dùng mật khẩu',
    disableDeviceFallback: false,
  });

  return result.success;
}

// Store Supabase session securely
export async function saveSessionSecurely(session: string) {
  await SecureStore.setItemAsync('supabase-session', session);
}

export async function getSessionSecurely(): Promise<string | null> {
  return await SecureStore.getItemAsync('supabase-session');
}
```

### 5. Camera & Media (QR/Barcode for Check-in)

```typescript
// components/QRScanner.tsx
import { CameraView, useCameraPermissions } from 'expo-camera';

export function QRScanner({ onScan }: { onScan: (data: string) => void }) {
  const [permission, requestPermission] = useCameraPermissions();

  if (!permission?.granted) {
    return (
      <View className="flex-1 items-center justify-center">
        <Text className="text-lg mb-4">Cần quyền truy cập camera</Text>
        <Button title="Cấp quyền" onPress={requestPermission} />
      </View>
    );
  }

  return (
    <CameraView
      className="flex-1"
      facing="back"
      barcodeScannerSettings={{ barcodeTypes: ['qr'] }}
      onBarcodeScanned={(result) => onScan(result.data)}
    />
  );
}
```

### 6. OTA Updates (EAS Update)

```typescript
// app/_layout.tsx
import * as Updates from 'expo-updates';

export function useOTAUpdate() {
  useEffect(() => {
    async function checkForUpdate() {
      if (__DEV__) return; // Skip in development

      try {
        const update = await Updates.checkForUpdateAsync();
        if (update.isAvailable) {
          await Updates.fetchUpdateAsync();
          // Optionally prompt user or auto-restart
          Alert.alert(
            'Cập nhật mới',
            'Ứng dụng đã được cập nhật. Khởi động lại?',
            [
              { text: 'Sau', style: 'cancel' },
              { text: 'Khởi động lại', onPress: () => Updates.reloadAsync() },
            ]
          );
        }
      } catch (error) {
        console.error('OTA update check failed:', error);
      }
    }

    checkForUpdate();
  }, []);
}
```

### 7. Performance Optimization

```typescript
// Optimized FlatList for athlete list
import { FlashList } from '@shopify/flash-list';

function AthleteList({ athletes }: { athletes: Athlete[] }) {
  const renderItem = useCallback(({ item }: { item: Athlete }) => (
    <AthleteCard athlete={item} />
  ), []);

  const keyExtractor = useCallback((item: Athlete) => item.id, []);

  return (
    <FlashList
      data={athletes}
      renderItem={renderItem}
      keyExtractor={keyExtractor}
      estimatedItemSize={80}
      // Performance optimizations
      removeClippedSubviews={true}
      maxToRenderPerBatch={10}
      windowSize={5}
    />
  );
}
```

### 8. App Store Optimization (ASO)

| Element | iOS (App Store) | Android (Play Store) |
|---------|----------------|---------------------|
| App Name | VCT Platform - Cycling & Triathlon | Same |
| Subtitle | Quản lý vận động viên xe đạp & 3 môn | Short description |
| Keywords | cycling, triathlon, sports, Vietnam, VĐV | Tags |
| Screenshots | 6.7" (iPhone 15 Pro Max) + 12.9" (iPad) | Phone + 7" + 10" tablet |
| Languages | Vietnamese (primary), English | Same |
| Category | Sports | Sports |
| Rating | 4+ | Everyone |

### 9. Mobile Developer Checklist
- [ ] Expo Router navigation works (all routes, deep links)
- [ ] Offline mode: cached data loads without network
- [ ] Push notifications: permission request, token storage, delivery
- [ ] Biometric auth: Face ID / Touch ID / fingerprint
- [ ] Camera: QR scan for check-in works
- [ ] OTA updates: EAS Update configured
- [ ] Performance: FlatList/FlashList for large lists
- [ ] Secure storage: Supabase session in SecureStore
- [ ] App size: < 50MB
- [ ] Accessibility: VoiceOver (iOS), TalkBack (Android)
- [ ] Vietnamese diacritics render correctly
- [ ] Background/foreground state handled properly
- [ ] Error boundaries for crash recovery
