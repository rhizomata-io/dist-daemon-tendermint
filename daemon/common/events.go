package common

import "github.com/rhizomata-io/dist-daemon-tendermint/types"

const (
	EventScopeDaemon  = types.EventScope("daemon")
)

var (
	daemonEventBus = types.RegisterEventBus(EventScopeDaemon)
)

type DaemonEvent types.Event

func PublishDaemonEvent(event DaemonEvent) {
	daemonEventBus.Publish(event)
}

func SubscribeDaemonEvent(path types.EventPath, name string, handler types.EventHandler) {
	daemonEventBus.Subscribe(path, name, handler)
}

func UnsubscribeDaemonEvent(path types.EventPath, name string) {
	daemonEventBus.Unsubscribe(path, name)
}
