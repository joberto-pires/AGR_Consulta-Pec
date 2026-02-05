// front-end/static/js/main.js
// JavaScript principal - carregado apenas uma vez

class AgroApp {
    constructor() {
        this.initSidebar();
        this.initTheme();
        this.initHTMX();
        this.initGlobalListeners();
    }

    initSidebar() {
        const menuToggle = document.getElementById('menu-toggle');
        const sidebar = document.querySelector('.sidebar');
        const sidebarOverlay = document.getElementById('sidebar-overlay');
        
        if (menuToggle && sidebar) {
            menuToggle.addEventListener('click', () => {
                sidebar.classList.toggle('mobile-show');
                if (sidebarOverlay) {
                    sidebarOverlay.classList.toggle('active');
                }
                document.body.style.overflow = sidebar.classList.contains('mobile-show') ? 'hidden' : '';
            });
        }
        
        if (sidebarOverlay) {
            sidebarOverlay.addEventListener('click', () => {
                sidebar?.classList.remove('mobile-show');
                sidebarOverlay.classList.remove('active');
                document.body.style.overflow = '';
            });
        }
        
        // Fechar com ESC
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                sidebar?.classList.remove('mobile-show');
                sidebarOverlay?.classList.remove('active');
                document.body.style.overflow = '';
            }
        });
    }

    initTheme() {
        const savedTheme = localStorage.getItem('agro-theme');
        const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        const theme = savedTheme || (systemPrefersDark ? 'dark' : 'light');
        
        document.documentElement.setAttribute('data-theme', theme);
        
        // Atualizar ícones
        document.querySelectorAll('.theme-toggle i').forEach(icon => {
            icon.className = theme === 'dark' ? 'fas fa-sun' : 'fas fa-moon';
        });
        
        // Botão de toggle
        document.querySelectorAll('.theme-toggle').forEach(btn => {
            btn.onclick = () => this.toggleTheme();
        });
    }

    toggleTheme() {
        const html = document.documentElement;
        const currentTheme = html.getAttribute('data-theme');
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        
        html.setAttribute('data-theme', newTheme);
        localStorage.setItem('agro-theme', newTheme);
        
        document.querySelectorAll('.theme-toggle i').forEach(icon => {
            icon.className = newTheme === 'dark' ? 'fas fa-sun' : 'fas fa-moon';
        });
    }

    initHTMX() {
        // Configurar HTMX globalmente
        if (window.htmx) {
            htmx.config.defaultSwapStyle = 'innerHTML';
            htmx.config.defaultSwapDelay = 0;
            htmx.config.scrollIntoViewOnBoost = true;
            
            // Adicionar classe durante carregamento
            document.body.addEventListener('htmx:beforeRequest', (e) => {
                const target = e.detail.target;
                if (target) {
                    target.classList.add('htmx-request');
                }
            });
            
            document.body.addEventListener('htmx:afterRequest', (e) => {
                const target = e.detail.target;
                if (target) {
                    target.classList.remove('htmx-request');
                }
            });
            
            // Executar scripts em conteúdo carregado via HTMX
            document.body.addEventListener('htmx:afterSwap', (e) => {
                if (e.detail.target.id === 'main-content') {
                    this.initFormFields();
                    this.initSidebar(); // Re-inicializar sidebar para novos elementos
                }
            });
        }
    }

    initFormFields() {
        // Formatação de CPF/CNPJ
        const cpfCnpjInput = document.getElementById('cpf_cnpj');
        if (cpfCnpjInput) {
            cpfCnpjInput.addEventListener('input', (e) => {
                this.formatDocument(e.target);
            });
        }
        
        // Formatação de telefone
        const telefoneInput = document.getElementById('telefone');
        if (telefoneInput) {
            telefoneInput.addEventListener('input', (e) => {
                this.formatPhone(e.target);
            });
        }
    }

    formatDocument(input) {
        let value = input.value.replace(/\D/g, '');
        
        if (value.length <= 11) {
            // CPF: 000.000.000-00
            value = value.replace(/(\d{3})(\d)/, '$1.$2');
            value = value.replace(/(\d{3})(\d)/, '$1.$2');
            value = value.replace(/(\d{3})(\d{1,2})$/, '$1-$2');
        } else {
            // CNPJ: 00.000.000/0000-00
            value = value.replace(/^(\d{2})(\d)/, '$1.$2');
            value = value.replace(/^(\d{2})\.(\d{3})(\d)/, '$1.$2.$3');
            value = value.replace(/\.(\d{3})(\d)/, '.$1/$2');
            value = value.replace(/(\d{4})(\d)/, '$1-$2');
        }
        
        input.value = value;
    }

    formatPhone(input) {
        let value = input.value.replace(/\D/g, '');
        
        if (value.length === 11) {
            // (11) 99999-9999
            value = value.replace(/^(\d{2})(\d)/, '($1) $2');
            value = value.replace(/(\d{5})(\d)/, '$1-$2');
        } else if (value.length === 10) {
            // (11) 9999-9999
            value = value.replace(/^(\d{2})(\d)/, '($1) $2');
            value = value.replace(/(\d{4})(\d)/, '$1-$2');
        } else if (value.length > 2) {
            value = value.replace(/^(\d{2})(\d)/, '($1) $2');
        }
        
        input.value = value;
    }

    initGlobalListeners() {
        // Header mobile responsivo
        this.updateMobileHeader();
        window.addEventListener('resize', () => this.updateMobileHeader());
    }

    updateMobileHeader() {
        const mobileHeader = document.querySelector('.mobile-header');
        if (mobileHeader) {
            mobileHeader.style.display = window.innerWidth <= 768 ? 'block' : 'none';
        }
    }
}

// Inicializar quando o DOM estiver pronto
document.addEventListener('DOMContentLoaded', () => {
    window.app = new AgroApp();
});
