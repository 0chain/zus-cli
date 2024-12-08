package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/0chain/errors"
	"github.com/0chain/gosdk/constants"
	"github.com/0chain/gosdk/zboxcore/logger"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zboxcore/zboxutil"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

func printTable(files []sdk.FileDiff) {
	header := []string{"Operation", "Path"}
	data := make([][]string, len(files))
	for idx, child := range files {
		data[idx] = []string{child.Op, child.Path}
	}
	util.WriteTable(os.Stdout, header, []string{}, data)
	fmt.Println("")
}

func saveCache(allocationObj *sdk.Allocation, path string, exclPath []string) {
	if len(path) > 0 {
		err := allocationObj.SaveRemoteSnapshot(path, exclPath)
		if err != nil {
			util.PrintError("Failed to save local cache.", err)
			os.Exit(1)
		}
		fmt.Println("Local cache saved.")
	}
}

func filterOperations(lDiff []sdk.FileDiff) (filterDiff []sdk.FileDiff, exclPath []string) {
	for _, f := range lDiff {
		if f.Op == sdk.Update || f.Op == sdk.Upload {
			filterDiff = append(filterDiff, f)
		} else {
			exclPath = append(exclPath, f.Path)
		}
	}
	return
}

func isEmptyUploadOrUpdate(operation string, localPath string) bool {

	if operation != sdk.Update && operation != sdk.Upload {
		return false
	}
	localPath = strings.TrimRight(localPath, "/")
	fileReader, err := os.Open(localPath)
	if err != nil {
		return false
	}
	defer fileReader.Close()

	fileInfo, err := fileReader.Stat()
	if err != nil {
		return false
	}
	if fileInfo.Size() == 0 {
		return true
	}
	return false
}

func filterEmptyFiles(localPath string, lDiff []sdk.FileDiff) (filterDiff []sdk.FileDiff) {
	localPath = strings.TrimRight(localPath, "/")
	for _, f := range lDiff {
		if !isEmptyUploadOrUpdate(f.Op, localPath+f.Path) {
			filterDiff = append(filterDiff, f)
		}
	}
	return
}

func startMultiUploadUpdate(allocationObj *sdk.Allocation, argsSlice []chunkedUploadArgs) error {
	totalOperations := len(argsSlice)
	if totalOperations == 0 {
		return nil
	}
	operationRequests := make([]sdk.OperationRequest, totalOperations)
	// wg := &sync.WaitGroup{}
	waitTimes := len(argsSlice)
	statusBar := util.NewStatusBar(waitTimes)
	for idx, args := range argsSlice {
		// statusBar := &StatusBar{wg: wg}
		// wg.Add(1)
		fileReader, err := os.Open(args.localPath)
		if err != nil {
			return err
		}
		defer fileReader.Close()

		fileInfo, err := fileReader.Stat()
		if err != nil {
			return err
		}

		mimeType, err := zboxutil.GetFileContentType(path.Ext(args.localPath), fileReader)
		if err != nil {
			return err
		}

		remotePath, fileName, err := fullPathAndFileNameForUpload(args.localPath, args.remotePath)
		if err != nil {
			return err
		}

		fileMeta := sdk.FileMeta{
			Path:       args.localPath,
			ActualSize: fileInfo.Size(),
			MimeType:   mimeType,
			RemoteName: fileName,
			RemotePath: remotePath,
		}
		options := []sdk.ChunkedUploadOption{
			sdk.WithThumbnailFile(args.thumbnailPath),
			sdk.WithEncrypt(args.encrypt),
			sdk.WithStatusCallback(statusBar),
			sdk.WithChunkNumber(args.chunkNumber),
		}
		operationRequests[idx] = sdk.OperationRequest{
			FileMeta:      fileMeta,
			FileReader:    fileReader,
			OperationType: constants.FileOperationInsert,
			Opts:          options,
		}
		if args.isUpdate {
			operationRequests[idx].OperationType = constants.FileOperationUpdate
		}
	}

	err := allocationObj.DoMultiOperation(operationRequests)
	if err != nil {
		logger.Logger.Error("[update/upload]", err)
		return err
	}
	// wg.Wait()
	for i := 0; i < waitTimes; i++ {
		err := <-statusBar.CmdErr
		if err != nil {
			os.Exit(1)
		}
	}
	return nil

}

// syncCmd represents sync command
var syncCmd = &cobra.Command{
	Use:    "sync",
	Short:  "Sync files to/from blobbers",
	Long:   `Sync all files to/from blobbers from/to a localpath`,
	Args:   cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags() // fflags is a *flag.FlagSet
		if fflags.Changed("localpath") == false {
			util.PrintError("Error: localpath flag is missing")
			os.Exit(1)
		}

		localpath := cmd.Flag("localpath").Value.String()
		remotepath := cmd.Flag("remotepath").Value.String()
		if !isAbsolutePathAndDirectory(remotepath) {
			util.PrintError("Error: remotepath should be absolute, and path of directory")
			os.Exit(1)
		}
		if len(localpath) == 0 {
			util.PrintError("Error: localpath flag is missing")
			os.Exit(1)
		}

		encryptpath := cmd.Flag("encryptpath").Value.String()

		if fflags.Changed("allocation") == false { // check if the flag "path" is set
			util.PrintError("Error: allocation flag is missing") // If not, we'll let the user know
			os.Exit(1)                                           // and os.Exit(1)
		}
		allocationID := cmd.Flag("allocation").Value.String()

		localcache := ""
		if fflags.Changed("localcache") {
			localcache = cmd.Flag("localcache").Value.String()
		}
		exclPath := []string{}
		if fflags.Changed("excludepath") {
			exclPath, _ = cmd.Flags().GetStringArray("excludepath")
		}

		verifyDownload, err := cmd.Flags().GetBool("verifydownload")
		if err != nil {
			util.PrintError("Error: ", err)
			os.Exit(1)
		}

		allocationObj, err := sdk.GetAllocation(allocationID)
		if err != nil {
			util.PrintError("Error fetching the allocation", err)
			os.Exit(1)
		}

		fileMetas := make(map[string]*sdk.ConsolidatedFileMeta)
		// wg := &sync.WaitGroup{}
		// Create filter
		filter := []string{".DS_Store", ".git"}

		uploadOnly, _ := cmd.Flags().GetBool("uploadonly")

		lDiff, err := allocationObj.GetAllocationDiff(localcache, localpath, filter, exclPath, remotepath)
		if err != nil {
			util.PrintError("Error getting diff.", err)
			os.Exit(1)
		}

		if uploadOnly {
			var otherPaths []string
			lDiff, otherPaths = filterOperations(lDiff)
			exclPath = append(exclPath, otherPaths...)
		}

		lDiff = filterEmptyFiles(localpath, lDiff)

		if len(lDiff) > 0 {
			printTable(lDiff)
		} else {
			fmt.Println("Already up to date")
			saveCache(allocationObj, localcache, exclPath)
			return
		}
		argsSlice := make([]chunkedUploadArgs, 0)
		downloadSlice := make([]MultiDownloadOption, 0)
		// downloadStatusBars := make([]*util.StatusBar, 0)
		waitTimes := 0
		for _, f := range lDiff {
			localpath = strings.TrimRight(localpath, "/")
			lPath := filepath.Join(localpath, f.Path)
			fileRemotePath, err := getFullRemotePath(f.Path, remotepath)
			if err != nil {
				return
			}

			switch f.Op {
			case sdk.Download:
				waitTimes++
				option := MultiDownloadOption{
					RemotePath: f.Path,
					LocalPath:  lPath,
					DownloadOp: 1,
				}
				downloadSlice = append(downloadSlice, option)
				// statusBar := util.NewStatusBar(1)
				// downloadStatusBars = append(downloadStatusBars, statusBar)

			case sdk.Upload:

				encrypt := len(encryptpath) != 0 && strings.Contains(lPath, encryptpath)
				argsSlice = append(argsSlice, chunkedUploadArgs{
					localPath:     lPath,
					thumbnailPath: "",
					remotePath:    fileRemotePath,
					encrypt:       encrypt,
					chunkNumber:   syncChunkNumber,
				})

			case sdk.Update:

				encrypt := len(encryptpath) != 0 && strings.Contains(lPath, encryptpath)
				argsSlice = append(argsSlice, chunkedUploadArgs{
					localPath:     lPath,
					thumbnailPath: "",
					remotePath:    fileRemotePath,
					encrypt:       encrypt,
					chunkNumber:   syncChunkNumber,
					isUpdate:      true,
				})

			case sdk.Delete:
				fileMeta, err := allocationObj.GetFileMeta(f.Path)
				if err != nil {
					util.PrintError("Error fetching metaData :", err.Error())
				}
				fileMetas[f.Path] = fileMeta
				// TODO: User confirm??
				fmt.Printf("Deleting remote %s...\n", f.Path)
				err = allocationObj.DeleteFile(f.Path)
				if err != nil {
					util.PrintError("Error deleting remote file,", err.Error())
				}
				continue
			case sdk.LocalDelete:
				// TODO: User confirm??
				fmt.Printf("Deleting local %s...\n", lPath)
				err = os.RemoveAll(lPath)
				if err != nil {
					util.PrintError("Error deleting local file.", err.Error())
				}
				continue
			}
			if err != nil {
				util.PrintError(err.Error())
			}
		}
		err = startMultiUploadUpdate(allocationObj, argsSlice)
		if err != nil {
			util.PrintError("\nSync Failed", err.Error())
			return
		}
		statusBar := util.NewStatusBar(waitTimes)
		for i := 0; i < len(downloadSlice); i++ {
			err = allocationObj.DownloadFile(downloadSlice[i].LocalPath, downloadSlice[i].RemotePath, verifyDownload, statusBar, i == len(downloadSlice)-1)
			if err != nil {
				util.PrintError("\nSync Download Failed", err.Error())
				return
			}
		}
		for i := 0; i < waitTimes; i++ {
			err := <-statusBar.CmdErr
			if err != nil {
				util.PrintError("\nSync Download Failed for file with remote path", downloadSlice[i].RemotePath)
			}
		}
		// wg.Wait()
		// for ind, bar := range downloadStatusBars {
		// 	if !bar.success {
		// 		util.PrintError("\nSync Download Failed for file with remote path", downloadSlice[ind].RemotePath)
		// 	}
		// }
		fmt.Println("\nSync Complete")
		saveCache(allocationObj, localcache, exclPath)
		return
	},
}

// statsCmd.Flags().Bool("json", false, "pass this option to print response as json data")

var syncChunkNumber int

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.PersistentFlags().String("allocation", "", "Allocation ID")
	syncCmd.PersistentFlags().String("localpath", "", "Local dir path to sync")
	syncCmd.PersistentFlags().String("remotepath", "/", `Remote dir path from where it sync`)
	syncCmd.PersistentFlags().String("encryptpath", "", "Local dir path to upload as encrypted")
	syncCmd.PersistentFlags().String("localcache", "", `Local cache of remote snapshot.
If file exists, this will be used for comparison with remote.
After sync complete, remote snapshot will be updated to the same file for next use.`)
	syncCmd.PersistentFlags().StringArray("excludepath", []string{}, "Remote folder paths exclude to sync")
	syncCmd.Flags().BoolP("verifydownload", "v", true, "pass this option to verify downloaded blocks")

	syncCmd.MarkFlagRequired("allocation")
	syncCmd.MarkFlagRequired("localpath")
	syncCmd.Flags().Bool("uploadonly", false, "(default false) pass this option to only upload/update the files")

	syncCmd.Flags().IntVarP(&syncChunkNumber, "chunknumber", "", 1, "how many chunks should be uploaded in a http request")

}

func isAbsolutePathAndDirectory(remotePath string) bool {
	if (strings.HasSuffix(remotePath, "/")) &&
		(strings.HasPrefix(remotePath, "/")) {
		return true
	}
	return false
}

func getFullRemotePath(localPath string, remotePath string) (string, error) {
	if !isAbsolutePathAndDirectory(remotePath) {
		return "", errors.New("invalid_path", "Remote path should be absolute, and path of directory")
	}
	localPath = strings.TrimLeft(localPath, "/")
	fullRemotePath := remotePath + localPath
	return fullRemotePath, nil
}
