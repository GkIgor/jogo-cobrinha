	package main

	import (
		"fmt"
		"math/rand"
		"time"

		"github.com/nsf/termbox-go"
	)

	// Constantes para direções
	const (
		Cima    = 0
		Direita = 1
		Baixo   = 2
		Esquerda = 3
	)

	// Constantes para cores
	const (
		ColorSnake  = termbox.ColorGreen
		ColorFood   = termbox.ColorRed
		ColorBorder = termbox.ColorYellow
		ColorText   = termbox.ColorWhite
	)

	// Constantes para níveis de dificuldade
	const (
		Facil   = 150
		Medio   = 100
		Dificil = 70
		Expert  = 40
	)

	// Estrutura para representar uma coordenada no jogo
	type Coordenada struct {
		x, y int
	}

	// Estrutura principal do jogo
	type Jogo struct {
		cobra        []Coordenada
		comida       Coordenada
		direcao      int
		largura      int
		altura       int
		pontuacao    int
		nivel        string
		velocidade   time.Duration
		jogoAtivo    bool
		pausado      bool
		msgGameOver  bool
	}

	// Inicializa um novo jogo
	func NovoJogo(largura, altura int, nivel string) *Jogo {
		// Define a posição inicial da cobra
		cobra := []Coordenada{
			{largura / 2, altura / 2},
			{largura/2 - 1, altura / 2},
			{largura/2 - 2, altura / 2},
		}

		// Inicializa o jogo com velocidade baseada no nível de dificuldade
		j := &Jogo{
			cobra:     cobra,
			direcao:   Direita,
			largura:   largura,
			altura:    altura,
			pontuacao: 0,
			nivel:     nivel,
			jogoAtivo: true,
			pausado:   false,
		}

		// Define a velocidade de acordo com o nível
		switch nivel {
		case "Fácil":
			j.velocidade = Facil
		case "Médio":
			j.velocidade = Medio
		case "Difícil":
			j.velocidade = Dificil
		case "Expert":
			j.velocidade = Expert
		default:
			j.velocidade = Medio
		}

		// Gera a primeira comida
		j.gerarComida()

		return j
	}

	// Gera uma nova comida em uma posição aleatória
	func (j *Jogo) gerarComida() {
		rand.Seed(time.Now().UnixNano())
		
		// Garante que a comida não seja gerada onde a cobra está
		for {
			j.comida = Coordenada{
				x: rand.Intn(j.largura - 4) + 2,
				y: rand.Intn(j.altura - 4) + 2,
			}
			
			// Verifica se a comida não está na posição da cobra
			valida := true
			for _, parte := range j.cobra {
				if parte.x == j.comida.x && parte.y == j.comida.y {
					valida = false
					break
				}
			}
			
			if valida {
				break
			}
		}
	}

	// Desenha o jogo no terminal
	func (j *Jogo) desenhar() {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		// Desenha a borda
		for i := 0; i < j.largura; i++ {
			termbox.SetCell(i, 0, '█', ColorBorder, termbox.ColorDefault)
			termbox.SetCell(i, j.altura-1, '█', ColorBorder, termbox.ColorDefault)
		}
		for i := 0; i < j.altura; i++ {
			termbox.SetCell(0, i, '█', ColorBorder, termbox.ColorDefault)
			termbox.SetCell(j.largura-1, i, '█', ColorBorder, termbox.ColorDefault)
		}

		// Desenha a cobra
		for i, c := range j.cobra {
			char := '█'
			if i == 0 {
				// Cabeça da cobra
				switch j.direcao {
				case Cima:
					char = '▲'
				case Baixo:
					char = '▼'
				case Esquerda:
					char = '◄'
				case Direita:
					char = '►'
				}
			}
			termbox.SetCell(c.x, c.y, char, ColorSnake, termbox.ColorDefault)
		}

		// Desenha a comida
		termbox.SetCell(j.comida.x, j.comida.y, '●', ColorFood, termbox.ColorDefault)

		// Desenha a pontuação e nível
		desenharTexto(2, j.altura+1, fmt.Sprintf("Pontuação: %d", j.pontuacao), ColorText)
		desenharTexto(j.largura-20, j.altura+1, fmt.Sprintf("Nível: %s", j.nivel), ColorText)
		
		// Instruções
		desenharTexto(2, j.altura+3, "Controles: W/↑ (Cima), A/← (Esquerda), S/↓ (Baixo), D/→ (Direita), P (Pausar), Q (Sair)", ColorText)

		// Mensagens de pausa ou game over
		if j.pausado {
			desenharTexto(j.largura/2-10, j.altura/2, "JOGO PAUSADO - Pressione P para continuar", termbox.ColorYellow)
		}
		
		if j.msgGameOver {
			desenharTexto(j.largura/2-15, j.altura/2-2, "GAME OVER! A COBRA BATEU!", termbox.ColorRed)
			desenharTexto(j.largura/2-15, j.altura/2, fmt.Sprintf("Sua pontuação final: %d", j.pontuacao), termbox.ColorYellow)
			desenharTexto(j.largura/2-15, j.altura/2+2, "Pressione ENTER para jogar novamente ou Q para sair", termbox.ColorWhite)
		}

		termbox.Flush()
	}

	// Função auxiliar para desenhar texto
	func desenharTexto(x, y int, texto string, cor termbox.Attribute) {
		for i, char := range texto {
			termbox.SetCell(x+i, y, char, cor, termbox.ColorDefault)
		}
	}

	// Atualiza o estado do jogo
	func (j *Jogo) atualizar() {
		if j.pausado || !j.jogoAtivo {
			return
		}

		// Obtém a posição da cabeça
		cabeca := j.cobra[0]
		
		// Calcula a nova posição da cabeça com base na direção
		var novaCabeca Coordenada
		switch j.direcao {
		case Cima:
			novaCabeca = Coordenada{cabeca.x, cabeca.y - 1}
		case Baixo:
			novaCabeca = Coordenada{cabeca.x, cabeca.y + 1}
		case Esquerda:
			novaCabeca = Coordenada{cabeca.x - 1, cabeca.y}
		case Direita:
			novaCabeca = Coordenada{cabeca.x + 1, cabeca.y}
		}

		// Verificar colisão com a parede
		if novaCabeca.x <= 0 || novaCabeca.x >= j.largura-1 || novaCabeca.y <= 0 || novaCabeca.y >= j.altura-1 {
			j.jogoAtivo = false
			j.msgGameOver = true
			return
		}

		// Verificar colisão com o próprio corpo
		for _, parte := range j.cobra {
			if novaCabeca.x == parte.x && novaCabeca.y == parte.y {
				j.jogoAtivo = false
				j.msgGameOver = true
				return
			}
		}

		// Adiciona a nova cabeça
		j.cobra = append([]Coordenada{novaCabeca}, j.cobra...)

		// Verifica se comeu a comida
		if novaCabeca.x == j.comida.x && novaCabeca.y == j.comida.y {
			j.pontuacao++
			j.gerarComida()
			
			// Aumenta a velocidade a cada 5 pontos
			if j.pontuacao % 5 == 0 && j.velocidade > 30 {
				j.velocidade -= 5
			}
		} else {
			// Remove a última parte da cobra
			j.cobra = j.cobra[:len(j.cobra)-1]
		}
	}

	// Processa as teclas pressionadas
	func (j *Jogo) processarTecla(key termbox.Key, ch rune) {
		if !j.jogoAtivo && j.msgGameOver {
			// Processamento de teclas na tela de game over
			if key == termbox.KeyEnter {
				// Reinicia o jogo
				*j = *NovoJogo(j.largura, j.altura, j.nivel)
			}
			return
		}

		// Teclas durante o jogo
		switch {
		// Teclas de direção (setas)
		case key == termbox.KeyArrowUp && j.direcao != Baixo:
			j.direcao = Cima
		case key == termbox.KeyArrowDown && j.direcao != Cima:
			j.direcao = Baixo
		case key == termbox.KeyArrowLeft && j.direcao != Direita:
			j.direcao = Esquerda
		case key == termbox.KeyArrowRight && j.direcao != Esquerda:
			j.direcao = Direita
			
		// Teclas WASD para direção
		case (ch == 'w' || ch == 'W') && j.direcao != Baixo:
			j.direcao = Cima
		case (ch == 's' || ch == 'S') && j.direcao != Cima:
			j.direcao = Baixo
		case (ch == 'a' || ch == 'A') && j.direcao != Direita:
			j.direcao = Esquerda
		case (ch == 'd' || ch == 'D') && j.direcao != Esquerda:
			j.direcao = Direita
			
		// Pausa o jogo
		case ch == 'p' || ch == 'P':
			j.pausado = !j.pausado
			
		// Sai do jogo
		case ch == 'q' || ch == 'Q':
			termbox.Close()
			fmt.Println("Jogo encerrado. Obrigado por jogar!")
			fmt.Printf("Pontuação final: %d\n", j.pontuacao)
			fmt.Println("Volte sempre!")
			time.Sleep(1 * time.Second)
			panic("saída normal")
		}
	}

	// Tela de menu para selecionar a dificuldade
	func exibirMenu() string {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		
		options := []string{"Fácil", "Médio", "Difícil", "Expert"}
		selected := 1 // Médio é a opção padrão
		
		for {
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			
			// Título
			titulo := "JOGO DA COBRINHA"
			desenharTexto(40-len(titulo)/2, 5, titulo, termbox.ColorGreen | termbox.AttrBold)
			
			// Subtítulo
			subtitulo := "Escolha a dificuldade:"
			desenharTexto(40-len(subtitulo)/2, 8, subtitulo, termbox.ColorWhite)
			
			// Opções
			for i, opt := range options {
				cor := termbox.ColorWhite
				prefixo := "  "
				
				if i == selected {
					cor = termbox.ColorGreen
					prefixo = "> "
				}
				
				texto := fmt.Sprintf("%s%s", prefixo, opt)
				desenharTexto(40-len(texto)/2, 10+i*2, texto, cor)
			}
			
			// Instruções
			instrucoes := "Use as setas ↑/↓ para selecionar e ENTER para confirmar"
			desenharTexto(40-len(instrucoes)/2, 20, instrucoes, termbox.ColorYellow)
			
			// Mensagem de boas-vindas brasileira
			msgBr := "Bora jogar esse clássico, meu consagrado!"
			desenharTexto(40-len(msgBr)/2, 22, msgBr, termbox.ColorCyan)
			
			termbox.Flush()
			
			// Processamento de teclas
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyArrowUp:
					selected = (selected - 1 + len(options)) % len(options)
				case termbox.KeyArrowDown:
					selected = (selected + 1) % len(options)
				case termbox.KeyEnter:
					return options[selected]
				case termbox.KeyEsc:
					termbox.Close()
					panic("saída normal")
				}
			}
		}
	}

	func main() {
		// Inicializa o termbox
		err := termbox.Init()
		if err != nil {
			panic(err)
		}
		defer termbox.Close()

		// Configura o terminal
		termbox.SetInputMode(termbox.InputEsc)

		// Obtém o tamanho do terminal
		largura, altura := termbox.Size()

		// Exibe o menu e obtém o nível de dificuldade escolhido
		nivelEscolhido := exibirMenu()

		// Cria um novo jogo com o tamanho ajustado para o terminal
		jogo := NovoJogo(largura, altura-5, nivelEscolhido)

		// Cria canais para eventos e atualização do jogo
		eventos := make(chan termbox.Event)
		atualizacao := time.NewTicker(time.Millisecond * jogo.velocidade)
		defer atualizacao.Stop()

		// Inicia a goroutine para capturar eventos do teclado
		go func() {
			for {
				eventos <- termbox.PollEvent()
			}
		}()

		// Loop principal do jogo
		for {
			// Desenha o estado atual do jogo
			jogo.desenhar()

			// Seleciona entre eventos do teclado e atualizações do jogo
			select {
			case ev := <-eventos:
				// Processa eventos do teclado
				if ev.Type == termbox.EventKey {
					// Se pressionar ESC, sai do jogo
					if ev.Key == termbox.KeyEsc {
						termbox.Close()
						fmt.Println("Jogo encerrado. Obrigado por jogar!")
						return
					}
					
					// Processa outras teclas
					jogo.processarTecla(ev.Key, ev.Ch)
				}

			case <-atualizacao.C:
				// Atualiza o estado do jogo
				jogo.atualizar()
				
				// Se a velocidade mudou, atualiza o ticker
				if !jogo.pausado && jogo.jogoAtivo {
					atualizacao.Reset(time.Millisecond * jogo.velocidade)
				}
			}
		}
		}
