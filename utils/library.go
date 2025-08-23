package utils

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func BuildLibrary(root string) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	artMap := map[string]*Artist{}
	albMap := map[string]*Album{}
	TrackByID = map[string]*Track{}

	err = filepath.WalkDir(absRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.ToLower(filepath.Ext(d.Name())) != ".flac" {
			return nil
		}

		rel, err := filepath.Rel(absRoot, path)

		if err != nil {
			return err
		}

		segs := strings.Split(rel, string(filepath.Separator))

		if len(segs) < 3 {
			return nil
		}

		artistName := segs[0]
		albumName := segs[1]
		fileName := segs[len(segs)-1]

		trackNo, trackTitle := ParseTrackInfo(fileName)
		if trackTitle == "" {
			trackTitle = strings.TrimSuffix(fileName, filepath.Ext(fileName))
		}

		artist := artMap[artistName]
		if artist == nil {
			artist = &Artist{Name: artistName}
			artMap[artistName] = artist
		}

		albKey := artistName + "\u0000" + albumName
		album := albMap[albKey]

		if album == nil {
			album = &Album{Name: albumName, Artist: artistName}
			albMap[albKey] = album
			artist.Albums = append(artist.Albums, album)
		}

		id := HashString(path)
		track := &Track{
			ID:     id,
			Name:   trackTitle,
			Track:  trackNo,
			Album:  albumName,
			Artist: artistName,
			Path:   path,
		}
		album.Tracks = append(album.Tracks, track)
		TrackByID[id] = track

		return nil
	})
	if err != nil {
		return err
	}

	var artists []*Artist
	for _, a := range artMap {
		sort.Slice(a.Albums, func(i, j int) bool { return a.Albums[i].Name < a.Albums[j].Name })
		for _, alb := range a.Albums {
			sort.Slice(alb.Tracks, func(i, j int) bool {
				if alb.Tracks[i].Track == alb.Tracks[j].Track {
					return alb.Tracks[i].Name < alb.Tracks[j].Name
				}
				if alb.Tracks[i].Track == 0 || alb.Tracks[j].Track == 0 {
					return alb.Tracks[i].Name < alb.Tracks[j].Name
				}
				return alb.Tracks[i].Track < alb.Tracks[j].Track
			})
		}
		artists = append(artists, a)
	}

	sort.Slice(artists, func(i, j int) bool { return artists[i].Name < artists[j].Name })
	TheLibrary = Library{Artists: artists}
	log.Printf("library loaded: %d artists, %d tracks", len(TheLibrary.Artists), len(TrackByID))
	return nil
}
