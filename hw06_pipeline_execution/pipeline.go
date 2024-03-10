package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = stage(processStage(out, done))
	}
	return out
}

func processStage(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case item, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- item:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()
	return out
}
