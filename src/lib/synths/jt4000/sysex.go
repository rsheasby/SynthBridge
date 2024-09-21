package jt4000

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"time"
)

func (s *Synth) parseIncomingSysex(data []byte) (err error) {
	defer s.wg.Done()
	if len(data) < 64 {
		return errors.New("sysex message is too short")
	}
	if bytes.HasPrefix(data, []byte{0x00, 0x20, 0x32, 0x00, 0x01, 0x38, 0x15}) {
		log.Println("Received current patch")
		err := s.parseCurrentPatchData(data[7:71])
		if err != nil {
			return err
		}
	} else if bytes.HasPrefix(data, []byte{0x00, 0x20, 0x32, 0x00, 0x01, 0x38, 0x10}) {
		log.Println("Received all patch names")
		patchNames := []string{}
		for i := 7; i+64 <= len(data); i += 64 {
			patchName, err := parsePatchName(data[i : i+64])
			if err != nil {
				return err
			}
			patchNames = append(patchNames, patchName)
		}
		s.PatchNames = patchNames
	}
	time.Sleep(10 * time.Millisecond)

	return
}

func (s *Synth) parseCurrentPatchData(patchData []byte) error {
	if len(patchData) != 64 {
		return errors.New("patch data length is not 64")
	}
	s.CurrentPatchName = strings.TrimSpace(string(patchData[55:64]))
	for id, param := range s.SelectionParams {
		param.StoredValueIndex = int(patchData[param.bytePosition])
		param.StoredValue = param.PossibleValues[param.StoredValueIndex]
		param.CurrentValueIndex = param.StoredValueIndex
		param.CurrentValue = param.StoredValue
		s.SelectionParams[id] = param
	}
	for id, param := range s.IntParams {
		param.StoredValue = int(patchData[param.bytePosition])
		param.CurrentValue = param.StoredValue
		s.IntParams[id] = param
	}
	return nil
}

func parsePatchName(patchData []byte) (patchName string, err error) {
	if len(patchData) != 64 {
		return patchName, errors.New("patch data length is not 64")
	}
	return strings.TrimSpace(string(patchData[55:64])), nil
}

func (s *Synth) GetAllPatchNames() (err error) {
	s.wg.Wait()
	s.wg.Add(1)
	log.Println("Requesting all patch names")
	allPatchesCmd := []byte{0xF0, 0x00, 0x20, 0x32, 0x00, 0x01, 0x38, 0x00, 0x20, 0xF7}

	s.outPort.Send(allPatchesCmd)
	s.wg.Wait()
	return
}

func (s *Synth) GetCurrentPatchDetails() (err error) {
	s.wg.Wait()
	s.wg.Add(1)
	log.Println("Requesting current patch details")
	getCurrentPatchCmd := []byte{0xF0, 0x00, 0x20, 0x32, 0x00, 0x01, 0x38, 0x00, 0x00, 0xF7}

	s.outPort.Send(getCurrentPatchCmd)
	s.wg.Wait()
	return
}

func (s *Synth) guessPatchIndex() {
	for i, name := range s.PatchNames {
		if name == s.CurrentPatchName {
			s.CurrentPatchIndex = i
			return
		}
	}
}
