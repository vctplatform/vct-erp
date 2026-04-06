import { useState, useEffect } from 'react';
import { Activity, GitBranch, Terminal, ShieldAlert, Cpu } from 'lucide-react';
import './index.css';

export default function App() {
  const [telemetry, setTelemetry] = useState<any>(null);

  useEffect(() => {
    // Fetch telemetry data periodically
    const interval = setInterval(() => {
      // In a real prod setup, this would be a WebSocket or API endpoint.
      // For this demo, we simulate the fetch or expect it to be served via Vite public
      fetch('/.telemetry.json').then(res => res.json()).then(data => setTelemetry(data)).catch(() => {});
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  // Mock telemetry if JSON fetch fails (for demonstration)
  const data = telemetry || {
    status: 'idle',
    active_tasks: [],
    completed_tasks: [{ id: 1, name: "Initialize DB schemas", assignee: "brian", status: "done" }],
    last_update: Date.now()
  };

  return (
    <div className="dashboard-container">
      {/* Sidebar */}
      <nav className="sidebar">
        <div className="brand">
          <Cpu className="brand-icon" />
          <h1>VCT Omniscience</h1>
        </div>
        <div className="nav-items">
          <a href="#" className="active"><Activity /> Live Swarm</a>
          <a href="#"><GitBranch /> Branches</a>
          <a href="#"><Terminal /> Logs</a>
          <a href="#"><ShieldAlert /> Vector Config</a>
        </div>
        <div className="system-status">
          <div className={`status-indicator ${data.status === 'processing' ? 'pulse' : ''}`}></div>
          <span>System: {data.status.toUpperCase()}</span>
        </div>
      </nav>

      {/* Main Content */}
      <main className="main-content">
        <header>
          <h2>Chairman Overview</h2>
          <div className="metrics">
            <div className="metric-card">
              <span className="label">Active Threads</span>
              <span className="value">{data.active_tasks.length}</span>
            </div>
            <div className="metric-card">
              <span className="label">Tasks Completed</span>
              <span className="value">{data.completed_tasks.length}</span>
            </div>
          </div>
        </header>

        <section className="board">
          {/* Active Work */}
          <div className="board-column">
            <h3>Active Agent Threads</h3>
            {data.active_tasks.length === 0 ? <p className="empty-state">No active tasks. Swarm is resting.</p> : null}
            {data.active_tasks.map((task: any) => (
              <div key={task.id} className="task-card processing">
                <div className="card-header">
                  <span className="task-id">#{task.id}</span>
                  <span className="assignee gradient-text">@{task.assignee}</span>
                </div>
                <h4 className="task-name">{task.name}</h4>
                <div className="progress-bar-container">
                  <div className="progress-bar indeterminate"></div>
                </div>
                <p className="status-text">AST/Vector Mapping & Docker Testing...</p>
              </div>
            ))}
          </div>

          {/* Merge Requests */}
          <div className="board-column">
            <h3>Ready for Merge (Done)</h3>
            {data.completed_tasks.map((task: any) => (
              <div key={task.id} className="task-card done">
                <div className="card-header">
                  <span className="task-id">#{task.id}</span>
                  <span className="assignee">@{task.assignee}</span>
                </div>
                <h4 className="task-name">{task.name}</h4>
                <button className="merge-btn">Merge PR</button>
              </div>
            ))}
          </div>
        </section>
      </main>
    </div>
  );
}
