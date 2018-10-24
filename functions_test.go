package main

import (
	"testing"
	"time"

	igc "github.com/marni/goigc"
)

func Test_parseTimeDifference(t *testing.T) {

	secondsArray := []int{0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576,
		2097152, 4194304, 8388608, 16777216, 33554432, 67108864, 134217728, 268435456, 536870912, 1073741824, 2147483648, 4294967296}

	returnValueArray := []string{"P", "PTS1", "PTS2", "PTS4", "PTS8", "PTS16", "PTS32", "PTM1S4", "PTM2S8", "PTM4S16", "PTM8S32", "PTM17S4", "PTM34S8",
		"PTH1M8S16", "PTH2M16S32", "PTH4M33S4", "PTH9M6S8", "PTH18M12S16", "PD1TH12M24S32", "PD3TM49S4", "PD6TH1M38S8", "PW1D5TH3M16S16",
		"PW3D3TH6M32S32", "PM1W2D4TH13M5S4", "PM3W1TH2M10S8", "PM6W2TH4M20S16", "PY1W3D2TH2M40S32", "PY2M1W2D2TH5M21S4",
		"PY4M3D2TH10M42S8", "PY8M6D4TH21M24S16", "PY17D4TH12M48S32", "PY34W1D2TH1M37S4", "PY68W2D4TH3M14S8", "PY136M1D6TH6M28S16"}

	// Check whether the seconds in secondsArray correspond to the formatted
	for i := 0; i < len(returnValueArray); i++ {
		if parseTimeDifference(secondsArray[i]) != returnValueArray[i] {
			t.Errorf("Time duration format is not correct\n%s\n%s\n", parseTimeDifference(secondsArray[i]), returnValueArray[i])
		}
	}
}

func Test_calculateTotalDistance(t *testing.T) {

	urlDistance := map[string]string{
		`http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc`:  `443.26`,
		`http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc`: `2885.44`,
		`http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc`:  `2001.11`,
		`http://skypolaris.org/wp-content/uploads/IGS%20Files/Medellin%20Guatemala.igc`: `2068.56`,
	}

	for key, val := range urlDistance {
		track, _ := igc.ParseLocation(key)
		dist := calculateTotalDistance(track.Points)
		if dist != val {
			t.Error("Distance is not the one wanted")
		}
	}

}

func Test_latestTimestamp(t *testing.T) {
	igcTracks := []igcTrack{
		igcTrack{TimeRecorded: time.Date(2018, 4, 25, 12, 32, 1, 0, time.UTC)},
		igcTrack{TimeRecorded: time.Now()},
		igcTrack{TimeRecorded: time.Date(2019, 4, 25, 12, 32, 1, 0, time.UTC)},
	}

	latestTS := latestTimestamp(igcTracks)
	if latestTS != igcTracks[2].TimeRecorded {
		t.Error("Not the latest timestamp")
	}
}

func Test_oldestTimestamp(t *testing.T) {
	igcTracks := []igcTrack{
		igcTrack{TimeRecorded: time.Date(2018, 4, 25, 12, 32, 1, 0, time.UTC)},
		igcTrack{TimeRecorded: time.Now()},
		igcTrack{TimeRecorded: time.Date(2019, 4, 25, 12, 32, 1, 0, time.UTC)},
	}

	oldestTS := oldestTimestamp(igcTracks)
	if oldestTS != igcTracks[0].TimeRecorded {
		t.Error("Not the oldest timestamp")
	}
}

func Test_oldestNewerTimestamp(t *testing.T) {
	igcTracks := []igcTrack{
		igcTrack{TimeRecorded: time.Date(2018, 4, 25, 12, 32, 1, 0, time.UTC)},
		igcTrack{TimeRecorded: time.Date(2018, 4, 26, 12, 32, 1, 0, time.UTC)},
		igcTrack{TimeRecorded: time.Date(2019, 4, 25, 12, 32, 1, 0, time.UTC)},
	}

	oldestNewTS := oldestNewerTimestamp("25.04.2018 12:34:30.314", igcTracks)

	if oldestNewTS != igcTracks[1].TimeRecorded {
		t.Error("Not the right timestamp")
	}
}

func Test_tickerTimestamps(t *testing.T) {
	igcTracks := []igcTrack{
		igcTrack{TimeRecorded: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	// No connection to the DB :(

	tickerTS := tickerTimestamps("25.04.2018 12:34:30.314")

	if tickerTS.oldestTimestamp != igcTracks[0].TimeRecorded {
		t.Error("Not the right timestamp")
	}
	if tickerTS.oldestNewerTimestamp != igcTracks[0].TimeRecorded {
		t.Error("Not the right timestamp")
	}
	if tickerTS.latestTimestamp != igcTracks[0].TimeRecorded {
		t.Error("Not the right timestamp")
	}
}
