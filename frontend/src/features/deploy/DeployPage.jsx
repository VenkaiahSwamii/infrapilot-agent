import { useState, useEffect } from 'react';
import { listGitRepositories, listPipelines, listBuilds, listDeployments, getDeployStats } from '../../api/deploy';
import './DeployPage.css';

export default function DeployPage() {
  const [activeTab, setActiveTab] = useState('overview');
  const [repositories, setRepositories] = useState([]);
  const [pipelines, setPipelines] = useState([]);
  const [builds, setBuilds] = useState([]);
  const [deployments, setDeployments] = useState([]);
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchData();
  }, [activeTab]);

  const fetchData = async () => {
    setLoading(true);
    try {
      switch (activeTab) {
        case 'repositories':
          const repos = await listGitRepositories();
          setRepositories(repos);
          break;
        case 'pipelines':
          const pipes = await listPipelines();
          setPipelines(pipes);
          break;
        case 'builds':
          const buildsData = await listBuilds();
          setBuilds(buildsData);
          break;
        case 'deployments':
          const deploys = await listDeployments();
          setDeployments(deploys);
          break;
        case 'overview':
          const statsData = await getDeployStats();
          setStats(statsData);
          const recentBuilds = await listBuilds('', 5);
          setBuilds(recentBuilds);
          const recentDeploys = await listDeployments('', 5);
          setDeployments(recentDeploys);
          break;
        default:
          break;
      }
    } catch (error) {
      console.error('Failed to fetch data:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status) => {
    const statusMap = {
      pending: 'badge-pending',
      running: 'badge-running',
      success: 'badge-success',
      failed: 'badge-failed',
      cancelled: 'badge-cancelled',
      deploying: 'badge-running',
      'rolled-back': 'badge-rolled-back',
    };
    return statusMap[status] || 'badge-default';
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleString();
  };

  return (
    <div className="deploy-page">
      <div className="deploy-header">
        <h1>InfraDeploy - CI/CD Platform</h1>
        <p className="deploy-subtitle">Enterprise-grade continuous integration and deployment</p>
      </div>

      <div className="deploy-tabs">
        <button
          className={activeTab === 'overview' ? 'tab-button active' : 'tab-button'}
          onClick={() => setActiveTab('overview')}
        >
          Overview
        </button>
        <button
          className={activeTab === 'repositories' ? 'tab-button active' : 'tab-button'}
          onClick={() => setActiveTab('repositories')}
        >
          Git Repositories
        </button>
        <button
          className={activeTab === 'pipelines' ? 'tab-button active' : 'tab-button'}
          onClick={() => setActiveTab('pipelines')}
        >
          Pipelines
        </button>
        <button
          className={activeTab === 'builds' ? 'tab-button active' : 'tab-button'}
          onClick={() => setActiveTab('builds')}
        >
          Builds
        </button>
        <button
          className={activeTab === 'deployments' ? 'tab-button active' : 'tab-button'}
          onClick={() => setActiveTab('deployments')}
        >
          Deployments
        </button>
      </div>

      <div className="deploy-content">
        {loading && <div className="loading">Loading...</div>}

        {!loading && activeTab === 'overview' && (
          <div className="overview-section">
            {stats && (
              <div className="stats-grid">
                <div className="stat-card">
                  <div className="stat-value">{stats.pipelines || 0}</div>
                  <div className="stat-label">Pipelines</div>
                </div>
                <div className="stat-card">
                  <div className="stat-value text-success">{stats.builds_success || 0}</div>
                  <div className="stat-label">Successful Builds</div>
                </div>
                <div className="stat-card">
                  <div className="stat-value text-danger">{stats.builds_failed || 0}</div>
                  <div className="stat-label">Failed Builds</div>
                </div>
                <div className="stat-card">
                  <div className="stat-value text-primary">{stats.deployments || 0}</div>
                  <div className="stat-label">Deployments</div>
                </div>
              </div>
            )}

            <div className="recent-section">
              <h3>Recent Builds</h3>
              {builds.length === 0 ? (
                <p className="no-data">No builds yet</p>
              ) : (
                <div className="table-container">
                  <table className="data-table">
                    <thead>
                      <tr>
                        <th>Build #</th>
                        <th>Pipeline</th>
                        <th>Status</th>
                        <th>Commit</th>
                        <th>Branch</th>
                        <th>Duration</th>
                        <th>Created</th>
                      </tr>
                    </thead>
                    <tbody>
                      {builds.map((build) => (
                        <tr key={build.id}>
                          <td>#{build.build_number}</td>
                          <td>{build.pipeline?.name || 'N/A'}</td>
                          <td>
                            <span className={`badge ${getStatusBadge(build.status)}`}>
                              {build.status}
                            </span>
                          </td>
                          <td className="font-mono text-sm">{build.commit_hash?.substring(0, 7)}</td>
                          <td>{build.branch}</td>
                          <td>{build.duration_ms ? `${build.duration_ms}ms` : 'N/A'}</td>
                          <td>{formatDate(build.created_at)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>

            <div className="recent-section">
              <h3>Recent Deployments</h3>
              {deployments.length === 0 ? (
                <p className="no-data">No deployments yet</p>
              ) : (
                <div className="table-container">
                  <table className="data-table">
                    <thead>
                      <tr>
                        <th>Name</th>
                        <th>Environment</th>
                        <th>Strategy</th>
                        <th>Image</th>
                        <th>Status</th>
                        <th>Deployed At</th>
                      </tr>
                    </thead>
                    <tbody>
                      {deployments.map((deploy) => (
                        <tr key={deploy.id}>
                          <td>{deploy.name}</td>
                          <td>
                            <span className={`badge badge-${deploy.environment}`}>
                              {deploy.environment}
                            </span>
                          </td>
                          <td>{deploy.strategy}</td>
                          <td className="font-mono text-sm">{deploy.image?.split(':').pop()}</td>
                          <td>
                            <span className={`badge ${getStatusBadge(deploy.status)}`}>
                              {deploy.status}
                            </span>
                          </td>
                          <td>{formatDate(deploy.deployed_at)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          </div>
        )}

        {!loading && activeTab === 'repositories' && (
          <div className="repositories-section">
            <div className="section-header">
              <h3>Git Repositories</h3>
              <button className="btn btn-primary">+ Add Repository</button>
            </div>
            {repositories.length === 0 ? (
              <p className="no-data">No repositories configured</p>
            ) : (
              <div className="repositories-grid">
                {repositories.map((repo) => (
                  <div key={repo.id} className="repo-card">
                    <div className="repo-header">
                      <h4>{repo.name}</h4>
                      <span className={`badge badge-${repo.provider}`}>{repo.provider}</span>
                    </div>
                    <div className="repo-info">
                      <p className="repo-url">{repo.url}</p>
                      <p className="repo-branch">Branch: {repo.branch}</p>
                      <p className="repo-auth">Auth: {repo.auth_type}</p>
                    </div>
                    <div className="repo-actions">
                      <button className="btn btn-sm btn-secondary">Configure</button>
                      <button className="btn btn-sm btn-danger">Delete</button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {!loading && activeTab === 'pipelines' && (
          <div className="pipelines-section">
            <div className="section-header">
              <h3>Pipelines</h3>
              <button className="btn btn-primary">+ Create Pipeline</button>
            </div>
            {pipelines.length === 0 ? (
              <p className="no-data">No pipelines configured</p>
            ) : (
              <div className="table-container">
                <table className="data-table">
                  <thead>
                    <tr>
                      <th>Name</th>
                      <th>Environment</th>
                      <th>Trigger</th>
                      <th>Repository</th>
                      <th>Status</th>
                      <th>Last Run</th>
                    </tr>
                  </thead>
                  <tbody>
                    {pipelines.map((pipeline) => (
                      <tr key={pipeline.id}>
                        <td>{pipeline.name}</td>
                        <td>
                          <span className={`badge badge-${pipeline.environment}`}>
                            {pipeline.environment}
                          </span>
                        </td>
                        <td>{pipeline.trigger_type}</td>
                        <td>{pipeline.repository?.name || 'N/A'}</td>
                        <td>
                          <span className={`badge ${pipeline.enabled ? 'badge-success' : 'badge-default'}`}>
                            {pipeline.enabled ? 'Enabled' : 'Disabled'}
                          </span>
                        </td>
                        <td>{formatDate(pipeline.last_run_at)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {!loading && activeTab === 'builds' && (
          <div className="builds-section">
            <div className="section-header">
              <h3>Builds</h3>
              <button className="btn btn-primary">+ Trigger Build</button>
            </div>
            {builds.length === 0 ? (
              <p className="no-data">No builds yet</p>
            ) : (
              <div className="table-container">
                <table className="data-table">
                  <thead>
                    <tr>
                      <th>Build #</th>
                      <th>Pipeline</th>
                      <th>Status</th>
                      <th>Commit</th>
                      <th>Branch</th>
                      <th>Author</th>
                      <th>Duration</th>
                      <th>Created</th>
                    </tr>
                  </thead>
                  <tbody>
                    {builds.map((build) => (
                      <tr key={build.id}>
                        <td>#{build.build_number}</td>
                        <td>{build.pipeline?.name || 'N/A'}</td>
                        <td>
                          <span className={`badge ${getStatusBadge(build.status)}`}>
                            {build.status}
                          </span>
                        </td>
                        <td className="font-mono text-sm">{build.commit_hash?.substring(0, 7)}</td>
                        <td>{build.branch}</td>
                        <td>{build.author}</td>
                        <td>{build.duration_ms ? `${build.duration_ms}ms` : 'N/A'}</td>
                        <td>{formatDate(build.created_at)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {!loading && activeTab === 'deployments' && (
          <div className="deployments-section">
            <div className="section-header">
              <h3>Deployments</h3>
              <button className="btn btn-primary">+ New Deployment</button>
            </div>
            {deployments.length === 0 ? (
              <p className="no-data">No deployments yet</p>
            ) : (
              <div className="table-container">
                <table className="data-table">
                  <thead>
                    <tr>
                      <th>Name</th>
                      <th>Environment</th>
                      <th>Strategy</th>
                      <th>Image</th>
                      <th>Replicas</th>
                      <th>Status</th>
                      <th>Deployed At</th>
                      <th>Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {deployments.map((deploy) => (
                      <tr key={deploy.id}>
                        <td>{deploy.name}</td>
                        <td>
                          <span className={`badge badge-${deploy.environment}`}>
                            {deploy.environment}
                          </span>
                        </td>
                        <td>{deploy.strategy}</td>
                        <td className="font-mono text-sm">{deploy.image}</td>
                        <td>{deploy.replicas}</td>
                        <td>
                          <span className={`badge ${getStatusBadge(deploy.status)}`}>
                            {deploy.status}
                          </span>
                        </td>
                        <td>{formatDate(deploy.deployed_at)}</td>
                        <td>
                          <button className="btn btn-sm btn-secondary">Rollback</button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}