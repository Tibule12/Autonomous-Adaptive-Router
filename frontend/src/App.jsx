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
  Zap,
  Server,
  Smartphone,
  Cloud
} from 'lucide-react';

function App() {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [history, setHistory] = useState([]);
  const [logs, setLogs] = useState([]);
  const [devices, setDevices] = useState([]);

  // Fetch History ONCE on load
  useEffect(() => {
    fetch('http://localhost:8080/history')
        .then(res => res.json())
        .then(jsonHistory => {
            if (jsonHistory && jsonHistory.length > 0) {
                setHistory(jsonHistory);
            }
        })
        .catch(console.error);
  }, []);

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
           // If we already have history from the API that is fresher, rely on that
           // Only append if the last timestamp is different
           const now = new Date();
           const timeLabel = `${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`;
           
           // Simple duplicate check or just append
           const newData = [...prev, { 
               time: timeLabel, 
               latency: jsonData.latency_ms,
               loss: jsonData.packet_loss,
               packet_loss: jsonData.packet_loss
           }];
           return newData.slice(-50); // Keep last 50 points
        });
      })
      .catch(err => {
        setError(err.message);
        setData(null);
      });
      
    // Fetch Logs
    fetch('http://localhost:8080/logs')
      .then(res => res.json())
      .then(jsonLogs => setLogs(jsonLogs || []))
      .catch(console.error);
      
    // Fetch Devices (New)
    fetch('http://localhost:8080/api/devices')
      .then(res => res.json())
      .then(jsonDevices => setDevices(jsonDevices || []))
      .catch(console.error);
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
    },
    terminal: {
        fontFamily: "'Courier New', monospace",
        background: '#0d1117',
        border: '1px solid #333',
        borderRadius: '8px',
        padding: '15px',
        height: '200px',
        overflowY: 'auto',
        fontSize: '0.85rem',
        marginTop: '20px',
        boxShadow: 'inset 0 0 10px rgba(0,0,0,0.5)'
    },
    logEntry: {
        marginBottom: '4px',
        display: 'flex',
        gap: '10px'
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
                        WAN: <span style={{color: '#fff'}}>{data.active_wan}</span> | Loss: <span style={{color: data.packet_loss > 0 ? '#ff1744' : '#00e676'}}>{data.packet_loss}%</span>
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
                    <div style={{display: 'flex', gap: '5px'}}>
                        <button onClick={() => toggleChaos('lag')} style={{...styles.btn, ...styles.btnChaos, flex: 1}}>
                            DoS (Lag)
                        </button>
                        <button onClick={() => toggleChaos('loss')} style={{...styles.btn, ...styles.btnChaos, flex: 1, borderColor: '#ff1744', color: '#ff1744'}}>
                            CUT (Loss)
                        </button>
                        <button onClick={() => toggleChaos('interference')} style={{...styles.btn, ...styles.btnChaos, flex: 1, borderColor: '#FFD740', color: '#FFD740'}}>
                            JAM (Wifi)
                        </button>
                    </div>
                </Card>

                {/* AI / Traffic CARD */}
                <Card style={styles.card}>
                     <div style={styles.label}><Cpu size={16}/> Traffic Intelligence</div>
                     <div style={{...styles.value, fontSize: '1.2rem', color: '#aa00ff'}}>
                        {data.traffic_type || "DEFAULT"}
                     </div>
                     <div style={{display: 'flex', gap: '5px', marginTop: '10px'}}>
                        <button onClick={() => setTraffic('Gaming')} style={{...styles.btn, flex: 1, fontSize: '0.7rem'}}>GAME</button>
                        <button onClick={() => setTraffic('Streaming')} style={{...styles.btn, flex: 1, fontSize: '0.7rem'}}>STREAM</button>
                        <button onClick={() => setTraffic('Default')} style={{...styles.btn, flex: 1, fontSize: '0.7rem'}}>IDLE</button>
                     </div>
                </Card>
            </div>

            {/* Connected Devices Grid (NEW) */}
            <h3 style={{...styles.title, fontSize: '1.2rem', marginBottom: '1rem'}}><Smartphone size={20} style={{verticalAlign: 'bottom'}}/> Connected Clients</h3>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '20px', marginBottom: '20px' }}>
                {devices.map((dev, i) => (
                    <div key={i} style={{
                        ...styles.card, 
                        display: 'flex', 
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        borderColor: dev.is_blocked ? '#ff1744' : '#333',
                        opacity: dev.is_blocked ? 0.6 : 1
                    }}>
                        <div>
                            <div style={{fontWeight: 'bold', color: '#fff'}}>{dev.name}</div>
                            <div style={{fontSize: '0.8rem', color: '#666'}}>{dev.ip}</div>
                            <div style={{fontSize: '0.7rem', color: '#444', fontFamily: 'monospace'}}>{dev.mac}</div>
                        </div>
                        <button 
                            onClick={() => toggleBlockDevice(dev.mac, dev.is_blocked)}
                            style={{
                                background: 'transparent',
                                border: '1px solid',
                                borderColor: dev.is_blocked ? '#00e676' : '#ff1744',
                                color: dev.is_blocked ? '#00e676' : '#ff1744',
                                padding: '5px 10px',
                                cursor: 'pointer',
                                fontSize: '0.7rem'
                            }}
                        >
                            {dev.is_blocked ? 'UNBLOCK' : 'BLOCK'}
                        </button>
                    </div>
                ))}
            </div>

            {/* Network Topology Map */}
            <div style={{...styles.card, marginBottom: '20px'}}>
                <div style={styles.label}><Globe size={16}/> Active Topology</div>
                <TopologyMap 
                    activeWan={data.active_wan} 
                    trafficType={data.traffic_type} 
                    packetLoss={data.packet_loss}
                />
            </div>

            {/* Real-time Graph */}
            <div style={{...styles.card, height: '320px', marginBottom: '20px', display: 'flex', flexDirection: 'column'}}>
                <div style={styles.label}><Activity size={16}/> Real-Time Metrics</div>
                <div style={{flexGrow: 1, height: '100%', minHeight: '200px', width: '100%', marginTop: '10px'}}>
                    <ResponsiveContainer width="100%" height="100%">
                        <LineChart data={history}>
                            <CartesianGrid strokeDasharray="3 3" stroke="#333" />
                            <XAxis dataKey="time" stroke="#666" tick={{fill: '#666', fontSize: 10}} />
                            <YAxis yAxisId="left" stroke="#666" tick={{fill: '#666', fontSize: 10}} label={{ value: 'ms', angle: -90, position: 'insideLeft', fill: '#666' }}/>
                            <YAxis yAxisId="right" orientation="right" stroke="#666" tick={{fill: '#666', fontSize: 10}} domain={[0, 100]} label={{ value: '%', angle: 90, position: 'insideRight', fill: '#666' }}/>
                            <Tooltip 
                                contentStyle={{backgroundColor: '#0d1117', borderColor: '#333', color: '#fff'}}
                                itemStyle={{color: '#00f2ff'}}
                            />
                            <Line 
                                yAxisId="left"
                                type="monotone" 
                                dataKey="latency" 
                                stroke="#00f2ff" 
                                name="Latency (ms)"
                                strokeWidth={2} 
                                dot={false} 
                                activeDot={{ r: 6, fill: '#fff' }}
                                isAnimationActive={false}
                            />
                             <Line 
                                yAxisId="right"
                                type="monotone" 
                                dataKey="packet_loss"
                                stroke="#ff1744" 
                                name="Packet Loss (%)"
                                strokeWidth={2} 
                                dot={false} 
                                activeDot={{ r: 6, fill: '#fff' }}
                                isAnimationActive={false}
                            />
                        </LineChart>
                    </ResponsiveContainer>
                </div>
            </div>

            {/* System Intelligence Console */}
            <div style={{...styles.card, height: '250px', marginBottom: '20px', display: 'flex', flexDirection: 'column', overflow: 'hidden'}}>
                <div style={styles.label}><Cpu size={16}/> System Intelligence Console</div>
                <div style={styles.terminal}>
                    {logs.length === 0 && <div style={{color: '#444'}}>Waiting for system events...</div>}
                    {logs.map((log, i) => (
                        <div key={i} style={styles.logEntry}>
                            <span style={{color: '#555'}}>[{log.timestamp}]</span>
                            <span style={{
                                color: log.level === 'ERROR' ? '#ff1744' : 
                                       log.level === 'WARN' ? '#ffea00' : 
                                       log.level === 'SUCCESS' ? '#00e676' : '#00b0ff',
                                fontWeight: 'bold'
                            }}>{log.level}</span>
                            <span style={{color: '#ddd'}}>{log.message}</span>
                        </div>
                    )).reverse()} 
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

  function toggleChaos(type) {
    fetch('http://localhost:8080/chaos/lag', { 
      method: 'POST',
      body: JSON.stringify({ enable: true, type: type })
    }).catch(err => console.error("Chaos Error", err));

    setTimeout(() => {
        fetch('http://localhost:8080/chaos/lag', { 
            method: 'POST',
            body: JSON.stringify({ enable: false, type: type })
        }).catch(err => console.error("Chaos Reset Error", err));
    }, 12000);
  }

  function setTraffic(type) {
    fetch('http://localhost:8080/simulation/traffic', { 
      method: 'POST',
      body: JSON.stringify({ type: type })
    }).catch(err => console.error("Traffic Error", err));
  }
  
  function toggleBlockDevice(mac, isCurrentlyBlocked) {
    const action = isCurrentlyBlocked ? 'unblock' : 'block';
    fetch('http://localhost:8080/api/devices', { 
        method: 'POST',
        body: JSON.stringify({ mac: mac, action: action })
    })
    .then(() => fetchData()) // Refresh list
    .catch(err => console.error("Block Error", err));
  }
}

function TopologyMap({ activeWan, trafficType, packetLoss }) {
    const isWan1 = activeWan && (activeWan.includes('wan1') || activeWan.includes('primary'));
    const isWan2 = activeWan && (activeWan.includes('wan2') || activeWan.includes('latency'));
    
    // Animation color based on status
    const pathColor = packetLoss > 10 ? '#ff1744' : '#00e676';
    const idleColor = '#333';

    return (
        <div style={{position: 'relative', height: '150px', display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '0 40px'}}>
            {/* Device Node */}
            <div style={{zIndex: 10, textAlign: 'center'}}>
                <Smartphone size={32} color="#fff" />
                <div style={{fontSize: '0.7rem', marginTop: '5px', color: '#888'}}>LOCAL</div>
            </div>

            {/* SVG Lines */}
            <div style={{position: 'absolute', top: 0, left: 0, width: '100%', height: '100%', pointerEvents: 'none'}}>
                <svg width="100%" height="100%" viewBox="0 0 100 100" preserveAspectRatio="none">
                    {/* Path to Router */}
                    <line x1="10" y1="50" x2="45" y2="50" stroke={pathColor} strokeWidth="0.5" strokeDasharray="2,2">
                         <animate attributeName="stroke-dashoffset" from="10" to="0" dur="1s" repeatCount="indefinite" />
                    </line>

                    {/* Path Router -> WAN 1 (Top) */}
                    <path d="M 55 50 C 65 50, 75 30, 90 30" stroke={isWan1 ? pathColor : idleColor} strokeWidth={isWan1 ? "1" : "0.5"} fill="none" />
                    
                    {/* Path Router -> WAN 2 (Bottom) */}
                    <path d="M 55 50 C 65 50, 75 70, 90 70" stroke={isWan2 ? pathColor : idleColor} strokeWidth={isWan2 ? "1" : "0.5"} fill="none" />
                </svg>
            </div>

            {/* Router Node */}
            <div style={{zIndex: 10, textAlign: 'center'}}>
                <div style={{background: '#0d1117', padding: '10px', borderRadius: '50%', border: `2px solid ${pathColor}`}}>
                    <Cpu size={32} color={pathColor} />
                </div>
                <div style={{fontSize: '0.7rem', marginTop: '5px', color: '#888'}}>AAR CORE</div>
            </div>

            {/* WAN Nodes */}
            <div style={{display: 'flex', flexDirection: 'column', gap: '40px', zIndex: 10}}>
                <div style={{textAlign: 'center', opacity: isWan1 ? 1 : 0.3}}>
                    <Cloud size={24} color={isWan1 ? "#2196F3" : "#555"} />
                    <div style={{fontSize: '0.6rem', color: '#888'}}>WAN 1 (DOCSIS)</div>
                </div>
                <div style={{textAlign: 'center', opacity: isWan2 ? 1 : 0.3}}>
                    <Cloud size={24} color={isWan2 ? "#ab47bc" : "#555"} />
                    <div style={{fontSize: '0.6rem', color: '#888'}}>WAN 2 (FIBER)</div>
                </div>
            </div>
        </div>
    );
}

function Card({children, style}) {
    return (
        <div style={style}>
            {children}
        </div>
    )
}

export default App;
