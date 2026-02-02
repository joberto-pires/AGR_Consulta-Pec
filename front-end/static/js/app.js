// AgroConsultoria - Frontend JavaScript

class AgroConsultoriaApp {
    constructor() {
        this.initHTMXExtensions();
        this.initEventListeners();
        this.initToastSystem();
    }

    initHTMXExtensions() {
        // Configuração global do HTMX
        htmx.config.defaultSwapStyle = 'innerHTML';
        htmx.config.scrollIntoViewOnBoost = true;
        
        // Adicionar classes durante requisições
        document.body.addEventListener('htmx:beforeRequest', (e) => {
            const target = e.detail.target;
            if (target) {
                target.classList.add('loading');
            }
        });

        document.body.addEventListener('htmx:afterRequest', (e) => {
            const target = e.detail.target;
            if (target) {
                target.classList.remove('loading');
            }
        });
    }

    initEventListeners() {
        // Menu responsivo
        document.addEventListener('click', (e) => {
            if (e.target.closest('.mobile-menu-toggle')) {
                document.querySelector('.sidebar').classList.toggle('active');
            }
        });

        // Fechar menu ao clicar fora
        document.addEventListener('click', (e) => {
            if (!e.target.closest('.sidebar') && !e.target.closest('.mobile-menu-toggle')) {
                document.querySelector('.sidebar').classList.remove('active');
            }
        });

        // Validação de formulários
        document.addEventListener('submit', (e) => {
            const form = e.target;
            if (form.hasAttribute('data-validate')) {
                e.preventDefault();
                if (this.validateForm(form)) {
                    htmx.ajax('POST', form.action, {
                        values: this.getFormData(form),
                        target: '#main-content'
                    });
                }
            }
        });
    }

    initToastSystem() {
        window.showToast = (message, type = 'success', duration = 5000) => {
            const container = document.getElementById('toast-container');
            const toast = document.createElement('div');
            toast.className = `toast toast-${type}`;
            toast.innerHTML = `
                <div class="toast-content">
                    <i class="fas fa-${type === 'success' ? 'check-circle' : type === 'error' ? 'exclamation-circle' : 'info-circle'}"></i>
                    <span>${message}</span>
                </div>
                <button class="toast-close" onclick="this.parentElement.remove()">
                    <i class="fas fa-times"></i>
                </button>
            `;
            
            container.appendChild(toast);
            
            setTimeout(() => {
                if (toast.parentElement) {
                    toast.remove();
                }
            }, duration);
        };
    }

    validateForm(form) {
        let isValid = true;
        const inputs = form.querySelectorAll('[required]');
        
        inputs.forEach(input => {
            if (!input.value.trim()) {
                input.classList.add('error');
                this.showToast(`Campo "${input.previousElementSibling?.textContent || input.name}" é obrigatório`, 'error');
                isValid = false;
            } else {
                input.classList.remove('error');
            }
        });
        
        return isValid;
    }

    getFormData(form) {
        const formData = new FormData(form);
        const data = {};
        
        for (let [key, value] of formData.entries()) {
            data[key] = value;
        }
        
        return data;
    }

    // Funções utilitárias
    formatCurrency(value) {
        return new Intl.NumberFormat('pt-BR', {
            style: 'currency',
            currency: 'BRL'
        }).format(value);
    }

    formatDate(dateString) {
        return new Date(dateString).toLocaleDateString('pt-BR');
    }

    formatArea(hectares) {
        return `${parseFloat(hectares).toFixed(2)} ha`;
    }

    // Modal system
    openModal(content, title = '') {
        const modal = document.createElement('div');
        modal.className = 'modal';
        modal.innerHTML = `
            <div class="modal-overlay" onclick="this.parentElement.remove()"></div>
            <div class="modal-content">
                <div class="modal-header">
                    <h3>${title}</h3>
                    <button class="modal-close" onclick="this.closest('.modal').remove()">
                        <i class="fas fa-times"></i>
                    </button>
                </div>
                <div class="modal-body">
                    ${content}
                </div>
            </div>
        `;
        
        document.body.appendChild(modal);
    }

    // Chart functions for reports
    createChart(canvasId, type, data, options = {}) {
        const ctx = document.getElementById(canvasId).getContext('2d');
        return new Chart(ctx, {
            type: type,
            data: data,
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'top',
                    },
                    title: {
                        display: true,
                        text: options.title || ''
                    }
                },
                ...options
            }
        });
    }
}

// Initialize app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.app = new AgroConsultoriaApp();
    
    // Update stats on dashboard
    if (window.location.pathname === '/') {
        htmx.ajax('GET', '/api/dashboard/stats', {
            target: '#dashboard-stats'
        });
    }
});

// Global utility functions
function confirmAction(message, callback) {
    if (confirm(message)) {
        callback();
    }
}

function showLoading(element) {
    element.innerHTML = '<div class="spinner"><i class="fas fa-spinner fa-spin"></i></div>';
}

function hideLoading(element, originalContent) {
    element.innerHTML = originalContent;
}

// Dark Mode
class DarkMode {
    constructor() {
        this.theme = localStorage.getItem('theme') || 'light';
        this.init();
    }

    init() {
        this.setTheme(this.theme);
        document.getElementById('theme-toggle').addEventListener('click', () => this.toggle());
    }

    setTheme(theme) {
        document.documentElement.setAttribute('data-theme', theme);
        localStorage.setItem('theme', theme);
        this.updateIcon(theme);
    }

    toggle() {
        const newTheme = this.theme === 'light' ? 'dark' : 'light';
        this.setTheme(newTheme);
        this.theme = newTheme;
    }

    updateIcon(theme) {
        const icon = document.querySelector('#theme-toggle i');
        icon.className = theme === 'light' ? 'fas fa-moon' : 'fas fa-sun';
    }
}

// Inicializar quando o DOM estiver carregado
document.addEventListener('DOMContentLoaded', () => {
    // ... inicialização existente ...
    new DarkMode();
});


