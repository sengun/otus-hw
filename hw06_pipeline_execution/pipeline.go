package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, terminate In, stages ...Stage) Out {
	if in == nil {
		panic("No input channel specified")
	}

	var out Out
	out = RunStage(in, terminate, func(in In) Out { return in })

	for _, stage := range stages {
		if stage == nil {
			continue
		}

		out = RunStage(out, terminate, stage)
	}

	return out
}

func RunStage(in In, terminate In, stage Stage) Out {
	stageOutputChannel := stage(in)
	out := make(Bi)

	go func() {
		defer func() {
			close(out)
		}()

		for {
			select {
			case <-terminate:
				return
			default:
			}

			select {
			case <-terminate:
				return
			case v, ok := <-stageOutputChannel:
				if !ok {
					return
				}

				out <- v
			}
		}
	}()

	return out
}
