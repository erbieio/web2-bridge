package bot

type BotInterface interface {
	Do() error
	App() string
}
