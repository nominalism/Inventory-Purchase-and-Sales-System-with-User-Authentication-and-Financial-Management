<script>
	import { browser } from '$app/environment';
	import BarraLateral from '$lib/components/BarraLateral.svelte';
	async function handleSubmit(event) {
		event.preventDefault(); // Prevent default form submission

		const data = {
			ID: document.getElementById('ID').value,
			codigo: document.getElementById('codigo').value,
			nome: document.getElementById('nome').value,
			local: document.getElementById('local').value,
			dataCompra: document.getElementById('dataCompra').value,
			valorCompra: document.getElementById('valorCompra').value,
			valorVenda: document.getElementById('valorVenda').value,
			estoque: document.getElementById('estoque').value,
			estoqueCritico: document.getElementById('estoqueCritico').value,
			fornecedor: document.getElementById('fornecedor').value
		};

		try {
			const response = await fetch('http://localhost:5000/add_material', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(data)
			});

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			const result = await response.json();
			alert(result.message);
		} catch (error) {
			console.error('Error submitting form:', error);
			alert('Erro ao enviar o formulário: ' + error.message);
		}
	}

	// Use onMount instead of DOMContentLoaded
	import { onMount } from 'svelte';

	onMount(() => {
		if (browser) {
			const form = document.getElementById('materialForm');
			if (form) {
				form.addEventListener('submit', handleSubmit);
			}
		}
	});
</script>

<div class="container">
	<BarraLateral />
	<div class="content">
		<div class="form-container">
			<h1>Registrar Material</h1>
			<form id="materialForm" onsubmit={handleSubmit}>
				<div class="form-group">
					<label for="ID">ID:</label>
					<input type="text" id="ID" placeholder="Digite o ID do material" required />
				</div>

				<div class="form-group">
					<label for="codigo">Código:</label>
					<input type="text" id="codigo" placeholder="Digite o código do material" required />
				</div>

				<div class="form-group">
					<label for="nome">Nome:</label>
					<input type="text" id="nome" placeholder="Digite o nome do material" required />
				</div>

				<div class="form-group">
					<label for="local">Local:</label>
					<input type="text" id="local" placeholder="Digite o local de armazenamento" required />
				</div>

				<div class="form-group">
					<label for="dataCompra">Data de Compra:</label>
					<input type="date" id="dataCompra" required />
				</div>

				<div class="form-group">
					<label for="valorCompra">Valor de Compra:</label>
					<input
						type="number"
						id="valorCompra"
						placeholder="Digite o valor de compra do material"
						step="0.01"
						required
					/>
				</div>

				<div class="form-group">
					<label for="valorVenda">Valor de Venda:</label>
					<input
						type="number"
						id="valorVenda"
						placeholder="Digite o valor de venda do material"
						step="0.01"
						required
					/>
				</div>

				<div class="form-group">
					<label for="estoque">Estoque:</label>
					<input type="number" id="estoque" placeholder="Digite a quantidade em estoque" required />
				</div>

				<div class="form-group">
					<label for="estoqueCritico">Estoque Crítico:</label>
					<input
						type="number"
						id="estoqueCritico"
						placeholder="Digite o estoque crítico"
						required
					/>
				</div>

				<div class="form-group">
					<label for="fornecedor">Fornecedor:</label>
					<input type="text" id="fornecedor" placeholder="Digite o nome do fornecedor" required />
				</div>

				<div class="button-group">
					<button type="submit" class="save-btn">Registrar Material</button>
					<button type="button" class="cancel-btn" onclick={() => goto('/material')}
						>Cancelar</button
					>
				</div>
			</form>
		</div>
	</div>
</div>

<style>
	.container {
		display: flex;
		height: 100vh;
		background-color: #f5f5f5;
		overflow: hidden; /* Previne scroll na página inteira */
	}

	.content {
		flex: 1;
		padding: 20px;
		overflow-y: auto; /* Permite scroll apenas no conteúdo */
	}

	.form-container {
		background: white;
		padding: 30px;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
		max-width: 800px;
		margin: 0 auto;
	}

	.form-group {
		margin-bottom: 15px;
	}

	label {
		display: block;
		margin-bottom: 8px;
		font-weight: 500;
		color: #333;
	}

	input {
		width: 100%;
		padding: 10px;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 14px;
	}

	.button-group {
		display: flex;
		gap: 10px;
		margin-top: 20px;
	}

	button {
		flex: 1;
		padding: 12px;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 14px;
	}

	.save-btn {
		background-color: #4caf50;
		color: white;
	}

	.cancel-btn {
		background-color: #f44336;
		color: white;
	}

	h1 {
		text-align: center;
		color: #333;
		margin-bottom: 30px;
		font-size: 24px;
	}
</style>
