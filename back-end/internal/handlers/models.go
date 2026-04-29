package handlers
 
import "time"
 
// Cliente already defined in clientes.go
 
type Propriedade struct {
	ID               int       `json:"id"`
	ClienteID        int       `json:"cliente_id"`
	ClienteNome      string    `json:"cliente_nome"`
	Nome             string    `json:"nome"`
	Hectares         float64   `json:"hectares"`
	Municipio        string    `json:"municipio"`
	Estado           string    `json:"estado"`
	Coordenadas      string    `json:"coordenadas"`
	CaracteristicaSolo string  `json:"caracteristica_solo"`
	Infraestrutura   string    `json:"infraestrutura"`
	HistoricoCulturas string   `json:"historico_culturas"`
	DataCadastro     time.Time `json:"data_cadastro"`
}
 
type AnaliseSolo struct {
	ID               int       `json:"id"`
	PropriedadeID    int       `json:"propriedade_id"`
	PropriedadeNome  string    `json:"propriedade_nome"`
	ClienteNome      string    `json:"cliente_nome"`
	Talhao           string    `json:"talhao"`
	DataAmostra      string    `json:"data_amostra"`
	PH               float64   `json:"ph"`
	MateriaOrganica  float64   `json:"materia_organica"`
	Nitrogenio       float64   `json:"nitrogenio"`
	Fosforo          float64   `json:"fosforo"`
	Potassio         float64   `json:"potassio"`
	CTC              float64   `json:"ctc"`
	SaturacaoBases   float64   `json:"saturacao_bases"`
	RecCalcario      string    `json:"rec_calcario"`
	RecAdubacao      string    `json:"rec_adubacao"`
	PlanoCorrecao    string    `json:"plano_correcao"`
	Status           string    `json:"status"`
	DataCadastro     time.Time `json:"data_cadastro"`
}
 
type Consulta struct {
	ID              int       `json:"id"`
	ClienteID       int       `json:"cliente_id"`
	ClienteNome     string    `json:"cliente_nome"`
	PropriedadeID   int       `json:"propriedade_id"`
	PropriedadeNome string    `json:"propriedade_nome"`
	DataConsulta    string    `json:"data_consulta"`
	Tipo            string    `json:"tipo"`
	Diagnostico     string    `json:"diagnostico"`
	Recomendacoes   string    `json:"recomendacoes"`
	PlanejamentoSafra string  `json:"planejamento_safra"`
	ControlePragas  string    `json:"controle_pragas"`
	GestaoIrrigacao string    `json:"gestao_irrigacao"`
	AnaliseCustos   string    `json:"analise_custos"`
	Status          string    `json:"status"`
	DataCadastro    time.Time `json:"data_cadastro"`
}
 
type Monitoramento struct {
	ID              int       `json:"id"`
	PropriedadeID   int       `json:"propriedade_id"`
	PropriedadeNome string    `json:"propriedade_nome"`
	ClienteNome     string    `json:"cliente_nome"`
	DataRegistro    string    `json:"data_registro"`
	TipoRegistro    string    `json:"tipo_registro"`
	Descricao       string    `json:"descricao"`
	RegistroAplicacao string  `json:"registro_aplicacao"`
	Produtividade   float64   `json:"produtividade"`
	Alertas         string    `json:"alertas"`
	DataCadastro    time.Time `json:"data_cadastro"`
}
 

