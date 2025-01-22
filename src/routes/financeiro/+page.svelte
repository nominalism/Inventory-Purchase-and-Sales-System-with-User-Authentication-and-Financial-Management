<script>
	import { onMount } from 'svelte';
	import Chart from 'chart.js/auto';
	import BarraLateral from '$lib/components/BarraLateral.svelte';

	let caixa = $state(0);
	let canvas = $state();
	let chart;

	onMount(async () => {
		await buscarCaixa();
		await criarGrafico();
	});

	async function buscarCaixa() {
		try {
			const user = JSON.parse(localStorage.getItem('user'));
			const response = await fetch('http://localhost:5000/admin/caixa', {
				headers: {
					Authorization: user.token
				}
			});

			if (!response.ok) {
				throw new Error('Erro ao buscar saldo');
			}

			const data = await response.json();
			caixa = data.saldo;
		} catch (error) {
			console.error('Erro ao buscar caixa:', error);
			alert('Erro ao buscar saldo do caixa');
		}
	}

	async function atualizarCaixa(operacao) {
		const valor = parseFloat(prompt(`Valor para ${operacao}:`));
		if (!valor || isNaN(valor)) return;

		try {
			const user = JSON.parse(localStorage.getItem('user'));
			const response = await fetch('http://localhost:5000/admin/atualizar_caixa', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: user.token
				},
				body: JSON.stringify({
					valor: operacao === 'adicionar' ? valor : -valor
				})
			});

			if (!response.ok) throw new Error('Erro ao atualizar caixa');
			await buscarCaixa();
			await criarGrafico();
			alert('Caixa atualizado com sucesso!');
		} catch (error) {
			alert('Erro: ' + error.message);
		}
	}

	async function criarGrafico() {
		try {
			const user = JSON.parse(localStorage.getItem('user'));
			const response = await fetch('http://localhost:5000/lucros_mensais', {
				headers: {
					Authorization: user.token
				}
			});

			if (!response.ok) {
				throw new Error('Erro ao buscar dados do gráfico');
			}

			const data = await response.json();

			if (chart) chart.destroy();

			chart = new Chart(canvas, {
				type: 'line',
				data: {
					labels: data.meses,
					datasets: [
						{
							label: 'Lucro Mensal',
							data: data.valores,
							borderColor: 'rgb(75, 192, 192)',
							tension: 0.1
						}
					]
				},
				options: {
					responsive: true,
					scales: {
						y: {
							beginAtZero: true,
							title: {
								display: true,
								text: 'Lucro (R$)'
							}
						}
					}
				}
			});
		} catch (error) {
			console.error('Erro ao criar gráfico:', error);
		}
	}
</script>

<div class="container">
	<BarraLateral />
	<div class="content">
		<h1>Financeiro</h1>
		<h2>Caixa Atual: R$ {caixa.toFixed(2)}</h2>

		<div class="buttons">
			<button onclick={() => atualizarCaixa('adicionar')}>Adicionar ao Caixa</button>
			<button onclick={() => atualizarCaixa('retirar')}>Retirar do Caixa</button>
		</div>

		<div class="graph">
			<canvas bind:this={canvas}></canvas>
		</div>
	</div>
</div>

<style>
	.container {
		display: flex;
		height: 100vh;
	}

	.content {
		flex: 1;
		padding: 20px;
	}

	.buttons {
		margin: 20px 0;
	}

	.graph {
		margin-top: 20px;
		height: 400px;
	}

	button {
		margin: 0 10px;
		padding: 10px 20px;
	}
</style>
