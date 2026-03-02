# Omada Controller Prometheus Exporter

Prometheus exporter for TP-Link Omada Controller WiFi station metrics. Fork of [jamessanford/omada-controller-exporter](https://github.com/jamessanford/omada-controller-exporter) with Omada Controller v6 compatibility.

## Controller v6 Compatibility

The upstream exporter ([jamessanford/omada-controller-exporter](https://github.com/jamessanford/omada-controller-exporter)) targets Controller v5 which uses the site **name** or **key** in API paths. Controller v6 returns API version 3 and requires the site **ID** instead. This fork adds a `SiteID()` method that prefers `id` over `key` over `name`, making the exporter compatible with both Controller v5 and v6.

## Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `omada_station_wireless` | gauge | Is wireless station (1/0) |
| `omada_station_channel` | gauge | Wireless channel |
| `omada_station_wireless_mode` | gauge | Wireless mode |
| `omada_station_signal_rssi` | gauge | Signal strength (dBm) |
| `omada_station_signal_level_pct` | gauge | Signal level (%) |
| `omada_station_powersaving` | gauge | Power-save mode (1/0) |
| `omada_station_transmit_bitrate_mbps` | gauge | Transmit bitrate (Mbps) |
| `omada_station_receive_bitrate_mbps` | gauge | Receive bitrate (Mbps) |
| `omada_station_transmit_bytes_total` | counter | Bytes transmitted |
| `omada_station_receive_bytes_total` | counter | Bytes received |
| `omada_station_transmit_packets_total` | counter | Packets transmitted |
| `omada_station_receive_packets_total` | counter | Packets received |
| `omada_station_last_seen_time_seconds` | counter | Last seen (unix epoch) |
| `omada_station_uptime_seconds` | counter | Station uptime |
| `omada_collector_requests_total` | counter | Total collection requests |
| `omada_collector_errors_total` | counter | Total collection errors |
| `omada_collector_duration_seconds` | gauge | Collection duration |

Labels on station metrics: `station` (MAC), `name`, `ap_address`, `ap_name`, `ssid`, `site`.

## Configuration

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `OMADA_URL` | `https://10.0.0.3:30077` | no | Omada Controller URL |
| `OMADA_USER` | — | yes | Controller username |
| `OMADA_PASS` | — | yes | Controller password |
| `OMADA_INSECURE` | `true` | no | Skip TLS verification |
| `LISTEN_PORT` | `6779` | no | Exporter listen port |
| `LOG_LEVEL` | `info` | no | Log level (debug, info, warn, error) |

The Omada user must have **Admin** role with **All Sites** site privileges.

## Run locally

Requires [1Password CLI](https://developer.1password.com/docs/cli/) for credential injection:

```
make run
```

Then verify:

```
make query
```

## Deploy

```yaml
# Prometheus scrape config
scrape_configs:
  - job_name: omada
    static_configs:
      - targets: ['10.0.0.3:6779']
```

## Tested with

- Omada Controller **v6.1.0.19**
- TP-Link **EAP610** access point

## License

MIT — see [LICENSE](LICENSE). Based on [jamessanford/omada-controller-exporter](https://github.com/jamessanford/omada-controller-exporter).
