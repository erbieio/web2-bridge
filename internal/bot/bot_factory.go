package bot

var (
	factory *Factory
)

type Factory struct {
	Jobs map[string]BotInterface
}

func (f *Factory) Register(j ...BotInterface) {
	for _, v := range j {
		f.Jobs[v.App()] = v
	}
}

func (f *Factory) Do() {
	for _, v := range f.Jobs {
		go v.Do()
	}
}

func GetFacotory() *Factory {
	if factory == nil {
		factory = &Factory{
			Jobs: make(map[string]BotInterface),
		}
	}
	return factory
}
