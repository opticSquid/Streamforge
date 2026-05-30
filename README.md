# Streamforge
Share what your heart desires

## Description
Streamforge is a live streaming infrastructure platform built from the ground up. It handles real-time RTMP ingest, adaptive multi-bitrate HLS transcoding via FFMPEG, and concurrent multi-stream delivery — all containerised and orchestrated via Docker Compose.
Beyond delivery, Streamforge includes a QoE (Quality of Experience) observability stack that collects real player-side signals such as

- startup time
- rebuffer events
- bitrate switches
- dropped frames

and surfaces them on a live Grafana dashboard via OpenTelemetry. Ad breaks are served mid-stream through a VAST-compliant mock ad server with server-side scheduling, giving the platform a complete content-to-monetisation pipeline.

## Tech Stack
- Go
- FFMPEG
- nginx-rtmp
- HLS
- Shaka Player
- OpenTelemetry
- Prometheus
- Grafana
- Docker Compose

## Architecture
```mermaid
flowchart TB
    subgraph INGEST["📡  INGEST"]
        direction LR
        broadcaster["📡 Broadcaster\nOBS / FFMPEG CLI"]
        nginx["⚡ nginx-rtmp\nRTMP Ingest · :1935"]
    end

    subgraph PROCESSING["🎬  PROCESSING"]
        direction LR
        ffmpeg["🎬 FFMPEG\nReal-time Transcode\n1080p · 720p · 480p"]
        segments["💾 HLS Segments\n.m3u8 + .ts · Shared Volume"]
    end

    subgraph DELIVERY["▶️  DELIVERY"]
        direction LR
        goserver["⚙️ Go Origin Server\nManifests · Stream Keys\nAd Scheduler · QoE API · :8080"]
        shaka["▶️ Shaka Player\nABR · HLS Playback\nQoE Signal Emitter"]
        viewers["👥 Viewers\nMultiple Concurrent\nBrowser Sessions"]
    end

    subgraph ADSTACK["📢  AD STACK"]
        direction LR
        vastserver["📢 Mock VAST Server\nGo HTTP · VAST 2.0 XML\nAd Break Scheduler · :8090"]
        adcreative["🎯 Ad Creative\n15s Test Video · MP4"]
    end

    subgraph OBSERVABILITY["📊  OBSERVABILITY"]
        direction LR
        otel["🔭 OTel Collector\nTraces · Metrics · :4317"]
        prometheus["🔥 Prometheus\nMetrics Store · :9090"]
        grafana["📊 Grafana Dashboard\nQoE Panels · :3000"]
    end

    %% Ingest flow
    broadcaster -->|RTMP| nginx
    nginx -->|RTMP stream| ffmpeg

    %% Processing flow
    ffmpeg -->|write segments| segments
    segments -->|read| goserver

    %% Delivery flow
    goserver -->|HLS manifest + segments| shaka
    shaka -->|renders| viewers

    %% QoE signals — async
    shaka -.->|"QoE signals · POST /qoe"| goserver

    %% Ad flow — async
    goserver -.->|ad break trigger| vastserver
    vastserver -.->|VAST XML| shaka
    adcreative -.->|ad video| shaka

    %% Observability flow
    goserver -.->|spans| otel
    goserver -.->|metrics| prometheus
    otel -->|export| prometheus
    prometheus -->|query| grafana
    otel -.->|query| grafana

    %% Styles
    classDef ingest       fill:#2D0A12,stroke:#E94560,color:#FECDD3
    classDef processing   fill:#2D1A00,stroke:#F59E0B,color:#FDE68A
    classDef delivery     fill:#0A1628,stroke:#3B82F6,color:#BFDBFE
    classDef adstack      fill:#1A0F2D,stroke:#8B5CF6,color:#DDD6FE
    classDef observability fill:#042F1E,stroke:#10B981,color:#A7F3D0

    class broadcaster,nginx ingest
    class ffmpeg,segments processing
    class goserver,shaka,viewers delivery
    class vastserver,adcreative adstack
    class otel,prometheus,grafana observability
```
## Local Setup
### Start Local Stream Ingest Server
```bash
docker compose up -d ingest-server
```
This would start the ingest server on port 1935 (RTMP).
### Start Local Stream
1. Have `ffmpeg` installed on your system.
2. Give execute permissions to the `start-stream.sh` script:
```bash
chmod +x scripts/start-stream.sh
```
3. Then run the following script to start the stream:
```bash
./scripts/start-stream.sh
```
3. When the script asks for `Stream server address: `, enter the RTMP server URL (e.g., `rtmp://localhost/live` or `rtmp://localhost:1935/live`).
4. When the script asks for `Stream key: `, enter the stream key (e.g., `test`).
5. The script will then start streaming to the specified server and key.
> WARNING: This script will start your system's default camera and microphone and stream their feed
