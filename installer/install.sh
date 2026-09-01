#!/bin/bash
# InfraPilot Enterprise Installation Script

set -euo pipefail

INFRAPILOT_VERSION="1.0.0"
INFRAPILOT_USER="infrapilot"
INFRAPILOT_HOME="/opt/infrapilot"

log_info() { echo -e "\033[0;32m[INFO]\033[0m $1"; }
log_error() { echo -e "\033[0;31m[ERROR]\033[0m $1"; }

# Check root
if [ "$EUID" -ne 0 ]; then
    log_error "Please run as root"
    exit 1
fi

# Install dependencies
log_info "Installing dependencies..."
apt-get update
apt-get install -y ca-certificates curl git ufw

# Install Docker
if ! command -v docker &> /dev/null; then
    log_info "Installing Docker..."
    curl -fsSL https://get.docker.com | sh
    systemctl enable --now docker
fi

# Install Docker Compose
if ! command -v docker-compose &> /dev/null; then
    log_info "Installing Docker Compose..."
    curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" \
        -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
fi

# Create user
if ! id "$INFRAPILOT_USER" &>/dev/null; then
    log_info "Creating user..."
    useradd -r -s /bin/false -m -d $INFRAPILOT_HOME "$INFRAPILOT_USER"
    usermod -aG docker "$INFRAPILOT_USER"
fi

# Download InfraPilot
log_info "Downloading InfraPilot..."
git clone https://github.com/yourusername/infrapilot-enterprise.git /tmp/infrapilot
cp -r /tmp/infrapilot/* $INFRAPILOT_HOME/
rm -rf /tmp/infrapilot
chown -R "$INFRAPILOT_USER:$INFRAPILOT_USER" $INFRAPILOT_HOME

# Configure environment
log_info "Configuring environment..."
DB_PASSWORD=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 64)
GRAFANA_PASSWORD=$(openssl rand -base64 32)

cat > $INFRAPILOT_HOME/backend/.env << EOF
DB_PASSWORD=$DB_PASSWORD
JWT_SECRET=$JWT_SECRET
GRAFANA_PASSWORD=$GRAFANA_PASSWORD
EOF

chmod 600 $INFRAPILOT_HOME/backend/.env

# Configure firewall
log_info "Configuring firewall..."
ufw --force disable
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
echo "y" | ufw enable

# Start services
log_info "Starting services..."
cd $INFRAPILOT_HOME
docker-compose pull
docker-compose up -d

# Create systemd service
cat > /etc/systemd/system/infrapilot.service << EOF
[Unit]
Description=InfraPilot Enterprise
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=$INFRAPILOT_HOME
ExecStart=/usr/local/bin/docker-compose up -d
ExecStop=/usr/local/bin/docker-compose down
User=$INFRAPILOT_USER

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable infrapilot

log_info "Installation complete!"
echo "Dashboard: http://$(hostname -I | awk '{print $1}')"
echo "Default login: admin@infrapilot.io / admin123"