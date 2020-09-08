/*
   Copyright The containerd-compose Authors.

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

/*
   file created by mc256.com in 2020
*/

package logger

import (
	"log"
	"strings"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

func N(format string, v ...interface{}) {
	log.Printf(" [containerd-compose] "+format, v...)
}

type StreamLogger struct {
	containerId string
	label       string
}

func NewStreamLogger(containerId string, label string) (s *StreamLogger) {
	return &StreamLogger{
		containerId: containerId,
		label:       label,
	}
}

func (s StreamLogger) Write(data []byte) (n int, err error) {
	for _, t := range strings.Split(string(data), "\n") {
		if len(t) != 0 {
			log.Printf(" (%3s) [%s] %s", s.label, s.containerId, t)
		}
	}
	return len(data), nil
}
