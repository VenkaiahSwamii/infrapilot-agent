import React, { useEffect, useState } from 'react';

const identityRows = (payload) => payload;

export default function TelemetryTable({ machineId, load, selectRows = identityRows }) {
  const [rows, setRows] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!machineId) return undefined;
    let active = true;
    const refresh = async () => {
      try {
        const payload = await load(machineId);
        const selected = selectRows(payload);
        if (active) {
          setRows(Array.isArray(selected) ? selected : selected ? [selected] : []);
          setError('');
        }
      } catch (requestError) {
        if (active) setError(requestError.message);
      }
    };
    refresh();
    const timer = setInterval(refresh, 5000);
    return () => {
      active = false;
      clearInterval(timer);
    };
  }, [load, machineId, selectRows]);

  if (error) return <div className="tab-empty" style={{ textAlign: 'center', padding: '40px', color: '#64748b' }}>Unable to load telemetry: {error}</div>;
  if (rows.length === 0) return <div className="tab-empty" style={{ textAlign: 'center', padding: '40px', color: '#64748b' }}>No telemetry data collected for this machine.</div>;

  const columns = Object.keys(rows[0]).filter((key) => !['id', 'machine_id'].includes(key));
  return (
    <div className="tab-container asset-table">
      <table>
        <thead><tr>{columns.map((column) => <th key={column}>{column.replaceAll('_', ' ')}</th>)}</tr></thead>
        <tbody>{rows.map((row, index) => (
          <tr key={row.id || `${row.machine_id || 'row'}-${index}`}>
            {columns.map((column) => {
              const value = row[column];
              return <td key={column}>{value == null ? 'N/A' : typeof value === 'object' ? JSON.stringify(value) : String(value)}</td>;
            })}
          </tr>
        ))}</tbody>
      </table>
    </div>
  );
}
