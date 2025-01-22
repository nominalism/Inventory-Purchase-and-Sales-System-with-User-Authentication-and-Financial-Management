package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Material represents the materials table structure
type Material struct {
	ID             string  `json:"id"`
	Codigo         string  `json:"codigo"`
	Nome           string  `json:"nome"`
	Local          string  `json:"local"`
	DataCompra     string  `json:"dataCompra"`
	ValorCompra    float64 `json:"valorCompra"`    // Changed from string to float64
	ValorVenda     float64 `json:"valorVenda"`     // Changed from string to float64
	Estoque        int     `json:"estoque"`        // Changed from string to int
	EstoqueCritico int     `json:"estoqueCritico"` // Changed from string to int
	Fornecedor     string  `json:"fornecedor"`
}

// Adicionar após a struct Material
type Boleto struct {
	ID             int     `json:"id"`
	Nome           string  `json:"nome"`
	DataVencimento string  `json:"dataVencimento"`
	Valor          float64 `json:"valor"`
	PdfPath        string  `json:"pdfPath"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

var db *sql.DB
var jwtKey = []byte("your_secret_key")

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		claims := &jwt.StandardClaims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		var user User
		err = db.QueryRow("SELECT id, username, role FROM usuarios WHERE username = $1",
			claims.Subject).Scan(&user.ID, &user.Username, &user.Role)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func adminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			c.Abort()
			return
		}

		user, ok := userInterface.(User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	err := db.QueryRow("SELECT id, username, password, role FROM usuarios WHERE username = $1",
		credentials.Username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   user.Username,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  user,
	})
}

func main() {
	// Database connection - corrigindo a string de conexão
	connStr := "host=localhost port=5432 user=postgres dbname=sistema_materiais sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	// Public routes - sem autenticação
	router.POST("/login", login)

	// Rotas que exigem autenticação mas não precisam ser admin
	authorized := router.Group("/")
	authorized.Use(authMiddleware())
	{
		authorized.GET("/busca_material", searchMaterial)
		authorized.GET("/busca_material_id/:id", getMaterialByID)
		authorized.POST("/vender_material/:id", venderMaterial)
		authorized.GET("/lucros_mensais", getLucrosMensais)
		authorized.GET("/buscar_boletos", getBoletos) // Move this here

		// Admin routes
		admin := authorized.Group("/admin")
		admin.Use(adminOnly())
		{
			admin.POST("/comprar_material/:id", comprarMaterial)
			admin.POST("/add_material", addMaterial)
			admin.DELETE("/excluir_material/:id", deleteMaterial)
			admin.PUT("/atualizar_material/:id", updateMaterial)
			admin.GET("/caixa", getCaixa)
			admin.POST("/atualizar_caixa", atualizarCaixa)
			admin.GET("/usuarios", getUsers)
			admin.POST("/criar_usuario", createUser)
			admin.POST("/alterar_senha", updatePassword)

			// Boletos
			admin.POST("/adicionar_boleto", addBoleto)
			admin.PUT("/atualizar_boleto/:id", updateBoleto)
			admin.DELETE("/excluir_boleto/:id", deleteBoleto)
			admin.POST("/pagar_boleto/:id", pagarBoleto)
			admin.POST("/upload_boleto", uploadBoleto)
		}
	}

	// Start server
	router.Run(":5000")
}

// Função auxiliar para converter string para float64
func strToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// Função auxiliar para converter string para int
func strToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func addMaterial(c *gin.Context) {
	var material Material
	if err := c.BindJSON(&material); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro nos dados: " + err.Error()})
		return
	}

	query := `
        INSERT INTO materiais (
            codigo, nome, local, data_compra, valor_compra,
            valor_venda, estoque, estoque_critico, fornecedor
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	_, err := db.Exec(query,
		material.Codigo,
		material.Nome,
		material.Local,
		material.DataCompra,
		material.ValorCompra,    // Remove strToFloat64 conversion
		material.ValorVenda,     // Remove strToFloat64 conversion
		material.Estoque,        // Remove strToInt conversion
		material.EstoqueCritico, // Remove strToInt conversion
		material.Fornecedor,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar material"})
		log.Printf("Erro ao inserir no banco: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Material adicionado com sucesso!"})
}

func searchMaterial(c *gin.Context) {
	busca := c.Query("busca")

	query := `
        SELECT id, codigo, nome, local, valor_venda::text, estoque::text, data_compra::text
        FROM materiais
        WHERE (CAST(codigo AS TEXT) ILIKE $1 OR nome ILIKE $1)
        AND ativo = true  -- Adicionar esta condição
        ORDER BY nome
    `

	rows, err := db.Query(query, "%"+busca+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar materiais"})
		log.Printf("Erro ao buscar materiais no banco: %v", err)
		return
	}
	defer rows.Close()

	var materials []Material
	for rows.Next() {
		var m Material
		err := rows.Scan(&m.ID, &m.Codigo, &m.Nome, &m.Local, &m.ValorVenda, &m.Estoque, &m.DataCompra)
		if err != nil {
			log.Printf("Erro ao ler linha: %v", err)
			continue
		}
		log.Printf("Material encontrado: ID=%s", m.ID) // Debug log
		materials = append(materials, m)
	}

	c.JSON(http.StatusOK, gin.H{"materiais": materials})
}

func deleteMaterial(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Tentando desativar material ID: %s", id)

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Erro ao iniciar transação: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar processo de desativação"})
		return
	}

	// Verificar se o material existe e está ativo
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM materiais WHERE id = $1 AND ativo = true)", id).Scan(&exists)
	if err != nil {
		tx.Rollback()
		log.Printf("Erro ao verificar existência do material: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar material"})
		return
	}

	if !exists {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Material não encontrado ou já desativado"})
		return
	}

	// Apenas marcar como inativo em vez de excluir
	result, err := tx.Exec("UPDATE materiais SET ativo = false WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		log.Printf("Erro ao desativar material: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao desativar material"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Printf("Erro ao verificar desativação: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar desativação"})
		return
	}

	if rowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Material não encontrado"})
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("Erro ao finalizar transação: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao finalizar desativação"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Material desativado com sucesso"})
}

func updateMaterial(c *gin.Context) {
	id := c.Param("id")
	var material Material
	if err := c.BindJSON(&material); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro nos dados: " + err.Error()})
		return
	}

	query := `
        UPDATE materiais
        SET codigo = $1, nome = $2, local = $3, data_compra = $4,
            valor_compra = $5, valor_venda = $6, estoque = $7,
            estoque_critico = $8, fornecedor = $9
        WHERE id = $10
    `

	result, err := db.Exec(query,
		material.Codigo,
		material.Nome,
		material.Local,
		material.DataCompra,
		material.ValorCompra,    // No need for conversion anymore
		material.ValorVenda,     // No need for conversion anymore
		material.Estoque,        // No need for conversion anymore
		material.EstoqueCritico, // No need for conversion anymore
		material.Fornecedor,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao alterar material"})
		log.Printf("Erro ao alterar material: %v", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar atualização"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Material atualizado com sucesso"})
}

func getMaterialByID(c *gin.Context) {
	id := c.Param("id")

	query := `
        SELECT 
            id, 
            codigo, 
            nome, 
            local, 
            data_compra,
            CAST(valor_compra AS TEXT),
            CAST(valor_venda AS TEXT),
            CAST(estoque AS TEXT),
            CAST(estoque_critico AS TEXT),
            fornecedor
        FROM materiais 
        WHERE id = $1 AND ativo = true  -- Adicionar esta condição
    `

	var material Material
	var rawDataCompra string

	err := db.QueryRow(query, id).Scan(
		&material.ID,
		&material.Codigo,
		&material.Nome,
		&material.Local,
		&rawDataCompra,
		&material.ValorCompra,
		&material.ValorVenda,
		&material.Estoque,
		&material.EstoqueCritico,
		&material.Fornecedor,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material não encontrado"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar material por ID"})
		log.Printf("Erro ao buscar material por ID no banco: %v", err)
		return
	}
	if rawDataCompra != "" {
		t, err := time.Parse("2006-01-02T15:04:05Z", rawDataCompra)
		if err == nil {
			material.DataCompra = t.Format("2006-01-02")
		} else {
			material.DataCompra = rawDataCompra // fallback to original format
		}
	}
	c.JSON(http.StatusOK, gin.H{"materiais": []Material{material}})
}

func comprarMaterial(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Comprando material ID: %s", id) // Debug log

	var data struct {
		Quantidade int `json:"quantidade"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Iniciar transação
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar transação"})
		return
	}

	// Buscar material e caixa
	var material Material
	var valorCompra float64
	var saldoCaixa float64
	err = tx.QueryRow(`
        SELECT 
            CAST(valor_compra AS FLOAT8),
            CAST(estoque AS INTEGER)
        FROM materiais 
        WHERE id = $1 AND ativo = true  -- Adicionar esta condição
    `, id).Scan(&valorCompra, &material.Estoque)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Material não encontrado"})
		return
	}

	err = tx.QueryRow("SELECT saldo FROM caixa ORDER BY id DESC LIMIT 1").Scan(&saldoCaixa)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar caixa"})
		return
	}

	custoTotal := valorCompra * float64(data.Quantidade)
	if custoTotal > saldoCaixa {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Saldo insuficiente"})
		return
	}

	// Atualizar estoque e caixa
	_, err = tx.Exec("UPDATE materiais SET estoque = estoque + $1 WHERE id = $2", data.Quantidade, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar estoque"})
		return
	}

	_, err = tx.Exec("INSERT INTO caixa (saldo) VALUES ($1)", saldoCaixa-custoTotal)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar caixa"})
		return
	}

	// Registrar transação
	_, err = tx.Exec("INSERT INTO transacoes (tipo, material_id, quantidade, valor, data) VALUES ($1, $2, $3, $4, $5)",
		"compra", id, data.Quantidade, -custoTotal, time.Now())
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar transação"})
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao finalizar transação"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compra realizada com sucesso"})
}

func venderMaterial(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Vendendo material ID: %s", id) // Debug log

	var data struct {
		Quantidade int `json:"quantidade"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Iniciar transação
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar transação"})
		return
	}

	// Buscar material e estoque atual
	var estoque int
	var valorVenda float64
	err = tx.QueryRow("SELECT CAST(estoque AS INTEGER), CAST(valor_venda AS FLOAT8) FROM materiais WHERE id = $1 AND ativo = true", id).Scan(&estoque, &valorVenda)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Material não encontrado"})
		return
	}

	// Verificar se há estoque suficiente
	if estoque < data.Quantidade {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Estoque insuficiente"})
		return
	}

	// Calcular valor total da venda
	valorTotal := valorVenda * float64(data.Quantidade)

	// Atualizar estoque
	_, err = tx.Exec("UPDATE materiais SET estoque = estoque - $1 WHERE id = $2", data.Quantidade, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar estoque"})
		return
	}

	// Buscar saldo atual do caixa
	var saldoCaixa float64
	err = tx.QueryRow("SELECT saldo FROM caixa ORDER BY id DESC LIMIT 1").Scan(&saldoCaixa)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar caixa"})
		return
	}

	// Atualizar caixa
	_, err = tx.Exec("INSERT INTO caixa (saldo) VALUES ($1)", saldoCaixa+valorTotal)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar caixa"})
		return
	}

	// Registrar transação
	_, err = tx.Exec(
		"INSERT INTO transacoes (tipo, material_id, quantidade, valor, data) VALUES ($1, $2, $3, $4, $5)",
		"venda", id, data.Quantidade, valorTotal, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar transação"})
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao finalizar transação"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Venda realizada com sucesso"})
}

func getCaixa(c *gin.Context) {
	var saldo float64
	err := db.QueryRow("SELECT saldo FROM caixa ORDER BY id DESC LIMIT 1").Scan(&saldo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar saldo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"saldo": saldo})
}

func atualizarCaixa(c *gin.Context) {
	var data struct {
		Valor float64 `json:"valor"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var saldoAtual float64
	err := db.QueryRow("SELECT saldo FROM caixa ORDER BY id DESC LIMIT 1").Scan(&saldoAtual)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar saldo atual"})
		return
	}

	_, err = db.Exec("INSERT INTO caixa (saldo) VALUES ($1)", saldoAtual+data.Valor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar caixa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Caixa atualizado com sucesso"})
}

func getLucrosMensais(c *gin.Context) {
	rows, err := db.Query(`
        SELECT DATE_TRUNC('month', data) as mes, SUM(valor) as lucro
        FROM transacoes
        GROUP BY mes
        ORDER BY mes
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar lucros"})
		return
	}
	defer rows.Close()

	var meses []string
	var valores []float64
	for rows.Next() {
		var mes time.Time
		var lucro float64
		if err := rows.Scan(&mes, &lucro); err != nil {
			continue
		}
		meses = append(meses, mes.Format("01/2006"))
		valores = append(valores, lucro)
	}

	c.JSON(http.StatusOK, gin.H{
		"meses":   meses,
		"valores": valores,
	})
}

// Adicionar as novas funções
func addBoleto(c *gin.Context) {
	var boleto Boleto
	if err := c.BindJSON(&boleto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
        INSERT INTO boletos (nome, data_vencimento, valor, pdf_path)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	err := db.QueryRow(query,
		boleto.Nome,
		boleto.DataVencimento,
		boleto.Valor,
		boleto.PdfPath,
	).Scan(&boleto.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar boleto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Boleto adicionado com sucesso", "id": boleto.ID})
}

func getBoletos(c *gin.Context) {
	log.Println("========= INICIANDO BUSCA DE BOLETOS =========")

	query := `
        SELECT id, nome, data_vencimento::text, valor::text, pdf_path
        FROM boletos
        ORDER BY data_vencimento ASC
    `

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Erro na query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar boletos"})
		return
	}
	defer rows.Close()

	boletos := []Boleto{}
	for rows.Next() {
		var b Boleto
		err := rows.Scan(&b.ID, &b.Nome, &b.DataVencimento, &b.Valor, &b.PdfPath)
		if err != nil {
			log.Printf("Erro ao ler linha: %v", err)
			continue
		}
		log.Printf("Boleto encontrado: ID=%d, Nome=%s", b.ID, b.Nome) // Changed %s to %d for b.ID
		boletos = append(boletos, b)
	}

	log.Printf("Total de boletos encontrados: %d", len(boletos))
	c.JSON(http.StatusOK, gin.H{"boletos": boletos})
	log.Println("========= FIM DA BUSCA DE BOLETOS =========")
}

func updateBoleto(c *gin.Context) {
	id := c.Param("id")
	var boleto Boleto
	if err := c.BindJSON(&boleto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
        UPDATE boletos
        SET nome = $1, data_vencimento = $2, valor = $3, pdf_path = $4
        WHERE id = $5
    `

	result, err := db.Exec(query,
		boleto.Nome,
		boleto.DataVencimento,
		boleto.Valor,
		boleto.PdfPath,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar boleto"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Boleto não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Boleto atualizado com sucesso"})
}

func deleteBoleto(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM boletos WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir boleto"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Boleto não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Boleto excluído com sucesso"})
}

func uploadBoleto(c *gin.Context) {
	// Recebe o arquivo
	file, err := c.FormFile("pdf")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhum arquivo recebido"})
		return
	}

	// Gera um nome único para o arquivo
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)

	// Define o caminho do arquivo
	uploadDir := "../jessica/static/boletos"
	filepath := filepath.Join(uploadDir, filename)

	// Garante que o diretório existe
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar diretório"})
		return
	}

	// Salva o arquivo
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar arquivo"})
		return
	}

	// Retorna o caminho relativo do arquivo
	relativePath := "/boletos/" + filename
	c.JSON(http.StatusOK, gin.H{
		"message":  "Arquivo enviado com sucesso",
		"filepath": relativePath,
	})
}

func pagarBoleto(c *gin.Context) {
	id := c.Param("id")

	// Iniciar transação
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar transação"})
		return
	}

	// Buscar boleto e saldo do caixa
	var valorBoleto float64
	var saldoCaixa float64

	err = tx.QueryRow("SELECT valor FROM boletos WHERE id = $1", id).Scan(&valorBoleto)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Boleto não encontrado"})
		return
	}

	err = tx.QueryRow("SELECT saldo FROM caixa ORDER BY id DESC LIMIT 1").Scan(&saldoCaixa)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar caixa"})
		return
	}

	// Verificar se há saldo suficiente
	if valorBoleto > saldoCaixa {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Saldo insuficiente para pagar o boleto"})
		return
	}

	// Atualizar caixa
	_, err = tx.Exec("INSERT INTO caixa (saldo) VALUES ($1)", saldoCaixa-valorBoleto)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar caixa"})
		return
	}

	// Registrar transação
	_, err = tx.Exec(
		"INSERT INTO transacoes (tipo, valor, data) VALUES ($1, $2, $3)",
		"pagamento_boleto", -valorBoleto, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar transação"})
		return
	}

	// Excluir boleto
	_, err = tx.Exec("DELETE FROM boletos WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir boleto"})
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao finalizar transação"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Boleto pago com sucesso"})
}

func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, username, role FROM usuarios")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuários"})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Role); err != nil {
			continue
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var newUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Only allow creating 'funcionario' users
	if newUser.Role != "funcionario" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Só é permitido criar usuários funcionários"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
		return
	}

	_, err = db.Exec("INSERT INTO usuarios (username, password, role) VALUES ($1, $2, $3)",
		newUser.Username, string(hashedPassword), newUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário criado com sucesso"})
}

func updatePassword(c *gin.Context) {
	var data struct {
		Username    string `json:"username"`
		NewPassword string `json:"newPassword"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
		return
	}

	result, err := db.Exec("UPDATE usuarios SET password = $1 WHERE username = $2",
		string(hashedPassword), data.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar senha"})
		return
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Senha atualizada com sucesso"})
}
