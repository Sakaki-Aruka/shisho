package util

import "github.com/spf13/cobra"

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