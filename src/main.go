package main

import (
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	router := NewRouter()
	router.KnobEventHandler = router.KnobEventDispatcher([]KnobEventHandler{
		router.KnobIntSpeedScaler(router.KnobIntParamController("osc1Adjustment"), 100),
		router.KnobIntSpeedScaler(router.KnobIntParamController("osc1Coarse"), 25),
		router.KnobIntSpeedScaler(router.KnobIntParamController("osc1Fine"), 100),
		router.KnobIntSpeedScaler(router.KnobIntParamController("osc2Adjustment"), 100),
		router.KnobIntSpeedScaler(router.KnobIntParamController("osc2Coarse"), 100),
		router.KnobIntSpeedScaler(router.KnobIntParamController("osc2Fine"), 100),
		router.KnobIntSpeedScaler(router.KnobIntParamController("vcaAttack"), 100),
		router.KnobIntSpeedScaler(router.KnobIntParamController("vcaDecay"), 100),
	})
	router.Run()
}
