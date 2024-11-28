package vant_util

import (
	"sync"
	"time"
)

// Estrutura do RateLimiter
type RateLimiter struct {
	mu        sync.Mutex
	limit     int64 // Limite em bits por janela
	used      int64 // Quantidade de bits usados na janela atual
	window    time.Duration
	lastReset int64 // Usando Unix() para controle de tempo preciso
}

// Função para criar um novo RateLimiter
func NewRateLimiter(limit int64, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:     limit,
		window:    window,
		lastReset: time.Now().Unix(), // Usando Unix() para pegar tempo em segundos
	}
}

func (rl *RateLimiter) AllowSoft(bits int64) (bool, int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Verificar se a janela de tempo expirou (em segundos)
	now := time.Now().Unix() // Obtém tempo em segundos
	//fmt.Println(now - rl.lastReset)
	if now-rl.lastReset >= int64(rl.window.Seconds()) { // Verifica se passou 1 segundo
		rl.used = 0
		rl.lastReset = now
	}

	// Verificar se o uso atual mais o pedido excede o limite
	available := rl.limit - rl.used
	if rl.used+bits > rl.limit {
		// Se os bits excedem o limite e a quantidade disponível é 0, retorna false
		if available <= 0 {
			return false, available
		}

		rl.used = available
		// Caso contrário, permite e retorna a quantidade de bits restantes
		return true, available
	}

	// Consumir os bits e permitir a operação
	rl.used += bits
	// Retorna verdadeiro e a quantidade de dados restantes
	return true, available
}
