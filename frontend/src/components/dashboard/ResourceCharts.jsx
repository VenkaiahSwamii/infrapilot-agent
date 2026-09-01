import React from 'react';
import HistoryChart from '../../charts/HistoryChart.jsx';

export default function ResourceCharts({ cpuHistory }) {
  const chartData = cpuHistory.map((item) => ({
    timestamp: item.timestamp,
    cpu: item.value,
  }));

  return (
    <section style={{ background: '#0b1120', border: '1px solid #1f2937', borderRadius: 16, padding: 20, marginBottom: 20 }}>
      <h3 style={{ fontSize: '1rem', fontWeight: 700, marginBottom: 15, color: 'var(--foreground)' }}>Live CPU Utilization History</h3>
      <HistoryChart samples={chartData} range="1h" />
    </section>
  );
}
