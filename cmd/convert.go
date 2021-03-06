/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package cmd

import (
	"os"
	"path/filepath"

	"github.com/loadimpact/k6/converter/har"
	"github.com/spf13/cobra"
)

var output = "har-script.js"

var (
	enableChecks bool
	threshold    uint
	only         []string
	skip         []string
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert a HAR file to a k6 script",
	Long:  "Convert a HAR (HTTP Archive) file to a k6 script",
	Example: `
  # Convert a HAR file to a k6 script.
  k6 convert -O har-session.js session.har

  # Convert a HAR file to a k6 script creating requests only for the given domain/s.
  k6 convert -O har-session.js --only yourdomain.com,additionaldomain.com session.har

  # Run the k6 script.
  k6 run har-session.js`[1:],
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse the HAR file
		filePath, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		r, err := os.Open(filePath)
		if err != nil {
			return err
		}
		h, err := har.Decode(r)
		if err != nil {
			return err
		}
		if err := r.Close(); err != nil {
			return err
		}

		script, err := har.Convert(h, enableChecks, threshold, only, skip)
		if err != nil {
			return err
		}

		// Write script content to output
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		if _, err := f.WriteString(script); err != nil {
			return err
		}
		if err := f.Sync(); err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(convertCmd)
	convertCmd.Flags().SortFlags = false
	convertCmd.Flags().StringVarP(&output, "output", "O", output, "k6 script output filename")
	convertCmd.Flags().StringSliceVarP(&only, "only", "", []string{}, "include only requests from the given domains")
	convertCmd.Flags().StringSliceVarP(&skip, "skip", "", []string{}, "skip requests from the given domains")
	convertCmd.Flags().UintVarP(&threshold, "batch-threshold", "", 500, "split requests in different batch statements when the start time difference between subsequent requests is smaller than the given value in ms. A sleep will be added between the batch statements.")
	convertCmd.Flags().BoolVarP(&enableChecks, "enable-status-code-checks", "", false, "add a check for each http status response")
}
