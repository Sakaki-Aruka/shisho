package cmd

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var systemdTemplateCmd = &cobra.Command{
	Use:   "systemd",
	Short: "Print systemd unit file",
	RunE:  systemd,
}

//go:embed shisho.service.template
var unitTemplate string
var db string
var username string

func systemd(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(db) == "" || strings.TrimSpace(username) == "" {
		return errors.New("'dbpath' and 'username' must not be empty.")
	}

	v := map[string]string{
		"DBPATH":   db,
		"USERNAME": username,
	}
	t, _ := template.New("").Parse(unitTemplate)
	buf := bytes.NewBufferString("")
	if err := t.Execute(buf, v); err != nil {
		return errors.New("'failed to replace values'")
	}

	fmt.Println(buf.String())
	return nil
}

func init() {
	rootCmd.AddCommand(systemdTemplateCmd)

	systemdTemplateCmd.Flags().StringVar(&db, "dbpath", "", "database path")
	systemdTemplateCmd.Flags().StringVar(&username, "username", "", "username")
}
