---
name: realtime-engineer
description: Realtime Engineer role - Supabase Realtime, WebSocket patterns, live scoring state machine, pub/sub, and event-driven architecture for VCT Platform.
---

# Realtime Engineer - VCT Platform

## Role Overview
Designs and implements real-time features using Supabase Realtime (Postgres Changes, Broadcast, Presence) and WebSocket patterns. Primary responsibility is the live scoring engine, real-time notifications, and collaborative features.

## Technology Stack
- **Primary**: Supabase Realtime (Channels, Postgres Changes, Broadcast, Presence)
- **WebSocket**: Supabase JS client (`@supabase/supabase-js`)
- **Backend Events**: PostgreSQL NOTIFY/LISTEN, pg_notify
- **State Machine**: XState v5 (TypeScript) for scoring state management
- **Message Queue**: Redis Pub/Sub (for backend-to-backend events)
- **Mobile**: Expo + Supabase Realtime client

## Core Patterns

### 1. Supabase Realtime Channels

#### Postgres Changes (Database CDC)
```typescript
// Listen to real-time database changes
import { supabase } from '@/lib/supabase';

const channel = supabase
  .channel('athlete-changes')
  .on(
    'postgres_changes',
    {
      event: '*',       // INSERT, UPDATE, DELETE
      schema: 'public',
      table: 'athletes',
      filter: 'organization_id=eq.${orgId}',
    },
    (payload) => {
      console.log('Change:', payload.eventType, payload.new);
      // Update local state / TanStack Query cache
      queryClient.invalidateQueries({ queryKey: ['athletes'] });
    }
  )
  .subscribe();

// Cleanup
return () => { supabase.removeChannel(channel); };
```

#### Broadcast (Client-to-Client)
```typescript
// Live scoring broadcast (no database involved)
const scoringChannel = supabase.channel('race:${raceId}');

// Send score update
scoringChannel.send({
  type: 'broadcast',
  event: 'score_update',
  payload: {
    athleteId: '...',
    checkpoint: 3,
    time: '01:23:45.678',
    position: 2,
  },
});

// Receive score updates
scoringChannel
  .on('broadcast', { event: 'score_update' }, (payload) => {
    updateLeaderboard(payload.payload);
  })
  .subscribe();
```

#### Presence (Who's Online)
```typescript
// Track who's viewing the live scoring page
const presenceChannel = supabase.channel('race:${raceId}:presence');

presenceChannel
  .on('presence', { event: 'sync' }, () => {
    const state = presenceChannel.presenceState();
    setViewerCount(Object.keys(state).length);
  })
  .on('presence', { event: 'join' }, ({ key, newPresences }) => {
    console.log('Joined:', newPresences);
  })
  .on('presence', { event: 'leave' }, ({ key, leftPresences }) => {
    console.log('Left:', leftPresences);
  })
  .subscribe(async (status) => {
    if (status === 'SUBSCRIBED') {
      await presenceChannel.track({
        userId: user.id,
        role: user.role,
        online_at: new Date().toISOString(),
      });
    }
  });
```

### 2. Live Scoring State Machine (XState v5)

```typescript
import { createMachine, assign } from 'xstate';

type RaceContext = {
  raceId: string;
  athletes: AthleteScore[];
  currentCheckpoint: number;
  totalCheckpoints: number;
  status: 'pending' | 'in_progress' | 'finished' | 'cancelled';
};

type RaceEvent =
  | { type: 'START_RACE' }
  | { type: 'RECORD_TIME'; athleteId: string; time: string }
  | { type: 'ADVANCE_CHECKPOINT' }
  | { type: 'FINISH_RACE' }
  | { type: 'CANCEL_RACE'; reason: string }
  | { type: 'DNS'; athleteId: string }  // Did Not Start
  | { type: 'DNF'; athleteId: string }  // Did Not Finish
  | { type: 'DSQ'; athleteId: string; reason: string };  // Disqualified

const raceScoringMachine = createMachine({
  id: 'raceScoring',
  initial: 'pending',
  context: {
    raceId: '',
    athletes: [],
    currentCheckpoint: 0,
    totalCheckpoints: 5,
    status: 'pending',
  } as RaceContext,
  states: {
    pending: {
      on: {
        START_RACE: {
          target: 'racing',
          actions: assign({ status: 'in_progress', currentCheckpoint: 1 }),
        },
        DNS: {
          actions: assign({
            athletes: ({ context, event }) =>
              context.athletes.map(a =>
                a.id === event.athleteId ? { ...a, status: 'DNS' } : a
              ),
          }),
        },
      },
    },
    racing: {
      on: {
        RECORD_TIME: {
          actions: [
            assign({
              athletes: ({ context, event }) =>
                context.athletes.map(a =>
                  a.id === event.athleteId
                    ? { ...a, checkpoints: { ...a.checkpoints, [context.currentCheckpoint]: event.time } }
                    : a
                ),
            }),
            'broadcastScoreUpdate',
            'persistToDatabase',
          ],
        },
        ADVANCE_CHECKPOINT: {
          guard: ({ context }) => context.currentCheckpoint < context.totalCheckpoints,
          actions: assign({
            currentCheckpoint: ({ context }) => context.currentCheckpoint + 1,
          }),
        },
        DNF: {
          actions: assign({
            athletes: ({ context, event }) =>
              context.athletes.map(a =>
                a.id === event.athleteId ? { ...a, status: 'DNF' } : a
              ),
          }),
        },
        DSQ: {
          actions: assign({
            athletes: ({ context, event }) =>
              context.athletes.map(a =>
                a.id === event.athleteId ? { ...a, status: 'DSQ', dsqReason: event.reason } : a
              ),
          }),
        },
        FINISH_RACE: 'finishing',
        CANCEL_RACE: 'cancelled',
      },
    },
    finishing: {
      entry: ['calculateFinalRankings', 'persistFinalResults'],
      always: 'finished',
    },
    finished: {
      type: 'final',
      entry: ['broadcastFinalResults', 'updateRankings'],
    },
    cancelled: {
      type: 'final',
      entry: ['notifyCancellation'],
    },
  },
});
```

### 3. Backend Event System (PostgreSQL NOTIFY)

```sql
-- Trigger to notify on score changes
CREATE OR REPLACE FUNCTION notify_score_change()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify(
        'score_changes',
        json_build_object(
            'race_id', NEW.race_id,
            'athlete_id', NEW.athlete_id,
            'checkpoint', NEW.checkpoint,
            'time', NEW.time_recorded,
            'operation', TG_OP
        )::text
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_score_notify
    AFTER INSERT OR UPDATE ON race_results
    FOR EACH ROW EXECUTE FUNCTION notify_score_change();
```

```go
// Go listener for PostgreSQL NOTIFY
func (s *ScoringService) ListenForScoreChanges(ctx context.Context) error {
    conn, err := s.pool.Acquire(ctx)
    if err != nil {
        return err
    }
    defer conn.Release()

    _, err = conn.Exec(ctx, "LISTEN score_changes")
    if err != nil {
        return err
    }

    for {
        notification, err := conn.Conn().WaitForNotification(ctx)
        if err != nil {
            return err
        }
        // Process notification.Payload
        s.handleScoreNotification(ctx, notification.Payload)
    }
}
```

### 4. React Hook for Realtime Subscriptions

```typescript
// hooks/useRealtimeScores.ts
import { useEffect, useState } from 'react';
import { supabase } from '@/lib/supabase';
import { useQueryClient } from '@tanstack/react-query';

export function useRealtimeScores(raceId: string) {
  const queryClient = useQueryClient();
  const [isConnected, setIsConnected] = useState(false);
  const [viewerCount, setViewerCount] = useState(0);

  useEffect(() => {
    const channel = supabase
      .channel(`race:${raceId}`)
      // Database changes
      .on('postgres_changes', {
        event: '*',
        schema: 'public',
        table: 'race_results',
        filter: `race_id=eq.${raceId}`,
      }, () => {
        queryClient.invalidateQueries({ queryKey: ['race-results', raceId] });
      })
      // Broadcast events
      .on('broadcast', { event: 'score_update' }, (payload) => {
        queryClient.setQueryData(['live-scores', raceId], (old: any) => ({
          ...old,
          ...payload.payload,
        }));
      })
      // Presence
      .on('presence', { event: 'sync' }, () => {
        setViewerCount(Object.keys(channel.presenceState()).length);
      })
      .subscribe((status) => {
        setIsConnected(status === 'SUBSCRIBED');
      });

    return () => { supabase.removeChannel(channel); };
  }, [raceId, queryClient]);

  return { isConnected, viewerCount };
}
```

### 5. Connection Management

```typescript
// Reconnection strategy
const MAX_RETRIES = 5;
const BACKOFF_BASE = 1000;

function subscribeWithRetry(channelName: string, retryCount = 0) {
  const channel = supabase.channel(channelName);

  channel.subscribe((status, err) => {
    if (status === 'SUBSCRIBED') {
      console.log('Connected to', channelName);
    } else if (status === 'CHANNEL_ERROR') {
      if (retryCount < MAX_RETRIES) {
        const delay = BACKOFF_BASE * Math.pow(2, retryCount);
        console.warn(`Retry ${retryCount + 1} in ${delay}ms`);
        setTimeout(() => subscribeWithRetry(channelName, retryCount + 1), delay);
      }
    }
  });

  return channel;
}
```

### 6. Realtime Checklist
- [ ] Supabase project has Realtime enabled for required tables
- [ ] RLS policies allow SELECT for realtime listeners
- [ ] Channels cleaned up on component unmount
- [ ] Reconnection strategy implemented
- [ ] Presence tracking for live viewer count
- [ ] Score state machine handles all race states (DNS, DNF, DSQ)
- [ ] Backend NOTIFY triggers for server-side events
- [ ] Mobile (Expo) receives realtime updates in background
- [ ] Load tested for 10,000+ concurrent viewers
- [ ] Fallback polling for unreliable connections
