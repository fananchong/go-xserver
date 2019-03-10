```mermaid
sequenceDiagram
participant Gateway
participant LoginServer
Gateway-->>Gateway: 检查闲置连接
alt 有闲置连接
Gateway-->>Gateway: 本地闲置连接处理
Gateway-->>LobbyServer: 通知账号登出
LobbyServer-->>LobbyServer: 账号登出处理
end
```
