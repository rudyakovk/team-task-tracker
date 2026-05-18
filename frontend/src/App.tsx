import "./styles.css";

const columns = [
  { title: "Backlog", count: 0 },
  { title: "Todo", count: 0 },
  { title: "In progress", count: 0 },
  { title: "Done", count: 0 },
];

export function App() {
  return (
    <main className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-mark">TT</span>
          <div>
            <strong>Team Task Tracker</strong>
            <span>Local workspace</span>
          </div>
        </div>

        <nav className="nav-list" aria-label="Main navigation">
          <a aria-current="page" href="/">
            Dashboard
          </a>
          <a href="/">Projects</a>
          <a href="/">Issues</a>
          <a href="/">Team</a>
        </nav>
      </aside>

      <section className="workspace">
        <header className="topbar">
          <div>
            <p className="eyebrow">Phase 0</p>
            <h1>Project scaffold is ready</h1>
          </div>
          <div className="status-pill">localhost</div>
        </header>

        <section className="summary-grid" aria-label="Project summary">
          <article>
            <span>Projects</span>
            <strong>0</strong>
          </article>
          <article>
            <span>Open issues</span>
            <strong>0</strong>
          </article>
          <article>
            <span>Team members</span>
            <strong>1</strong>
          </article>
        </section>

        <section className="board" aria-label="Task board preview">
          {columns.map((column) => (
            <article className="board-column" key={column.title}>
              <header>
                <h2>{column.title}</h2>
                <span>{column.count}</span>
              </header>
              <div className="empty-state">No issues yet</div>
            </article>
          ))}
        </section>
      </section>
    </main>
  );
}

