package test

import (
	"mouniu/internal/services"
	"testing"
)

func TestGetShenZhenAnnouncement(t *testing.T) {
	t.Run("Fetch SZ Announcements", func(t *testing.T) {
		services.GetShenZhenAnnouncement("601360")
	})
}

func TestGetHongKongAnnouncement(t *testing.T) {
	t.Run("Fetch HK Announcements", func(t *testing.T) {
		services.GetHongKongAnnouncement("01810")
	})
}
