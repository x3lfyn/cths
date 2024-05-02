package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

var GlobalState State

func main() {
	rootCmd := &cobra.Command{
		Use:   "cths",
		Short: "simple test web server with tui",
	}

	rootCmd.PersistentFlags().StringVarP(&GlobalState.listenAddress, "address", "a", ":6969", "Address to listen on")

	fileserverCmd := &cobra.Command{
		Use:   "fileserver",
		Short: "Act like simple file server",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			println("run of fs called!!!!")

			path := args[0]

			handler := func(writer http.ResponseWriter, request *http.Request) {
				http.FileServer(http.Dir(path)).ServeHTTP(writer, request)
			}

			RunServerAndTui(handler)
		},
	}

	fileCmd := &cobra.Command{
		Use:   "file",
		Short: "Serve one file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]

			handler := func(writer http.ResponseWriter, request *http.Request) {
				file, err := os.ReadFile(path)
				if err != nil {
					panic(err)
				}

				writer.WriteHeader(http.StatusOK)

				_, err = writer.Write(file)
				if err != nil {
					panic(err)
				}
			}

			RunServerAndTui(handler)
		},
	}

	stringCmd := &cobra.Command{
		Use:   "string",
		Short: "Respond with static string",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			text := args[0]

			handler := func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Header().Set("Content-Type", "text/plain")

				_, err := writer.Write([]byte(text))
				if err != nil {
					panic(err)
				}
			}

			RunServerAndTui(handler)
		},
	}

	rootCmd.AddCommand(fileserverCmd)
	rootCmd.AddCommand(fileCmd)
	rootCmd.AddCommand(stringCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Some error occurred '%s'", err)
		os.Exit(1)
	}
}
