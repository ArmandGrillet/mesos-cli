package cmd

import (
	"errors"

	"github.com/ArmandGrillet/mesos-cli/mesos"
	"github.com/dcos/dcos-cli/api"
	"github.com/dcos/dcos-core-cli/pkg/pluginutil"
	gomesos "github.com/mesos/mesos-go/api/v1/lib"
	"github.com/spf13/cobra"
)

// newCmdMesosSandbox ataches the CLI to a cluster.
func newCmdMesosSandbox(ctx api.Context) *cobra.Command {
	// var removeDir bool
	cmd := &cobra.Command{
		Use:   "sandbox <task-id> <absolute-path-to-dir>",
		Short: "Downloads the sandbox of a given task in a given directory",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// We need to confirm that the task is in the cluster.
			// We then need the executor ID, the framework ID, and
			// the agent ID so that we can get a path to the sandbox.
			client := mesos.NewClient(pluginutil.HTTPClient(""), ctx.Fs())
			state, err := client.State()
			if err != nil {
				return err
			}
			tasks := state.GetTasks
			c := make(chan gomesos.Task, 2)
			go getTask(args[0], tasks.GetTasks(), c)
			go getTask(args[0], tasks.GetCompletedTasks(), c)
			x, y := <-c, <-c
			var task gomesos.Task
			if x.Name == "" && y.Name == "" {
				return errors.New("unable to find task, make sure it is running or completed")
			} else if x.Name == "" {
				task = y
			} else {
				task = x
			}

			lastStatus := task.Statuses[len(task.Statuses)-1]
			containerID := lastStatus.ContainerStatus.GetContainerID()

			var executorID string
			if task.ExecutorID != nil {
				executorID = (*task.ExecutorID).Value
			} else {
				executorID = task.TaskID.Value
			}
			path, err := client.Debug(task.AgentID.Value, task.FrameworkID.Value, executorID, containerID.GetParent().Value)
			if err != nil {
				return err
			}
			// if afero.Exists(ctx.Fs(), args[1]) {
			// 	return fmt.Errorf("unable to download sandbox in '%s', directory already exists", args[1])
			// }
			return client.Download(task.AgentID.Value, path, args[1])
		},
	}
	// cmd.Flags().BoolVar(&removeDir, "force", false, "Remove the local directory where the sandbox will be downloaded")
	return cmd
}

func getTask(taskID string, tasks []gomesos.Task, c chan gomesos.Task) {
	for _, task := range tasks {
		if task.GetTaskID().Value == taskID {
			c <- task
			return
		}
	}
	c <- gomesos.Task{}
}
