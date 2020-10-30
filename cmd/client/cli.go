package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/andrejtokarcik/jobworker/client"
	"github.com/andrejtokarcik/jobworker/mtls"
	pb "github.com/andrejtokarcik/jobworker/proto"
)

var (
	// Global flags
	serverAddress string
	timeout       time.Duration
	credsFiles    mtls.CredsFiles
)

func CLI(ctx context.Context) *cobra.Command {
	cmdRoot := &cobra.Command{
		Use:   "client",
		Short: "Client for the Job Worker service",
	}

	cmdRoot.PersistentFlags().StringVarP(&serverAddress, "server", "s", "127.0.0.1:50051", "address of the server to connect to")
	cmdRoot.PersistentFlags().DurationVar(&timeout, "timeout", 5*time.Second, "connection & RPC timeout")
	cmdRoot.PersistentFlags().StringVar(&credsFiles.Cert, "client-cert", "client.crt", "certificate file to use for the client")
	cmdRoot.PersistentFlags().StringVar(&credsFiles.Key, "client-key", "client.key", "private key file to use for the client")
	cmdRoot.PersistentFlags().StringVar(&credsFiles.PeerCACert, "server-ca-cert", "server-ca.crt", "certificate file of the CA to authenticate the server")

	cmdRoot.AddCommand(cmdStartJob(ctx), cmdStopJob(ctx), cmdGetJob(ctx))
	return cmdRoot
}

func cmdStartJob(ctx context.Context) *cobra.Command {
	var (
		jobDir string
		jobEnv []string
	)

	cmdStartJob := &cobra.Command{
		Use:   "start-job [flags] -- shell-cmd",
		Short: "Start a job process",
		Args:  cobra.MinimumNArgs(1),
		Run: withJobWorkerClient(ctx, func(ctx context.Context, client pb.JobWorkerClient, args []string) {
			cmdSpec := &pb.CommandSpec{
				Command: args[0], Args: args[1:], Env: jobEnv, Dir: jobDir,
			}
			req := &pb.StartJobRequest{Command: cmdSpec}
			log.Print("Requesting StartJob: ", req)

			resp, err := client.StartJob(ctx, req)
			if err != nil {
				log.Fatal("StartJob failed: ", err)
			}
			fmt.Println(resp)
		}),
	}

	cmdStartJob.Flags().StringVar(&jobDir, "dir", "", "working directory in which to run the command")
	cmdStartJob.Flags().StringArrayVar(&jobEnv, "env-var", []string{}, "declaration of an environment variable (can be specified multiple times)")
	return cmdStartJob
}

func cmdStopJob(ctx context.Context) *cobra.Command {
	cmdStopJob := &cobra.Command{
		Use:   "stop-job [flags] uuid",
		Short: "Stop a job process",
		Args:  cobra.ExactArgs(1),
		Run: withJobWorkerClient(ctx, func(ctx context.Context, client pb.JobWorkerClient, args []string) {
			req := &pb.StopJobRequest{JobUUID: args[0]}
			log.Print("Requesting StopJob: ", req)
			_, err := client.StopJob(ctx, req)
			if err != nil {
				log.Fatal("StopJob failed: ", err)
			}
		}),
	}
	return cmdStopJob
}

func cmdGetJob(ctx context.Context) *cobra.Command {
	var withLogs bool

	cmdGetJob := &cobra.Command{
		Use:   "get-job [flags] uuid",
		Short: "Get information about a job",
		Run: withJobWorkerClient(ctx, func(ctx context.Context, client pb.JobWorkerClient, args []string) {
			req := &pb.GetJobRequest{JobUUID: args[0], WithLogs: withLogs}
			log.Print("Requesting GetJob: ", req)
			resp, err := client.GetJob(ctx, req)
			if err != nil {
				log.Fatal("GetJob failed: ", err)
			}
			fmt.Println(resp)
		}),
	}

	cmdGetJob.Flags().BoolVar(&withLogs, "with-logs", false, "whether to include stdout/stderr outputs")
	return cmdGetJob
}

func withJobWorkerClient(ctx context.Context, callRPC func(context.Context, pb.JobWorkerClient, []string)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		creds, err := mtls.NewClientCreds(credsFiles)
		if err != nil {
			log.Fatal("Failed to load mTLS credentials: ", err)
		}

		conn, err := client.DialContextWithTimeout(
			context.Background(),
			timeout,
			serverAddress,
			grpc.WithTransportCredentials(creds),
		)
		if err != nil {
			log.Fatal("Failed to dial server: ", err)
		}
		defer conn.Close()
		log.Print("Successfully connected to server at ", serverAddress)

		callRPC(ctx, pb.NewJobWorkerClient(conn), args)
	}
}
