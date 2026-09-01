import React, { useEffect, useState, useRef } from 'react';
import { Terminal, RefreshCw, AlertTriangle, ShieldAlert } from 'lucide-react';
import { useDashboardStore } from '../../../store/dashboardStore.jsx';

export default function TerminalTab({ machine }) {
  const [output, setOutput] = useState('');
  const [input, setInput] = useState('');
  const [status, setStatus] = useState('Disconnected'); // Disconnected, Connecting, Connected, Error
  const [errorMsg, setErrorMsg] = useState('');
  const socketRef = useRef(null);
  const consoleBottomRef = useRef(null);
  const inputRef = useRef(null);
  const { addToast } = useDashboardStore();

  const connectTerminal = () => {
    if (socketRef.current) {
      socketRef.current.close();
    }

    setStatus('Connecting');
    setErrorMsg('');
    setOutput('// Establishing secure shell tunnel to target host...\n');

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    // Connect to backend port 8080
    const wsUrl = `${protocol}//${window.location.hostname}:8080/api/v1/terminal/connect`;

    try {
      const ws = new WebSocket(wsUrl);
      socketRef.current = ws;

      ws.onopen = () => {
        setStatus('Connected');
        addToast('success', 'Terminal Connected', 'Live interactive terminal session started.');
        // Press enter once to get the initial Windows cmd.exe prompt
        ws.send('\n');
      };

      ws.onmessage = async (event) => {
        let text = '';
        if (typeof event.data === 'string') {
          text = event.data;
        } else if (event.data instanceof Blob) {
          text = await event.data.text();
        }
        setOutput((prev) => prev + text);
      };

      ws.onclose = (e) => {
        setStatus('Disconnected');
        setOutput((prev) => prev + '\n// Tunnel closed. Session terminated.\n');
      };

      ws.onerror = (err) => {
        console.error('WS Terminal Error:', err);
        setStatus('Error');
        setErrorMsg('Failed to establish tunnel connection.');
        setOutput((prev) => prev + '\n// Error: Connection failed.\n');
        addToast('critical', 'Tunnel Error', 'WebSocket connection failed.');
      };
    } catch (err) {
      setStatus('Error');
      setErrorMsg(err.message);
    }
  };

  useEffect(() => {
    connectTerminal();
    return () => {
      if (socketRef.current) {
        socketRef.current.close();
      }
    };
  }, [machine?.id]);

  useEffect(() => {
    consoleBottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [output]);

  const handleSend = (e) => {
    e.preventDefault();
    if (!input.trim() && e.type === 'submit') return;

    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(input + '\r\n');
      setInput('');
    } else {
      addToast('warning', 'Session Offline', 'Cannot send input. Terminal is disconnected.');
    }
  };

  const handleClear = () => {
    setOutput('// Console output cleared.\n');
  };

  return (
    <div className="tab-container" style={{ display: 'flex', flexDirection: 'column', gap: '12px', height: 'calc(100vh - 280px)' }}>
      {/* Header status bar */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', backgroundColor: '#0d1220', border: '1px solid var(--border-soft)', padding: '10px 16px', borderRadius: '10px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '10px', fontSize: '13px' }}>
          <Terminal size={16} color="#06b6d4" />
          <span style={{ color: 'var(--text)', fontWeight: 700 }}>Interactive Shell Session</span>
          <span style={{ color: 'var(--muted)' }}>|</span>
          <span style={{ color: status === 'Connected' ? '#22c55e' : status === 'Connecting' ? '#f59e0b' : '#ef4444', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '6px' }}>
            <span style={{ width: '8px', height: '8px', borderRadius: '50%', backgroundColor: status === 'Connected' ? '#22c55e' : status === 'Connecting' ? '#f59e0b' : '#ef4444' }} />
            {status.toUpperCase()}
          </span>
        </div>

        <div style={{ display: 'flex', gap: '8px' }}>
          <button
            onClick={handleClear}
            style={{
              padding: '6px 12px',
              backgroundColor: 'rgba(255, 255, 255, 0.03)',
              border: '1px solid var(--border-soft)',
              borderRadius: '6px',
              color: 'var(--text)',
              fontSize: '11px',
              cursor: 'pointer',
              fontWeight: 700
            }}
          >
            Clear Console
          </button>
          <button
            onClick={connectTerminal}
            disabled={status === 'Connecting'}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '6px',
              padding: '6px 12px',
              backgroundColor: '#06b6d4',
              border: 'none',
              borderRadius: '6px',
              color: '#080c14',
              fontSize: '11px',
              cursor: 'pointer',
              fontWeight: 700
            }}
          >
            <RefreshCw size={12} className={status === 'Connecting' ? 'spin' : ''} />
            Reconnect
          </button>
        </div>
      </div>

      {/* Console area */}
      <div
        style={{
          flex: 1,
          backgroundColor: '#05070c',
          border: '1px solid #1e293b',
          borderRadius: '12px',
          padding: '16px',
          fontFamily: '"JetBrains Mono", "Courier New", monospace',
          fontSize: '13px',
          color: '#cbd5e1',
          overflowY: 'auto',
          display: 'flex',
          flexDirection: 'column',
          boxShadow: 'inset 0 4px 20px rgba(0, 0, 0, 0.6)'
        }}
        onClick={() => inputRef.current?.focus()}
      >
        <pre style={{ margin: 0, whiteSpace: 'pre-wrap', flex: 1, lineHeight: '1.5' }}>{output}</pre>
        <div ref={consoleBottomRef} />
      </div>

      {/* Input prompt */}
      <form onSubmit={handleSend} style={{ display: 'flex', alignItems: 'center', gap: '10px', padding: '12px', backgroundColor: '#0d1220', border: '1px solid var(--border-soft)', borderRadius: '10px' }}>
        <span style={{ color: '#06b6d4', fontFamily: 'monospace', fontWeight: 800, fontSize: '15px' }}>$</span>
        <input
          type="text"
          ref={inputRef}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder='Type a command (e.g., "dir", "ipconfig", "ping google.com") and press Enter...'
          disabled={status !== 'Connected'}
          style={{
            flex: 1,
            background: 'none',
            border: 'none',
            color: '#f1f5f9',
            fontFamily: '"JetBrains Mono", monospace',
            fontSize: '13px',
            outline: 'none'
          }}
          autoFocus
        />
      </form>

      <style>{`
        @keyframes spin { to { transform: rotate(360deg); } }
        .spin { animation: spin 0.8s linear infinite; }
      `}</style>
    </div>
  );
}