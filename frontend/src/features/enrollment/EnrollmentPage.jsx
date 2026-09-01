import React, { useState } from 'react';

export default function EnrollmentPage() {
  const [serverUrl, setServerUrl] = useState(`${window.location.protocol}//${window.location.host}`);
  const [enrollmentToken, setEnrollmentToken] = useState('');
  const [status, setStatus] = useState('idle'); // idle | enrolling | success | error
  const [message, setMessage] = useState('');

  const handleEnroll = async () => {
    if (!serverUrl || !enrollmentToken) {
      setStatus('error');
      setMessage('Server URL and enrollment token are required');
      return;
    }

    setStatus('enrolling');
    setMessage('Enrolling agent with the backend...');

    try {
      // Step 1: Get system info from backend config endpoint
      const configResp = await fetch(`${serverUrl}/api/v1/agent/config`, {
        method: 'GET',
        headers: { 'Accept': 'application/json' },
      });

      if (!configResp.ok) {
        throw new Error(`Failed to fetch agent config: ${configResp.status}`);
      }

      const config = await configResp.json();

      // Step 2: Collect local system info
      const sysInfo = {
        enrollment_token: enrollmentToken,
        hostname: window.location.hostname || 'unknown',
        os: navigator.platform || 'Unknown',
        platform: 'linux',
        ip_address: '',
        kernel: '',
        architecture: navigator.userAgent.includes('x64') ? 'amd64' : 'unknown',
        cpu_model: '',
        memory: 0,
      };

      // Step 3: Enroll with the backend
      const enrollResp = await fetch(`${serverUrl}/api/v1/agent/enroll`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(sysInfo),
      });

      if (!enrollResp.ok) {
        const errData = await enrollResp.json().catch(() => ({}));
        throw new Error(errData.error || `Enrollment failed: ${enrollResp.status}`);
      }

      const data = await enrollResp.json();

      // Step 4: Build local config.json content
      const localConfig = {
        server: serverUrl,
        machine_id: data.machine_id,
        api_key: data.api_key,
        organization: data.organization_id || 'default',
        metrics_interval: config.metrics_interval_seconds || 5,
        heartbeat_interval: config.heartbeat_interval_seconds || 15,
      };

      // In a real Electron/Tauri app, you'd write this to disk.
      // For demo, we display the JSON that the agent should persist locally.
      console.log('Local config to save:', JSON.stringify(localConfig, null, 2));

      setStatus('success');
      setMessage(JSON.stringify(localConfig, null, 2));
    } catch (err) {
      setStatus('error');
      setMessage(err.message || 'Enrollment failed');
    }
  };

  return (
    <div className="page enrollment-page">
      <h1>Agent Enrollment</h1>
      <p className="subtitle">One-time secure enrollment for InfraPilot Enterprise Agent</p>

      <div className="enrollment-card">
        <div className="field">
          <label>Backend Server URL</label>
          <input
            type="text"
            value={serverUrl}
            onChange={(e) => setServerUrl(e.target.value)}
            placeholder="http://192.168.1.41:8080"
          />
        </div>

        <div className="field">
          <label>Enrollment Token</label>
          <input
            type="text"
            value={enrollmentToken}
            onChange={(e) => setEnrollmentToken(e.target.value)}
            placeholder="Paste token from dashboard"
          />
        </div>

        <button onClick={handleEnroll} disabled={status === 'enrolling'}>
          {status === 'enrolling' ? 'Enrolling...' : 'Enroll Agent'}
        </button>

        {status === 'success' && (
          <div className="success notice">
            <strong>Enrollment Successful</strong>
            <pre>{message}</pre>
            <p>Save this JSON into the agent's <code>config.json</code> to complete setup.</p>
          </div>
        )}

        {status === 'error' && (
          <div className="error notice">
            <strong>Enrollment Failed</strong>
            <p>{message}</p>
          </div>
        )}
      </div>

      <div className="steps">
        <h3>How it works</h3>
        <ol>
          <li>Generate an enrollment token from the InfraPilot Dashboard.</li>
          <li>Paste the server URL and token above.</li>
          <li>Click <strong>Enroll Agent</strong>.</li>
          <li>Save the returned JSON into the agent's <code>config.json</code>.</li>
          <li>Agent auto-connects and starts sending telemetry.</li>
        </ol>
      </div>
    </div>
  );
}