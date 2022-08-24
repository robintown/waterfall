/*
Copyright 2022 The Matrix.org Foundation C.I.C.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"sync"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

type Subscriber struct {
	Track *webrtc.TrackLocalStaticRTP

	mutex     sync.RWMutex
	logger    *logrus.Entry
	call      *Call
	sender    *webrtc.RTPSender
	publisher *Publisher
}

func NewSubscriber(call *Call) *Subscriber {
	subscriber := new(Subscriber)

	subscriber.call = call
	subscriber.logger = call.logger

	return subscriber
}

func (s *Subscriber) initLoggingWithTrack(track *webrtc.TrackRemote) {
	s.logger = s.call.logger.WithFields(logrus.Fields{
		"track_id":   (*track).ID(),
		"track_kind": (*track).Kind(),
		"stream_id":  (*track).StreamID(),
	})
}

func (s *Subscriber) Subscribe(publisher *Publisher) {
	s.initLoggingWithTrack(publisher.Track)

	track, err := webrtc.NewTrackLocalStaticRTP(
		publisher.Track.Codec().RTPCodecCapability,
		publisher.Track.ID(),
		publisher.Track.StreamID(),
	)
	if err != nil {
		s.logger.WithError(err).Error("failed to create local static RTP track")
	}

	sender, err := s.call.PeerConnection.AddTrack(track)
	if err != nil {
		s.logger.WithError(err).Error("failed to add track to peer connection")
	}

	s.mutex.Lock()
	s.Track = track
	s.sender = sender
	s.publisher = publisher
	s.mutex.Unlock()

	publisher.AddSubscriber(s)

	s.logger.Info("subscribed")
}

func (s *Subscriber) Unsubscribe() {
	if s.call.PeerConnection.ConnectionState() != webrtc.PeerConnectionStateClosed {
		err := s.call.PeerConnection.RemoveTrack(s.sender)
		if err != nil {
			s.logger.WithError(err).Error("failed to remove track")
		}
	}

	s.call.RemoveSubscriber(s)

	s.logger.Info("unsubscribed")
}
