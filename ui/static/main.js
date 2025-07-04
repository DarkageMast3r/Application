const API_BASE = 'http://localhost:4000';

const API_ENDPOINTS = {
    dashboardStats: `${API_BASE}/dashboardStats`,
    recentActivities: `${API_BASE}/recentActivities`,
    signaalStats: `${API_BASE}/signaalStats`,
    recentCases: `${API_BASE}/recentCases`,
    activeTechList: `${API_BASE}/activeTechList`,
    availableTechGrid: `${API_BASE}/availableTechGrid`,
    budgetStats: `${API_BASE}/budgetStats`,
    financieringOverview: `${API_BASE}/financieringOverview`,
    rapportageMaand: `${API_BASE}/rapportageMaand`
};

window.SmartCare = {
    init: function () {
        const page = SmartCare.currentPage;
        loadPage(page);

        document.querySelectorAll('.nav-tab').forEach(tab => {
            tab.addEventListener('click', e => {
                e.preventDefault();
                const targetPage = tab.getAttribute('data-page');
                if (targetPage && targetPage !== SmartCare.currentPage) {
                    SmartCare.currentPage = targetPage;
                    loadPage(targetPage);
                }
            });
        });
    }
};

function loadPage(page) {
    const mainContent = document.getElementById('main-content');

    let html = '';
    switch (page) {
        case 'dashboard':
            html = `
                <h2 class="section-title text-center">📊 Overzicht: Dashboard</h2>
                <div class="stats-grid" id="dashboardStats">
                    <div class="stat-card"><div class="stat-number" data-stat="activeClients">-</div><div class="stat-label">Actieve Cases</div></div>
                    <div class="stat-card"><div class="stat-number" data-stat="pendingSignals">-</div><div class="stat-label">In Behandeling</div></div>
                    <div class="stat-card"><div class="stat-number" data-stat="activeTech">-</div><div class="stat-label">Beschikbare Technologieën</div></div>
                    <div class="stat-card"><div class="stat-number" data-stat="budgetUsed">-</div><div class="stat-label">Budget Gebruikt</div></div>
                </div>
                <div class="card">
                    <h3 class="card-title text-center">🕒 Laatste Activiteiten</h3>
                    <div class="timeline" id="recentActivities"></div>
                </div>
            `;
            break;

        case 'signalen':
            html = `
                <h2 class="section-title text-center">🚨 Overzicht: Signalen</h2>
                <div class="card">
                    <h3 class="card-title text-center">📈 Statistieken</h3>
                    <div id="signaalStats"></div>
                </div>
                <div class="card">
                    <h3 class="card-title text-center">📋 Recente Cases</h3>
                    <div class="client-list" id="recentCases"></div>
                </div>
            `;
            break;

        case 'zorgtechnologie':
            html = `
                <h2 class="section-title text-center">⚙️ Overzicht: Zorgtechnologie</h2>
                <div class="card">
                    <h3 class="card-title text-center">🟢 Actieve Technologieën</h3>
                    <div id="activeTechList"></div>
                </div>
                <div class="card">
                    <h3 class="card-title text-center">🧪 Beschikbare Technologieën</h3>
                    <div class="tech-grid" id="availableTechGrid"></div>
                </div>
            `;
            break;

        case 'financiering':
            html = `
                <h2 class="section-title text-center">💰 Overzicht: Financiering</h2>
                <div class="card">
                    <h3 class="card-title text-center">📊 Budgetstatus</h3>
                    <div id="budgetStats"></div>
                </div>
                <div class="card">
                    <h3 class="card-title text-center">💼 Financieel Overzicht</h3>
                    <div class="stats-grid" id="financieringOverview"></div>
                </div>
            `;
            break;

        case 'rapportage':
            html = `
                <h2 class="section-title text-center">📈 Overzicht: Rapportage</h2>
                <div class="card">
                    <h3 class="card-title text-center">📆 Rapportage Maandelijks</h3>
                    <div id="main-rapportage"></div>
                </div>
            `;
            break;
    }

    mainContent.innerHTML = html;

    switch (page) {
        case 'dashboard': loadDashboard(); break;
        case 'signalen': loadSignalen(); break;
        case 'zorgtechnologie': loadZorgtechnologie(); break;
        case 'financiering': loadFinanciering(); break;
        case 'rapportage': loadRapportage(); break;
    }
}

function loadDashboard() {
    fetch(API_ENDPOINTS.dashboardStats)
        .then(res => res.json())
        .then(data => {
            document.querySelector('[data-stat="activeClients"]').textContent = data.activeClients;
            document.querySelector('[data-stat="pendingSignals"]').textContent = data.pendingSignals;
            document.querySelector('[data-stat="activeTech"]').textContent = data.activeTech;
            document.querySelector('[data-stat="budgetUsed"]').textContent = data.budgetUsed;
        });

    fetch(API_ENDPOINTS.recentActivities)
        .then(res => res.json())
        .then(activities => {
            const container = document.getElementById('recentActivities');
            container.innerHTML = activities.map(a => `
                <div class="timeline-item">
                    <strong>${new Date(a.timestamp).toLocaleString()}</strong>: ${a.activity}
                </div>
            `).join('');
        });
}

function loadSignalen() {
    fetch(API_ENDPOINTS.signaalStats)
        .then(res => res.json())
        .then(data => {
            const statsHTML = `
                <div class="stat-card"><div class="stat-number">${data.total}</div><div class="stat-label">Totaal</div></div>
                <div class="stat-card"><div class="stat-number">${data.newToday}</div><div class="stat-label">Nieuw Vandaag</div></div>
                <div class="stat-card"><div class="stat-number">${data.inProgress}</div><div class="stat-label">In Behandeling</div></div>
                <div class="stat-card"><div class="stat-number">${data.closed}</div><div class="stat-label">Afgesloten</div></div>
            `;
            document.getElementById('signaalStats').innerHTML = `<div class="stats-grid">${statsHTML}</div>`;
        });

    fetch(API_ENDPOINTS.recentCases)
        .then(res => res.json())
        .then(cases => {
            const container = document.getElementById('recentCases');
            container.innerHTML = cases.map(c => `
                <div class="case-item">👤 ${c.client} – <em>${c.status}</em></div>
            `).join('');
        });
}

function loadZorgtechnologie() {
    fetch(API_ENDPOINTS.activeTechList)
        .then(res => res.json())
        .then(list => {
            const html = list.map(item => `
                <div class="stat-card">
                    <div class="stat-number">${item.name}</div>
                    <div class="stat-label">${item.status}</div>
                </div>
            `).join('');
            document.getElementById('activeTechList').innerHTML = `<div class="stats-grid">${html}</div>`;
        });

    fetch(API_ENDPOINTS.availableTechGrid)
        .then(res => res.json())
        .then(tech => {
            const container = document.getElementById('availableTechGrid');
            container.innerHTML = tech.map(t => `
                <div class="tech-card">
                    <div class="tech-icon">🛠️</div>
                    <h4>${t.naam}</h4>
                    <p>${t.beschrijving}</p>
                    <small><strong>Leverancier:</strong> ${t.leverancier} – €${t.prijs}</small>
                </div>
            `).join('');
        });
}

function loadFinanciering() {
    fetch(API_ENDPOINTS.budgetStats)
        .then(res => res.json())
        .then(data => {
            const html = `
                <div class="stat-card"><div class="stat-number">${data.totalBudget}</div><div class="stat-label">Totaal Budget</div></div>
                <div class="stat-card"><div class="stat-number">${data.usedBudget}</div><div class="stat-label">Gebruikt</div></div>
                <div class="stat-card"><div class="stat-number">${data.remaining}</div><div class="stat-label">Beschikbaar</div></div>
            `;
            document.getElementById('budgetStats').innerHTML = `<div class="stats-grid">${html}</div>`;
        });

    fetch(API_ENDPOINTS.financieringOverview)
        .then(res => res.json())
        .then(items => {
            const html = items.map(item => `
                <div class="stat-card">
                    <div class="stat-number">€${item.Budget.MaxBedrag.toLocaleString()}</div>
                    <div class="stat-label">ZorgTech #${item.ZorgTechID} – ${item.Budget.BudgetStatus}</div>
                </div>
            `).join('');
            document.getElementById('financieringOverview').innerHTML = `<div class="stats-grid">${html}</div>`;
        });
}

function loadRapportage() {
    fetch(API_ENDPOINTS.rapportageMaand)
        .then(res => res.json())
        .then(data => {
            const html = `
                <div class="stat-card"><div class="stat-number">${data.signalenHandled}</div><div class="stat-label">Afgehandelde signalen</div></div>
                <div class="stat-card"><div class="stat-number">${data.techDeployments}</div><div class="stat-label">Tech Deployments</div></div>
                <div class="stat-card"><div class="stat-number">${data.budgetChange}</div><div class="stat-label">Budget verandering</div></div>
            `;
            document.getElementById('main-rapportage').innerHTML = `<div class="stats-grid">${html}</div>`;
        });
}
