package main

type nodesResponse struct {
	Result nodesResponseResult `json:"result"`
	Code   int                 `json:"code"`
}

type nodesResponseResult struct {
	NodeName           string `json:"name"`
	Release            string `json:"otp_release"`
	Status             string `json:"node_status"`
	MemoryTotal        string `json:"memory_total"`
	MemoryUsed         string `json:"memory_used"`
	ProcessesAvailable int    `json:"process_available"`
	ProcessesUsed      int    `json:"process_used"`
	MaxFds             int    `json:"max_fds"`
	Clients            int    `json:"clients"`
	Load1              string `json:"load1"`
	Load5              string `json:"load5"`
	Load15             string `json:"load15"`
}

type metricsResponse struct {
	Result metricsResponseResult `json:"result"`
	Code   int                   `json:"code"`
}

type metricsResponseResult struct {
	MessagesDropped        int `json:"messages/dropped"`
	PacketsReceived        int `json:"packets/received"`
	PacketsPubcompReceived int `json:"packets/pubcomp/received"`
	PacketsUnsuback        int `json:"packets/unsuback"`
	PacketsPingresp        int `json:"packets/pingresp"`
	PacketsPingreq         int `json:"packets/pingreq"`
	MessagesQos0Sent       int `json:"messages/qos0/sent"`
	MessagesQos2Received   int `json:"messages/qos2/received"`
	PacketsPubcompMissed   int `json:"packets/pubcomp/missed"`
	MessagesRetained       int `json:"messages/retained"`
	PacketsSuback          int `json:"packets/suback"`
	BytesSent              int `json:"bytes/sent"`
	PacketsPubackReceived  int `json:"packets/puback/received"`
	PacketsPubrecReceived  int `json:"packets/pubrec/received"`
	MessagesQos2Sent       int `json:"messages/qos2/sent"`
	PacketsPubrecSent      int `json:"packets/pubrec/sent"`
	PacketsPubackSent      int `json:"packets/puback/sent"`
	PacketsPubrelMissed    int `json:"packets/pubrel/missed"`
	PacketsConnect         int `json:"packets/connect"`
	MessagesQos1Sent       int `json:"messages/qos1/sent"`
	PacketsConnack         int `json:"packets/connack"`
	PacketsPubrelReceived  int `json:"packets/pubrel/received"`
	PacketsPublishReceived int `json:"packets/publish/received"`
	BytesReceived          int `json:"bytes/received"`
	PacketsPubrelSent      int `json:"packets/pubrel/sent"`
	PacketsPubrecMissed    int `json:"packets/pubrec/missed"`
	PacketsSent            int `json:"packets/sent"`
	MessagesQos0Received   int `json:"messages/qos0/received"`
	PacketsPubcompSent     int `json:"packets/pubcomp/sent"`
	MessagesReceived       int `json:"messages/received"`
	MessagesSent           int `json:"messages/sent"`
	PacketsSubscribe       int `json:"packets/subscribe"`
	MessagesQos2Dropped    int `json:"messages/qos2/dropped"`
	PacketsUnsubscribe     int `json:"packets/unsubscribe"`
	MessagesQos1Received   int `json:"messages/qos1/received"`
	PacketsDisconnect      int `json:"packets/disconnect"`
	PacketsPublishSent     int `json:"packets/publish/sent"`
	PacketsPubackMissed    int `json:"packets/puback/missed"`
}

type statsResponse struct {
	Result statsResponseResult `json:"result"`
	Code   int                 `json:"code"`
}

type statsResponseResult struct {
	ClientsCount       int `json:"clients/count"`
	ClientsMax         int `json:"clients/max"`
	RetainedCount      int `json:"retained/count"`
	RetainedMax        int `json:"retained/max"`
	RoutesCount        int `json:"routes/count"`
	RoutesMax          int `json:"routes/max"`
	SessionsCount      int `json:"sessions/count"`
	SessionsMax        int `json:"sessions/max"`
	SubscribersCount   int `json:"subscribers/count"`
	SubscribersMax     int `json:"subscribers/max"`
	SubscriptionsCount int `json:"subscriptions/count"`
	SubscriptionsMax   int `json:"subscriptions/max"`
	TopicsCount        int `json:"topics/count"`
	TopicsMax          int `json:"topics/max"`
}

type managementResponse struct {
	Result []ManagementResponseResult `json:"result"`
	Code   int                        `json:"code"`
}

// ManagementResponseResult contains the management data for a single node
type ManagementResponseResult struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Sysdescr   string `json:"sysdescr"`
	Uptime     string `json:"uptime"`
	Datetime   string `json:"datetime"`
	OtpRelease string `json:"otp_release"`
	NodeStatus string `json:"node_status"`
}

type combinedResponse struct {
	nodes       nodesResponse
	metrics     metricsResponse
	stats       statsResponse
	ClusterSize int
}
