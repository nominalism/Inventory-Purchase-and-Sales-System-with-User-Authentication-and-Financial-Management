<script>
	import { preventDefault } from 'svelte/legacy';

	import { goto } from '$app/navigation';
	import BarraLateral from '$lib/components/BarraLateral.svelte';

	let material = $state({
		codigo: '',
		nome: '',
		local: '',
		dataCompra: '',
		valorCompra: '',
		valorVenda: '',
		estoque: '',
		estoqueCritico: '',
		fornecedor: ''
	});

	let error = $state('');
	let success = $state('');

	async function handleSubmit() {
		try {
			const user = JSON.parse(localStorage.getItem('user'));
			// Converter os valores numéricos antes de enviar
			const materialToSend = {
				...material,
				valorCompra: parseFloat(material.valorCompra),
				valorVenda: parseFloat(material.valorVenda),
				estoque: parseInt(material.estoque),
				estoqueCritico: parseInt(material.estoqueCritico)
			};

			const response = await fetch('http://localhost:5000/admin/add_material', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: user.token
				},
				body: JSON.stringify(materialToSend)
			});

			if (response.ok) {
				success = 'Material registrado com sucesso!';
				setTimeout(() => goto('/material'), 1500);
			} else {
				const data = await response.json();
				error = data.error || 'Erro ao registrar material';
			}
		} catch (err) {
			error = 'Erro ao conectar com o servidor';
			console.error(err);
		}
	}
</script>

<div class="page">
	<BarraLateral />
	<div class="content">
		<h1>Registrar Novo Material</h1>

		{#if error}
			<div class="error">{error}</div>
		{/if}
		{#if success}
			<div class="success">{success}</div>
		{/if}

		<form onsubmit={preventDefault(handleSubmit)} class="form-container">
			<div class="input-group">
				<label for="codigo">Código:</label>
				<input type="text" id="codigo" bind:value={material.codigo} required />
			</div>

			<div class="input-group">
				<label for="nome">Nome:</label>
				<input type="text" id="nome" bind:value={material.nome} required />
			</div>

			<div class="input-group">
				<label for="local">Local:</label>
				<input type="text" id="local" bind:value={material.local} required />
			</div>

			<div class="input-group">
				<label for="dataCompra">Data da Compra:</label>
				<input type="date" id="dataCompra" bind:value={material.dataCompra} required />
			</div>

			<div class="input-group">
				<label for="valorCompra">Valor de Compra:</label>
				<input
					type="number"
					step="0.01"
					id="valorCompra"
					bind:value={material.valorCompra}
					required
				/>
			</div>

			<div class="input-group">
				<label for="valorVenda">Valor de Venda:</label>
				<input
					type="number"
					step="0.01"
					id="valorVenda"
					bind:value={material.valorVenda}
					required
				/>
			</div>

			<div class="input-group">
				<label for="estoque">Estoque:</label>
				<input type="number" id="estoque" bind:value={material.estoque} required />
			</div>

			<div class="input-group">
				<label for="estoqueCritico">Estoque Crítico:</label>
				<input type="number" id="estoqueCritico" bind:value={material.estoqueCritico} required />
			</div>

			<div class="input-group">
				<label for="fornecedor">Fornecedor:</label>
				<input type="text" id="fornecedor" bind:value={material.fornecedor} required />
			</div>

			<button type="submit">Registrar Material</button>
		</form>
	</div>
</div>

<style>
	.page {
		display: flex;
		min-height: 100vh;
	}

	.content {
		flex: 1;
		padding: 20px;
	}

	h1 {
		margin-bottom: 20px;
		color: #333;
	}

	.form-container {
		max-width: 600px;
		background: white;
		padding: 20px;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.input-group {
		margin-bottom: 15px;
	}

	label {
		display: block;
		margin-bottom: 5px;
		color: #333;
	}

	input {
		width: 100%;
		padding: 8px;
		border: 1px solid #ddd;
		border-radius: 4px;
	}

	button {
		background: #4caf50;
		color: white;
		padding: 10px 20px;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		width: 100%;
		margin-top: 10px;
	}

	button:hover {
		background: #45a049;
	}

	.error {
		color: red;
		margin-bottom: 10px;
	}

	.success {
		color: green;
		margin-bottom: 10px;
	}
</style>
