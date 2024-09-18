package jt4000

import (
	"bytes"
	"errors"
	"log"
)

func (s *Synth) parseIncomingSysex(data []byte) (err error) {
	defer s.Done()
	if len(data) < 64 {
		return errors.New("sysex message is too short")
	}
	if bytes.HasPrefix(data, []byte{0x00, 0x20, 0x32, 0x00, 0x01, 0x38, 0x15}) {
		log.Println("Received current patch")
		patch, err := parsePatch(data[7:71])
		if err != nil {
			return err
		}
		s.CurrentPatch = patch
	} else if bytes.HasPrefix(data, []byte{0x00, 0x20, 0x32, 0x00, 0x01, 0x38, 0x10}) {
		log.Println("Received all patches")
		patches := []Patch{}
		for i := 7; i+64 <= len(data); i += 64 {
			patch, err := parsePatch(data[i : i+64])
			if err != nil {
				return err
			}
			patches = append(patches, patch)
		}
		s.Patches = patches
	}

	return
}

func parsePatch(data []byte) (patch Patch, err error) {
	if len(data) != 64 {
		return patch, errors.New("patch data length is not 64")
	}
	patch.Name = string(data[55:64])
	patch.Osc1Wave = OscWave(data[0])
	patch.Osc1Adj = data[2]
	patch.Osc1Coarse = data[4]
	patch.Osc1Fine = data[5]

	patch.Osc2Wave = OscWave(data[1])
	if patch.Osc2Wave == OscWaveSuperSaw {
		patch.Osc2Wave = OscWaveNoise
	}
	patch.Osc2Adj = data[3]
	patch.Osc2Coarse = data[6]
	patch.Osc2Fine = data[7]

	patch.OscBalance = data[8] - 0x40

	patch.PortamentoAmount = data[46]
	patch.PortamentoMode = PortamentoMode(data[45])

	patch.RingModEnabled = data[43] == 0x01
	patch.RingModAmount = data[44]

	patch.LFO1Wave = LFOWave(data[47])
	patch.LFO1Destination = LFODestination(data[53])
	patch.LFO1Speed = data[49]
	patch.LFO1Amount = data[50]

	patch.LFO2Wave = LFOWave(data[48])
	patch.LFO2Speed = data[51]
	patch.LFO2Amount = data[52]

	patch.VCFAttack = data[14]
	patch.VCFDecay = data[15]
	patch.VCFSustain = data[16]
	patch.VCFRelease = data[17]
	patch.VCFAmount = data[22]

	patch.VCFCutoff = data[12]
	patch.VCFResonance = data[13]

	patch.VCAAttack = data[18]
	patch.VCADecay = data[19]
	patch.VCASustain = data[20]
	patch.VCARelease = data[21]

	return
}

func (s *Synth) GetAllPatches() (err error) {
	s.Wait()
	s.Add(1)
	log.Println("Requesting all patches")
	allPatchesCmd := []byte{0xF0, 0x00, 0x20, 0x32, 0x00, 0x01, 0x38, 0x00, 0x20, 0xF7}

	s.outPort.Send(allPatchesCmd)
	s.Wait()
	return
}

func (s *Synth) GetCurrentPatch() (err error) {
	s.Wait()
	s.Add(1)
	log.Println("Requesting current patch")
	getCurrentPatchCmd := []byte{0xF0, 0x00, 0x20, 0x32, 0x00, 0x01, 0x38, 0x00, 0x00, 0xF7}

	s.outPort.Send(getCurrentPatchCmd)
	s.Wait()
	return
}
