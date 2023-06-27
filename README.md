# Prometheus exporter for gym utilization

This is a prometheus exporter, that exports the current utilization of you local gym in percent. At the moment only McFit is supported.

## How to run:
Before you can use the exporter, you have to find out the studio code of you local McFIT. One way to do it, is to login to you McFIT account and open the developer tools of your browser. Then search for "/nox/public/v1/studios/" in the network tab and you should see a GET request with your studio code.

### Prerequisites
- Docker
- Docker Compose
- Prometheus
- Grafana

### Docker
1. Open a terminal or command prompt in the same directory where the Dockerfile is located.
2. Build the Docker image using the following command:

```bash
docker build -t gym_util_exporter .
```

3. Run the Docker container, specifying the port you want to use. For example, if you want to use port 8080, execute the following command:

```
docker run -p 8080:2112 gym_util_exporter
```

4. Now you can access the metrics under ```http://localhost:8080/metrics```

Alternatively you can use the following command:

```bash
docker build -t gym_util_exporter . && docker run -p 8080:2112 --restart always gym_util_exporter
```

### Docker Compose
If you want to use Docker Compose, you can execute the following command:
```bash
docker-compose up -d
```

This will start the container in detached mode.

### Add exporter to Prometheus
Open you Prometheus config file `/etc/prometheus/prometheus.yaml` and add a new job. It should look similar to this:
```
scrape_configs:
- job_name: prometheus
  honor_timestamps: true
  scrape_interval: 5s
  scrape_timeout: 5s
  metrics_path: /metrics
  scheme: http
  static_configs:
  - targets:
    - localhost:9090
- job_name: mcfit
  honor_timestamps: true
  scrape_interval: 15s
  scrape_timeout: 10s
  metrics_path: /metrics
  scheme: http
  static_configs:
  - targets:
    - localhost:2112
```

### Import dashboard to Grafana
Perform the following steps in Grafana, to import the dashboard:
1. Click on "Dashboards"
2. Click on "New"
3. Click on "Import"
4. Drag the "Grafana_dashboard.json" file or paste the JSON into the textbox
5. Click on "Load"