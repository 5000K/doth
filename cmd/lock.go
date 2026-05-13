package cmd

import (
	"fmt"
	"os"

	"github.com/5000K/doth/model"
	"github.com/5000K/doth/util"
	"github.com/spf13/cobra"
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Locks the project to a specific version of doth",
	Long: `Locks the project to a specific version of doth (defaults to the current version).
	
This only affects the wrapper script. If you installed doth using another source, this won't have any effect.`,
	Run: func(cmd *cobra.Command, args []string) {
		version, err := cmd.Flags().GetString("version")
		if err != nil {
			fmt.Printf("Error getting version flag: %v\n", err)
			return
		}

		loc := util.CleanPath(model.DothLockFileLocation)

		if ok, _ := util.Exists(loc); ok {
			err := util.Delete(loc)
			if err != nil {
				fmt.Printf("Failed to overwrite existing lock: %v\n", err)
				return
			}
		}

		err = os.WriteFile(loc, []byte(version), 0644)
		if err != nil {
			fmt.Printf("Failed to lock the project: %v\n", err)
			return
		}

		fmt.Println("Project locked.")
	},
}

// lockCmd represents the lock command
var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlocks the project",
	Long: `Unlocks the project.
	
This only affects the wrapper script. If you installed doth using another source, this won't have any effect.`,
	Run: func(cmd *cobra.Command, args []string) {

		loc := util.CleanPath(model.DothLockFileLocation)

		if ok, _ := util.Exists(loc); !ok {
			fmt.Println("Project wasn't locked.")
			return
		}

		err := util.Delete(loc)
		if err != nil {
			fmt.Printf("Failed to unlock the project: %v\n", err)
			return
		}

		fmt.Println("Project unlocked.")
	},
}

func init() {
	versionCmd.AddCommand(lockCmd)
	versionCmd.AddCommand(unlockCmd)

	lockCmd.Flags().StringP("version", "v", model.Version, "The version of doth to lock to (defaults to the current version).")
}
