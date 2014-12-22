package volume

/*
	volume.Manager providers interfaces for both provisioning volume backends, and then creating volumes using them.

	There is one volume.Manager per host daemon process (though of course it's not an enforced singleton, because tests behave otherwise).
*/
type Manager struct {
	defaultProvider Provider

	// `map[providerName]provider`
	//
	// It's possible to configure multiple volume providers for a flynn-host daemon.
	// This can be used to create volumes using providers backed by different storage resources,
	// or different volume backends entirely.
	providers map[string]Provider

	// `map[volume.Id]volume`
	volumes map[string]Volume
}

func NewManager(p Provider) *Manager {
	return &Manager{
		defaultProvider: p,
		providers: map[string]Provider{"default": p},
		volumes: map[string]Volume{},
	}
}

/*
	volume.Manager implements the volume.Provider interface by
	delegating NewVolume requests to the default Provider.
*/
func (m *Manager) NewVolume() (Volume, error) {
	return managerProviderProxy{m.defaultProvider, m}.NewVolume()
}

func (m *Manager) GetVolume(id string) Volume {
	return m.volumes[id]
}

/*
	Proxies `volume.Provider` while making sure the manager remains
	apprised of all volume lifecycle events.

	@implements Provider
*/
type managerProviderProxy struct {
	Provider
	m *Manager
}

func (p managerProviderProxy) NewVolume() (Volume, error) {
	v, err := p.Provider.NewVolume() // please don't be a loop
	if err != nil {
		return v, err
	}
	p.m.volumes[v.ID()] = v
	return v, err
}

// TODO: lots of options for configuring, adding, and later selecting providers for use.  will take shape along with API.
