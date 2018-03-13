// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// fingerprintCmd represents the fingerprint command
var fingerprintCmd = &cobra.Command{
	Use:   "fingerprint <private-key-file>",
	Short: "Generate github-style fingerprint from private key file",
	Long:  `Generate github-style fingerprint from private key file`,
	RunE:  fingerprint,
}

func init() {
	rootCmd.AddCommand(fingerprintCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fingerprintCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fingerprintCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func fingerprint(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("must specify a private key file")
	}

	privateKeyFile, err := filepath.Abs(args[0])
	if err != nil {
		return errors.New("invalid keyfile path")
	}
	publicKeyFile := privateKeyFile + ".pub"

	if _, err = os.Stat(publicKeyFile); err == nil {
		fmt.Printf("found public keyfile %s\n", publicKeyFile)
	} else {
		fmt.Printf("Generating public keyfile %s\n\n", publicKeyFile)
		err = generatePublicKey(privateKeyFile, publicKeyFile)
		if err != nil {
			return err
		}
	}

	fingerprintKeyCmd := exec.Command("ssh-keygen",
		"-l",
		"-E", "md5",
		"-f", publicKeyFile,
	)

	fingerprint, err := fingerprintKeyCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to create public key from private key file %s", privateKeyFile)
	}

	fmt.Printf("Github-style fingerprint: %s\n", fingerprint)

	return nil
}

func generatePublicKey(privateKeyFile, destinationFilePath string) error {
	// open the out file for writing
	outfile, err := os.Create(destinationFilePath)
	if err != nil {
		return fmt.Errorf("failed to create public key from private key file %s", destinationFilePath)
	}
	defer outfile.Close()

	publicKeyCmd := exec.Command("ssh-keygen",
		"-y",
		"-f", privateKeyFile,
	)
	publicKeyCmd.Stdout = outfile

	err = publicKeyCmd.Start()
	if err != nil {
		return fmt.Errorf("failed to generate public key from private key file %s", destinationFilePath)
	}

	publicKeyCmd.Wait()

	return nil
}
