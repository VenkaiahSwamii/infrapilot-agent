# Public Cloud Deployment Guide

Deploy InfraPilot Enterprise to public cloud providers for production use.

## Table of Contents

1. [DigitalOcean](#digitalocean)
2. [AWS EC2](#aws-ec2)
3. [Google Cloud Platform](#google-cloud-platform)
4. [Azure](#azure)
5. [Vercel + Railway](#vercel--railway)
6. [Security Considerations](#security-considerations)

---

## DigitalOcean

Easiest and most cost-effective for small to medium deployments.

### Prerequisites

- DigitalOcean account
- Domain name (optional but recommended)
- SSH key configured

### Step 1: Create Droplet

```bash
# Recommended specs:
# - 4 vCPUs, 8GB RAM, 160GB SSD
# - Ubuntu 22.04 LTS
# - Regions: NYC, SFO, or LON (closest to users)

# Create via CLI
doctl compute droplet create infrapilot \
  --size s-2vcpu-4gb \
  --image ubuntu-22-04-x64 \
  --region nyc1 \
  --ssh-keys <your-ssh-key-id> \
  --wait
```

### Step 2: Configure DNS (Optional)

```bash
# Point domain to droplet IP
doctl compute domain create infrapilot.example.com
doctl compute domain records create infrapilot.example.com \
  --record-type A \
  --record-data <droplet-ip>
```

### Step 3: Deploy Application

```bash
# SSH into droplet
ssh root@<droplet-ip>

# Clone repository
git clone https://github.com/venkaiswami/infrapilot-enterprise.git
cd infrapilot-enterprise

# Run deployment script
chmod +x scripts/deploy-cloud.sh
./scripts/deploy-cloud.sh infrapilot.example.com

# Or manual deployment:
docker-compose up -d
```

### Step 4: Configure SSL with Let's Encrypt

```bash
# Install Certbot
apt install certbot

# Obtain certificate
certbot certbot certonly --standalone -d infrapilot.example.com

# Configure nginx with SSL
# See nginx/ssl.conf
```

### Cost Estimate

| Resource | Monthly Cost |
|----------|--------------|
| Droplet (4 vCPU, 8GB) | $48 |
| Load Balancer | $12 |
| Floating IP | $0 (included) |
| **Total** | **~$60/month** |

---

## AWS EC2

Best for enterprise with multi-region requirements.

### Prerequisites

- AWS account with billing enabled
- AWS CLI configured
- Domain name in Route53 (optional)

### Step 1: Launch EC2 Instance

```bash
# Create security group
aws ec2 create-security-group \
  --group-name infrapilot-sg \
  --description "InfraPilot security group" \
  --vpc-id vpc-xxxxx

# Allow ports: 22 (SSH), 80 (HTTP), 443 (HTTPS), 8080 (API)
aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxx \
  --protocol tcp \
  --port 22 \
  --cidr 0.0.0.0/0

aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxx \
  --protocol tcp \
  --port 443 \
  --cidr 0.0.0.0/0

# Launch instance
aws ec2 run-instances \
  --image-id ami-0c7217cdde317cfec \
  --instance-type t3.xlarge \
  --key-name infrapilot-key \
  --security-group-ids sg-xxxxx \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=infrapilot}]'
```

### Step 2: Configure RDS (PostgreSQL)

```bash
# Create RDS instance
aws rds create-db-instance \
  --db-instance-identifier infrapilot-db \
  --db-instance-class db.t3.micro \
  --engine postgres \
  --engine-version 15.4 \
  --master-username infrapilot \
  --master-user-password <secure-password> \
  --allocated-storage 20 \
  --storage-type gp3 \
  --backup-retention-period 7
```

### Step 3: Configure ElastiCache (Redis)

```bash
# Create Redis cluster
aws elasticache create-cache-cluster \
  --cache-cluster-id infrapilot-redis \
  --cache-node-type cache.t3.micro \
  --engine redis \
  --num-cache-nodes 1
```

### Step 4: Deploy Application

```bash
# SSH into instance
ssh -i infrapilot-key.pem ec2-user@<instance-public-ip>

# Clone and deploy
git clone https://github.com/venkaiswami/infrapilot-enterprise.git
cd infrapilot-enterprise

# Configure environment
cp backend/.env.example backend/.env
# Update with RDS and ElastiCache endpoints

# Start with Docker Compose
docker-compose -f docker-compose.prod.yml up -d
```

### Step 5: Configure ALB with SSL

```bash
# Create Application Load Balancer
aws elbv2 create-load-balancer \
  --name infrapilot-alb \
  --subnets subnet-xxxxx subnet-yyyyy \
  --security-groups sg-xxxxx

# Create target group
aws elbv2 create-target-group \
  --name infrapilot-tg \
  --protocol HTTP \
  --port 8080 \
  --vpc-id vpc-xxxxx \
  --target-type instance

# Register targets
aws elbv2 register-targets \
  --target-group-arn <tg-arn> \
  --targets Id=i-xxxxx

# Create HTTPS listener with ACM certificate
aws elbv2 create-listener \
  --load-balancer-arn <alb-arn> \
  --protocol HTTPS \
  --port 443 \
  --certificates CertificateArn=<acm-cert-arn> \
  --default-actions Type=forward,TargetGroupArn=<tg-arn>
```

### Cost Estimate (Monthly)

| Resource | Monthly Cost |
|----------|--------------|
| EC2 t3.xlarge | $120 |
| RDS db.t3.micro | $15 |
| ElastiCache t3.micro | $10 |
| ALB | $25 |
| Data Transfer | $10 |
| **Total** | **~$180/month** |

---

## Google Cloud Platform (GCP)

Best for Kubernetes-native deployments.

### Step 1: Create GKE Cluster

```bash
# Create cluster
gcloud container clusters create infrapilot \
  --zone us-central1-a \
  --num-nodes 3 \
  --machine-type e2-standard-4 \
  --disk-size 100GB \
  --enable-autoscaling --min-nodes 2 --max-nodes 10

# Get credentials
gcloud container clusters get-credentials infrapilot --zone us-central1-a
```

### Step 2: Create Cloud SQL (PostgreSQL)

```bash
# Create PostgreSQL instance
gcloud sql instances create infrapilot-db \
  --database-version=POSTGRES_15 \
  --tier=db-f1-micro \
  --region=us-central1 \
  --storage-size=20GB \
  --storage-type=SSD

# Create database
gcloud sql databases create infrapilot --instance=infrapilot-db

# Create user
gcloud sql users create infrapilot \
  --instance=infrapilot-db \
  --password=<secure-password>
```

### Step 3: Create Memorystore (Redis)

```bash
gcloud redis instances create infrapilot-redis \
  --size=1 \
  --region=us-central1 \
  --redis-version=redis_7_0
```

### Step 4: Deploy to GKE

```yaml
# k8s/gcp-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: infrapilot-config
data:
  DATABASE_URL: "cloudsql://<project-id>:us-central1:infrapilot-db/infrapilot"
  REDIS_URL: "redis://10.0.0.3:6379"
```

```bash
# Deploy application
kubectl apply -f k8s/gcp-configmap.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

### Step 5: Configure Cloud CDN and SSL

```bash
# Create static IP
gcloud compute addresses create infrapilot-ip --global

# Create SSL certificate
gcloud compute ssl-certificates create infrapilot-cert \
  --domains infrapilot.example.com

# Create HTTPS load balancer (via GKE Ingress)
kubectl apply -f k8s/gcp-ingress.yaml
```

### Cost Estimate (Monthly)

| Resource | Monthly Cost |
|----------|--------------|
| GKE cluster (3 nodes) | $120 |
| Cloud SQL db-f1-micro | $15 |
| Memorystore 1GB | $35 |
| Load Balancer | $18 |
| Cloud CDN | $10 |
| **Total** | **~$198/month** |

---

## Azure

Best for enterprise Microsoft shops.

### Step 1: Create AKS Cluster

```bash
# Create resource group
az group create --name infrapilot-rg --location eastus

# Create AKS cluster
az aks create \
  --resource-group infrapilot-rg \
  --name infrapilot \
  --node-count 3 \
  --node-vm-size Standard_D4s_v3 \
  --enable-addons monitoring \
  --generate-ssh-keys

# Get credentials
az aks get-credentials --resource-group infrapilot-rg --name infrapilot
```

### Step 2: Create Azure Database (PostgreSQL)

```bash
# Create PostgreSQL server
az postgres server create \
  --name infrapilot-db \
  --resource-group infrapilot-rg \
  --location eastus \
  --admin-user infrapilot \
  --admin-password <secure-password> \
  --sku-name B_Gen5_2 \
  --storage-size 51200

# Create database
az postgres db create \
  --name infrapilot \
  --server-name infrapilot-db \
  --resource-group infrapilot-rg
```

### Step 3: Create Azure Cache (Redis)

```bash
az redis create \
  --name infrapilot-redis \
  --resource-group infrapilot-rg \
  --location eastus \
  --sku Basic \
  --vm-size C1
```

### Step 4: Deploy Application

```bash
# Update secrets in k8s/secret.yaml with Azure credentials
kubectl apply -f k8s/azure-secret.yaml
kubectl apply -f k8s/deployment.yaml
```

### Cost Estimate (Monthly)

| Resource | Monthly Cost |
|----------|--------------|
| AKS cluster (3 nodes) | $120 |
| Azure Database (Basic) | $15 |
| Azure Cache (Basic) | $15 |
| Application Gateway | $25 |
| **Total** | **~$175/month** |

---

## Vercel + Railway

Best for frontend + backend separation with managed services.

### Step 1: Deploy Frontend to Vercel

```bash
# Install Vercel CLI
npm i -g vercel

# Deploy frontend
cd frontend
vercel --prod

# Configure environment variables in Vercel dashboard:
# - VITE_API_URL=https://infrapilot-api.example.com
# - VITE_WS_URL=wss://infrapilot-api.example.com/ws
```

### Step 2: Deploy Backend to Railway

```bash
# Install Railway CLI
npm i -g @railway/cli

# Login and initialize
cd backend
railway login
railway init

# Deploy
railway up

# Add PostgreSQL plugin
railway add postgresql

# Add Redis plugin
railway add redis
```

### Step 3: Deploy Agent

```bash
# Agent runs on monitored machines, not cloud
# Download from releases page
curl -LO https://github.com/venkaiswami/infrapilot-enterprise/releases/download/v1.0.0/infrapilot-agent-linux-amd64
chmod +x infrapilot-agent-linux-amd64

# Run with Railway API URL
./infrapilot-agent-linux-amd64 --server https://infrapilot-api.railway.app --token <token>
```

### Cost Estimate (Monthly)

| Resource | Monthly Cost |
|----------|--------------|
| Vercel Pro | $20 |
| Railway Backend | $20 |
| PostgreSQL (Railway) | $9 |
| Redis (Railway) | $5 |
| **Total** | **~$54/month** |

---

## Security Considerations

### 1. Firewall Configuration

```bash
# UFW (Ubuntu)
ufw enable
ufw allow 22/tcp      # SSH
ufw allow 80/tcp      # HTTP
ufw allow 443/tcp     # HTTPS
ufw deny 8080/tcp     # Block direct API access (use reverse proxy)
ufw deny 6379/tcp     # Block Redis
ufw deny 5432/tcp     # Block PostgreSQL
```

### 2. SSL/TLS Configuration

```nginx
# nginx/ssl.conf
server {
    listen 443 ssl http2;
    server_name infrapilot.example.com;

    ssl_certificate /etc/letsencrypt/live/infrapilot.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/infrapilot.example.com/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    location / {
        proxy_pass http://localhost:5173;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /ws/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
    }
}
```

### 3. SSH Hardening

```bash
# /etc/ssh/sshd_config
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
PermitEmptyPasswords no
ChallengeResponseAuthentication no
UsePAM yes
X11Forwarding no
Protocol 2

# Allow only specific users
AllowUsers deploy

# Restart SSH
systemctl restart sshd
```

### 4. Fail2ban Configuration

```bash
# Install Fail2ban
apt install fail2ban

# Create jail for SSH
cat > /etc/fail2ban/jail.local << EOF
[sshd]
enabled = true
port = ssh
logpath = /var/log/auth.log
maxretry = 3
bantime = 86400
findtime = 600
EOF

# Start Fail2ban
systemctl enable fail2ban
systemctl start fail2ban
```

### 5. Automatic Security Updates

```bash
# Install unattended-upgrades
apt install unattended-upgrades

# Enable automatic updates
dpkg-reconfigure -plow unattended-upgrades
```

### 6. Monitoring & Alerts

```bash
# Set up CloudWatch (AWS), Stackdriver (GCP), or Azure Monitor
# Configure alerts for:
# - High CPU/Memory usage
# - Disk space < 20%
# - Service downtime
# - Failed login attempts
```

### 7. Backup Strategy

```bash
# Automated daily backups
cat > /etc/cron.daily/infrapilot-backup << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR=/backup/infrapilot
DB_NAME=infrapilot_enterprise
DB_USER=postgres

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup database
docker exec infrapilot-postgres pg_dump -U $DB_USER $DB_NAME | gzip > $BACKUP_DIR/db_$DATE.sql.gz

# Backup uploads
tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz -C /data/infrapilot uploads/

# Delete backups older than 30 days
find $BACKUP_DIR -name "*.gz" -mtime +30 -delete

# Upload to S3 (optional)
# aws s3 sync $BACKUP_DIR s3://infrapilot-backups/
EOF

chmod +x /etc/cron.daily/infrapilot-backup
```

---

## Troubleshooting

### Issue: Services won't start

```bash
# Check logs
docker-compose logs -f
docker-compose logs -f backend

# Verify environment variables
docker-compose config

# Check disk space
df -h

# Check memory
free -m
```

### Issue: Database connection refused

```bash
# Verify PostgreSQL is running
docker-compose ps postgres

# Check credentials
docker-compose exec postgres psql -U postgres -c "SELECT 1"

# Verify network connectivity
docker-compose exec backend ping postgres
```

### Issue: WebSocket connection fails

```bash
# Check nginx proxy configuration
nginx -t

# Verify proxy headers
curl -I https://infrapilot.example.com/ws/

# Check firewall
ufw status
```

### Issue: High CPU/Memory usage

```bash
# Check running processes
docker stats

# Scale backend
docker-compose up -d --scale backend=3

# Optimize database queries
docker-compose exec postgres psql -U postgres -d infrapilot -c "SELECT * FROM pg_stat_activity;"
```

---

## Performance Optimization

### 1. Enable Gzip Compression

```bash
# nginx/nginx.conf
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css text/xml text/javascript application/javascript application/json application/xml+rss application/rss+xml font/truetype font/opentype application/vnd.ms-fontobject image/svg+xml;
```

### 2. Configure Redis Persistence

```conf
# redis.conf
save 900 1
save 300 10
save 60 10000
appendonly yes
appendfsync everysec
```

### 3. Database Connection Pooling

```go
// backend/internal/database/database.go
import "github.com/jackc/pgx/v5/pgxpool"

config.MaxConns = 50
config.MinConns = 10
config.MaxConnLifetime = time.Hour
config.MaxConnIdleTime = 30 * time.Minute
```

### 4. CDN for Static Assets

```bash
# CloudFlare setup:
# 1. Add site to CloudFlare
# 2. Update nameservers
# 3. Enable CDN, WAF, DDoS protection
# 4. Create page rules for caching
```

---

## Next Steps

1. **Configure monitoring** with Prometheus and Grafana
2. **Set up alerting** for critical issues
3. **Configure backups** and test restoration
4. **Implement CI/CD** for automated deployments
5. **Deploy to staging** environment first
6. **Load test** before production traffic
7. **Document runbooks** for common issues

---

## Support

- **Issues:** [GitHub Issues](https://github.com/venkaiswami/infrapilot-enterprise/issues)
- **Email:** venkaiswami@pm.me
- **Docs:** [Deployment Guide](docs/deployment.md)

---

<p align="center">
  Made with ❤️ by the InfraPilot Team
</p>