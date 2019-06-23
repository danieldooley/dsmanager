package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)


const (
	plexToken = "hiiBVvyBfvJ64_W5HUqy"
	plexUrl = "dooley-server.local:32400"

	hibernateStart = 22 //10:00pm //Hibernate start MUST BE greater than hibernateEnd
	hibernateEnd = 10 //10:00am
)

func hibernate() error {
	var till time.Time
	if time.Now().Hour() > hibernateStart {
		till = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, hibernateEnd, 0, 0, 0, time.Local)
	} else {
		till = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hibernateEnd, 0, 0, 0, time.Local)
	}

	return shutDownTill(till)
}

func plexHibernate() scheduledTask {
	return scheduledTask{
		interval: 10,
		timeout: time.Minute,

		task: func() error {
			if time.Since(startTime) < time.Hour {
				webLogf("PlexHibernate: Server has not been up an hour")
				return nil //If server has been turned on wait at least an hour before turning off
			}

			if time.Now().Hour() < hibernateStart && time.Now().Hour() >= hibernateEnd {
				webLogf("PlexHibernate: Hour %v is not between %v and %v", time.Now().Hour(), hibernateStart, hibernateEnd)
				return nil //If time is between start&end don't try to turn off
			}

			idle, err := isIdle()
			if err != nil {
				return fmt.Errorf("could not determine idle status: %v", err)
			}

			if idle {
				webLogf("PlexHibernate: Shutting down till: %v", hibernateEnd)
				err := hibernate()
				if err != nil {
					return fmt.Errorf("hibernate() failed: %v", err)
				}
			} else {
				webLogf("PlexHibernate: Not Idle")
			}

			return nil
		},
	}
}

func isIdle() (bool, error) {
	cmc, err := getClients()
	if err != nil {
		return false, fmt.Errorf("getClients() failed: %v", err)
	}

	ssmc, err := getStatusSessions()
	if err != nil {
		return false, fmt.Errorf("getStatusSessions() failed: %v", err)
	}

	return len(cmc.Server) == 0 && len(ssmc.Video) == 0, nil
}

func getClients() (ClientsMediaContainer, error) {
	var mc ClientsMediaContainer

	res, err := http.DefaultClient.Get(fmt.Sprintf("http://%v/clients?X-Plex-Token=%v", plexUrl, plexToken))
	if err != nil {
		return mc, fmt.Errorf("failed to GET from plex server: %v", err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return mc, fmt.Errorf("failed to read from response body: %v", err)
	}

	err = xml.Unmarshal(b, &mc)
	if err != nil {
		return mc, fmt.Errorf("failed to unmarshal xml: %v", err)
	}

	return mc, nil
}

func getStatusSessions() (StatusSessionsMediaContainer, error) {
	var mc StatusSessionsMediaContainer

	res, err := http.DefaultClient.Get(fmt.Sprintf("http://%v/status/sessions?X-Plex-Token=%v", plexUrl, plexToken))
	if err != nil {
		return mc, fmt.Errorf("failed to GET from plex server: %v", err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return mc, fmt.Errorf("failed to read from response body: %v", err)
	}

	err = xml.Unmarshal(b, &mc)
	if err != nil {
		return mc, fmt.Errorf("failed to unmarshal xml: %v", err)
	}

	return mc, nil
}


/*
	Types
 */

type ClientsMediaContainer struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Text    string   `xml:",chardata"`
	Size    string   `xml:"size,attr"`
	Server  []struct {
		Text                 string `xml:",chardata"`
		Name                 string `xml:"name,attr"`
		Host                 string `xml:"host,attr"`
		Address              string `xml:"address,attr"`
		Port                 string `xml:"port,attr"`
		MachineIdentifier    string `xml:"machineIdentifier,attr"`
		Version              string `xml:"version,attr"`
		Protocol             string `xml:"protocol,attr"`
		Product              string `xml:"product,attr"`
		DeviceClass          string `xml:"deviceClass,attr"`
		ProtocolVersion      string `xml:"protocolVersion,attr"`
		ProtocolCapabilities string `xml:"protocolCapabilities,attr"`
	} `xml:"Server"`
}

type StatusSessionsMediaContainer struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Text    string   `xml:",chardata"`
	Size    string   `xml:"size,attr"`
	Video   []struct {
		Text                  string `xml:",chardata"`
		AddedAt               string `xml:"addedAt,attr"`
		Art                   string `xml:"art,attr"`
		ContentRating         string `xml:"contentRating,attr"`
		Duration              string `xml:"duration,attr"`
		Guid                  string `xml:"guid,attr"`
		Key                   string `xml:"key,attr"`
		LibrarySectionID      string `xml:"librarySectionID,attr"`
		LibrarySectionKey     string `xml:"librarySectionKey,attr"`
		LibrarySectionTitle   string `xml:"librarySectionTitle,attr"`
		OriginallyAvailableAt string `xml:"originallyAvailableAt,attr"`
		Rating                string `xml:"rating,attr"`
		RatingImage           string `xml:"ratingImage,attr"`
		RatingKey             string `xml:"ratingKey,attr"`
		SessionKey            string `xml:"sessionKey,attr"`
		Studio                string `xml:"studio,attr"`
		Summary               string `xml:"summary,attr"`
		Tagline               string `xml:"tagline,attr"`
		Thumb                 string `xml:"thumb,attr"`
		Title                 string `xml:"title,attr"`
		Type                  string `xml:"type,attr"`
		UpdatedAt             string `xml:"updatedAt,attr"`
		ViewOffset            string `xml:"viewOffset,attr"`
		Year                  string `xml:"year,attr"`
		ChapterSource         string `xml:"chapterSource,attr"`
		GrandparentArt        string `xml:"grandparentArt,attr"`
		GrandparentKey        string `xml:"grandparentKey,attr"`
		GrandparentRatingKey  string `xml:"grandparentRatingKey,attr"`
		GrandparentTheme      string `xml:"grandparentTheme,attr"`
		GrandparentThumb      string `xml:"grandparentThumb,attr"`
		GrandparentTitle      string `xml:"grandparentTitle,attr"`
		Index                 string `xml:"index,attr"`
		LastViewedAt          string `xml:"lastViewedAt,attr"`
		ParentIndex           string `xml:"parentIndex,attr"`
		ParentKey             string `xml:"parentKey,attr"`
		ParentRatingKey       string `xml:"parentRatingKey,attr"`
		ParentThumb           string `xml:"parentThumb,attr"`
		ParentTitle           string `xml:"parentTitle,attr"`
		TitleSort             string `xml:"titleSort,attr"`
		Media                 struct {
			Text                  string `xml:",chardata"`
			AspectRatio           string `xml:"aspectRatio,attr"`
			AudioChannels         string `xml:"audioChannels,attr"`
			AudioCodec            string `xml:"audioCodec,attr"`
			AudioProfile          string `xml:"audioProfile,attr"`
			Bitrate               string `xml:"bitrate,attr"`
			Container             string `xml:"container,attr"`
			Duration              string `xml:"duration,attr"`
			Has64bitOffsets       string `xml:"has64bitOffsets,attr"`
			Height                string `xml:"height,attr"`
			ID                    string `xml:"id,attr"`
			OptimizedForStreaming string `xml:"optimizedForStreaming,attr"`
			VideoCodec            string `xml:"videoCodec,attr"`
			VideoFrameRate        string `xml:"videoFrameRate,attr"`
			VideoProfile          string `xml:"videoProfile,attr"`
			VideoResolution       string `xml:"videoResolution,attr"`
			Width                 string `xml:"width,attr"`
			Selected              string `xml:"selected,attr"`
			Protocol              string `xml:"protocol,attr"`
			Part                  struct {
				Text                  string `xml:",chardata"`
				AudioProfile          string `xml:"audioProfile,attr"`
				Container             string `xml:"container,attr"`
				Duration              string `xml:"duration,attr"`
				File                  string `xml:"file,attr"`
				Has64bitOffsets       string `xml:"has64bitOffsets,attr"`
				ID                    string `xml:"id,attr"`
				Key                   string `xml:"key,attr"`
				OptimizedForStreaming string `xml:"optimizedForStreaming,attr"`
				Size                  string `xml:"size,attr"`
				VideoProfile          string `xml:"videoProfile,attr"`
				Decision              string `xml:"decision,attr"`
				Selected              string `xml:"selected,attr"`
				Bitrate               string `xml:"bitrate,attr"`
				Height                string `xml:"height,attr"`
				Protocol              string `xml:"protocol,attr"`
				Width                 string `xml:"width,attr"`
				Stream                []struct {
					Text              string `xml:",chardata"`
					BitDepth          string `xml:"bitDepth,attr"`
					Bitrate           string `xml:"bitrate,attr"`
					ChromaLocation    string `xml:"chromaLocation,attr"`
					ChromaSubsampling string `xml:"chromaSubsampling,attr"`
					Codec             string `xml:"codec,attr"`
					Default           string `xml:"default,attr"`
					DisplayTitle      string `xml:"displayTitle,attr"`
					FrameRate         string `xml:"frameRate,attr"`
					HasScalingMatrix  string `xml:"hasScalingMatrix,attr"`
					Height            string `xml:"height,attr"`
					ID                string `xml:"id,attr"`
					Index             string `xml:"index,attr"`
					Level             string `xml:"level,attr"`
					Profile           string `xml:"profile,attr"`
					RefFrames         string `xml:"refFrames,attr"`
					StreamIdentifier  string `xml:"streamIdentifier,attr"`
					StreamType        string `xml:"streamType,attr"`
					Width             string `xml:"width,attr"`
					Location          string `xml:"location,attr"`
					Channels          string `xml:"channels,attr"`
					Language          string `xml:"language,attr"`
					LanguageCode      string `xml:"languageCode,attr"`
					SamplingRate      string `xml:"samplingRate,attr"`
					Selected          string `xml:"selected,attr"`
					Decision          string `xml:"decision,attr"`
					BitrateMode       string `xml:"bitrateMode,attr"`
				} `xml:"Stream"`
			} `xml:"Part"`
		} `xml:"Media"`
		Genre []struct {
			Text   string `xml:",chardata"`
			Count  string `xml:"count,attr"`
			Filter string `xml:"filter,attr"`
			ID     string `xml:"id,attr"`
			Tag    string `xml:"tag,attr"`
		} `xml:"Genre"`
		Director struct {
			Text   string `xml:",chardata"`
			Filter string `xml:"filter,attr"`
			ID     string `xml:"id,attr"`
			Tag    string `xml:"tag,attr"`
		} `xml:"Director"`
		Writer struct {
			Text   string `xml:",chardata"`
			Filter string `xml:"filter,attr"`
			ID     string `xml:"id,attr"`
			Tag    string `xml:"tag,attr"`
		} `xml:"Writer"`
		Producer []struct {
			Text   string `xml:",chardata"`
			Count  string `xml:"count,attr"`
			Filter string `xml:"filter,attr"`
			ID     string `xml:"id,attr"`
			Tag    string `xml:"tag,attr"`
		} `xml:"Producer"`
		Country struct {
			Text   string `xml:",chardata"`
			Count  string `xml:"count,attr"`
			Filter string `xml:"filter,attr"`
			ID     string `xml:"id,attr"`
			Tag    string `xml:"tag,attr"`
		} `xml:"Country"`
		Role []struct {
			Text   string `xml:",chardata"`
			Count  string `xml:"count,attr"`
			Filter string `xml:"filter,attr"`
			ID     string `xml:"id,attr"`
			Role   string `xml:"role,attr"`
			Tag    string `xml:"tag,attr"`
			Thumb  string `xml:"thumb,attr"`
		} `xml:"Role"`
		Similar []struct {
			Text   string `xml:",chardata"`
			Count  string `xml:"count,attr"`
			Filter string `xml:"filter,attr"`
			ID     string `xml:"id,attr"`
			Tag    string `xml:"tag,attr"`
		} `xml:"Similar"`
		User struct {
			Text  string `xml:",chardata"`
			ID    string `xml:"id,attr"`
			Thumb string `xml:"thumb,attr"`
			Title string `xml:"title,attr"`
		} `xml:"User"`
		Player struct {
			Text                string `xml:",chardata"`
			Address             string `xml:"address,attr"`
			Device              string `xml:"device,attr"`
			MachineIdentifier   string `xml:"machineIdentifier,attr"`
			Model               string `xml:"model,attr"`
			Platform            string `xml:"platform,attr"`
			PlatformVersion     string `xml:"platformVersion,attr"`
			Product             string `xml:"product,attr"`
			Profile             string `xml:"profile,attr"`
			RemotePublicAddress string `xml:"remotePublicAddress,attr"`
			State               string `xml:"state,attr"`
			Title               string `xml:"title,attr"`
			Vendor              string `xml:"vendor,attr"`
			Version             string `xml:"version,attr"`
			Local               string `xml:"local,attr"`
			Relayed             string `xml:"relayed,attr"`
			Secure              string `xml:"secure,attr"`
			UserID              string `xml:"userID,attr"`
		} `xml:"Player"`
		Session struct {
			Text      string `xml:",chardata"`
			ID        string `xml:"id,attr"`
			Bandwidth string `xml:"bandwidth,attr"`
			Location  string `xml:"location,attr"`
		} `xml:"Session"`
		TranscodeSession struct {
			Text                    string `xml:",chardata"`
			Key                     string `xml:"key,attr"`
			Throttled               string `xml:"throttled,attr"`
			Complete                string `xml:"complete,attr"`
			Progress                string `xml:"progress,attr"`
			Speed                   string `xml:"speed,attr"`
			Duration                string `xml:"duration,attr"`
			Remaining               string `xml:"remaining,attr"`
			Context                 string `xml:"context,attr"`
			SourceVideoCodec        string `xml:"sourceVideoCodec,attr"`
			SourceAudioCodec        string `xml:"sourceAudioCodec,attr"`
			VideoDecision           string `xml:"videoDecision,attr"`
			AudioDecision           string `xml:"audioDecision,attr"`
			Protocol                string `xml:"protocol,attr"`
			Container               string `xml:"container,attr"`
			VideoCodec              string `xml:"videoCodec,attr"`
			AudioCodec              string `xml:"audioCodec,attr"`
			AudioChannels           string `xml:"audioChannels,attr"`
			TranscodeHwRequested    string `xml:"transcodeHwRequested,attr"`
			TranscodeHwFullPipeline string `xml:"transcodeHwFullPipeline,attr"`
			TimeStamp               string `xml:"timeStamp,attr"`
			MaxOffsetAvailable      string `xml:"maxOffsetAvailable,attr"`
			MinOffsetAvailable      string `xml:"minOffsetAvailable,attr"`
		} `xml:"TranscodeSession"`
	} `xml:"Video"`
}