package util

import (
	"fmt"
	"sync"

	"gopkg.in/cheggaaa/pb.v1"
)

// CHECK: Verify

type StatusBar struct {
	B  *pb.ProgressBar
	Wg *sync.WaitGroup
}

func NewStatusBar() *StatusBar {
	return &StatusBar{
		Wg: &sync.WaitGroup{},
	}
}

func (s *StatusBar) Started(allocationId, filePath string, op int, totalBytes int) {
	s.B = pb.StartNew(totalBytes)
	s.B.Set(0)
}
func (s *StatusBar) InProgress(allocationId, filePath string, op int, completedBytes int, data []byte) {
	s.B.Set(completedBytes)
}

func (s *StatusBar) Completed(allocationId, filePath string, filename string, mimetype string, size int, op int) {
	defer s.Wg.Done()
	if s.B != nil {
		s.B.Finish()
	}

	// s.Success = true
	// if !allocUnderRepair {
	// 	defer s.wg.Done()
	// }
	fmt.Println("Status completed callback. Type = " + mimetype + ". Name = " + filename)
}

func (s *StatusBar) Error(allocationID string, filePath string, op int, err error) {
	defer s.Wg.Done()
	if s.B != nil {
		s.B.Finish()
	}
	// s.Success = false
	// if !allocUnderRepair {
	// 	defer s.wg.Done()
	// }

	var errDetail interface{} = "Unknown Error"
	if err != nil {
		errDetail = err.Error()
	}

	PrintError("Error in file operation:", errDetail)
	// s.C <- errDetail
}

func (s *StatusBar) RepairCompleted(filesRepaired int) {
	// defer s.wg.Done()
	// allocUnderRepair = false
	// s.Success = true
	defer s.Wg.Done()
	fmt.Println("Repair file completed, Total files repaired: ", filesRepaired)
}
