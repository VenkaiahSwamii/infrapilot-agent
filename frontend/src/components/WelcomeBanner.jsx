import React, { useState } from 'react';
import { 
  Rocket, 
  Monitor, 
  Shield, 
  Zap, 
  BookOpen, 
  Github,
  ExternalLink,
  ChevronRight,
  Play,
  X
} from 'lucide-react';
import './WelcomeBanner.css';

const WelcomeBanner = ({ onDismiss }) => {
  const [isVisible, setIsVisible] = useState(true);

  if (!isVisible) return null;

  const handleDismiss = () => {
    setIsVisible(false);
    if (onDismiss) onDismiss();
  };

  const features = [
    {
      icon: <Monitor size={24} />,
      title: 'Real-Time Monitoring',
      description: 'Live metrics, alerts, and remote management'
    },
    {
      icon: <Shield size={24} />,
      title: 'Enterprise Security',
      description: 'RBAC, JWT auth, rate limiting, audit logs'
    },
    {
      icon: <Zap size={24} />,
      title: 'AI-Powered Insights',
      description: 'Anomaly detection and predictive analytics'
    },
    {
      icon: <Rocket size={24} />,
      title: 'Production Ready',
      description: 'Docker, K8s, Prometheus, Grafana included'
    }
  ];

  return (
    <div className="welcome-banner">
      <button className="welcome-banner__close" onClick={handleDismiss}>
        <X size={20} />
      </button>

      <div className="welcome-banner__content">
        <div className="welcome-banner__header">
          <div className="welcome-banner__logo">
            <Rocket size={48} />
          </div>
          <div>
            <h1 className="welcome-banner__title">
              Welcome to InfraPilot Enterprise
            </h1>
            <p className="welcome-banner__subtitle">
              Production-ready infrastructure monitoring & management platform
            </p>
          </div>
        </div>

        <div className="welcome-banner__features">
          {features.map((feature, index) => (
            <div key={index} className="welcome-banner__feature">
              <div className="welcome-banner__feature-icon">
                {feature.icon}
              </div>
              <div>
                <h3>{feature.title}</h3>
                <p>{feature.description}</p>
              </div>
            </div>
          ))}
        </div>

        <div className="welcome-banner__actions">
          <a 
            href="https://github.com/venkaiswami/infrapilot-enterprise" 
            className="welcome-banner__button welcome-banner__button--primary"
            target="_blank"
            rel="noopener noreferrer"
          >
            <Github size={18} />
            Star on GitHub
            <ExternalLink size={14} />
          </a>
          
          <a 
            href="https://github.com/venkaiswami/infrapilot-enterprise/blob/main/docs/installation.md" 
            className="welcome-banner__button welcome-banner__button--secondary"
            target="_blank"
            rel="noopener noreferrer"
          >
            <BookOpen size={18} />
            Read Documentation
            <ExternalLink size={14} />
          </a>

          <button className="welcome-banner__button welcome-banner__button--outline">
            <Play size={18} />
            Launch Demo
            <ChevronRight size={14} />
          </button>
        </div>

        <div className="welcome-banner__quick-start">
          <h3>Quick Start</h3>
          <div className="welcome-banner__code">
            <code>
              {`# Clone and start\n`}
              {`git clone https://github.com/venkaiswami/infrapilot-enterprise.git\n`}
              {`cd infrapilot-enterprise\n`}
              {`docker-compose up -d`}
            </code>
          </div>
        </div>

        <div className="welcome-banner__footer">
          <p>
            Need help? Check out our{' '}
            <a 
              href="https://github.com/venkaiswami/infrapilot-enterprise/discussions"
              target="_blank"
              rel="noopener noreferrer"
            >
              Discussions forum
            </a>
            {' '}or{' '}
            <a 
              href="https://github.com/venkaiswami/infrapilot-enterprise/issues"
              target="_blank"
              rel="noopener noreferrer"
            >
              report an issue
            </a>
            .
          </p>
        </div>
      </div>
    </div>
  );
};

export default WelcomeBanner;