package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/0chain/gosdk/zboxcore/fileref"
	"github.com/0chain/zus-cli/util"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/spf13/cobra"
)

// downloadCmd represents download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download file from blobbers",
	Long:  `download file from blobbers`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags() // fflags is a *flag.FlagSet
		if !(fflags.Changed("remotepath") || fflags.Changed("authticket") || fflags.Changed("multidownloadjson")) {
			util.PrintError("Error: remotepath / authticket flag is missing")
			os.Exit(1)
		}

		remotePath := cmd.Flag("remotepath").Value.String()
		authTicket := cmd.Flag("authticket").Value.String()
		lookupHash := cmd.Flag("lookuphash").Value.String()
		verifyDownload, err := cmd.Flags().GetBool("verifydownload")
		if err != nil {
			util.PrintError("Error: ", err)
			os.Exit(1)
		}

		thumbnail, err := cmd.Flags().GetBool("thumbnail")
		if err != nil {
			util.PrintError("Error: ", err)
			os.Exit(1)
		}

		localPath := cmd.Flag("localpath").Value.String()
		allocationID := cmd.Flag("allocation").Value.String()

		live, _ := cmd.Flags().GetBool("live")

		if live {
			delay, _ := cmd.Flags().GetInt("delay")

			m3u8, err := util.CreateM3u8Downloader(localPath, remotePath, authTicket, allocationID, lookupHash, delay)

			if err != nil {
				util.PrintError("Error: download files and build playlist: ", err)
				os.Exit(1)
			}

			err = m3u8.Start()

			if err != nil {
				util.PrintError("Error: download files and build playlist: ", err)
				os.Exit(1)
			}

			return

		}

		numBlocks, _ := cmd.Flags().GetInt("blockspermarker")
		if numBlocks == 0 {
			numBlocks = 100
		}

		startBlock, _ := cmd.Flags().GetInt64("startblock")
		if startBlock < 1 {
			util.PrintError("Error: start block should not be less than 1")
		}
		endBlock, _ := cmd.Flags().GetInt64("endblock")

		sdk.SetNumBlockDownloads(numBlocks)
		// wg := &sync.WaitGroup{}
		// statusBar := &StatusBar{wg: wg}
		// wg.Add(1)
		statusBar := util.NewStatusBar(1)
		var errE error
		waitTimes := 1
		var allocationObj *sdk.Allocation

		var multidownloadJSON string
		if fflags.Changed("multidownloadjson") {
			multidownloadJSON = cmd.Flag("multidownloadjson").Value.String()
		}

		if len(authTicket) > 0 {
			at, err := sdk.InitAuthTicket(authTicket).Unmarshall()

			if err != nil {
				util.PrintError(err)
				os.Exit(1)
			}

			allocationObj, err = sdk.GetAllocationFromAuthTicket(authTicket)
			if err != nil {
				util.PrintError("Error fetching the allocation", err)
				os.Exit(1)
			}

			var fileName string

			if at.RefType == fileref.FILE {
				fileName = at.FileName
				lookupHash = at.FilePathHash
			} else if len(lookupHash) > 0 {
				fileMeta, err := allocationObj.GetFileMetaFromAuthTicket(authTicket, lookupHash)
				if err != nil {
					util.PrintError("Either remotepath or lookuphash is required when using authticket of directory type")
					os.Exit(1)
				}
				fileName = fileMeta.Name
			} else if len(remotePath) > 0 {
				lookupHash = fileref.GetReferenceLookup(allocationObj.Tx, remotePath)

				pathNames := strings.Split(remotePath, "/")
				fileName = pathNames[len(pathNames)-1]
			} else {
				util.PrintError("Either remotepath or lookuphash is required when using authticket of directory type")
				os.Exit(1)
			}

			if thumbnail {
				errE = allocationObj.DownloadThumbnailFromAuthTicket(localPath,
					authTicket, lookupHash, fileName, verifyDownload, statusBar, true)
			} else {
				if startBlock != 0 || endBlock != 0 {
					errE = allocationObj.DownloadFromAuthTicketByBlocks(
						localPath, authTicket, startBlock, endBlock, numBlocks,
						lookupHash, fileName, verifyDownload, statusBar, true)
				} else {
					errE = allocationObj.DownloadFromAuthTicket(localPath,
						authTicket, lookupHash, fileName, verifyDownload, statusBar, true)
				}
			}
		} else if len(remotePath) > 0 {
			if !fflags.Changed("allocation") { // check if the flag "path" is set
				util.PrintError("Error: allocation flag is missing") // If not, we'll let the user know
				os.Exit(1)                                           // and return
			}
			allocationID := cmd.Flag("allocation").Value.String()
			allocationObj, err = sdk.GetAllocation(allocationID)

			if err != nil {
				util.PrintError("Error fetching the allocation", err)
				os.Exit(1)
			}

			if thumbnail {
				errE = allocationObj.DownloadThumbnail(localPath, remotePath, verifyDownload, statusBar, true)
			} else {
				if startBlock == 1 && endBlock == 0 {
					ds := sdk.CreateFsDownloadProgress()
					errE = allocationObj.DownloadFile(localPath, remotePath, verifyDownload, statusBar, true, sdk.WithDownloadProgressStorer(ds), sdk.WithWorkDir(util.GetHomeDir()))
				} else {
					errE = allocationObj.DownloadFileByBlock(localPath, remotePath, startBlock, endBlock, numBlocks, verifyDownload, statusBar, true)
				}
			}
		} else if len(multidownloadJSON) > 0 {
			if !fflags.Changed("allocation") { // check if the flag "path" is set
				util.PrintError("Error: allocation flag is missing") // If not, we'll let the user know
				os.Exit(1)                                           // and return
			}
			allocationID := cmd.Flag("allocation").Value.String()
			allocationObj, err = sdk.GetAllocation(allocationID)
			if err != nil {
				util.PrintError("Error: getting allocation", err)
				os.Exit(1)
			}

			waitTimes, errE = MultiDownload(allocationObj, multidownloadJSON, statusBar, verifyDownload)
		}

		if errE == nil {
			// wg.Wait()
			for i := 0; i < waitTimes; i++ {
				err := <-statusBar.CmdErr
				if err != nil {
					// util.PrintError("Multidownload error ", i+1, " : ", err)
					os.Exit(1)
				}
			}
		} else {
			util.PrintError("Download failed.", errE.Error())
			os.Exit(1)
		}

	},
}

type MultiDownloadOption struct {
	RemotePath       string `json:"remotePath"`
	LocalPath        string `json:"localPath"`
	DownloadOp       int    `json:"downloadOp"`
	RemoteFileName   string `json:"remoteFileName,omitempty"`   //Required only for file download with auth ticket
	RemoteLookupHash string `json:"remoteLookupHash,omitempty"` //Required only for file download with auth ticket
}

func MultiDownload(a *sdk.Allocation, jsonMultiDownloadOptions string, statusBar *util.StatusBar, verifyDownload bool) (int, error) {
	var options []MultiDownloadOption
	file, err := os.Open(jsonMultiDownloadOptions)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&options)
	if err != nil {
		return 0, err
	}

	lastOp := len(options) - 1
	for i := 0; i <= len(options)-1; i++ {
		// if i > 0 {
		// 	statusBar.wg.Add(1)
		// }
		ds := sdk.CreateFsDownloadProgress()
		if options[i].DownloadOp == 1 {
			err = a.DownloadFile(options[i].LocalPath, options[i].RemotePath, verifyDownload, statusBar, i == lastOp, sdk.WithDownloadProgressStorer(ds), sdk.WithWorkDir(util.GetHomeDir()))
		} else {
			err = a.DownloadThumbnail(options[i].LocalPath, options[i].RemotePath, false, statusBar, i == lastOp)
		}
		if err != nil {
			return 0, err
		}
	}

	return len(options), err
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.PersistentFlags().String("allocation", "", "Allocation ID")
	downloadCmd.PersistentFlags().String("remotepath", "", "Remote path to download")
	downloadCmd.PersistentFlags().String("localpath", "", "Local path of file to download")
	downloadCmd.PersistentFlags().String("authticket", "", "Auth ticket fot the file to download if you dont own it")
	downloadCmd.PersistentFlags().String("lookuphash", "", "The remote lookuphash of the object retrieved from the list")
	downloadCmd.PersistentFlags().String("multidownloadjson", "", "A JSON file containing multi download options")
	downloadCmd.Flags().BoolP("thumbnail", "t", false, "(default false) pass this option to download only the thumbnail")

	downloadCmd.Flags().Int64P("startblock", "s", 1,
		"Pass this option to download from specific block number. It should not be less than 1")
	downloadCmd.Flags().Int64P("endblock", "e", 0, "pass this option to download till specific block number")
	downloadCmd.Flags().IntP("blockspermarker", "b", 10, "pass this option to download multiple blocks per marker")
	downloadCmd.Flags().BoolP("verifydownload", "v", false, "(default false) pass this option to verify downloaded blocks")

	downloadCmd.Flags().Bool("live", false, "(default false) start m3u8 downloader,and automatically generate media playlist(m3u8) on --localpath")
	downloadCmd.Flags().Int("delay", 5, "pass segment duration to generate media playlist(m3u8). only works with --live. default duration is 5s.")

	downloadCmd.MarkFlagRequired("allocation")
	downloadCmd.MarkFlagRequired("localpath")
}
