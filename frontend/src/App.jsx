import React, { useEffect, useState } from 'react';
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts';
import { 
  Activity, 
  Wifi, 
  Shield, 
  Globe, 
  Cpu, 
  Radio, 
  AlertTriangle,
  Zap
} from 'lucide-react';

function App() {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [history, setHistory] = useState([]);

  const fetchData = () => {
    fetch('http://localhost:8080/health')
      .then(res => {
        if (res.ok) return res.json();
        throw new Error('Backend not reachable');
      })
      .then(jsonData => {
        setData(jsonData);
        setError(null);
        
        // Update History for Graph
        setHistory(prev => {
           const now = new Date();
           const timeLabel = `${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`;
           const newData = [...prev, { time: timeLabel, latency: jsonData.latency_ms }];
           return newData.slice(-20); // Keep last 20 points
        });
      })
      .catch(err => {
        setError(err.message);
        setData(null);
      });
  };

  useEffect(() => {
    fetchData(); // Initial fetch
    const interval = setInterval(fetchData, 1000); // Faster polling for smooth graph
    return () => clearInterval(interval);
  }, []);

  // Cyberpunk Styles
  const styles = {
    container: {
      width: '100vw',
      minHeight: '100vh',
      backgroundColor: '#0a0a0f',
      color: '#e0e0e0',
      padding: '2rem',
      boxSizing: 'border-box',
    },
    header: {
      marginBottom: '2rem',
      borderBottom: '1px solid #333',
      paddingBottom: '1rem',
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'center'
    },
    title: {
      fontSize: '2rem',
      fontWeight: 'bold',
      background: 'linear-gradient(45deg, #00f2ff, #00c3ff)',
      WebkitBackgroundClip: 'text',
      WebkitTextFillColor: 'transparent',
      margin: 0,
      letterSpacing: '1px'
    },
    grid: {
      display: 'grid',
      gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))',
      gap: '20px',
      marginBottom: '20px'
    },
    card: {
      background: 'rgba(20, 20, 30, 0.6)',
      border: '1px solid #333',
      borderRadius: '12px',
      padding: '20px',
      backdropFilter: 'blur(10px)',
      boxShadow: '0 4px 6px rgba(0,0,0,0.3)',
      transition: 'transform 0.2s, border-color 0.2s',
    },
    value: {
      fontSize: '1.8rem',
      fontWeight: 'bold',
      margin: '10px 0',
      fontFamily: 'monospace'
    },
    label: {
      color: '#888',
      fontSize: '0.9rem',
      textTransform: 'uppercase',
      letterSpacing: '1px',
      display: 'flex',
      alignItems: 'center',
      gap: '8px'
    },
    btn: {
        background: 'linear-gradient(45deg, #2196F3, #21CBF3)',
        border: 'none',
        borderRadius: '4px',
        padding: '8px 16px',
        color: 'white',
        cursor: 'pointer',
        fontWeight: 'bold',
        marginTop: '10px',
        width: '100%',
        textTransform: 'uppercase',
        letterSpacing: '1px'
    },
    btnDanger: {
        background: 'linear-gradient(45deg, #FF5252, #FF1744)',
    },
    btnChaos: {
        background: 'linear-gradient(45deg, #9C27B0, #E040FB)',
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <div>
           <h1 style={styles.title}>AAR <span style={{fontSize: '0.5em', color: '#555'}}>// AUTONOMOUS ADAPTIVE ROUTER</span></h1>
           <div style={{color: '#666', fontSize: '0.8rem', marginTop: '5px'}}>SYSTEM STATUS: ONLINE</div>
        </div>
        <div style={{textAlign: 'right'}}>
            <Activity color="#00f2ff" />
        </div>
      </div>

      {error && (
        <div style={{...styles.card, borderColor: '#ff4444', color: '#ff4444', marginBottom: '20px'}}>
          <AlertTriangle size={20} style={{marginRight: '10px', verticalAlign: 'middle'}}/>
          <strong>CONNECTION LOST:</strong> {error}
        </div>
      )}

      {data && (
        <>
            {/* Top Metrics Row */}
            <div style={styles.grid}>
                <Card style={styles.card}>
                    <div style={styles.label}><Globe size={16}/> Internet Status</div>
                    <div style={{...styles.value, color: data.connectivity ? '#00e676' : '#ff1744'}}>
                        {data.connectivity ? "ONLINE" : "OFFLINE"}
                    </div>
                    <div style={{fontSize: '0.8rem', color: '#666'}}>
                        WAN Interface: <span style={{color: '#fff'}}>{data.active_wan}</span>
                    </div>
                </Card>

                <Card style={styles.card}>
                    <div style={styles.label}><Wifi size={16}/> Wi-Fi Health</div>
                    <div style={{...styles.value, color: getQualityColor(data.wifi_quality)}}>
                        {data.wifi_quality}%
                    </div>
                    <div style={{fontSize: '0.8rem', color: '#666'}}>
                        Channel: <span style={{color: '#fff'}}>{data.wifi_channel}</span>
                    </div>
                </Card>

                <Card style={styles.card}>
                    <div style={styles.label}><Shield size={16}/> VPN Tunnel</div>
                    <div style={{...styles.value, color: data.vpn_status === "Connected" ? '#2979ff' : '#757575'}}>
                        {data.vpn_status.toUpperCase()}
                    </div>
                    <button 
                        onClick={toggleVPN}
                        style={{...styles.btn, ...(data.vpn_status === "Connected" ? styles.btnDanger : {})}}
                    >
                        {data.vpn_status === "Connected" ? 'DISENGAGE' : 'INITIALIZE'}
                    </button>
                </Card>

                <Card style={styles.card}>
                    <div style={styles.label}><Zap size={16}/> Chaos Engineering</div>
                    <div style={{...styles.value, color: '#e040fb'}}>TEST MODE</div>
                    <button onClick={toggleChaos} style={{...styles.btn, ...styles.btnChaos}}>
                        INJECT LATENCY (500ms)
                    </button>
                </Card>
            </div>

            {/* Real-time Graph */}
            <div style={{...styles.card, height: '300px', marginBottom: '20px'}}>
                <div style={styles.label}><Activity size={16}/> Real-Time Latency (ms)</div>
                <div style={{width: '100%', height: '100%', marginTop: '10px'}}>
                    <ResponsiveContainer width="100%" height="90%">
                        <LineChart data={history}>
                            <CartesianGrid strokeDasharray="3 3" stroke="#333" />
                            <XAxis dataKey="time" stroke="#666" tick={{fill: '#666', fontSize: 10}} />
                            <YAxis stroke="#666" tick={{fill: '#666', fontSize: 10}}/>
                            <Tooltip 
                                contentStyle={{backgroundColor: '#0d1117', borderColor: '#333', color: '#fff'}}
                                itemStyle={{color: '#00f2ff'}}
                            />
                            <Line 
                                type="monotone" 
                                dataKey="latency" 
                                stroke="#00f2ff" 
                                strokeWidth={2} 
                                dot={false} 
                                activeDot={{ r: 6, fill: '#fff' }}
                                isAnimationActive={false}
                            />
                        </LineChart>
                    </ResponsiveContainer>
                </div>
            </div>

            {/* System Log / Footer */}
             <div style={{...styles.card, fontSize: '0.8rem', color: '#666'}}>
                <div style={styles.label}><Cpu size={16}/> System Engine</div>
                <div style={{marginTop: '10px'}}>
                    Auto-Pilot Version: <strong>v0.1.0-alpha</strong> | Strategy: <strong>Latency-Drift-Correction</strong>
                    <br/>
                    Current Latency: <strong style={{color: '#fff'}}>{data.latency_ms}ms</strong>
                </div>
             </div>
        </>
      )}

      {!data && !error && (
        <div style={{textAlign: 'center', marginTop: '50px', color: '#666'}}>
            <div className="loading-spinner"></div>
            CONNECTING TO NEURAL NET...
        </div>
      )}
    </div>
  );

  function getQualityColor(q) {
      if (q > 70) return '#00e676';
      if (q > 40) return '#ffea00';
      return '#ff1744';
  }

  function toggleVPN() {
    fetch('http://localhost:8080/vpn/toggle', { method: 'POST' })
      .then(() => fetchData()) 
      .catch(err => console.error("VPN Error", err));
  }

  function toggleChaos() {
    fetch('http://localhost:8080/chaos/lag', { 
      method: 'POST',
      body: JSON.stringify({ enable: true })
    }).catch(err => console.error("Chaos Error", err));

    setTimeout(() => {
        fetch('http://localhost:8080/chaos/lag', { 
            method: 'POST',
            body: JSON.stringify({ enable: false })
        }).catch(err => console.error("Chaos Reset Error", err));
    }, 12000);
  }
}

function Card({children, style}) {
    return (
        <div style={style}>
            {children}
        </div>
    )
}

export default App;
