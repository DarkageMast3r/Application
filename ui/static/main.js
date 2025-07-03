window.SmartCare = {
    currentPage: 'dashboard',
    
    init() {
        this.setupNavigation();
        console.log('SmartCare Assist initialized ✅');
    },

    setupNavigation() {
        const navTabs = document.querySelectorAll('.nav-tab');
        navTabs.forEach(tab => {
            tab.addEventListener('click', (e) => this.handleNavClick(e));
        });
    },

    handleNavClick(event) {
        event.preventDefault();
        const tab = event.target;
        const page = tab.getAttribute('data-page') || this.extractPageFromText(tab.textContent);
        
        if (page === this.currentPage) return;
        
        document.querySelectorAll('.nav-tab').forEach(t => t.classList.remove('active'));
        tab.classList.add('active');
        
        this.loadPage(page);
    },

    extractPageFromText(text) {
        const pageMap = {
            'Dashboard': 'dashboard',
            'Signalen': 'signalen', 
            'Zorgtechnologie': 'zorgtechnologie',
            'Financiering': 'financiering',
            'Rapportage': 'rapportage'
        };
        
        for (const [key, value] of Object.entries(pageMap)) {
            if (text.includes(key)) return value;
        }
        return 'dashboard';
    },

    async loadPage(page) {
        try {
            const response = await fetch(`/App/pages/${page}`, {
                headers: {
                    'X-Requested-With': 'XMLHttpRequest'
                }
            });
            
            if (response.ok) {
                const html = await response.text();
                document.getElementById('main-content').innerHTML = html;
                this.currentPage = page;
            }
        } catch (error) {
            console.error('Error loading page:', error);
        }
    }
};