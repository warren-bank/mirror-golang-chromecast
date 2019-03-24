// Copyright © 2019 Jonathan Pentecost <pentecostjonathan@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/vishen/go-chromecast/tts"
)

// ttsCmd represents the tts command
var ttsCmd = &cobra.Command{
	Use:   "tts <message>",
	Short: "text-to-speech",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Printf("expected exactly one argument to convert to speech\n")
			return
		}

		googleServiceAccount, _ := cmd.Flags().GetString("google-service-account")
		if googleServiceAccount == "" {
			fmt.Printf("--google-service-account is required\n")
			return
		}

		b, err := ioutil.ReadFile(googleServiceAccount)
		if err != nil {
			fmt.Printf("unable to open google service account file: %v\n", err)
			return
		}

		app, err := castApplication(cmd, args)
		if err != nil {
			fmt.Printf("unable to get cast application: %v\n", err)
			return
		}

		data, err := tts.Create("Hello, World!", b)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		f, err := ioutil.TempFile("", "go-chromecast-tts")
		if err != nil {
			fmt.Printf("unable to create temp file: %v", err)
			return
		}
		defer os.Remove(f.Name())

		if _, err := f.Write(data); err != nil {
			fmt.Printf("unable to write to temp file: %v\n", err)
			return
		}
		if err := f.Close(); err != nil {
			fmt.Printf("unable to close temp file: %v\n", err)
			return
		}

		if err := app.Load(f.Name(), "audio/mp3", false); err != nil {
			fmt.Printf("unable to load media to device: %v\n", err)
			return
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(ttsCmd)
	ttsCmd.Flags().String("google-service-account", "", "google service account JSON file")
}
