package workerlist

type Entry struct {
	Path string
}

type Worklist struct {
	jobs chan Entry
}

func (w *Worklist) Add(work Entry) {
	w.jobs <- work
}
func (w *Worklist) Next() Entry {
	return <-w.jobs
}

func New(bufSize int) Worklist {
	return Worklist{jobs: make(chan Entry, bufSize)}
}

func NewJob(path string) Entry {
	return Entry{Path: path}
}

func (w *Worklist) Finilize(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		w.Add(Entry{Path:  ""})
	}
}
