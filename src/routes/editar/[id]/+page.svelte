<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import BarraLateral from '$lib/components/BarraLateral.svelte';

	let material = $state({
		codigo: '',
		nome: '',
		local: '',
		dataCompra: '',
		valorCompra: 0, // Changed from string to number
		valorVenda: 0, // Changed from string to number
		estoque: 0, // Changed from string to number
		estoqueCritico: 0, // Changed from string to number
		fornecedor: ''
	});

	let error = $state('');
	let success = $state('');

	let id = $derived($page.params.id);

	onMount(async () => {
		if (id) {
			try {
				const user = JSON.parse(localStorage.getItem('user'));
				console.log('Fetching material with ID:', id);
				const response = await fetch(`http://localhost:5000/busca_material_id/${id}`, {
					headers: {
						Authorization: user.token
					}
				});
				console.log('Response status:', response.status);

				if (response.ok) {
					const data = await response.json();
					console.log('Received data:', data);

					const rawMaterial = data.materiais[0];
					console.log('Raw material:', rawMaterial);

					// Keep everything as string to match Go struct
					material = {
						id: String(rawMaterial.id || ''),
						codigo: String(rawMaterial.codigo || ''),
						nome: String(rawMaterial.nome || ''),
						local: String(rawMaterial.local || ''),
						dataCompra: String(rawMaterial.dataCompra || ''),
						valorCompra: Number(rawMaterial.valorCompra || 0),
						valorVenda: Number(rawMaterial.valorVenda || 0),
						estoque: Number(rawMaterial.estoque || 0),
						estoqueCritico: Number(rawMaterial.estoqueCritico || 0),
						fornecedor: String(rawMaterial.fornecedor || '')
					};
					console.log('Transformed material:', material);
				} else {
					console.error('Response not OK:', await response.text());
					error = 'Erro ao carregar material';
				}
			} catch (err) {
				console.error('Fetch error:', err);
				error = 'Erro ao conectar com o servidor';
			}
		}
	});

	async function handleSubmit(event) {
		event.preventDefault();
		error = '';
		success = '';

		try {
			const user = JSON.parse(localStorage.getItem('user'));
			// No need for conversion since the values are already numbers
			const materialToSend = {
				...material
			};

			const response = await fetch(`http://localhost:5000/admin/atualizar_material/${id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: user.token
				},
				body: JSON.stringify(materialToSend)
			});

			if (response.ok) {
				success = 'Material atualizado com sucesso!';
				setTimeout(() => goto('/material'), 1500);
			} else {
				const data = await response.json();
				error = data.error || 'Erro ao atualizar material';
			}
		} catch (err) {
			console.error('Submit error:', err);
			error = 'Erro ao conectar com o servidor';
		}
	}
</script>

<div class="container">
	<BarraLateral class="w-64 bg-gray-800 text-white p-4" />

	<div class="content">
		{#if error}
			<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
				{error}
			</div>
		{/if}

		{#if success}
			<div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
				{success}
			</div>
		{/if}

		<div class="form-container">
			<h1>Editar Material</h1>
			<form onsubmit={handleSubmit} class="formulario">
				<div class="form-group">
					<label for="codigo" class="block mb-1">Código:</label>
					<input
						type="text"
						id="codigo"
						bind:value={material.codigo}
						placeholder={material.codigo || ''}
						class="w-full p-2 border rounded"
						required
					/>
				</div>

				<div class="form-group">
					<label for="nome" class="block mb-1">Nome:</label>
					<input
						type="text"
						id="nome"
						bind:value={material.nome}
						placeholder={material.nome || ''}
						class="w-full p-2 border rounded"
						required
					/>
				</div>

				<div class="form-group">
					<label for="local" class="block mb-1">Local:</label>
					<input
						type="text"
						id="local"
						bind:value={material.local}
						placeholder={material.local || ''}
						class="w-full p-2 border rounded"
						required
					/>
				</div>

				<div class="form-group">
					<label for="dataCompra" class="block mb-1">Data de Compra:</label>
					<input
						type="date"
						id="dataCompra"
						bind:value={material.dataCompra}
						class="w-full p-2 border rounded"
						required
					/>
				</div>

				<div class="form-group">
					<label for="valorCompra" class="block mb-1">Valor de Compra:</label>
					<input
						type="number"
						id="valorCompra"
						bind:value={material.valorCompra}
						placeholder={material.valorCompra || ''}
						step="0.01"
						class="w-full p-2 border rounded"
						required
					/>
				</div>

				<div class="form-group">
					<label for="valorVenda" class="block mb-1">Valor de Venda:</label>
					<input
						type="number"
						id="valorVenda"
						bind:value={material.valorVenda}
						placeholder={material.valorVenda || ''}
						step="0.01"
						class="w-full p-2 border rounded"
						required
					/>
				</div>

				<div class="form-group">
					<label for="estoque" class="block mb-1">Estoque:</label>
					<input
						type="number"
						id="estoque"
						bind:value={material.estoque}
						placeholder={material.estoque || ''}
						class="w-full p-2 border rounded"
						required
					/>
				</div>

				<div class="form-group">
					<label for="estoqueCritico" class="block mb-1">Estoque Crítico:</label>
					<input
						type="number"
						id="estoqueCritico"
						bind:value={material.estoqueCritico}
						placeholder={material.estoqueCritico || ''}
						class="w-full p-2 border rounded"
						required
					/>
				</div>

				<div class="form-group">
					<label for="fornecedor" class="block mb-1">Fornecedor:</label>
					<input
						type="text"
						id="fornecedor"
						bind:value={material.fornecedor}
						placeholder={material.fornecedor || ''}
						class="w-full p-2 border rounded"
						required
					/>
				</div>

				<div class="button-group">
					<button type="submit" class="save-btn"> Salvar Alterações </button>
					<button type="button" onclick={() => goto('/material')} class="cancel-btn">
						Cancelar
					</button>
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
		margin-bottom: 20px;
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
		font-size: 16px;
	}

	input:focus {
		border-color: #4caf50;
		outline: none;
	}

	.button-group {
		display: flex;
		gap: 10px;
		margin-top: 20px;
	}

	button {
		flex: 1;
		padding: 12px 20px;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 16px;
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
	}
</style>
