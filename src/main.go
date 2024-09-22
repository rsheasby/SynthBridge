package main

import (
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	router := NewRouter()
	router.KnobEventHandler = router.KnobEventDispatcher([]KnobEventHandler{
		router.KnobIntDualSpeedScaler(router.KnobDualIntParamController("osc1Coarse", "osc1Fine", "Osc 1 Tune"), 60, 1200),
		router.KnobIntSpeedScaler(router.KnobIntParamController("osc2Adjustment"), 100),
		router.KnobIntSpeedScaler(router.KnobIntParamController("osc2Coarse"), 100),
		router.KnobIntSpeedScaler(router.KnobIntParamController("osc2Fine"), 100),
		router.KnobIntSpeedScaler(router.KnobIntParamController("vcaAttack"), 100),
		router.KnobIntSpeedScaler(router.KnobIntParamController("vcaDecay"), 100),
		router.KnobIntSpeedScaler(router.KnobIntParamController("vcaSustain"), 100),
		router.KnobIntSpeedScaler(router.KnobIntParamController("vcaRelease"), 100),
	})
	router.Run()
}
