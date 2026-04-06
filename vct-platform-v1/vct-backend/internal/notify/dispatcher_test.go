package notify

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

// ── Mock channel ─────────────────────

type mockChannel struct {
	chType ChannelType
	sent   atomic.Int32
	failAt int32
}

func (m *mockChannel) Type() ChannelType { return m.chType }
func (m *mockChannel) Send(_ context.Context, msg *Message) error {
	n := m.sent.Add(1)
	if m.failAt > 0 && n >= m.failAt {
		return errors.New("delivery failed")
	}
	return nil
}

// ── Mock preferences ─────────────────

type mockPrefs struct {
	prefs map[string]*UserPreferences
}

func (m *mockPrefs) Get(_ context.Context, userID string) (*UserPreferences, error) {
	if p, ok := m.prefs[userID]; ok {
		return p, nil
	}
	return nil, nil
}

func TestDispatch_SingleChannel(t *testing.T) {
	d := NewDispatcher(nil)
	email := &mockChannel{chType: ChannelEmail}
	d.RegisterChannel(email)

	results := d.Dispatch(context.Background(), &Message{
		Type:      "test",
		Recipient: "user-1",
		Title:     "Test",
		Body:      "Hello",
		Channels:  []ChannelType{ChannelEmail},
	})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Delivered {
		t.Error("should be delivered")
	}
	if email.sent.Load() != 1 {
		t.Error("email should have been sent")
	}
}

func TestDispatch_MultiChannel(t *testing.T) {
	d := NewDispatcher(nil)
	email := &mockChannel{chType: ChannelEmail}
	push := &mockChannel{chType: ChannelPush}
	d.RegisterChannel(email)
	d.RegisterChannel(push)

	results := d.Dispatch(context.Background(), &Message{
		Type:     "match.started",
		Channels: []ChannelType{ChannelEmail, ChannelPush},
	})

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestDispatch_ChannelFailure(t *testing.T) {
	d := NewDispatcher(nil)
	failing := &mockChannel{chType: ChannelSMS, failAt: 1}
	d.RegisterChannel(failing)

	results := d.Dispatch(context.Background(), &Message{
		Channels: []ChannelType{ChannelSMS},
	})

	if results[0].Delivered {
		t.Error("should not be delivered")
	}
	if results[0].Error == "" {
		t.Error("should have error message")
	}

	stats := d.Stats()
	if stats.Failed != 1 {
		t.Errorf("expected 1 failure, got %d", stats.Failed)
	}
}

func TestDispatch_UserPreferences(t *testing.T) {
	prefs := &mockPrefs{prefs: map[string]*UserPreferences{
		"user-1": {
			UserID: "user-1",
			EnabledChannels: map[ChannelType]bool{
				ChannelEmail: true,
				ChannelPush:  false,
				ChannelInApp: true,
			},
		},
	}}

	d := NewDispatcher(prefs)
	email := &mockChannel{chType: ChannelEmail}
	push := &mockChannel{chType: ChannelPush}
	inApp := &mockChannel{chType: ChannelInApp}
	d.RegisterChannel(email)
	d.RegisterChannel(push)
	d.RegisterChannel(inApp)

	results := d.Dispatch(context.Background(), &Message{
		Recipient: "user-1",
		Channels:  []ChannelType{ChannelEmail, ChannelPush, ChannelInApp},
	})

	// push disabled → only email + in_app
	if len(results) != 2 {
		t.Errorf("expected 2 (email+in_app), got %d", len(results))
	}
	if push.sent.Load() != 0 {
		t.Error("push should not have been sent")
	}
}

func TestDispatch_QuietMode(t *testing.T) {
	prefs := &mockPrefs{prefs: map[string]*UserPreferences{
		"user-quiet": {
			UserID: "user-quiet",
			Quiet:  true,
			EnabledChannels: map[ChannelType]bool{
				ChannelEmail: true,
			},
		},
	}}

	d := NewDispatcher(prefs)
	d.RegisterChannel(&mockChannel{chType: ChannelEmail})

	// Normal priority → filtered
	results := d.Dispatch(context.Background(), &Message{
		Recipient: "user-quiet",
		Priority:  PriorityNormal,
		Channels:  []ChannelType{ChannelEmail},
	})

	if len(results) != 0 {
		t.Error("quiet mode should filter normal priority")
	}

	// Critical → goes through
	results = d.Dispatch(context.Background(), &Message{
		Recipient: "user-quiet",
		Priority:  PriorityCritical,
		Channels:  []ChannelType{ChannelEmail},
	})

	if len(results) != 1 {
		t.Error("critical should bypass quiet mode")
	}
}

func TestAutoID(t *testing.T) {
	d := NewDispatcher(nil)
	d.RegisterChannel(&mockChannel{chType: ChannelInApp})

	d.Dispatch(context.Background(), &Message{Channels: []ChannelType{ChannelInApp}})
	d.Dispatch(context.Background(), &Message{Channels: []ChannelType{ChannelInApp}})

	history := d.History()
	if len(history) != 2 {
		t.Errorf("expected 2 history entries, got %d", len(history))
	}
}

func TestStats(t *testing.T) {
	d := NewDispatcher(nil)
	d.RegisterChannel(&mockChannel{chType: ChannelEmail})
	d.RegisterChannel(&mockChannel{chType: ChannelPush})

	d.Dispatch(context.Background(), &Message{Channels: []ChannelType{ChannelEmail}})
	d.Dispatch(context.Background(), &Message{Channels: []ChannelType{ChannelEmail, ChannelPush}})

	stats := d.Stats()
	if stats.Sent != 3 {
		t.Errorf("expected 3 sent, got %d", stats.Sent)
	}
	if stats.Channels != 2 {
		t.Errorf("expected 2 channels, got %d", stats.Channels)
	}
}

func TestUnregisteredChannel(t *testing.T) {
	d := NewDispatcher(nil)
	// SMS channel not registered
	results := d.Dispatch(context.Background(), &Message{Channels: []ChannelType{ChannelSMS}})

	if len(results) != 0 {
		t.Error("unregistered channel should be skipped")
	}
}
