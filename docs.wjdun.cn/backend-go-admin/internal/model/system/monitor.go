package system

type MonitorOverview struct {
	AppName          string `json:"appName"`
	Environment      string `json:"environment"`
	Address          string `json:"address"`
	StartTime        string `json:"startTime"`
	Uptime           string `json:"uptime"`
	RequestIDEnabled bool   `json:"requestIdEnabled"`
	CurrentTime      string `json:"currentTime"`
}

type MonitorDatabase struct {
	Status             string `json:"status"`
	Driver             string `json:"driver"`
	OpenConnections    int    `json:"openConnections"`
	InUseConnections   int    `json:"inUseConnections"`
	IdleConnections    int    `json:"idleConnections"`
	WaitCount          int64  `json:"waitCount"`
	WaitDuration       string `json:"waitDuration"`
	MaxOpenConnections int    `json:"maxOpenConnections"`
}

type MonitorRedis struct {
	Status       string `json:"status"`
	Address      string `json:"address"`
	DB           int    `json:"db"`
	PoolSize     int    `json:"poolSize"`
	PoolHits     uint32 `json:"poolHits"`
	PoolMisses   uint32 `json:"poolMisses"`
	Timeouts     uint32 `json:"timeouts"`
	TotalConns   uint32 `json:"totalConns"`
	IdleConns    uint32 `json:"idleConns"`
	StaleConns   uint32 `json:"staleConns"`
	PingDuration string `json:"pingDuration"`
}

type MonitorRuntime struct {
	GoVersion     string `json:"goVersion"`
	Goroutines    int    `json:"goroutines"`
	CPUs          int    `json:"cpus"`
	MemoryAllocMB uint64 `json:"memoryAllocMB"`
	MemorySysMB   uint64 `json:"memorySysMB"`
	HeapAllocMB   uint64 `json:"heapAllocMB"`
	HeapSysMB     uint64 `json:"heapSysMB"`
	LastGC        string `json:"lastGC"`
	NumGC         uint32 `json:"numGC"`
}

type MonitorHost struct {
	Hostname      string  `json:"hostname"`
	OS            string  `json:"os"`
	Arch          string  `json:"arch"`
	CPULoad       float64 `json:"cpuLoad"`
	TotalMemoryMB uint64  `json:"totalMemoryMB"`
	FreeMemoryMB  uint64  `json:"freeMemoryMB"`
	TotalDiskGB   uint64  `json:"totalDiskGB"`
	FreeDiskGB    uint64  `json:"freeDiskGB"`
}

type MonitorResponse struct {
	Overview MonitorOverview `json:"overview"`
	Database MonitorDatabase `json:"database"`
	Redis    MonitorRedis    `json:"redis"`
	Runtime  MonitorRuntime  `json:"runtime"`
	Host     MonitorHost     `json:"host"`
}
