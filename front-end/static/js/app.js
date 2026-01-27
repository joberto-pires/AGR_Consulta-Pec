// static/js/app.js
// Funções utilitárias para o sistema

// Formatação de dados
const Formatadores = {
    data: (dataString) => {
        if (!dataString) return '--';
        const data = new Date(dataString);
        return data.toLocaleDateString('pt-BR');
    },

    numero: (valor, decimais = 2) => {
        if (!valor && valor !== 0) return '--';
        return Number(valor).toLocaleString('pt-BR', {
            minimumFractionDigits: decimais,
            maximumFractionDigits: decimais
        });
    },

    moeda: (valor) => {
        if (!valor && valor !== 0) return 'R$ --';
        return Number(valor).toLocaleString('pt-BR', {
            style: 'currency',
            currency: 'BRL'
        });
    },

    area: (hectares) => {
        if (!hectares) return '-- ha';
        return `${Number(hectares).toLocaleString('pt-BR', {
            minimumFractionDigits: 1,
            maximumFractionDigits: 1
        })} ha`;
    }
};

// Modal System
const Modal = {
    open: (content, options = {}) => {
        const modal = document.createElement('div');
        modal.className = 'modal-overlay';
        modal.innerHTML = `
            <div class="modal">
                <div class="modal-header">
                    <h3>${options.title || ''}</h3>
                    <button class="modal-close" onclick="Modal.close()">
                        <i class="fas fa-times"></i>
                    </button>
                </div>
                <div class="modal-body">
                    ${content}
                </div>
            </div>
        `;
        
        document.body.appendChild(modal);
        setTimeout(() => modal.classList.add('active'), 10);
    },

    close: () => {
        const modal = document.querySelector('.modal-overlay');
        if (modal) {
            modal.classList.remove('active');
            setTimeout(() => modal.remove(), 300);
        }
    }
};

// Fechar modal com ESC
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') Modal.close();
});

// Dashboard Charts
const Graficos = {
    produtividade: (data) => {
        // Implementar gráfico de produtividade
        console.log('Gerando gráfico de produtividade:', data);
    },

    analiseSolo: (data) => {
        // Implementar gráfico de análise de solo
        console.log('Gerando gráfico de análise de solo:', data);
    }
};

// Event Listeners
document.addEventListener('DOMContentLoaded', () => {
    // Adicionar estilos para modais dinamicamente
    if (!document.querySelector('#modal-styles')) {
        const style = document.createElement('style');
        style.id = 'modal-styles';
        style.textContent = `
            .modal-overlay {
                position: fixed;
                top: 0;
                left: 0;
                right: 0;
                bottom: 0;
                background: rgba(0,0,0,0.5);
                display: flex;
                align-items: center;
                justify-content: center;
                opacity: 0;
                transition: opacity 0.3s ease;
                z-index: 1000;
            }
            
            .modal-overlay.active {
                opacity: 1;
            }
            
            .modal {
                background: white;
                border-radius: var(--border-radius);
                width: 90%;
                max-width: 600px;
                max-height: 90vh;
                overflow-y: auto;
                transform: translateY(-20px);
                transition: transform 0.3s ease;
            }
            
            .modal-overlay.active .modal {
                transform: translateY(0);
            }
            
            .modal-header {
                padding: 1.5rem;
                border-bottom: 1px solid #eee;
                display: flex;
                justify-content: space-between;
                align-items: center;
            }
            
            .modal-body {
                padding: 1.5rem;
            }
            
            .modal-close {
                background: none;
                border: none;
                font-size: 1.2rem;
                cursor: pointer;
                color: #666;
            }
        `;
        document.head.appendChild(style);
    }
});

// Exportar dados para CSV
function exportarParaCSV(dados, nomeArquivo = 'dados.csv') {
    let csvContent = "data:text/csv;charset=utf-8,";
    
    // Cabeçalho
    const cabecalho = Object.keys(dados[0]).join(';');
    csvContent += cabecalho + "\r\n";
    
    // Dados
    dados.forEach(item => {
        const linha = Object.values(item).join(';');
        csvContent += linha + "\r\n";
    });
    
    const link = document.createElement('a');
    link.setAttribute('href', encodeURI(csvContent));
    link.setAttribute('download', nomeArquivo);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}
