// Copyright © 2018 Andrew Edstrom
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
	"os/exec"

	"github.com/spf13/cobra"
)

var email string

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Short: "Generate ssh key pair",
	Long:  "Generate ssh key pair",
	Use:   "ssh <new-keyfile-name>",
	RunE:  ssh,
}

func init() {
	rootCmd.AddCommand(sshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	sshCmd.Flags().StringVarP(&email, "email", "e", "", "email to associate with this ssh key")
}

func ssh(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("must specify a name for output key files")
	}

	privateKeyFile := args[0]
	publicKeyFile := privateKeyFile + ".pub"

	fmt.Printf("Generating keys...\n")
	publicKeyCmd := exec.Command("ssh-keygen",
		"-t", "rsa",
		"-b", "4096",
		"-N", "",
		"-C", email,
		"-f", privateKeyFile,
	)

	_, err := publicKeyCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to create public key from private key file %s", privateKeyFile)
	}

	fmt.Printf("Success! Generated keys %s and %s\n", privateKeyFile, publicKeyFile)
	return nil
}
