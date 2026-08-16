package redisadapter

import "testing"

func TestPlayerKeyIsCaseInsensitiveAndURLSafe(t *testing.T) {
	first := playerKey("Player Name", "BR:1")
	second := playerKey("player name", "br:1")

	if first != second {
		t.Fatalf("playerKey() should be case insensitive: %q != %q", first, second)
	}
	if first != "player:account:player+name%3Abr%3A1" {
		t.Fatalf("playerKey() = %q, want URL-safe key", first)
	}
}
