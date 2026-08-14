/*
Copyright © 2026 Aruka

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"

	shishodb "github.com/Sakaki-Aruka/shisho/shisho-db"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.1.0",
	Use:     "shisho",
	Short:   "Manage your books",
	Long: ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return dbInit(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var dbPath string

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "", "Specify book DB file path")
}

func dbInit(cmd *cobra.Command) error {
	var err error
	dbPath, err = getDBPath(cmd);
	if err != nil {
		return err
	}

	if err := shishodb.Init(dbPath); err != nil {
		return err
	}
	return nil
} 

func getDBPath(cmd *cobra.Command) (string, error) {
	var path string
	v, err := cmd.Flags().GetString("db")
	if err != nil || !cmd.Flags().Changed("db") {
		if p, err := shishodb.GetDefaultPath(); err != nil {
			return "", err
		} else { 
			path = p
		}
	} else {
		path = v
	}
	return path, nil
}

func GetFlagOrNil[T any](cmd *cobra.Command, name string, getter func(string) (T, error)) (*T, error) {
	if !cmd.Flags().Changed(name) {
		return nil, nil
	}

	if v, err := getter(name); err != nil {
		return nil, err
	} else {
		return &v, nil
	}
}
