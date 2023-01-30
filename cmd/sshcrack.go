package cmd

import (
	"fmt"
	"gorun/ssh"

	"github.com/spf13/cobra"
)

var (
	path string
)
var sshCrack = &cobra.Command{
	Use:   "SshCrack",
	Short: "The realization of xshell password cracking",
	PreRun: func(cmd *cobra.Command, args []string) {
		localname, SID := ssh.GetUserSid()
		fmt.Println("[*]find this localname is :" + localname)
		fmt.Println("[*]find this   SID     is :" + SID)
		fmt.Println("[*]crack ssh loading...................")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Lookup("Path").Value.String() == "" {
			m := ssh.Run()
			for _, v := range m {
				ssh.InitXSh(v)
			}
		} else {
			File := cmd.Flags().Lookup("Path").Value.String()
			ssh.RunFile(File)
		}
		// fmt.Printf("cmd.Flags().Lookup(\"path\").Value: %v\n", cmd.Flags().Lookup("Path").Value)
	},
}

func init() {
	rootCmd.AddCommand(sshCrack)
	sshCrack.Flags().StringVarP(&path, "Path", "P", "", "eg -P C:\\xxxx\\xxx.xsh")
}
