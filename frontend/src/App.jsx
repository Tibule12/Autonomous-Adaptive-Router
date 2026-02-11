import React, { useEffect, useState } from 'react';

function App() {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);

  const fetchData = () => {
    fetch('http://localhost:8080/health')
      .then(res => {
        if (res.ok) return res.json();
        throw new Error('Backend not reachable');
      })
      .then(jsonData => {
        setData(jsonData);
        setError(null);
      })
      .catch(err => {
        setError(err.message);
        setData(null);
      });
  };

  useEffect(() => {
    fetchData(); // Initial fetch
    const interval = setInterval(fetchData, 2000); // Poll every 2s
    return () => clearInterval(interval);
  }, []);

  return (
    <div style={{ padding: '20px', fontFamily: 'Arial, sans-serif', maxWidth: '800px', margin: '0 auto' }}>
      <h1>AAR Dashboard <span style={{fontSize: '0.6em', color: '#666'}}>(MVP Phase 1)</span></h1>
      
      {error && (
        <div style={{ padding: '15px', background: '#ffebee', color: '#c62828', borderRadius: '5px', marginBottom: '20px' }}>
          <strong>Error:</strong> {error}
          <p style={{margin: '5px 0 0 0', fontSize: '0.9em'}}>Make sure the backend is running (<code>go run main.go</code>)</p>
        </div>
      )}

      {data && (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '20px' }}>
          
          <Card title="System Status" value={data.status} color="#2e7d32" />
          
          <Card 
            title="Internet Connectivity" 
            value={data.connectivity ? "Online" : "Offline"} 
            color={data.connectivity ? "#2e7d32" : "#c62828"} 
          />
          
          <Card 
            title="Latency (Google DNS)" 
            value={`${data.latency_ms} ms`} 
            color={data.latency_ms < 50 ? "#2e7d32" : data.latency_ms < 150 ? "#f57c00" : "#c62828"} 
          />
          
          <Card 
            title="Active WAN Interface" 
            value={data.active_wan} 
            color="#0288d1" 
          />

          <Card
            title="Wi-Fi Health"
            value={`Ch ${data.wifi_channel} (${data.wifi_quality}%)`}
            color={data.wifi_quality > 70 ? "#2e7d32" : data.wifi_quality > 40 ? "#f57c00" : "#c62828"}
          />

          <Card 
            title="VPN Status" 
            value={data.vpn_status} 
            color={data.vpn_status === "Connected" ? "#1565c0" : "#757575"}
            action={
              <button 
                onClick={toggleVPN}
                style={{
                  marginTop: '10px',
                  padding: '8px 16px',
                  backgroundColor: data.vpn_status === "Connected" ? '#ef5350' : '#4caf50',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: 'pointer'
                }}
              >
                {data.vpn_status === "Connected" ? 'Simulate Kill Switch' : 'Enable VPN'}
              </button>
            }
          />

          {/* Chaos Tool */}
          <Card 
            title="Chaos Engineering" 
            value="Simulate Lag" 
            color="#8e24aa"
            action={
              <button 
                onClick={toggleChaos}
                style={{
                  marginTop: '10px',
                  padding: '8px 16px',
                  backgroundColor: '#ab47bc',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: 'pointer'
                }}
              >
                Inject 500ms Latency
              </button>
            }
          />
        </div>
      )}

      {!data && !error && <p>Loading system metrics...</p>}
    </div>
  );

  function toggleVPN() {
    fetch('http://localhost:8080/vpn/toggle', { method: 'POST' })
      .then(() => fetchData()) // Refresh data immediately
      .catch(err => console.error("VPN Toggle Error", err));
  }

  function toggleChaos() {
    // Enable lag to force the engine to react
    fetch('http://localhost:8080/chaos/lag', { 
      method: 'POST',
      body: JSON.stringify({ enable: true })
    }).catch(err => console.error("Chaos Error", err));

    // Automatically disable it after 10 seconds so the system can "recover"
    setTimeout(() => {
        fetch('http://localhost:8080/chaos/lag', { 
            method: 'POST',
            body: JSON.stringify({ enable: false })
        }).catch(err => console.error("Chaos Reset Error", err));
    }, 12000);
  }
}

function Card({ title, value, color, action }) {
  return (
    <div style={{ padding: '20px', border: '1px solid #ddd', borderRadius: '8px', boxShadow: '0 2px 4px rgba(0,0,0,0.05)' }}>
      <h3 style={{ margin: '0 0 10px 0', fontSize: '1em', color: '#555' }}>{title}</h3>
      <div style={{ fontSize: '1.5em', fontWeight: 'bold', color: color }}>
        {value}
      </div>
      {action && <div>{action}</div>}
    </div>
  );
}

export default App;
